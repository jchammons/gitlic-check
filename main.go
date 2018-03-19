package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Follow instructions in README for proper authentication
func getTokens() []string {
	fh, err := ioutil.ReadFile("auth.txt")
	if err != nil {
		fmt.Printf("Failed to read PAT from auth.txt: %s", err)
	}
	tokens := strings.Split(string(fh), ",")
	return tokens
}

func getIgnoredOrgs() map[string]bool {
	fh, err := ioutil.ReadFile("ignore.txt")
	if err != nil {
		return nil
	}
	orgNames := strings.Split(string(fh), ",")
	ignoredOrgs := make(map[string]bool)
	for _, org := range orgNames {
		ignoredOrgs[org] = true
	}
	return ignoredOrgs
}

func prepareOutput() (*os.File, *os.File, *os.File) {
	os.RemoveAll("output")
	os.Mkdir("output", 0777)
	os.Chdir("output")
	fRepos, err := os.OpenFile("repos.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Couldn't create repos.csv. Failed with %s", err)
	}
	fUsers, err := os.OpenFile("users.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Couldn't create users.csv. Failed with %s", err)
	}
	fInvites, err := os.OpenFile("invites.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Couldn't create invites.csv. Failed with %s", err)
	}
	return fRepos, fUsers, fInvites
}

func main() {
	ctx := context.Background()
	ignoredOrgs := getIgnoredOrgs()
	auth := getTokens()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: auth[0],
		},
	)
	authClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(authClient)
	lo := &github.ListOptions{PerPage: 100}

	fmt.Print("Working...\n\n")

	fRepos, fUsers, fInvites := prepareOutput()

	if _, err := fRepos.WriteString("Org,Repo,Private,Fork,License\n"); err != nil {
		fmt.Printf("Initial save to repos CSV failed with %s\n", err)
	}

	if _, err := fUsers.WriteString("Org,User,Two-Factor Enabled\n"); err != nil {
		fmt.Printf("Initial save to users CSV failed with %s\n", err)
	}

	if _, err := fInvites.WriteString("Org,User,Date Sent,Invited By\n"); err != nil {
		fmt.Printf("Initial save to invites CSV failed with %s\n", err)
	}

	var orgs []*github.Organization

	for {
		partialOrgs, resp, err := client.Organizations.List(ctx, "", lo)
		if err != nil {
			fmt.Printf("Organizations.List failed with %s\n", err)
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
			fmt.Printf("Ignored %s, %d of %d\n", *org.Login, i+1, len(orgs))
			continue
		}

		var repos []*github.Repository
		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		var invites []*github.Invitation

		for {
			partialInvites, resp, err := client.Organizations.ListPendingOrgInvitations(ctx, *org.Login, lo)
			if err != nil {
				fmt.Printf("Repositories.List failed with %s\n", err)
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
				// TODO Check invite date; if older than 2 weeks, call Organization
				// Cancel Membership endpoint with user id
				inviteDate := fmt.Sprint(invite.CreatedAt.Year(), "-", invite.CreatedAt.Day(), "-", invite.CreatedAt.Month())
				if _, err := fInvites.WriteString(fmt.Sprint(*org.Login, ",", *invite.Login, ",", inviteDate, ",", *invite.Inviter.Login, "\n")); err != nil {
					fmt.Printf("Failed to write invite for %s\n", *invite.Login)
				}
			}
		}

		for {
			partialRepos, resp, err := client.Repositories.ListByOrg(ctx, *org.Login, opt)
			if err != nil {
				fmt.Printf("Repositories.List failed with %s\n", err)
				return
			}

			repos = append(repos, partialRepos...)

			if resp.NextPage == 0 {
				opt.Page = 1
				break
			}
			opt.Page = resp.NextPage
		}

		for _, repo := range repos {
			if _, err := fRepos.WriteString(fmt.Sprint(*org.Login, ",", *repo.Name, ",", *repo.Private, ",", *repo.Fork, ",")); err != nil {
				fmt.Printf("Failed to write to repos.csv on %s with %s\n", *org.Login, err)
			}

			lics, _, err := client.Repositories.License(ctx, *repo.Owner.Login, *repo.Name)
			if err != nil {
				fRepos.WriteString("None\n")
				continue
			}

			if _, err = fRepos.WriteString(fmt.Sprint(*lics.License.Name, "\n")); err != nil {
				fmt.Printf("Failed to write to repos.csv on %s with %s\n", *org.Login, err)
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
			partialMembers, resp, err := client.Organizations.ListMembers(ctx, *org.Login, memOpt)
			if err != nil {
				fmt.Printf("Organizations.ListMembers, no filter, failed with %s\n", err)
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
			partialNo2f, resp, err := client.Organizations.ListMembers(ctx, *org.Login, no2fOpt)
			if err != nil {
				fmt.Printf("Organizations.ListMembers, 2FA filter, failed with %s\n", err)
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
			if _, err := fUsers.WriteString(fmt.Sprint(*org.Login, ",", *member.Login, ",")); err != nil {
				fmt.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
			}

			if membersFilter[*member.Login] {
				if _, err := fUsers.WriteString("False\n"); err != nil {
					fmt.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
				}
				continue
			}

			if _, err := fUsers.WriteString("True\n"); err != nil {
				fmt.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
			}
		}

		fmt.Printf("Completed %d of %d\n", i+1, len(orgs))
	}

	if err := fRepos.Close(); err != nil {
		log.Fatal(err)
	}
	if err := fUsers.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Print("CSVs are now ready\n")
}
