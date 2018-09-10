package models

import "github.com/gobuffalo/pop"

type ServiceAccountAccessor interface {
	Create(*ServiceAccount) error
	Exists(string) (bool, error)
	FindByGithubID(string) (*ServiceAccount, error)
	List() ([]*ServiceAccount, error)
}

type ServiceAccountDB struct {
	tx *pop.Connection
}

func NewServiceAccountDB(cxn *pop.Connection) *ServiceAccountDB {
	return &ServiceAccountDB{cxn}
}

func (sadb *ServiceAccountDB) Create(acct *ServiceAccount) error {
	vErrs, err := sadb.tx.ValidateAndCreate(acct)
	if vErrs.HasAny() {
		return vErrs
	} else if err != nil {
		return err
	}
	return nil
}

func (sadb *ServiceAccountDB) FindByGithubID(ghID string) (*ServiceAccount, error) {
	foundAcct := &ServiceAccount{}
	return foundAcct, sadb.tx.Where("LOWER(github_id) = LOWER(?)", ghID).First(foundAcct)
}

func (sadb *ServiceAccountDB) List() ([]*ServiceAccount, error) {
	accts := []*ServiceAccount{}
	return accts, sadb.tx.All(&accts)
}

func (sadb *ServiceAccountDB) Exists(ghID string) (bool, error) {
	return sadb.tx.Where("LOWER(github_id) = LOWER(?)", ghID).Exists(&ServiceAccount{})
}
