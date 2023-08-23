package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	U "github.com/udon-code-sudios/vaidya-signal-service/utils"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/go-co-op/gocron"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// define constants
const EMA12_SMOOTHING float64 = 2
const EMA26_SMOOTHING float64 = 2
const RSI_PERIOD int = 14
const LOW_DETECTION int = 3 // # of days before and after a low that should have higher closes

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
	// Open PostgreSQL connection
	//--------------------------------------------------------------------------

	db := connectToPg()
	defer db.Close()

	//--------------------------------------------------------------------------
	// Create cron job to print status every hour
	//--------------------------------------------------------------------------

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Hour().Do(func() { fmt.Println("[ INFO ] Service is still running...") })
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

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received for URI:", r.RequestURI, "and method:", r.Method)
		fmt.Fprintf(w, `{"message": "hello world."}`)
	})

	// start server
	log.Fatal(http.ListenAndServe(":"+port, nil))

	//--------------------------------------------------------------------------
	// Generate data for each ticker
	//--------------------------------------------------------------------------

	today := time.Now().Add(-15 * time.Minute) // 15 minutes are subtracted due to Alpaca free-tier limitations
	fiveYearsAgo := today.AddDate(-5, 0, 0)
	tenYearsAgo := today.AddDate(-10, 0, 0)
	fmt.Println("[ DEBUG ] Today:", today.Format("2006/01/02"))
	fmt.Println("[ DEBUG ] 10 Years Ago:", tenYearsAgo.Format("2006/01/02"))

	for _, ticker := range tickers {
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("[ INFO ] Starting detection for ticker:", ticker)

		// get day bars for past 7 years from Alpaca
		bars, err := marketdata.GetBars(ticker, marketdata.GetBarsRequest{
			TimeFrame: marketdata.OneDay,
			Start:     tenYearsAgo,
			End:       today,
		})
		panicOnNotNil(err)

		// print first and last bars
		fmt.Println("[ DEBUG ] first bar:", bars[0])
		fmt.Println("[ DEBUG ] last bar:", bars[len(bars)-1])

		// loop over days to generate indicators
		var last12EMA float64
		var last26EMA float64
		var lastAvgGain float64 = 0 // for RSI
		var lastAvgLoss float64 = 0 // for RSI
		var days = make([]DaysTable, len(bars))
		for i, bar := range bars {
			// copy Alpaca Bars into DaysTable
			days[i].Date = bar.Timestamp
			days[i].Open = bar.Open
			days[i].High = bar.High
			days[i].Low = bar.Low
			days[i].Close = bar.Close
			days[i].Volume = bar.Volume

			// skip first 25 bars
			if i < 25 {
				continue
			}

			// initialize EMAs and Avg Gain/Loss on bar 26
			if i == 25 {
				last12EMA = U.CalcBarCloseSMA(bars[i-11 : i+1])
				last26EMA = U.CalcBarCloseSMA(bars[i-25 : i+1])

				lastAvgGain = U.CalcFirstAvgGainLoss(bars[i-RSI_PERIOD:i+1], true)
				lastAvgLoss = U.CalcFirstAvgGainLoss(bars[i-RSI_PERIOD:i+1], false)

				continue
			}

			last12EMA = U.CalcEMA(bar.Close, last12EMA, 12, EMA12_SMOOTHING)
			last26EMA = U.CalcEMA(bar.Close, last26EMA, 26, EMA26_SMOOTHING)
			days[i].MACD = last12EMA - last26EMA

			lastAvgGain = U.CalcAvgGainLoss(RSI_PERIOD, lastAvgGain, bars[i-1].Close, bar.Close, true)
			lastAvgLoss = U.CalcAvgGainLoss(RSI_PERIOD, lastAvgLoss, bars[i-1].Close, bar.Close, false)
			days[i].RSI = 100 - 100/(1+(lastAvgGain/lastAvgLoss))
		}

		//------------------------------------------------------------------------
		// write data file
		//------------------------------------------------------------------------

		// loop over days to generate output file
		for _, day := range days {
			// skip until 5 years ago date
			if day.Date.Before(fiveYearsAgo) {
				continue
			}

			query := `INSERT INTO days(ticker, date, open, high, low, close, volume, macd, rsi)
								VALUES(:ticker, :date, :open, :high, :low, :close, :volume, :macd, :rsi)`

			// insert row in db (doesn't check for errors)
			db.NamedExec(query, day)
		}

		//------------------------------------------------------------------------
		// find local lows
		//
		// NOTE: local lows are defined as having a lower close than the three
		//       previous days and the three following days
		//------------------------------------------------------------------------

		var lows []int // indexes of lows
	find_lows:
		for i, bar := range bars {
			// skip lows that are not from last 5 years
			if bar.Timestamp.Before(fiveYearsAgo) {
				continue
			}

			// NOTE: local lows are defined as having a lower close than the three
			// previous days and the three following days

			// skip first and last three bars
			if i < LOW_DETECTION || i >= len(bars)-LOW_DETECTION {
				continue
			}

			// skip if close is not a low
			for j := 1; j <= LOW_DETECTION; j++ {
				if bar.Close > bars[i-j].Close || bar.Close > bars[i+j].Close {
					continue find_lows
				}
			}

			lows = append(lows, i)
		}

		//------------------------------------------------------------------------
		// find signal triggers
		//------------------------------------------------------------------------

		for i, daysIdx := range lows {
			// skip first low, as a previous low is needed
			if i == 0 {
				continue
			}

			// skip if current low is not lower (i.e. is greater) than previous low
			// NOTE: if lows are equal, it will not skip
			if bars[daysIdx].Close > bars[lows[i-1]].Close {
				continue
			}

			// skip if current low's MACD and RSI are not both higher than the
			// previous low's
			// NOTE: if MACDs and RSIs are equal, it will not skip
			if days[daysIdx].MACD < days[lows[i-1]].MACD ||
				days[daysIdx].RSI < days[lows[i-1]].RSI {
				continue
			}

			// skip if current low's volume is not lower (i.e. is greater) than the
			// previous low's
			// NOTE: if volumes are equal, it will not skip
			if bars[daysIdx].Volume > bars[lows[i-1]].Volume {
				continue
			}

			query := `INSERT INTO vaidya_` + strings.ToLower(ticker) + ` (ticker, trigger_date, low_2_date, low_1_date) 
								VALUES(:ticker, :trigger_date, :low_2_date, :low_1_date)`

			// insert row in db (doesn't check for errors)
			db.NamedExec(query, VaidyaSignalsTable{
				TriggerDate: days[daysIdx+LOW_DETECTION].Date,
				Low2Date:    days[daysIdx].Date,
				Low1Date:    days[lows[i-1]].Date,
			})
		}
	}
}

//----------------------------------------------------------------------------
// types
//----------------------------------------------------------------------------

type DaysTable struct {
	Date   time.Time `db:"date"`
	Open   float64   `db:"open"`
	High   float64   `db:"high"`
	Low    float64   `db:"low"`
	Close  float64   `db:"close"`
	Volume uint64    `db:"volume"`
	MACD   float64   `db:"macd"`
	RSI    float64   `db:"rsi"`
}

type TickersTable struct {
	Ticker    string `db:"ticker"`
	FirstDate string `db:"first_date"`
	LastDate  string `db:"last_date"`
}

type VaidyaSignalsTable struct {
	TriggerDate time.Time `db:"trigger_date"` // day signal was triggered
	Low2Date    time.Time `db:"low_2_date"`   // current low
	Low1Date    time.Time `db:"low_1_date"`   // previous low
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
