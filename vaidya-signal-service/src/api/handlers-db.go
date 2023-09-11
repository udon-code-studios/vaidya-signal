package api

import (
	"fmt"
	"net/http"

	S "github.com/udon-code-sudios/vaidya-signal-service/service"

	"github.com/jmoiron/sqlx"
)

//---------------------------------------------------------------------------
// Types
//---------------------------------------------------------------------------

type Database struct {
	DB *sqlx.DB
}

//---------------------------------------------------------------------------
// Methods
//---------------------------------------------------------------------------

// will trigger S.ScanWatchlist()
func (db *Database) UpdateWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	// log request invocation
	fmt.Println("[ INFO ] Request received for URI:", r.RequestURI, "and method:", r.Method)

	// return wrong method if not POST
	if r.Method != "POST" {
		fmt.Println("[ INFO ] Method", r.Method, "is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// return 202 Accepted
	w.WriteHeader(http.StatusAccepted)

	go S.UpdateWatchlist(db.DB)
}

// will trigger S.ScanWatchlist() and S.EmailTodayWatchlistTriggers()
func (db *Database) UpdateWatchlistEmailTodayTriggersHandler(w http.ResponseWriter, r *http.Request) {
	// log request invocation
	fmt.Println("[ INFO ] Request received for URI:", r.RequestURI, "and method:", r.Method)

	// return wrong method if not POST
	if r.Method != "POST" {
		fmt.Println("[ INFO ] Method", r.Method, "is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// return 202 Accepted
	w.WriteHeader(http.StatusAccepted)

	S.UpdateWatchlist(db.DB)
	S.EmailTodayWatchlistTriggers(db.DB)
}
