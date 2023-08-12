package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

func main() {
	//--------------------------------------------------------------------------
	// Parse flags for tickers
	//--------------------------------------------------------------------------

	tickers := parseFlags()

	// print parsed flags
	fmt.Println("[ INFO ] Tickers:", tickers)

	//--------------------------------------------------------------------------
	// Ensure Alpaca environment variables are set
	//--------------------------------------------------------------------------

	checkAlpacaEnvVars()

	//--------------------------------------------------------------------------
	// Loop over tickers
	//--------------------------------------------------------------------------

	today := time.Now().Add(-15*time.Minute) // 15 minutes are subtracted due to Alpaca free-tier limitations
	sixYearsAgo := today.AddDate(-6, 0, 0)
	fmt.Println("[ DEBUG ] Today:", today.Format("2006/01/02"))
	fmt.Println("[ DEBUG ] Six Years Ago:", sixYearsAgo.Format("2006/01/02"))

	for _, ticker := range tickers {
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("[ INFO ] Starting detection for ticker:", ticker)

		// get day bars for past 6 years from Alpaca
		bars, err := marketdata.GetBars(ticker, marketdata.GetBarsRequest{
			TimeFrame: marketdata.OneDay,
			Start:     sixYearsAgo,
			End:       today,
		})
		panicOnNotNil(err)

		// print first and last bars
		fmt.Println("[ DEBUG ] first bar:", bars[0])
		fmt.Println("[ DEBUG ] last bar:", bars[len(bars)-1])
	}
}

//----------------------------------------------------------------------------
// types
//----------------------------------------------------------------------------

type DayMetadata struct {
	PrevDate   string  `json:"prev_date"`
	PrevOpen   float64 `json:"prev_open"`
	PrevHigh   float64 `json:"prev_high"`
	PrevLow    float64 `json:"prev_low"`
	PrevClose  float64 `json:"prev_close"`
	PrevVolume uint64  `json:"prev_volume"`
	PrevVWAP   float64 `json:"prev_vwap"`
}

//----------------------------------------------------------------------------
// helper functions
//----------------------------------------------------------------------------

func panicOnNotNil(value interface{}) {
	if value != nil {
		panic(value)
	}
}

func parseFlags() (tickers []string) {
	// define and parse flags
	tickersFlag := flag.String("t", "", "Comma-separated list of ticker symbols (format: SYMBOL1,SYMBOL2) (required)")
	flag.Parse()

	// default tickers to SPY if none specified
	if *tickersFlag == "" {
		fmt.Println("[ INFO ] Ticker symbol list flag (-t) is missing. Defaulting to ticker: SPY")
		tickers = []string{"SPY"}
	} else {
		tickers = strings.Split(*tickersFlag, ",")
	}

	return
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