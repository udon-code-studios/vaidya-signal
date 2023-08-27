package service

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func ExportedFunction() {
	fmt.Println("uhh ok")
}


// Creates day bar and vaidya signal tables in the database for the specified
// ticker. Sets first_date and last_date for the ticker as 5 years ago and 
// today (if after 4:15pm, otherwise yesterday), respectively.
func AddTicker(ticker string, db *sqlx.DB) {
	// determine first and last date
	now := time.Now()
	var lastDate time.Time
	if now.After(time.Date(now.Year(), now.Month(), now.Day(), 16, 15, 0, 0, time.UTC)) {
		lastDate = now
	} else {
		lastDate = now.AddDate(0, 0, -1)
	}
	firstDate := lastDate.AddDate(-5, 0, 0)

	// define query
	query := `INSERT INTO tickers (ticker, first_date, last_date)
						VALUES (:ticker, :first_date, :last_date)`

	// insert row in db (doesn't check for errors)
	// this will trigger the creation of days_TICKER and vaidya_TICKER
	db.NamedExec(query, TickersTable{
		Ticker:    ticker,
		FirstDate: firstDate,
		LastDate:  lastDate,
	})
}