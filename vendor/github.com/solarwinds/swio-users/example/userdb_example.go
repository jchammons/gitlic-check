package main

import (
	"errors"
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
	populator := swio.NewPopulator(id, secret)
	err := populator.Populate(&SampleUserDB{})
	if err != nil {
		return err
	}
}

func main() {
	err := PopulateDomainUsers()
	if err != nil {
		log.Fatalln(err)
	}
}
