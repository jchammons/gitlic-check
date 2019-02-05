package models

import (
	"github.com/gobuffalo/pop"
)

type AuditLogAccessor interface {
	Create(*AuditLog) error
	List() ([]*AuditLog, error)
}

type AuditLogDB struct {
	tx *pop.Connection
}

func NewAuditLogDB(tx *pop.Connection) *AuditLogDB {
	return &AuditLogDB{tx}
}

func (ldb *AuditLogDB) List() ([]*AuditLog, error) {
	entries := []*AuditLog{}
	return entries, ldb.tx.All(&entries)
}

func (ldb *AuditLogDB) Create(al *AuditLog) error {
	return ldb.tx.Create(al)
}
