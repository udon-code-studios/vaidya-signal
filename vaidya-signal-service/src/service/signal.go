package service

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	U "github.com/udon-code-sudios/vaidya-signal-service/utils"
)

// define constants
const EMA12_SMOOTHING float64 = 2
const EMA26_SMOOTHING float64 = 2
const RSI_PERIOD int = 14
const LOW_DETECTION int = 3 // # of days before and after a low that should have higher closes

func FindAllVaidyaSignalsForTickers(tickers []string) map[string][]VaidyaSignal {
	// define map to return
	signals := make(map[string][]VaidyaSignal)

	// loop through tickers 25 at a time
	n := 25
	for i := 0; i < len(tickers); i += n {
		// get tickers to process
		end := i + n
		if end > len(tickers) {
			end = len(tickers)
		}

		fmt.Println("[ INFO ] Processing tickers:", tickers[i:end])

		// get day bars for tickers (up to 10 years)
		today := time.Now().Add(-15 * time.Minute) // 15 minutes are subtracted due to Alpaca free-tier limitations
		tenYearsAgo := today.AddDate(-10, 0, 0)
		multiBars, _ := marketdata.GetMultiBars(tickers[i:end], marketdata.GetBarsRequest{
			TimeFrame: marketdata.OneDay,
			Start:     tenYearsAgo,
			End:       today,
		})

		// detect signals for each ticker
		for ticker, bars := range multiBars {
			tickerDays := generateDaysFromAlpacaDayBars(bars)
			tickerSignals := scanDaysForVaidyaSignals(tickerDays)
			signals[ticker] = tickerSignals
		}
	}

	return signals
}

//------------------------------------------------------------------------
// helper functions
//------------------------------------------------------------------------

func generateDaysFromAlpacaDayBars(bars []marketdata.Bar) []Day {
	// make sure there are enough bars (only 26th bar and on are used)
	if len(bars) < 27 {
		return make([]Day, 0)
	}

	var last12EMA float64
	var last26EMA float64
	var lastAvgGain float64 = 0 // for RSI
	var lastAvgLoss float64 = 0 // for RSI
	var days = make([]Day, len(bars))
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

	// only include days starting from 27th day
	return days[26:]
}

func scanDaysForVaidyaSignals(days []Day) []VaidyaSignal {
	//------------------------------------------------------------------------
	// find local lows
	//
	// NOTE: local lows are defined as having a lower close than the three
	//       previous days and the three following days
	//------------------------------------------------------------------------

	var lows []int // indexes of lows
find_lows:
	for i, day := range days {
		// // skip lows that are not from last 5 years
		// if bar.Timestamp.Before(fiveYearsAgo) {
		// 	continue
		// }

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

	signals := make([]VaidyaSignal, 0)
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

		signals = append(signals, VaidyaSignal{
			TriggerDate: days[daysIdx+LOW_DETECTION].Date,
			Low2Date:    days[daysIdx].Date,
			Low1Date:    days[lows[i-1]].Date,
		})
	}

	return signals
}

//------------------------------------------------------------------------
// types
//------------------------------------------------------------------------

type Day struct {
	Date   time.Time `db:"date"`
	Open   float64   `db:"open"`
	High   float64   `db:"high"`
	Low    float64   `db:"low"`
	Close  float64   `db:"close"`
	Volume uint64    `db:"volume"`
	MACD   float64   `db:"macd"`
	RSI    float64   `db:"rsi"`
}

type VaidyaSignal struct {
	TriggerDate time.Time `db:"trigger_date" json:"trigger_date"` // day signal was triggered
	Low2Date    time.Time `db:"low_2_date" json:"low_2_date"`     // current low
	Low1Date    time.Time `db:"low_1_date" json:"low_1_date"`     // previous low
}
