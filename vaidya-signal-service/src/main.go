package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	A "github.com/udon-code-sudios/vaidya-signal-service/api"
	S "github.com/udon-code-sudios/vaidya-signal-service/service"

	"github.com/go-co-op/gocron"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	//--------------------------------------------------------------------------
	// Ensure Alpaca environment variables are set
	//--------------------------------------------------------------------------

	checkAlpacaEnvVars()

	//--------------------------------------------------------------------------
	// Open PostgreSQL connection
	//--------------------------------------------------------------------------

	db := connectToPg()
	defer db.Close()
	database := A.Database{DB: db}

	//--------------------------------------------------------------------------
	// Create cron job to print status every hour
	//--------------------------------------------------------------------------

	// create scheduler
	newYork, err := time.LoadLocation("America/New_York")
	panicOnNotNil(err)
	scheduler := gocron.NewScheduler(newYork)

	// define jobs
	scheduler.Every(1).Hour().Do(func() {
		fmt.Println("[ INFO ] Service is still running...")
	})
	scheduler.Every(1).Day().At("12:00").Do(func() {
		fmt.Println("[ INFO ] Running daily watchlist scan...")
		S.ScanWatchlist(db)
	})

	// start scheduler
	scheduler.StartAsync()

	//--------------------------------------------------------------------------
	// Start HTTP server
	//--------------------------------------------------------------------------

	// get server port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("[ WARN ] PORT environment variable not set, defaulting to", port)
	}

	// define endpoints and handlers
	http.HandleFunc("/", A.UselessHandler)
	http.HandleFunc("/hello", A.UselessHandler)
	http.HandleFunc("/alive", A.UselessHandler)
	// http.HandleFunc("/api/v1/add-ticker", database.AddTickerHandler)
	// http.HandleFunc("/api/v1/remove-ticker", A.UselessHandler)
	// http.HandleFunc("/api/v1/add-email", A.UselessHandler)
	// http.HandleFunc("/api/v1/remove-email", A.UselessHandler)
	http.HandleFunc("/api/v1/get-signal-triggers", A.GetVaidyaSignalsHandler)
	http.HandleFunc("/api/v1/update-watchlist", database.UpdateWatchlistHandler)

	// start server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

//----------------------------------------------------------------------------
// types
//----------------------------------------------------------------------------

// type DaysTable struct {
// 	Date   time.Time `db:"date"`
// 	Open   float64   `db:"open"`
// 	High   float64   `db:"high"`
// 	Low    float64   `db:"low"`
// 	Close  float64   `db:"close"`
// 	Volume uint64    `db:"volume"`
// 	MACD   float64   `db:"macd"`
// 	RSI    float64   `db:"rsi"`
// }

// type TickersTable struct {
// 	Ticker    string `db:"ticker"`
// 	FirstDate string `db:"first_date"`
// 	LastDate  string `db:"last_date"`
// }

// type VaidyaSignalsTable struct {
// 	TriggerDate time.Time `db:"trigger_date"` // day signal was triggered
// 	Low2Date    time.Time `db:"low_2_date"`   // current low
// 	Low1Date    time.Time `db:"low_1_date"`   // previous low
// }

//----------------------------------------------------------------------------
// helper functions
//----------------------------------------------------------------------------

func panicOnNotNil(value interface{}) {
	if value != nil {
		panic(value)
	}
}

/*
Checks if all Alpaca environment variables are set, exiting with status 1 if
not.
*/
func checkAlpacaEnvVars() {
	alpacaEnvVars := []string{"APCA_API_KEY_ID", "APCA_API_SECRET_KEY"}

	for _, envVar := range alpacaEnvVars {
		if os.Getenv(envVar) == "" {
			fmt.Println("[ ERROR ] The environment variable", envVar, "is not set")
			os.Exit(1)
		}
	}
}

func connectToPg() *sqlx.DB {
	pgConnectionStringEnvVar := "PG_URI"
	pgConnectionString := os.Getenv(pgConnectionStringEnvVar)

	if pgConnectionString == "" {
		fmt.Println("[ ERROR ] The environment variable", pgConnectionStringEnvVar, "is not set")
		os.Exit(1)
	}

	// connect to db
	db, err := sqlx.Connect("pgx", pgConnectionString)
	panicOnNotNil(err)

	return db
}
