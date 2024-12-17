package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Resolver struct {
	Db *sql.DB
}
