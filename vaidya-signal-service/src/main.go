package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	A "github.com/udon-code-sudios/vaidya-signal-service/api"
	U "github.com/udon-code-sudios/vaidya-signal-service/utils"
)

func main() {
	//--------------------------------------------------------------------------
	// Ensure environment variables are set
	//--------------------------------------------------------------------------

	envVars := []string{"APCA_API_KEY_ID", "APCA_API_SECRET_KEY", "PG_URI", "EMAIL_PW"}

	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			fmt.Println("[ ERROR ] The environment variable", envVar, "is not set")
			os.Exit(1)
		}
	}

	//--------------------------------------------------------------------------
	// Open PostgreSQL connection
	//--------------------------------------------------------------------------

	db, err := sqlx.Connect("pgx", os.Getenv("PG_URI"))
	U.PanicOnNotNil(err)
	defer db.Close()

	database := A.Database{DB: db}

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
	http.HandleFunc("/api/v1/get-signal-triggers", A.GetVaidyaSignalsHandler)
	http.HandleFunc("/api/v1/update-watchlist", database.UpdateWatchlistHandler)
	http.HandleFunc("/api/v1/update-watchlist-email-today-triggers", database.UpdateWatchlistEmailTodayTriggersHandler)

	// start server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
