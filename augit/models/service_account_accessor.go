package models

import "github.com/gobuffalo/pop"

type ServiceAccountAccessor interface {
	Upsert(*ServiceAccount) error
	FindByGithubID(string) (*ServiceAccount, error)
	List() ([]*ServiceAccount, error)
}

type ServiceAccountDB struct {
	tx *pop.Connection
}

func NewServiceAccountDB(cxn *pop.Connection) *ServiceAccountDB {
	return &ServiceAccountDB{cxn}
}

func (sadb *ServiceAccountDB) Upsert(acct *ServiceAccount) error {
	serviceAcct := &ServiceAccount{}
	// I could not get Pop to let me update a record without first retrieving it, I assume
	// because it requires the primary key.
	err := sadb.tx.Where("github_id = ?", acct.GithubID).First(serviceAcct)
	if err != nil {
		return err
	}
	serviceAcct.AdminResponsible = acct.AdminResponsible
	vErrs, err := sadb.tx.ValidateAndUpdate(serviceAcct)
	if vErrs.HasAny() {
		return vErrs
	} else if err != nil {
		return err
	}
	return nil
}

func (sadb *ServiceAccountDB) FindByGithubID(ghID string) (*ServiceAccount, error) {
	foundAcct := &ServiceAccount{}
	return foundAcct, sadb.tx.Where("github_id = ?", ghID).First(foundAcct)
}

func (sadb *ServiceAccountDB) List() ([]*ServiceAccount, error) {
	accts := []*ServiceAccount{}
	return accts, sadb.tx.All(&accts)
}
