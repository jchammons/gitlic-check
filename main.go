package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func getTokens() []string {
	fh, err := ioutil.ReadFile("tokens.txt")
	if err != nil {
		log.Fatalf("Failed to read data from tokens.txt: %v\n", err)
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
	fRepos, err := os.OpenFile("repos.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Failed to create repos.csv with %s\n", err)
	}
	fUsers, err := os.OpenFile("users.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Failed to create users.csv with %s\n", err)
	}
	fInvites, err := os.OpenFile("invites.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Failed to create invites.csv with %s\n", err)
	}
	return fRepos, fUsers, fInvites
}

func uploadFiles(d *drive.Service, auth string, enableTeamDrive bool, files ...*os.File) {
	var parents []string
	parents = append(parents, auth)
	for _, file := range files {
		f := &drive.File{
			MimeType: "text/csv",
			Name:     file.Name(),
			Parents:  parents,
		}
		_, err := d.Files.Create(f).Media(file).SupportsTeamDrives(enableTeamDrive).Do()
		if err != nil {
			log.Printf("Failed to upload %s: %s\n", file.Name(), err)
		} else {
			log.Printf("Successfully uploaded %s", file.Name())
		}
	}
}

func main() {
	ctx := context.Background()
	ignoredOrgs := getIgnoredOrgs()
	auth := getTokens()

	ghClient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth[0]})))
	lo := &github.ListOptions{PerPage: 100}
	twoWeeks := time.Duration(336) * time.Hour

	log.Print("Working...\n\n")

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	fRepos, fUsers, fInvites := prepareOutput()
	defer func() {
		if err := fRepos.Close(); err != nil {
			log.Fatal(err)
		}
		if err := fUsers.Close(); err != nil {
			log.Fatal(err)
		}
		if err := fInvites.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if _, err := fRepos.WriteString("Org,Repo,Private,Fork,License\n"); err != nil {
		log.Printf("Initial save to repos CSV failed with %s\n", err)
	}

	if _, err := fUsers.WriteString("Org,User,Two-Factor Enabled\n"); err != nil {
		log.Printf("Initial save to users CSV failed with %s\n", err)
	}

	if _, err := fInvites.WriteString("Org,User,Date Sent,Invited By,Deleted\n"); err != nil {
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
				if _, err := fInvites.WriteString(fmt.Sprint(*org.Login, ",", *invite.Login, ",", inviteDate, ",", *invite.Inviter.Login)); err != nil {
					log.Printf("Failed to write invite data for %s from %s to invite.csv\n", *invite.Login, *org.Login)
				}
				if tSinceInvite > twoWeeks {
					_, err := ghClient.Organizations.RemoveOrgMembership(ctx, *invite.Login, *org.Login)
					if err != nil {
						log.Printf("Failed to remove flagged pending invitation for %s from org %s\n", *invite.Login, *org.Login)
					} else {
						fInvites.WriteString(",True\n")
					}
					continue
				}
				fInvites.WriteString(",\n")
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

		log.Printf("Completed %s, %d of %d\n", *org.Login, i+1, len(orgs))
	}

	log.Print("\nCSVs are ready!\n")

	if len(auth) > 1 {
		log.Print("Uploading to Google Drive...\n")
		enableTeamDrive := false
		if len(auth) > 2 {
			input, err := strconv.ParseBool(strings.TrimSpace(auth[2]))
			if err != nil {
				log.Printf("Failed to parse boolean value for Team Drive from third argument in tokens.txt (default - false). Error: %v\n", err)
			} else {
				enableTeamDrive = input
			}
		}

		secret, err := ioutil.ReadFile(filepath.Join(pwd, "config.json"))
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
		uploadFiles(drClient, auth[1], enableTeamDrive, fRepos, fUsers, fInvites)
	}
}
