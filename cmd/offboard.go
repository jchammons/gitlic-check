package cmd

import (
	"context"
	"os"
	"time"

	ao "github.com/appoptics/appoptics-api-go"
	"github.com/gobuffalo/pop"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/gitlic-check/config"
	"github.com/solarwinds/gitlic-check/swgithub"
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
		log.Println("couldn't connect to db: ", err)
		os.Exit(1)
	}
	offboardCmd.Flags().StringSliceVar(&orgsToProcess, "orgs", []string{}, "organization names to process")
	offboardCmd.Flags().BoolVar(&dryRun, "dry", false, "set if you just want to dry run")
	rootCmd.AddCommand(offboardCmd)
}

var offboardCmd = &cobra.Command{
	Use: "offboard",
	Run: func(cmd *cobra.Command, args []string) {
		aoToken := os.Getenv("AO_TOKEN")
		aoClient := ao.NewClient(aoToken)
		mService := aoClient.MeasurementsService()
		measurement := ao.Measurement{
			Name:  "augit.offboard.runs",
			Value: 1,
			Time:  time.Now().Unix(),
			Tags:  map[string]string{"environment": os.Getenv("ENVIRONMENT")},
		}
		log.Infoln("==========")
		log.Infoln("Starting offboard process with these orgs:")
		log.Infof("%+v", orgsToProcess)
		log.Infoln("Dry run is:")
		log.Infof("%+v", dryRun)
		aldb := models.NewAuditLogDB(db)
		offboard(aldb)
		log.Info(generateSuccessString("offboard"))
		_, err := mService.Create(ao.NewMeasurementsBatch([]ao.Measurement{measurement}, nil))
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func offboard(aldb models.AuditLogAccessor) {
	cf := config.GetConfig()
	client := github.NewClient(
		oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cf.Github.Token})))

	orgs, err := swgithub.GetSWOrgs(context.Background(), client, cf)
	if err != nil {
		log.WithError(err).Fatal("50011: Could not retrieve GitHub orgs")
	}

	for _, org := range orgs {
		if org.Login == nil || !isOrgToProcess(org) {
			continue
		}
		lo := &github.ListOptions{PerPage: 100}
		memOpt := &github.ListMembersOptions{ListOptions: *lo}
		members, err := swgithub.GetOrgMembers(context.Background(), client, org, memOpt)
		if err != nil {
			log.WithError(err).Errorf("50002: Could not get members for %s, continuing to next org", *org.Login)
		}
		for _, memb := range members {
			err := processMember(memb, client, org, aldb)
			if err != nil {
				log.WithError(err).Errorf("50001: Error processing member %s: %s", memb.GetLogin(), err)
			}
		}
	}
}

func processMember(member *github.User, client *github.Client, org *github.Organization, aldb models.AuditLogAccessor) error {
	swiUser := &models.GithubUser{}
	sa := &models.ServiceAccount{}
	exists, err := db.Where("(LOWER(github_id) = LOWER(?) AND username != '') OR (LOWER(github_id) = LOWER(?) AND email != '')", member.GetLogin(), member.GetLogin()).Exists(swiUser)
	if err != nil {
		return err
	}
	saExists, err := db.Where("LOWER(github_id) = LOWER(?)", member.GetLogin()).Exists(sa)
	if err != nil {
		return err
	}
	if !exists && !saExists {
		log.Infof("Did not find registered account for %s in org %s", member.GetLogin(), org.GetLogin())
		if !dryRun {
			log.Infof("Removing %s from %s", member.GetLogin(), org.GetLogin())
			_, err := client.Organizations.RemoveOrgMembership(context.Background(), member.GetLogin(), org.GetLogin())
			if err != nil {
				return err
			}
			al := &models.AuditLog{
				GithubID: member.GetLogin(),
			}
			err = aldb.Create(al)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{
					"github_id": member.GetLogin(),
				}).Warn("Could not create entry in audit log table")
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
