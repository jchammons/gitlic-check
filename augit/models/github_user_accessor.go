package models

import (
	"errors"
	"github.com/gobuffalo/pop"
)

type GithubUserAccessor interface {
	Create(*GithubUser) error
	Upsert(*GithubUser) error
	ReplaceGHRow(*GithubUser) error
	Find(string) (*GithubUser, error)
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

func (ghudb *GithubUserDB) Upsert(user *GithubUser) error {
	ghUser := &GithubUser{}
	// I could not get Pop to let me update a record without first retrieving it, I assume
	// because it requires the primary key.
	err := ghudb.tx.Where("email = ? OR github_id = ?", user.Email, user.GithubID).First(ghUser)
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

func (ghudb *GithubUserDB) ReplaceGHRow(inUser *GithubUser) error {
	return ghudb.tx.Transaction(func(tx *pop.Connection) error {
		if inUser.GithubID == "" {
			return errors.New("must provide a GitHub ID")
		}
		existingGHRow := &GithubUser{}
		err := ghudb.tx.Where("github_id = ?", inUser.GithubID).First(existingGHRow)
		if err != nil {
			return err
		}

		existingUser := &GithubUser{}
		err = tx.Where("email = ?", inUser.Email).First(existingUser)
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
	return foundUser, ghudb.tx.Where("email = ?", email).First(foundUser)
}

//
func (ghudb *GithubUserDB) ExistsByGithubID(ghID string) (bool, error) {
	return ghudb.tx.Where("github_id = ?", ghID).Exists(&GithubUser{})
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
	err := ghudb.tx.Where("github_id = ?", ghID).First(foundUser)
	if err != nil {
		return err
	}

	return ghudb.tx.Destroy(foundUser)
}

func (ghudb *GithubUserDB) AddAdmin(email string) error {
	foundUser := &GithubUser{}
	err := ghudb.tx.Where("email = ?", email).First(foundUser)
	if err != nil {
		return err
	}
	foundUser.Admin = true
	return ghudb.tx.Update(foundUser)
}

func (ghudb *GithubUserDB) RemoveAdmin(email string) error {
	foundUser := &GithubUser{}
	err := ghudb.tx.Where("email = ?", email).First(foundUser)
	if err != nil {
		return err
	}
	foundUser.Admin = false
	return ghudb.tx.Update(foundUser)
}

func (ghudb *GithubUserDB) MakeOwner(ghID string) error {
	foundUser := &GithubUser{}
	err := ghudb.tx.Where("github_id = ?", ghID).First(foundUser)
	if err != nil {
		return err
	}

	foundUser.Owner = true
	return ghudb.tx.Update(foundUser)
}
