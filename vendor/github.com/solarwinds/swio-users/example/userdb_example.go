package main

import (
	"errors"
	"fmt"
	swio "github.com/solarwinds/swio-users"
	"log"
	"os"
)

// DomainUser is an example type to illustrate how a swio-users caller would
// have their own type that they can use in their methods fulfilling the UserDatabase
// interface
type DomainUser struct {
	FirstName string
	LastName  string
	Email     string
	GithubID  string
}

func (du *DomainUser) Save() error {
	// Persist the thing to your DB however you want
	log.Printf("[-] Saving user %+v\n", du)
	return nil
}

type SampleUserDB struct{}

func (sudb *SampleUserDB) Create(inUser *swio.User) error {
	newDomainUser := &DomainUser{
		FirstName: inUser.FirstName,
		LastName:  inUser.LastName,
		Email:     inUser.Email,
	}
	return newDomainUser.Save()
}

func PopulateDomainUsers() error {
	id := os.Getenv("ad_cilent_id")
	secret := os.Getenv("ad_secret")
	if id == "" || secret == "" {
		return errors.New("must provide id and secret")
	}
	db := &SampleUserDB{}
	populator := swio.NewPopulator(id, secret)
	for populator.MoreGroups() {
		users, err := populator.GetUsersByGroup()
		if err != nil {
			log.Fatalln(err)
		}
		for _, user := range users {
			err := db.Create(user)
			if err != nil {
				fmt.Printf("err: %+v\n", err)
			}
		}
	}
}

func main() {
	err := PopulateDomainUsers()
	if err != nil {
		log.Fatalln(err)
	}
}
