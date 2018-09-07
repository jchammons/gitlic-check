package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type GithubOwner struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Org       string    `json:"org" db:"org"`
	GithubID  string    `json:"github_id" db:"github_id"`
}

// String is not required by pop and may be deleted
func (g GithubOwner) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// GithubOwners is not required by pop and may be deleted
type GithubOwners []GithubOwner

// String is not required by pop and may be deleted
func (g GithubOwners) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (g *GithubOwner) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: g.Org, Name: "Org"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (g *GithubOwner) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (g *GithubOwner) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
