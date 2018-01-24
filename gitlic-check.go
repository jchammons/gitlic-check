package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// you need to generate personal access token at
// https://github.com/settings/applications#personal-access-tokens
// Put your personal access token in pat.txt

type userdata struct {
	name string
	pat  string
}

func setUser() userdata {
	fh, err := ioutil.ReadFile("user.txt")
	if err != nil {
		fmt.Print(err)
	}
	userInfo := strings.Split(string(fh), ",")
	user := userdata{name: userInfo[0], pat: userInfo[1]}
	// fmt.Printf("Retrieved user info: %s\n", user)
	return user
}

func main() {
	ctx := context.Background()
	user := setUser()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: user.pat,
		},
	)
	authClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(authClient)

	orgs, _, err := client.Organizations.List(ctx, user.name, nil)
	if err != nil {
		fmt.Printf("Failed with %s\n", err)
		return
	}

	fmt.Print("Working...\n")
	var buffer bytes.Buffer
	for _, org := range orgs {
		if err != nil {
			fmt.Printf("Failed with %s\n", err)
			return
		}
		buffer.WriteString(fmt.Sprint("Org: ", *org.Login, "\n"))

		repos, _, err := client.Repositories.List(ctx, *org.Login, nil)
		if err != nil {
			fmt.Printf("Failed with %s\n", err)
			return
		}
		for _, repo := range repos {
			// fmt.Printf("  Repo: %s [%s] \n", *repo.Name, *lics.License.Name)
			if *repo.Private == false {
				buffer.WriteString(fmt.Sprint("  Repo: ", *repo.Name))

				lics, _, err := client.Repositories.License(ctx, *repo.Owner.Login, *repo.Name)
				if err != nil {
					buffer.WriteString(" - No License\n")
					continue
				}
				buffer.WriteString(fmt.Sprint(" [", *lics.License.Name, "]\n"))
			}
		}
	}

	ioutil.WriteFile("list.txt", []byte(buffer.String()), 0644)

	if err != nil {
		fmt.Printf("Failed to write to file with %s\n", err)
		return
	}

	fmt.Print("List now available in list.txt\n")
}
