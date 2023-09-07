package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ScanWatchlist(db *sqlx.DB) {
	// get watchlist from "watchlist" table and story in array
	var watchlist []WatchlistTable
	db.Select(&watchlist, "SELECT * FROM watchlist")

	fmt.Println("[ DEBUG ] watchlist legnth:", len(watchlist))

	// for each ticker in watchlist
	for _, ticker := range watchlist {
		signals := GetHistoricalVaidyaSignals(ticker.Ticker)

		// if signals are found take last signal and update "watchlist" table
		if len(signals) > 0 {
			lastSignal := signals[len(signals)-1]
			// fmt.Println("[ DEBUG ] lastSignal.TriggerDate:", lastSignal.TriggerDate)

			// update "watchlist" table
			db.Exec("UPDATE watchlist SET last_trigger = $1 WHERE ticker = $2", lastSignal.TriggerDate, ticker.Ticker)
		}
	}
}
