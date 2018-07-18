package cmd

import (
	"context"
	"log"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/google/go-github/github"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/gitlic-check/config"
	"github.com/solarwinds/gitlic-check/gitlic"
	"github.com/solarwinds/gitlic-check/swgithub"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var ghReportCmd = &cobra.Command{
	Use:   "gh-report",
	Short: "gh-report generates and uploads a list of GH users in SolarWinds organizations alongside their SolarWinds emails if we have them",
	Run: func(cmd *cobra.Command, args []string) {
		cxn, err := pop.Connect(os.Getenv("ENVIRONMENT"))
		if err != nil {
			log.Fatal(err)
		}
		ghudb := models.NewGithubUserDB(cxn)
		err = generateReport(ghudb)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ghReportCmd)
}

func generateReport(ghudb models.GithubUserAccessor) error {
	ctx := context.Background()
	cf := config.GetConfig()
	ghClient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cf.Github.Token})))
	orgs, err := swgithub.GetSWOrgs(ctx, ghClient, cf)
	if err != nil {
		return err
	}
	allMembers := []*github.User{}
	for _, org := range orgs {
		lo := &github.ListOptions{PerPage: 100}
		memOpt := &github.ListMembersOptions{ListOptions: *lo}
		members, err := swgithub.GetOrgMembers(ctx, ghClient, org, memOpt)
		if err != nil {
			log.Printf("Couldn't get org members, no filter, for %s: %s", *org.Login, err.Error())
		}
		allMembers = append(allMembers, members...)
	}

	dedupedMembers := map[string]string{}
	for _, member := range allMembers {
		if member.Login == nil {
			continue
		}
		swUser, err := ghudb.FindByGithubID(*member.Login)
		if err != nil {
			if !models.IsErrRecordNotFound(err) {
				return err
			}
		}
		dedupedMembers[*member.Login] = swUser.Email
	}

	swGhUserRows := [][]interface{}{}
	for ghId, email := range dedupedMembers {
		swGhUserRows = append(swGhUserRows, []interface{}{ghId, email})
	}
	err = gitlic.UploadToSheets(swGhUserRows, cf.Drive)
	return err
}
