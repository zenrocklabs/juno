package v3

import (
	"github.com/jmoiron/sqlx"

	"github.com/zenrocklabs/juno/database"
	"github.com/zenrocklabs/juno/database/postgresql"
)

var _ database.Migrator = &Migrator{}

// Migrator represents the database migrator that should be used to migrate from v2 of the database to v3
type Migrator struct {
	SQL *sqlx.DB
}

func NewMigrator(db *postgresql.Database) *Migrator {
	return &Migrator{
		SQL: db.SQL,
	}
}
