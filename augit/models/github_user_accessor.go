package models

import (
	"github.com/gobuffalo/pop"
)

type GithubUserAccessor interface {
	Upsert(*GithubUser) error
	Find(string) (*GithubUser, error)
	FindByGithubID(string) (*GithubUser, error)
}

type GithubUserDB struct {
	tx *pop.Connection
}

func NewGithubUserDB(tx *pop.Connection) *GithubUserDB {
	return &GithubUserDB{tx}
}

func (ghudb *GithubUserDB) Upsert(user *GithubUser) error {
	ghUser := &GithubUser{}
	// I could not get Pop to let me update a record without first retrieving it, I assume
	// because it requires the primary key.
	err := ghudb.tx.Where("email = ?", user.Email).First(ghUser)
	if err != nil {
		return err
	}
	if user.GithubID != "" {
		ghUser.GithubID = user.GithubID
	}
	vErrs, err := ghudb.tx.ValidateAndUpdate(ghUser)
	if vErrs.HasAny() {
		return vErrs
	} else if err != nil {
		return err
	}
	return nil
}

// Find returns the user with the given email
func (ghudb *GithubUserDB) Find(email string) (*GithubUser, error) {
	foundUser := &GithubUser{}
	return foundUser, ghudb.tx.Where("email = ?", email).First(foundUser)
}

// Find returns the user with the given GitHub ID
func (ghudb *GithubUserDB) FindByGithubID(ghID string) (*GithubUser, error) {
	foundUser := &GithubUser{}
	return foundUser, ghudb.tx.Where("github_id = ?", ghID).First(foundUser)
}
