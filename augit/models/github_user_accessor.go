package models

import (
	"github.com/gobuffalo/pop"
)

type GithubUserAccessor interface {
	Upsert(*GithubUser) error
	Find(string) (*GithubUser, error)
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

func (ghudb *GithubUserDB) Find(email string) (*GithubUser, error) {
	foundUser := &GithubUser{}
	return foundUser, ghudb.tx.Where("email = ?", email).First(foundUser)
}
