package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/solarwinds/gitlic-check/swgithub"

	"github.com/gobuffalo/pop"
	"github.com/google/go-github/github"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/gitlic-check/config"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	db            *pop.Connection
	dryRun        bool
	orgsToProcess []string
)

func init() {
	var err error
	db, err = pop.Connect(os.Getenv("ENVIRONMENT"))
	if err != nil {
		fmt.Println("couldn't connect to db: ", err)
		os.Exit(1)
	}
	offboardCmd.Flags().StringSliceVar(&orgsToProcess, "orgs", []string{}, "organization names to process")
	offboardCmd.Flags().BoolVar(&dryRun, "dry", false, "set if you just want to dry run")
	rootCmd.AddCommand(offboardCmd)
}

var offboardCmd = &cobra.Command{
	Use: "offboard",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("==========")
		fmt.Println("Starting offboard process with these orgs:")
		fmt.Printf("%+v\n", orgsToProcess)
		fmt.Println("Dry run is:")
		fmt.Printf("%+v\n", dryRun)
		offboard()
	},
}

func offboard() {
	cf := config.GetConfig()
	client := github.NewClient(
		oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cf.Github.Token})))

	orgs, err := swgithub.GetSWOrgs(context.Background(), client, cf)
	if err != nil {
		log.Fatalln(err)
	}

	for _, org := range orgs {
		if org.Login == nil || !isOrgToProcess(org) {
			continue
		}
		lo := &github.ListOptions{PerPage: 100}
		memOpt := &github.ListMembersOptions{ListOptions: *lo}
		members, err := swgithub.GetOrgMembers(context.Background(), client, org, memOpt)
		if err != nil {
			fmt.Printf("Could not get members for %s, continuing to next org", *org.Login)
		}
		for _, memb := range members {
			processMember(memb, client, org)
		}
	}
}

func processMember(member *github.User, client *github.Client, org *github.Organization) error {
	swiUser := &models.GithubUser{}
	if member.Login == nil {
		fmt.Printf("no Login for member %+v\n", member)
		return nil
	}
	exists, err := db.Where("github_id = ? AND username != ''", member.Login).Exists(swiUser)
	if !exists {
		fmt.Printf("Did not find registered account for %s\n", *member.Login)
		if !dryRun {
			_, err := client.Organizations.RemoveOrgMembership(context.Background(), member.GetLogin(), org.GetLogin())
			if err != nil {
				return err
			}
		}
	}
	return err
}

func isOrgToProcess(org *github.Organization) bool {
	for _, orgName := range orgsToProcess {
		if *org.Login == orgName {
			return true
		}
	}
	return false
}
