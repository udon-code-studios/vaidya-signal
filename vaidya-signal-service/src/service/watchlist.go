package service

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	U "github.com/udon-code-sudios/vaidya-signal-service/utils"
)

func UpdateWatchlistThenEmailTodayTriggers(db *sqlx.DB) {
	UpdateWatchlist(db)
	EmailTodayWatchlistTriggers(db)
}

func UpdateWatchlist(db *sqlx.DB) {
	// get watchlist from "watchlist" table and story in array
	var watchlist []U.WatchlistTable
	db.Select(&watchlist, "SELECT * FROM watchlist")

	fmt.Println("[ DEBUG ] watchlist legnth:", len(watchlist))

	// create list of tickers from watchlist
	tickers := make([]string, 0)
	for _, ticker := range watchlist {
		tickers = append(tickers, ticker.Ticker)
	}

	// get signals for tickers
	signalsMap := FindAllVaidyaSignalsForTickers(tickers)

	// for each ticker in signalsMap
	for ticker, signals := range signalsMap {
		// if signals are found take last signal and update "watchlist" table
		if len(signals) > 0 {
			lastSignal := signals[len(signals)-1]

			// update "watchlist" table
			db.Exec("UPDATE watchlist SET last_trigger = $1 WHERE ticker = $2", lastSignal.TriggerDate, ticker)
		}
	}

	fmt.Println("[ INFO ] Finished updating watchlist")
}

func EmailTodayWatchlistTriggers(db *sqlx.DB) {
	// get watchlist from "watchlist" table and story in array
	var watchlistTableRows []U.WatchlistTable
	db.Select(&watchlistTableRows, "SELECT * FROM watchlist")

	// for each ticker in watchlist check if signal was triggered today
	tickers := make([]string, 0)
	for _, ticker := range watchlistTableRows {
		// skip tickers that don't have a last trigger date
		if !ticker.LastTrigger.Valid {
			continue
		}

		if ticker.LastTrigger.Time.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour)) {
			tickers = append(tickers, ticker.Ticker)
		}
	}
	// fmt.Println("[ DEBUG ] Watchlist tickers with today as trigger date:", tickers)

	// get emails from "emails" table and store in array
	var emailsTableRows []U.EmailsTable
	db.Select(&emailsTableRows, "SELECT * FROM emails")

	// convert emailsTableRows to array of strings
	emails := make([]string, 0)
	for _, email := range emailsTableRows {
		emails = append(emails, email.Email)
	}
	// fmt.Println("[ DEBUG ] Emails:", emails)

	// send email notification if tickers and emails exist
	if len(tickers) > 0 && len(emails) > 0 {
		fmt.Println("[ INFO ] Sending email notification for today's watchlist triggers:", tickers)
		U.SendEmail(
			emails,
			fmt.Sprintf("[%s] Signals Triggered", time.Now().Format("2006-01-02")),
			fmt.Sprintf(
				"Hello traders,\n\nVaidya signal(s) triggered on %s for tickers: %s\n\nView entire watchlist here: https://vaidya.udon.studio/watchlist\n\n--\nUdon Code Studios",
				time.Now().Format("2006-01-02"),
				tickers,
			),
		)
	}

	fmt.Println("[ INFO ] Finished emailing today's watchlist triggers")
}
