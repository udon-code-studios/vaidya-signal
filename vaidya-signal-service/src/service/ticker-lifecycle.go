package service

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/jmoiron/sqlx"
	U "github.com/udon-code-sudios/vaidya-signal-service/utils"
)

func ExportedFunction() {
	fmt.Println("uhh ok")
}

// Creates day bar and vaidya signal tables in the database for the specified
// ticker. Sets first_date and last_date for the ticker as 5 years ago and
// today, respectively.
//
// Return all vaidya signals found for the specified ticker.
func AddTicker(ticker string, db *sqlx.DB) []VaidyaSignalsTable {
	// determine first and last date
	lastDate := time.Now()
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

	return InitializeDaysAndVaidyaSignalsTables(TickersTable{Ticker: ticker, FirstDate: firstDate, LastDate: lastDate}, db)
}

func InitializeDaysAndVaidyaSignalsTables(ticker TickersTable, db *sqlx.DB) []VaidyaSignalsTable {
	// get day bars for past 7 years from Alpaca
	fiveYearsBeforeFirstDate := ticker.FirstDate.AddDate(-5, 0, 0)
	bars, _ := marketdata.GetBars(ticker.Ticker, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start:     fiveYearsBeforeFirstDate,
		End:       time.Now().Add(-15 * time.Minute), // 15 minutes are subtracted due to Alpaca free-tier limitations
	})

	// print first and last bars
	// fmt.Println("[ DEBUG ] first bar:", bars[0])
	// fmt.Println("[ DEBUG ] last bar:", bars[len(bars)-1])

	// loop over days to generate indicators
	var last12EMA float64
	var last26EMA float64
	var lastAvgGain float64 = 0 // for RSI
	var lastAvgLoss float64 = 0 // for RSI
	var days = make([]DaysTable, len(bars))
	var startDateIdx int = 0
	for i, bar := range bars {
		// copy Alpaca Bars into DaysTable
		days[i].Date = bar.Timestamp
		days[i].Open = bar.Open
		days[i].High = bar.High
		days[i].Low = bar.Low
		days[i].Close = bar.Close
		days[i].Volume = bar.Volume

		// if date is start date, set firstDateIdx
		if days[i].Date.Year() == ticker.FirstDate.Year() && days[i].Date.YearDay() == ticker.FirstDate.YearDay() {
			fmt.Println("found start date:", days[i].Date)
			startDateIdx = i
		}

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

	// only include days starting from startDateIdx
	days = days[startDateIdx:]

	//------------------------------------------------------------------------
	// write days to database
	//------------------------------------------------------------------------

	// instead of the above for loop where one insert is done for each day, we
	// can do one insert for all days
	query := `INSERT INTO days_` + ticker.Ticker + ` (date, open, high, low, close, volume, macd, rsi)
							VALUES (:date, :open, :high, :low, :close, :volume, :macd, :rsi)`

	// insert row in db (doesn't check for errors)
	db.NamedExec(query, days)

	//------------------------------------------------------------------------
	// update last date in database
	//------------------------------------------------------------------------

	// update last date
	ticker.LastDate = days[len(days)-1].Date

	// define query
	query = `UPDATE tickers SET last_date = :last_date WHERE ticker = :ticker`

	// update row in db (doesn't check for errors)
	db.NamedExec(query, ticker)

	//------------------------------------------------------------------------
	// find local lows
	//
	// NOTE: local lows are defined as having a lower close than the three
	//       previous days and the three following days
	//------------------------------------------------------------------------

	var lows []int // indexes of lows
find_lows:
	for i, day := range days {
		// skip lows that are before first date
		if day.Date.Before(ticker.FirstDate) {
			continue
		}

		// skip first and last three bars
		if i < LOW_DETECTION || i >= len(days)-LOW_DETECTION {
			continue
		}

		// skip if close is not a low
		for j := 1; j <= LOW_DETECTION; j++ {
			if day.Close > days[i-j].Close || day.Close > days[i+j].Close {
				continue find_lows
			}
		}

		lows = append(lows, i)
	}

	//------------------------------------------------------------------------
	// find signal triggers
	//------------------------------------------------------------------------

	signals := make([]VaidyaSignalsTable, 0)
	for i, daysIdx := range lows {
		// skip first low, as a previous low is needed
		if i == 0 {
			continue
		}

		// skip if current low is not lower (i.e. is greater) than previous low
		// NOTE: if lows are equal, it will not skip
		if days[daysIdx].Close > days[lows[i-1]].Close {
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
		if days[daysIdx].Volume > days[lows[i-1]].Volume {
			continue
		}

		query := `INSERT INTO vaidya_` + ticker.Ticker + ` (trigger_date, low_2_date, low_1_date)
							VALUES (:trigger_date, :low_2_date, :low_1_date)`

		signal := VaidyaSignalsTable{
			TriggerDate: days[daysIdx+LOW_DETECTION].Date,
			Low2Date:    days[daysIdx].Date,
			Low1Date:    days[lows[i-1]].Date,
		}

		// insert row in db (doesn't check for errors)
		db.NamedExec(query, signal)

		signals = append(signals, signal)
	}

	return signals
}
