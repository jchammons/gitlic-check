package cmd

import (
	"context"
	"log"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/google/go-github/github"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/gitlic-check/config"
	"github.com/solarwinds/gitlic-check/swgithub"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var ghReportCmd = &cobra.Command{
	Use:   "gh-report",
	Short: "gh-report generates and persists a list of GH users in SolarWinds organizations",
	Run: func(cmd *cobra.Command, args []string) {
		cxn, err := pop.Connect(os.Getenv("ENVIRONMENT"))
		if err != nil {
			log.Fatal(err)
		}
		ghudb := models.NewGithubUserDB(cxn)
		err = persistUsers(ghudb)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ghReportCmd)
}

// persistUsers gets all of the users for all relevant organizations and saves them to the Augit db.
// If they already exist in the db it is a no-op. This is intended to be run regularly.
func persistUsers(ghudb models.GithubUserAccessor) error {
	ctx := context.Background()
	cf := config.GetConfig()
	ghClient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cf.Github.Token})))
	orgs, err := swgithub.GetSWOrgs(ctx, ghClient, cf)
	if err != nil {
		return err
	}
	allMembers := []*github.User{}
	allOwners := []*github.User{}
	for _, org := range orgs {
		lo := &github.ListOptions{PerPage: 100}
		memOpt := &github.ListMembersOptions{ListOptions: *lo}
		members, err := swgithub.GetOrgMembers(ctx, ghClient, org, memOpt)
		if err != nil {
			log.Printf("Couldn't get org members, no filter, for %s: %s", *org.Login, err.Error())
		}
		allMembers = append(allMembers, members...)

		owners, err := swgithub.GetOrgOwners(ctx, ghClient, org)
		if err != nil {
			log.Printf("Couldn't get org owners for %s: %s", *org.Login, err.Error())
		}
		allOwners = append(allOwners, owners...)
	}

	for _, member := range allMembers {
		if member.Login == nil {
			continue
		}
		exists, err := ghudb.ExistsByGithubID(*member.Login)
		if exists {
			continue
		}
		if err != nil {
			return err
		}
		err = ghudb.Create(&models.GithubUser{
			GithubID: *member.Login,
		})
		if err != nil {
			return err
		}
	}
	for _, owner := range allOwners {
		if owner.Login == nil {
			continue
		}
		err = ghudb.MakeOwner(*owner.Login)
		if err != nil {
			return err
		}
	}
	return err
}
