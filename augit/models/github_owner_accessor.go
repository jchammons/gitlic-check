package models

import (
	"github.com/gobuffalo/pop"
)

type GithubOwnerAccessor interface {
	Create(*GithubOwner) error
	ExistsByGithubID(string) (bool, error)
	ExistsByGithubIDInOrg(string, string) (bool, error)
	List() ([]*GithubOwner, error)
	Delete(string) error
}

type GithubOwnerDB struct {
	tx *pop.Connection
}

func NewGithubOwnerDB(tx *pop.Connection) *GithubOwnerDB {
	return &GithubOwnerDB{tx}
}

func (ghodb *GithubOwnerDB) Create(inUser *GithubOwner) error {
	return ghodb.tx.Create(inUser)
}

func (ghodb *GithubOwnerDB) ExistsByGithubID(ghID string) (bool, error) {
	return ghodb.tx.Where("LOWER(github_id = LOWER(?)", ghID).Exists(&GithubOwner{})
}

func (ghodb *GithubOwnerDB) ExistsByGithubIDInOrg(ghID, org string) (bool, error) {
	return ghodb.tx.Where("LOWER(github_id) = LOWER(?) AND LOWER(org) = LOWER(?)", ghID, org).Exists(&GithubOwner{})
}

func (ghodb *GithubOwnerDB) List() ([]*GithubOwner, error) {
	owners := []*GithubOwner{}
	return owners, ghodb.tx.All(&owners)
}

func (ghodb *GithubOwnerDB) Delete(ghID string) error {
	foundOwner := &GithubOwner{}
	err := ghodb.tx.Where("LOWER(github_id) = LOWER(?)", ghID).First(foundOwner)
	if err != nil {
		return err
	}

	return ghodb.tx.Destroy(foundOwner)
}
