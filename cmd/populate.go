package cmd

import (
	"fmt"

	"errors"
	"github.com/gobuffalo/pop"
	"github.com/solarwinds/gitlic-check/augit/models"
	swio "github.com/solarwinds/swio-users"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// populateCmd represents the populate command
var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "populate is the command used to populate a local Augit database with users from Azure AD",
	Run: func(cmd *cobra.Command, args []string) {
		err := PopulateDomainUsers()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

type AugitDB struct {
	db *pop.Connection
}

func (adb *AugitDB) Create(inUser *swio.User) error {
	if !inUser.Enabled {
		fmt.Printf("skipping or deleting %s for being disabled\n", inUser.Email)
		// Delete disabled users, so that the offboarding command can check for users existing within the
		// Augit DB associated w/ every GH user
		err := adb.checkForDeletion(inUser)
		if err != nil {
			return err
		}
		return nil
	}
	ghUser := &models.GithubUser{
		Name:  fmt.Sprintf("%s %s", inUser.FirstName, inUser.LastName),
		Email: inUser.Email,
	}
	if !adb.exists(inUser) {
		vErrs, err := adb.db.ValidateAndCreate(ghUser)
		if vErrs.HasAny() {
			return vErrs
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (adb *AugitDB) checkForDeletion(inUser *swio.User) error {
	queryUser := &models.GithubUser{}
	err := adb.db.Where("email = ?", inUser.Email).First(queryUser)
	if err != nil {
		if models.IsErrRecordNotFound(err) {
			return nil
		}
		return err
	}
	fmt.Printf("deleting %s for disabled or bad email\n", inUser.Email)
	return adb.db.Destroy(queryUser)
}

func (adb *AugitDB) exists(inUser *swio.User) bool {
	ghUser := &models.GithubUser{}
	exists, err := adb.db.Where("email = ?", inUser.Email).Exists(ghUser)
	if err != nil {
		fmt.Printf("error checking if user %s exists: %s\n", inUser.Email, err.Error())
		return false
	}
	return exists
}

func PopulateDomainUsers() error {
	cxn, err := pop.Connect(os.Getenv("ENVIRONMENT"))
	if err != nil {
		return err
	}
	id := os.Getenv("AD_CLIENT_ID")
	secret := os.Getenv("AD_SECRET")
	if id == "" || secret == "" {
		return errors.New("must provide id and secret")
	}
	augitDb := &AugitDB{cxn}
	populator := swio.NewPopulator(id, secret)
	for populator.MoreUsers() {
		users, err := populator.GetUsers()
		if err != nil {
			return err
		}
		for _, user := range users {
			err := augitDb.Create(user)
			if err != nil {
				// TODO: Create error type for array of errors to keep track of failures
				fmt.Printf("[ERROR] skipping user: %+v\n", user)
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(populateCmd)
}
