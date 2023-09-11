package utils

import (
	"database/sql"
)

type WatchlistTable struct {
	Ticker      string       `db:"ticker"`
	LastTrigger sql.NullTime `db:"last_trigger"`
}

type EmailsTable struct {
	Email string `db:"email"`
}
