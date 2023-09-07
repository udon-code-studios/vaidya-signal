package service

import (
	"fmt"
	"time"

	U "github.com/udon-code-sudios/vaidya-signal-service/utils"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

// define constants
const EMA12_SMOOTHING float64 = 2
const EMA26_SMOOTHING float64 = 2
const RSI_PERIOD int = 14
const LOW_DETECTION int = 3 // # of days before and after a low that should have higher closes

func GetHistoricalVaidyaSignals(ticker string) []VaidyaSignalsTable {
	today := time.Now().Add(-15 * time.Minute) // 15 minutes are subtracted due to Alpaca free-tier limitations
	fiveYearsAgo := today.AddDate(-5, 0, 0)
	tenYearsAgo := today.AddDate(-10, 0, 0)

	fmt.Println("[ INFO ] Starting detection for ticker:", ticker)

	// get day bars for past 7 years from Alpaca
	bars, _ := marketdata.GetBars(ticker, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start:     tenYearsAgo,
		End:       today,
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

	// // loop over days to generate output file
	// for _, day := range days {
	// 	// skip until 5 years ago date
	// 	if day.Date.Before(fiveYearsAgo) {
	// 		continue
	// 	}

	// 	query := `INSERT INTO days(ticker, date, open, high, low, close, volume, macd, rsi)
	// 						VALUES(:ticker, :date, :open, :high, :low, :close, :volume, :macd, :rsi)`

	// 	// insert row in db (doesn't check for errors)
	// 	db.NamedExec(query, day)
	// }

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

	signals := make([]VaidyaSignalsTable, 0)
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

		// query := `INSERT INTO vaidya_` + strings.ToLower(ticker) + ` (ticker, trigger_date, low_2_date, low_1_date)
		// 					VALUES(:ticker, :trigger_date, :low_2_date, :low_1_date)`

		// // insert row in db (doesn't check for errors)
		// db.NamedExec(query, VaidyaSignalsTable{
		// 	TriggerDate: days[daysIdx+LOW_DETECTION].Date,
		// 	Low2Date:    days[daysIdx].Date,
		// 	Low1Date:    days[lows[i-1]].Date,
		// })

		signals = append(signals, VaidyaSignalsTable{
			TriggerDate: days[daysIdx+LOW_DETECTION].Date,
			Low2Date:    days[daysIdx].Date,
			Low1Date:    days[lows[i-1]].Date,
		})
	}

	return signals
}
