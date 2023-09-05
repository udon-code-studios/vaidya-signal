package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ScanWatchlist(db *sqlx.DB) {
	// get watchlist from "watchlist" table and story in array
	var watchlist []WatchlistTable
	db.Select(&watchlist, "SELECT * FROM watchlist")

	fmt.Println("[ DEBUG ] watchlist:", watchlist)
}