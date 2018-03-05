package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Follow instructions in README for proper authentication

type userdata struct {
	name string
	pat  string
}

func setUser() userdata {
	fh, err := ioutil.ReadFile("auth.txt")
	if err != nil {
		fmt.Print(err)
	}
	userInfo := strings.Split(string(fh), ",")
	user := userdata{name: userInfo[0], pat: userInfo[1]}
	// fmt.Printf("Retrieved user info: %s\n", user)
	return user
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

func saveToFile(name string, data string) error {
	// From Go docs
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.WriteString(data)
	if err != nil {
		return err
	}
	f.Close()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	user := setUser()
	ignoredOrgs := getIgnoredOrgs()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: user.pat,
		},
	)
	authClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(authClient)

	orgs, _, err := client.Organizations.List(ctx, user.name, nil)
	if err != nil {
		fmt.Printf("Organizations.List failed with %s\n", err)
		return
	}

	fmt.Print("Working... This may take a minute or two.\n")

	os.Remove("repos.csv")
	os.Remove("users.csv")

	var repoBuffer bytes.Buffer
	var userBuffer bytes.Buffer

	repoBuffer.WriteString("Org,Repo,Private,Fork,License\n")
	if err != nil {
		fmt.Printf("Failed to write to repo buffer with %s\n", err)
	}
	userBuffer.WriteString("Org,User,Two-Factor Enabled\n")
	if err != nil {
		fmt.Printf("Failed to write to user buffer with %s\n", err)
	}

	saveToFile("repos.csv", repoBuffer.String())
	if err != nil {
		fmt.Printf("Initial save to repos CSV failed with %s\n", err)
	}

	saveToFile("users.csv", userBuffer.String())
	if err != nil {
		fmt.Printf("Initial save to users CSV failed with %s\n", err)
	}

	repoBuffer.Reset()
	userBuffer.Reset()

	for i, org := range orgs {
		if ignoredOrgs != nil && ignoredOrgs[*org.Login] {
			fmt.Printf("Ignored %s, %d of %d\n", *org.Login, i+1, len(orgs))
			continue
		}

		var repos []*github.Repository
		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
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
			// fmt.Printf("  Repo: %s [%s] \n", *repo.Name, *lics.License.Name)
			repoBuffer.WriteString(fmt.Sprint(*org.Login, ",", *repo.Name, ",", *repo.Private, ",", *repo.Fork, ","))
			if err != nil {
				fmt.Printf("Failed to write to repo buffer on %s with %s\n", *org.Login, err)
			}

			lics, _, err := client.Repositories.License(ctx, *repo.Owner.Login, *repo.Name)
			if err != nil {
				repoBuffer.WriteString("None\n")
				continue
			}
			repoBuffer.WriteString(fmt.Sprint(*lics.License.Name, "\n"))
			if err != nil {
				fmt.Printf("Failed to write to repo buffer on %s with %s\n", *org.Login, err)
			}
		}

		var members []*github.User
		memOpt := &github.ListMembersOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		var membersNo2f []*github.User
		no2fOpt := &github.ListMembersOptions{
			Filter:      "2fa_disabled",
			ListOptions: github.ListOptions{PerPage: 100},
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
			userBuffer.WriteString(fmt.Sprint(*org.Login, ",", *member.Login, ","))
			if err != nil {
				fmt.Printf("Failed to write to user buffer on %s with %s\n", *org.Login, err)
			}

			if membersFilter[*member.Login] {
				userBuffer.WriteString("False\n")
				if err != nil {
					fmt.Printf("Failed to write to user buffer on %s with %s\n", *org.Login, err)
				}
				continue
			}
			userBuffer.WriteString("True\n")
			if err != nil {
				fmt.Printf("Failed to write to user buffer on %s with %s\n", *org.Login, err)
			}
		}

		saveToFile("repos.csv", repoBuffer.String())
		if err != nil {
			fmt.Printf("Saving to repos CSV failed on %d with %s\n", i, err)
		}

		saveToFile("users.csv", userBuffer.String())
		if err != nil {
			fmt.Printf("Saving to users CSV failed on %d with %s\n", i, err)
		}

		repoBuffer.Reset()
		userBuffer.Reset()
		fmt.Printf("Completed %d of %d\n", i+1, len(orgs))
	}

	fmt.Print("CSVs are now ready\n")
}
