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
		// Continue without care
		return nil
	}
	ghUser := &models.GithubUser{
		Name:  fmt.Sprintf("%s %s", inUser.FirstName, inUser.LastName),
		Email: inUser.Email,
	}
	vErrs, err := adb.db.ValidateAndCreate(ghUser)
	if vErrs.HasAny() {
		return vErrs
	}
	if err != nil {
		return err
	}
	return nil
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
	populator := swio.NewPopulator(id, secret)
	err = populator.Populate(&AugitDB{cxn})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(populateCmd)
}
