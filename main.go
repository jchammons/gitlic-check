package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func getTokens() []string {
	fh, err := ioutil.ReadFile("tokens.txt")
	if err != nil {
		log.Printf("Failed to read data from tokens.txt: %v\n", err)
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
		log.Printf("Failed to create repos.csv with %s\n", err)
	}
	fUsers, err := os.OpenFile("users.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create users.csv with %s\n", err)
	}
	fInvites, err := os.OpenFile("invites.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create invites.csv with %s\n", err)
	}
	return fRepos, fUsers, fInvites
}

func uploadFiles(d *drive.Service, auth string, teamDrive bool, files ...*os.File) {
	var parents []string
	parents = append(parents, auth)
	for _, file := range files {
		f := &drive.File{
			MimeType: "text/csv",
			Name:     file.Name(),
			Parents:  parents,
		}
		if teamDrive {
			_, err := d.Files.Create(f).SupportsTeamDrives(true).Do()
			if err != nil {
				log.Printf("Failed to upload %s: %s\n", file.Name(), err)
			} else {
				log.Printf("Successfully uploaded %s", file.Name())
			}
		} else {
			_, err := d.Files.Create(f).Do()
			if err != nil {
				log.Printf("Failed to upload %s: %s\n", file.Name(), err)
			} else {
				log.Printf("Successfully uploaded %s\n", file.Name())
			}
		}
	}
}

func main() {
	ctx := context.Background()
	ignoredOrgs := getIgnoredOrgs()
	auth := getTokens()

	ghClient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth[0]})))
	lo := &github.ListOptions{PerPage: 100}

	fmt.Print("Working...\n\n")

	fRepos, fUsers, fInvites := prepareOutput()

	if _, err := fRepos.WriteString("Org,Repo,Private,Fork,License\n"); err != nil {
		log.Printf("Initial save to repos CSV failed with %s\n", err)
	}

	if _, err := fUsers.WriteString("Org,User,Two-Factor Enabled\n"); err != nil {
		log.Printf("Initial save to users CSV failed with %s\n", err)
	}

	if _, err := fInvites.WriteString("Org,User,Date Sent,Invited By\n"); err != nil {
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
			fmt.Printf("\nIgnored %s, %d of %d\n", *org.Login, i+1, len(orgs))
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
				// TODO Check invite date; if older than 2 weeks, call Organization
				// Cancel Membership endpoint with user id
				inviteDate := fmt.Sprint(invite.CreatedAt.Year(), "-", invite.CreatedAt.Day(), "-", invite.CreatedAt.Month())
				if _, err := fInvites.WriteString(fmt.Sprint(*org.Login, ",", *invite.Login, ",", inviteDate, ",", *invite.Inviter.Login, "\n")); err != nil {
					log.Printf("Failed to write invite for %s\n", *invite.Login)
				}
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
			if _, err := fRepos.WriteString(fmt.Sprint(*org.Login, ",", *repo.Name, ",", *repo.Private, ",", *repo.Fork, ",")); err != nil {
				log.Printf("Failed to write to repos.csv on %s with %s\n", *org.Login, err)
			}

			lics, _, err := ghClient.Repositories.License(ctx, *repo.Owner.Login, *repo.Name)
			if err != nil {
				fRepos.WriteString("None\n")
				continue
			}

			if _, err = fRepos.WriteString(fmt.Sprint(*lics.License.Name, "\n")); err != nil {
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
			if _, err := fUsers.WriteString(fmt.Sprint(*org.Login, ",", *member.Login, ",")); err != nil {
				log.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
			}

			if membersFilter[*member.Login] {
				if _, err := fUsers.WriteString("False\n"); err != nil {
					log.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
				}
				continue
			}

			if _, err := fUsers.WriteString("True\n"); err != nil {
				log.Printf("Failed to write to users.csv on %s with %s\n", *org.Login, err)
			}
		}

		fmt.Printf("\nCompleted %s, %d of %d\n", *org.Login, i+1, len(orgs))
	}

	fmt.Print("\nCSVs are ready!\n")

	if len(auth) > 1 {
		fmt.Print("Uploading to Google Drive...\n")
		teamDrive := false
		if len(auth) > 2 {
			input, err := strconv.ParseBool(auth[2])
			if err != nil {
				fmt.Print("Could not parse boolean value for Team Drive from third argument in tokens.txt. Please ensure you are using a boolean value.\n")
			}
			teamDrive = input
		}
		secret, err := ioutil.ReadFile("../config.json")
		if err != nil {
			log.Fatalf("Failed to read JSON config file: %v\n", err)
		}
		config, err := google.JWTConfigFromJSON(secret, drive.DriveFileScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v\n", err)
		}
		drClient, err := drive.New(config.Client(ctx))
		if err != nil {
			log.Fatalf("Failed to connect to Drive: %v\n", err)
		}
		uploadFiles(drClient, auth[1], teamDrive, fRepos, fUsers, fInvites)
	}

	if err := fRepos.Close(); err != nil {
		log.Print(err)
	}
	if err := fUsers.Close(); err != nil {
		log.Print(err)
	}
	if err := fInvites.Close(); err != nil {
		log.Print(err)
	}
}
