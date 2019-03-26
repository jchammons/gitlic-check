package swgithub

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/solarwinds/gitlic-check/config"
)

// getIncludedMap will turn the array of included orgs from the config file into a map that is easy to check against as we iterate over organizations returned from the GitHub API
func getIncludedMap(a []string) map[string]bool {
	includedOrgs := make(map[string]bool)
	for _, org := range a {
		includedOrgs[strings.ToLower(org)] = true
	}
	return includedOrgs
}

// RunGitlicCheck begins the process of querying the GitHub API. It will loop through your organizations and their repositories and pull info on configuration, license, and users, including invitations. It will output the results to respective CSV files in the output folder. See the README for an idea of what these CSV reports contain.
func RunGitlicCheck(ctx context.Context, cf config.Config, fo map[string]*os.File) {
	ghClient := GetGHClient(ctx, cf)
	lo := &github.ListOptions{PerPage: 100}
	maxInviteT := time.Duration(cf.Github.RmInvitesAfter) * time.Hour

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

	orgs, err := GetSWOrgs(ctx, ghClient, cf)
	if err != nil {
		log.Fatalf("Failed to get orgs: %s", err.Error())
		return
	}

	for i, org := range orgs {
		invites, err := GetOrgInvites(ctx, ghClient, org)
		if err != nil {
			log.Printf("Organizations.ListPendingOrgInvitations failed with %s\n", err)
			return
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

		repos, err := GetOrgRepositories(ctx, ghClient, org)
		if err != nil {
			log.Fatalf("Repositories.ListByOrg failed with %s\n", err)
			return
		}

		for _, repo := range repos {
			license := "None"
			if repo.License.GetName() != "" {
				license = repo.License.GetName()
			}
			if _, err := fo["repos.csv"].WriteString(fmt.Sprint(*org.Login, ",", *repo.Name, ",", *repo.Private, ",", *repo.Fork, ",", license, "\n")); err != nil {
				log.Printf("Failed to write to repos.csv on %s with %s\n", *org.Login, err)
			}
		}

		memOpt := &github.ListMembersOptions{ListOptions: *lo}
		members, err := GetOrgMembers(ctx, ghClient, org, memOpt)
		if err != nil {
			log.Printf("Couldn't get org members, no filter, for %s: %s", *org.Login, err.Error())
		}

		no2fOpt := &github.ListMembersOptions{
			Filter:      "2fa_disabled",
			ListOptions: *lo,
		}
		membersNo2f, err := GetOrgMembers(ctx, ghClient, org, no2fOpt)
		if err != nil {
			log.Printf("Couldn't get org members, 2fa filter, for %s: %s", *org.Login, err.Error())
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

// RunGitlicStatsReport will run a report across all repositories for the specified orgs that provides data on commits and lines of code
func RunGitlicStatsReport(ctx context.Context, cf config.Config, fo map[string]*os.File) {
	weekdays = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	ghClient := GetGHClient(ctx, cf)

	if _, err := fo["stats.csv"].WriteString("Org,Repo,Commits in Last Year,Most Commits in a Day,Most Commits in a Week,Most Committable Day,Total Lines Added,Total Lines Removed,Most Lines Added in a Week,Most Lines Removed in a Week\n"); err != nil {
		log.Printf("Initial save to stats CSV failed with %s\n", err)
	}

	log.Print("Working...\n\n")

	orgs, err := GetSWOrgs(ctx, ghClient, cf)
	if err != nil {
		log.Fatalf("Failed to get orgs: %s", err.Error())
		return
	}

	for i, org := range orgs {
		repos, err := GetOrgRepositories(ctx, ghClient, org)
		if err != nil {
			log.Fatalf("Repositories.ListByOrg failed with %s\n", err)
			return
		}

		var retries []string
		commitsMap := make(map[string][]*github.WeeklyCommitActivity)
		modsMap := make(map[string][]*github.WeeklyStats)

		for _, repo := range repos {
			commits, err, cmRetry := GetCommitActivity(ctx, ghClient, *org.Login, *repo.Name)
			if err != nil {
				log.Printf("Failed to get commit activity for %s. Error: %v\n", *repo.Name, err)
			} else {
				mods, err, mdRetry := GetAdditionsDeletions(ctx, ghClient, *org.Login, *repo.Name)
				if err != nil {
					log.Printf("Failed to get additions/deletions for %s. Error: %v\n", *repo.Name, err)
				}

				if cmRetry || mdRetry {
					retries = append(retries, *repo.Name)

					if !cmRetry {
						commitsMap[*repo.Name] = commits
					}
					if !mdRetry {
						modsMap[*repo.Name] = mods
					}
					continue
				}

				if _, err := fo["stats.csv"].WriteString(formatRepoStats(*org.Login, *repo.Name, commits, mods)); err != nil {
					log.Printf("Failed to write stats for %s repository in %s org to stats.csv. Error: %v\n", *repo.Name, *org.Login, err)
				}
			}
		}

		// If the stats have not been calculated, wait, then try again
		if len(retries) > 0 {
			time.Sleep(3000 * time.Millisecond)

			for _, repo := range retries {
				var commits []*github.WeeklyCommitActivity
				var mods []*github.WeeklyStats
				var err error

				if commitsMap[repo] == nil {
					commits, err, _ = GetCommitActivity(ctx, ghClient, *org.Login, repo)
					if err != nil {
						log.Printf("Failed to get commit activity for %s. Error: %v\n", repo, err)
						continue
					}
				} else {
					commits = commitsMap[repo]
				}

				if modsMap[repo] == nil {
					mods, err, _ = GetAdditionsDeletions(ctx, ghClient, *org.Login, repo)
					if err != nil {
						log.Printf("Failed to get additions/deletions for %s. Error: %v\n", repo, err)
					}
				} else {
					mods = modsMap[repo]
				}

				if _, err := fo["stats.csv"].WriteString(formatRepoStats(*org.Login, repo, commits, mods)); err != nil {
					log.Printf("Failed to write stats for %s repository in %s org to stats.csv. Error: %v\n", repo, *org.Login, err)
				}
			}
		}

		log.Printf("Completed %s, %d of %d\n", *org.Login, i+1, len(orgs))
	}
	log.Print("\nCSV is ready!\n")
	return
}
