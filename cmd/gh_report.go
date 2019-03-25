package cmd

import (
	"context"
	"os"
	"strings"
	"time"

	ao "github.com/appoptics/appoptics-api-go"
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
		aoToken := os.Getenv("AO_TOKEN")
		aoClient := ao.NewClient(aoToken)
		mService := aoClient.MeasurementsService()
		measurement := ao.Measurement{
			Name:  "augit.gh-report.runs",
			Value: 1,
			Time:  time.Now().Unix(),
			Tags:  map[string]string{"environment": os.Getenv("ENVIRONMENT")},
		}
		cxn, err := pop.Connect(os.Getenv("ENVIRONMENT"))
		if err != nil {
			log.Fatal(err)
		}
		ghudb := models.NewGithubUserDB(cxn)
		ghodb := models.NewGithubOwnerDB(cxn)
		sadb := models.NewServiceAccountDB(cxn)
		err = persistUsers(ghudb, ghodb, sadb)
		if err != nil {
			log.Fatalln(err)
		}
		log.Info(generateSuccessString("gh-report"))
		_, err = mService.Create(ao.NewMeasurementsBatch([]ao.Measurement{measurement}, nil))
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(ghReportCmd)
}

type orgOwners struct {
	org    string
	owners []*github.User
}

// persistUsers gets all of the users for all relevant organizations and saves them to the Augit db.
// If they already exist in the db it is a no-op. This is intended to be run regularly.
func persistUsers(ghudb models.GithubUserAccessor, ghodb models.GithubOwnerAccessor, sadb models.ServiceAccountAccessor) error {
	ctx := context.Background()
	cf := config.GetConfig()
	ghClient := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cf.Github.Token})))
	orgs, err := swgithub.GetSWOrgs(ctx, ghClient, cf)
	if err != nil {
		log.WithError(err).Fatal("50011: Could not retrieve GitHub orgs")
		return err
	}
	allMembers := []*github.User{}
	allOwners := []*orgOwners{}
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
		allOwners = append(allOwners, &orgOwners{
			org:    *org.Login,
			owners: owners,
		})
	}

	err = persistMembers(allMembers, ghudb, sadb)
	if err != nil {
		return err
	}

	err = persistOwners(allOwners, ghodb)
	if err != nil {
		return err
	}

	allOwnerUsers := []*github.User{}
	for _, org := range allOwners {
		allOwnerUsers = append(allOwnerUsers, org.owners...)
	}

	ghOwners := generateMembersMap(allOwnerUsers)
	err = purgeOldOwners(ghOwners, ghodb)
	if err != nil {
		return err
	}

	return err
}

// generateMembersMap turns an array of GitHub members into a map of their logins so we can check if they are a member without iterating
// through a full array every time
func generateMembersMap(members []*github.User) map[string]bool {
	newMap := map[string]bool{}
	for _, member := range members {
		if member.Login == nil {
			continue
		}
		newMap[strings.ToLower(*member.Login)] = true
	}
	return newMap
}

func purgeOldOwners(ghMembers map[string]bool, ghodb models.GithubOwnerAccessor) error {
	existingOwners, err := ghodb.List()
	for _, owner := range existingOwners {
		if _, ok := ghMembers[strings.ToLower(owner.GithubID)]; !ok {
			err = ghodb.Delete(owner.GithubID)
			if err != nil {
				return err
			}
			log.Printf("Deleted github_owner with ID: %s", owner.GithubID)
		}
	}
	return nil
}

func persistMembers(allMembers []*github.User, ghudb models.GithubUserAccessor, sadb models.ServiceAccountAccessor) error {
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
		saExists, err := sadb.Exists(*member.Login)
		if saExists {
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
	return nil
}

func persistOwners(allOwners []*orgOwners, ghodb models.GithubOwnerAccessor) error {
	for _, orgOwner := range allOwners {
		for _, owner := range orgOwner.owners {
			if owner.Login == nil {
				continue
			}
			exists, err := ghodb.ExistsByGithubIDInOrg(*owner.Login, orgOwner.org)
			if exists {
				continue
			}
			err = ghodb.Create(&models.GithubOwner{
				GithubID: *owner.Login,
				Org:      orgOwner.org,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
