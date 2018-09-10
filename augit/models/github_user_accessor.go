package models

import (
	"errors"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type GithubUserAccessor interface {
	Create(*GithubUser) error
	ReplaceGHRow(*GithubUser) error
	Find(string) (*GithubUser, error)
	FindByID(uuid.UUID) (*GithubUser, error)
	ExistsByGithubID(string) (bool, error)
	ListGHUsers() ([]*GithubUser, error)
	Delete(string) error
	AddAdmin(string) error
	RemoveAdmin(string) error
	MakeOwner(string) error
}

type GithubUserDB struct {
	tx *pop.Connection
}

func NewGithubUserDB(tx *pop.Connection) *GithubUserDB {
	return &GithubUserDB{tx}
}

func (ghudb *GithubUserDB) ReplaceGHRow(inUser *GithubUser) error {
	return ghudb.tx.Transaction(func(tx *pop.Connection) error {
		if inUser.GithubID == "" {
			return errors.New("must provide a GitHub ID")
		}
		alreadyExists, err := ghudb.tx.Where("LOWER(github_id) = LOWER(?) AND LOWER(email) = LOWER(?)", inUser.GithubID, inUser.Email).Exists(&GithubUser{})
		if err != nil {
			return err
		}
		if alreadyExists {
			return nil
		}
		existingGHRow := &GithubUser{}
		err = ghudb.tx.Where("LOWER(github_id) = LOWER(?)", inUser.GithubID).First(existingGHRow)
		if err != nil {
			return err
		}

		existingUser := &GithubUser{}
		err = tx.Where("LOWER(email) = LOWER(?)", inUser.Email).First(existingUser)
		if err != nil {
			return err
		}
		// Update the existing row with the GH ID
		existingUser.GithubID = inUser.GithubID

		vErrs, err := tx.ValidateAndUpdate(existingUser)
		if vErrs.HasAny() {
			return vErrs
		} else if err != nil {
			return err
		}
		// Delete the old row with the GH ID
		return tx.Destroy(existingGHRow)
	})
}

// Find returns the user with the given email
func (ghudb *GithubUserDB) Find(email string) (*GithubUser, error) {
	foundUser := &GithubUser{}
	return foundUser, ghudb.tx.Where("LOWER(email) = LOWER(?)", email).First(foundUser)
}

func (ghudb *GithubUserDB) FindByID(id uuid.UUID) (*GithubUser, error) {
	foundUser := &GithubUser{}
	return foundUser, ghudb.tx.Where("id = ?", id).First(foundUser)
}

//
func (ghudb *GithubUserDB) ExistsByGithubID(ghID string) (bool, error) {
	return ghudb.tx.Where("LOWER(github_id) = LOWER(?)", ghID).Exists(&GithubUser{})
}

func (ghudb *GithubUserDB) Create(inUser *GithubUser) error {
	return ghudb.tx.Create(inUser)
}

func (ghudb *GithubUserDB) ListGHUsers() ([]*GithubUser, error) {
	users := []*GithubUser{}
	return users, ghudb.tx.Where("github_id != ''").All(&users)
}

func (ghudb *GithubUserDB) Delete(ghID string) error {
	foundUser := &GithubUser{}
	err := ghudb.tx.Where("LOWER(github_id) = LOWER(?)", ghID).First(foundUser)
	if err != nil {
		return err
	}

	return ghudb.tx.Destroy(foundUser)
}

func (ghudb *GithubUserDB) AddAdmin(email string) error {
	foundUser := &GithubUser{}
	err := ghudb.tx.Where("LOWER(email) = LOWER(?)", email).First(foundUser)
	if err != nil {
		return err
	}
	foundUser.Admin = true
	return ghudb.tx.Update(foundUser)
}

func (ghudb *GithubUserDB) RemoveAdmin(email string) error {
	foundUser := &GithubUser{}
	err := ghudb.tx.Where("LOWER(email) = LOWER(?)", email).First(foundUser)
	if err != nil {
		return err
	}
	foundUser.Admin = false
	return ghudb.tx.Update(foundUser)
}

func (ghudb *GithubUserDB) MakeOwner(ghID string) error {
	foundUser := &GithubUser{}
	err := ghudb.tx.Where("LOWER(github_id) = LOWER(?)", ghID).First(foundUser)
	if err != nil {
		return err
	}

	foundUser.Owner = true
	return ghudb.tx.Update(foundUser)
}
