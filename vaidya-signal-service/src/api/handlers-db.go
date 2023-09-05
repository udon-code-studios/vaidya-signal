package api

import (
	"encoding/json"
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

func (db *Database) AddTickerHandler(w http.ResponseWriter, r *http.Request) {
	// log request invocation
	fmt.Println("[ INFO ] Request received for URI:", r.RequestURI, "and method:", r.Method)

	// return wrong method if not POST
	if r.Method != "POST" {
		fmt.Println("[ INFO ] Method", r.Method, "is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// parse request JSON body for ticker
	// body is in the form {"ticker": "SPY"}

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		fmt.Println("[ ERROR ] Error parsing request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify ticker exists
	ticker := body["ticker"].(string)
	if ticker == "" {
		fmt.Println("[ ERROR ] Missing ticker in request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// add ticker to database
	signals := S.AddTicker(ticker, db.DB)

	// return signals
	signalsJSON, _ := json.Marshal(signals)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(signalsJSON)
}

func (db *Database) UpdateWatchlistHandler(w http.ResponseWriter, r *http.Request) {
	// log request invocation
	fmt.Println("[ INFO ] Request received for URI:", r.RequestURI, "and method:", r.Method)

	// return wrong method if not POST
	if r.Method != "POST" {
		fmt.Println("[ INFO ] Method", r.Method, "is not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	S.ScanWatchlist(db.DB)

	// return 200 OK
	w.WriteHeader(http.StatusOK)
}
