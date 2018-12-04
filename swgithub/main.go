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
	"golang.org/x/oauth2"
)

var weekdays [7]string

// getIncludedMap will turn the array of included orgs from the config file into a map that is easy to check against as we iterate over organizations returned from the GitHub API
func getIncludedMap(a []string) map[string]bool {
	includedOrgs := make(map[string]bool)
	for _, org := range a {
		includedOrgs[strings.ToLower(org)] = true
	}
	return includedOrgs
}

func saveRepoStats(ctx context.Context, fo map[string]*os.File, org string, repo string, commits []*github.WeeklyCommitActivity, mods []*github.WeeklyStats) {
	totalCommits := 0
	linesAdded := 0
	linesRemoved := 0
	mostCommitsOneDay := 0
	mostCommitsOneWeek := 0
	var mostCommittable string
	highestDayCounts := map[string]int{"Sunday": 0, "Monday": 0, "Tuesday": 0, "Wednesday": 0, "Thursday": 0, "Friday": 0, "Saturday": 0}

	for _, commit := range commits {
		totalCommits += *commit.Total

		if *commit.Total > 0 {
			highestDay := 0
			highestDayInd := 0
			weekTotal := 0
			for i, dayTotal := range commit.Days {
				weekTotal += dayTotal

				if dayTotal > highestDay {
					highestDayInd = i
				}
				if dayTotal > mostCommitsOneDay {
					mostCommitsOneDay = dayTotal
				}
			}
			highestDayCounts[weekdays[highestDayInd]]++
			if weekTotal > mostCommitsOneWeek {
				mostCommitsOneWeek = weekTotal
			}
		}
	}

	if totalCommits > 0 {
		highestCount := 0
		for day, count := range highestDayCounts {
			if count > highestCount {
				mostCommittable = day
			}
		}
	} else {
		mostCommittable = "N/A"
	}

	mostLinesAddedOneWeek := 0
	mostLinesRemovedOneWeek := 0
	for _, mod := range mods {
		linesAdded += *mod.Additions
		linesRemoved += *mod.Deletions

		if *mod.Additions > mostLinesAddedOneWeek {
			mostLinesAddedOneWeek = *mod.Additions
		}

		if *mod.Deletions < mostLinesRemovedOneWeek {
			mostLinesRemovedOneWeek = *mod.Deletions
		}
	}

	fo["stats.csv"].WriteString(fmt.Sprint(org, ",", repo, ",", totalCommits, ",", mostCommitsOneDay, ",", mostCommitsOneWeek, ",", mostCommittable, ",", linesAdded, ",", linesRemoved, ",", mostLinesAddedOneWeek, ",", mostLinesRemovedOneWeek, "\n"))

	return
}

// RunGitlicCheck begins the process of querying the GitHub API. It will loop through your organizations and their repositories and pull info on configuration, license, and users, including invitations. It will output the results to respective CSV files in the output folder. See the README for an idea of what these CSV reports contain.
func RunGitlicCheck(ctx context.Context, cf config.Config, fo map[string]*os.File) {
	weekdays = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	ghClient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cf.Github.Token})))
	lo := &github.ListOptions{PerPage: 100}
	maxInviteT := time.Duration(cf.Github.RmInvitesAfter) * time.Hour

	log.Print("Working...\n\n")

	if _, err := fo["repos.csv"].WriteString("Org,Repo,Private,Fork,License\n"); err != nil {
		log.Printf("Initial save to repos CSV failed with %s\n", err)
	}

	if _, err := fo["stats.csv"].WriteString("Org,Repo,Commits in Last Year,Most Commits in a Day,Most Commits in a Week,Most Committable Day,Total Lines Added,Total Lines Removed,Most Lines Added in a Week,Most Lines Removed in a Week\n"); err != nil {
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
		log.Printf("Could not get orgs: %s", err.Error())
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

		var retries []string
		commitsMap := make(map[string][]*github.WeeklyCommitActivity)
		modsMap := make(map[string][]*github.WeeklyStats)
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

			commits, err, cmRetry := GetCommitActivity(ctx, ghClient, *org.Login, *repo.Name)
			if err != nil {
				log.Printf("Failed to get commit activity for %s. Error: %v\n", *repo.Name, err)
				continue
			}
			mods, modsErr, mdRetry := GetAdditionsDeletions(ctx, ghClient, *org.Login, *repo.Name)
			if modsErr != nil {
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

			saveRepoStats(ctx, fo, *org.Login, *repo.Name, commits, mods)
		}

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

				saveRepoStats(ctx, fo, *org.Login, repo, commits, mods)
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
