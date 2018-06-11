package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// getIgnoredMap will turn the array of ignored orgs from the config file into a map that is easy to check against as we iterate over organizations returned from the GitHub API
func getIgnoredMap(a []string) map[string]bool {
	ignoredOrgs := make(map[string]bool)
	for _, org := range a {
		ignoredOrgs[org] = true
	}
	return ignoredOrgs
}

// RunCheck begins the process of querying the GitHub API. It will loop through your organizations and their repositories and pull info on configuration, license, and users, including invitations. It will output the results to respective CSV files in the output folder. See the README for an idea of what these CSV reports contain.
func RunCheck(ctx context.Context, cf config, fo map[string]*os.File) {
	ghClient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cf.Github.Token})))
	lo := &github.ListOptions{PerPage: 100}
	maxInviteT := time.Duration(cf.Github.RmInvitesAfter) * time.Hour
	ignoredOrgs := getIgnoredMap(cf.Github.IgnoredOrgs)

	log.Print("Working...\n\n")

	if _, err := fo["repos.csv"].WriteString("Org,Repo,Private,Fork,License\n"); err != nil {
		log.Printf("Initial save to repos CSV failed with %s\n", err)
	}

	if _, err := fo["users.csv"].WriteString("Org,User,Two-Factor Enabled\n"); err != nil {
		log.Printf("Initial save to users CSV failed with %s\n", err)
	}

	if _, err := fo["invites.csv"].WriteString("Org,User,Date Sent,Invited By,Deleted\n"); err != nil {
		log.Printf("Initial save to invites CSV failed with %s\n", err)
	}

	var orgs []*github.Organization

	for {
		partialOrgs, resp, err := ghClient.Organizations.List(ctx, "", lo)
		if err != nil {
			log.Fatalf("Organizations.List failed with %s\n", err)
			return
		}

		orgs = append(orgs, partialOrgs...)

		if resp.NextPage == 0 {
			lo.Page = 1
			break
		}
		lo.Page = resp.NextPage
	}

	for i, org := range orgs {
		if ignoredOrgs != nil && ignoredOrgs[*org.Login] {
			log.Printf("Ignored %s, %d of %d\n", *org.Login, i+1, len(orgs))
			continue
		}

		var invites []*github.Invitation
		for {
			partialInvites, resp, err := ghClient.Organizations.ListPendingOrgInvitations(ctx, *org.Login, lo)
			if err != nil {
				log.Printf("Organizations.ListPendingOrgInvitations failed with %s\n", err)
				return
			}

			invites = append(invites, partialInvites...)

			if resp.NextPage == 0 {
				lo.Page = 1
				break
			}
			lo.Page = resp.NextPage
		}

		if len(invites) > 0 {
			for _, invite := range invites {
				inviteDate := fmt.Sprint(int(invite.CreatedAt.Month()), "/", fmt.Sprintf("%02d", invite.CreatedAt.Day()), "/", invite.CreatedAt.Year())
				tSinceInvite := time.Now().UTC().Sub(invite.CreatedAt.UTC())
				if _, err := fo["invites.csv"].WriteString(fmt.Sprint(*org.Login, ",", *invite.Login, ",", inviteDate, ",", *invite.Inviter.Login)); err != nil {
					log.Printf("Failed to write invite data for %s from %s to invite.csv\n", *invite.Login, *org.Login)
				}
				if cf.Github.RmInvitesAfter != 0 && tSinceInvite > maxInviteT {
					_, err := ghClient.Organizations.RemoveOrgMembership(ctx, *invite.Login, *org.Login)
					if err != nil {
						log.Printf("Failed to remove flagged pending invitation for %s from org %s\n", *invite.Login, *org.Login)
					} else {
						fo["invites.csv"].WriteString(",True\n")
					}
					continue
				}
				fo["invites.csv"].WriteString(",\n")
			}
		}

		var repos []*github.Repository
		repoOpt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			partialRepos, resp, err := ghClient.Repositories.ListByOrg(ctx, *org.Login, repoOpt)
			if err != nil {
				log.Fatalf("Repositories.ListByOrg failed with %s\n", err)
				return
			}

			repos = append(repos, partialRepos...)

			if resp.NextPage == 0 {
				repoOpt.Page = 1
				break
			}
			repoOpt.Page = resp.NextPage
		}

		for _, repo := range repos {
			if _, err := fo["repos.csv"].WriteString(fmt.Sprint(*org.Login, ",", *repo.Name, ",", *repo.Private, ",", *repo.Fork, ",")); err != nil {
				log.Printf("Failed to write to repos.csv on %s with %s\n", *org.Login, err)
			}

			lics, _, err := ghClient.Repositories.License(ctx, *repo.Owner.Login, *repo.Name)
			if err != nil {
				fo["repos.csv"].WriteString("None\n")
				continue
			}

			if _, err = fo["repos.csv"].WriteString(fmt.Sprint(*lics.License.Name, "\n")); err != nil {
				log.Printf("Failed to write to repos.csv on %s with %s\n", *org.Login, err)
			}
		}

		var members []*github.User
		memOpt := &github.ListMembersOptions{
			ListOptions: *lo,
		}

		var membersNo2f []*github.User
		no2fOpt := &github.ListMembersOptions{
			Filter:      "2fa_disabled",
			ListOptions: *lo,
		}

		for {
			partialMembers, resp, err := ghClient.Organizations.ListMembers(ctx, *org.Login, memOpt)
			if err != nil {
				log.Printf("Organizations.ListMembers, no filter, failed with %s\n", err)
				break
			}

			members = append(members, partialMembers...)

			if resp.NextPage == 0 {
				memOpt.Page = 1
				break
			}
			memOpt.Page = resp.NextPage
		}

		for {
			partialNo2f, resp, err := ghClient.Organizations.ListMembers(ctx, *org.Login, no2fOpt)
			if err != nil {
				log.Printf("Organizations.ListMembers, 2FA filter, failed with %s\n", err)
				break
			}

			membersNo2f = append(membersNo2f, partialNo2f...)

			if resp.NextPage == 0 {
				no2fOpt.Page = 1
				break
			}
			no2fOpt.Page = resp.NextPage
		}

		membersFilter := make(map[string]bool)
		for _, member := range membersNo2f {
			membersFilter[*member.Login] = true
		}

		for _, member := range members {
			if _, err := fo["users.csv"].WriteString(fmt.Sprint(*org.Login, ",", *member.Login, ",")); err != nil {
				log.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
			}

			if membersFilter[*member.Login] {
				if _, err := fo["users.csv"].WriteString("False\n"); err != nil {
					log.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
				}
				continue
			}

			if _, err := fo["users.csv"].WriteString("True\n"); err != nil {
				log.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
			}
		}

		log.Printf("Completed %s, %d of %d\n", *org.Login, i+1, len(orgs))
	}
	log.Print("\nCSVs are ready!\n")
	return
}
