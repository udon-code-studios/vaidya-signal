package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
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
		var indicators = make([]Indicators, len(bars))
		for i, bar := range bars {
			// skip first 25 bars
			if i < 25 {
				continue
			}

			// initialize EMAs and Avg Gain/Loss on bar 26
			if i == 25 {
				last12EMA = calcBarCloseSMA(bars[i-11 : i+1])
				last26EMA = calcBarCloseSMA(bars[i-25 : i+1])

				lastAvgGain = calcFirstAvgGainLoss(bars[i-RSI_PERIOD:i+1], true)
				lastAvgLoss = calcFirstAvgGainLoss(bars[i-RSI_PERIOD:i+1], false)

				continue
			}

			last12EMA = calcEMA(bar.Close, last12EMA, 12, EMA12_SMOOTHING)
			last26EMA = calcEMA(bar.Close, last26EMA, 26, EMA26_SMOOTHING)
			indicators[i].MACD = last12EMA - last26EMA

			lastAvgGain = calcAvgGainLoss(RSI_PERIOD, lastAvgGain, bars[i-1].Close, bar.Close, true)
			lastAvgLoss = calcAvgGainLoss(RSI_PERIOD, lastAvgLoss, bars[i-1].Close, bar.Close, false)
			indicators[i].RSI = 100 - 100/(1+(lastAvgGain/lastAvgLoss))
		}

		// define output directory and filenames
		outputDirectory := "tickers"
		outputDataFilename := fmt.Sprintf("%s_data.csv", ticker)
		outputMetaFilename := fmt.Sprintf("%s_meta.json", ticker)

		// create output directory
		err = os.MkdirAll(outputDirectory, 0755)
		panicOnNotNil(err)

		//------------------------------------------------------------------------
		// write data file
		//------------------------------------------------------------------------

		// create output data file and writer
		dataFile, err := os.Create(fmt.Sprintf("%s/%s", outputDirectory, outputDataFilename))
		panicOnNotNil(err)
		dataWriter := csv.NewWriter(dataFile)
		defer dataWriter.Flush()

		// write data columns header
		dataWriter.Write([]string{
			"date",
			"open",
			"high",
			"low",
			"close",
			"volume",
			"MACD",
			"RSI",
		})

		// loop over days to generate output file
		for i, bar := range bars {
			// skip until 5 years ago date
			if bar.Timestamp.Before(fiveYearsAgo) {
				continue
			}

			// write row
			dataWriter.Write([]string{
				bar.Timestamp.Format("2006-01-02"),
				fmt.Sprintf("%.3f", bar.Open),
				fmt.Sprintf("%.3f", bar.High),
				fmt.Sprintf("%.3f", bar.Low),
				fmt.Sprintf("%.3f", bar.Close),
				fmt.Sprintf("%d", bar.Volume),
				fmt.Sprintf("%.3f", indicators[i].MACD),
				fmt.Sprintf("%.3f", indicators[i].RSI),
			})
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

		var triggers []Vaidya
		for i, barIdx := range lows {
			// skip first low, as a previous low is needed
			if i == 0 {
				continue
			}

			// skip if current low is not lower (i.e. is greater) than previous low
			// NOTE: if lows are equal, it will not skip
			if bars[barIdx].Close > bars[lows[i-1]].Close {
				continue
			}

			// skip if current low's MACD and RSI are not both higher than the
			// previous low's
			// NOTE: if MACDs and RSIs are equal, it will not skip
			if indicators[barIdx].MACD < indicators[lows[i-1]].MACD ||
				indicators[barIdx].RSI < indicators[lows[i-1]].RSI {
				continue
			}

			// skip if current low's volume is not lower (i.e. is greater) than the
			// previous low's
			// NOTE: if volumes are equal, it will not skip
			if bars[barIdx].Volume > bars[lows[i-1]].Volume {
				continue
			}

			// save signal trigger
			triggers = append(triggers, Vaidya{
				TriggerDate: bars[barIdx+LOW_DETECTION].Timestamp.Format("2006-01-02"),
				Low2Date:    bars[barIdx].Timestamp.Format("2006-01-02"),
				Low1Date:    bars[lows[i-1]].Timestamp.Format("2006-01-02"),
			})
		}

		//------------------------------------------------------------------------
		// write signals to metadata file
		//------------------------------------------------------------------------

		// create output metadata file
		metaFile, err := os.Create(fmt.Sprintf("%s/%s", outputDirectory, outputMetaFilename))
		panicOnNotNil(err)
		defer metaFile.Close()

		// encode the metadata to json
		jsonMetadata, err := json.Marshal(map[string]interface{}{"vaidyaSignals": triggers})
		panicOnNotNil(err)

		// write json to file
		metaFile.Write(jsonMetadata)
	}
}

//----------------------------------------------------------------------------
// types
//----------------------------------------------------------------------------

type Indicators struct {
	MACD float64
	RSI  float64
}

type Vaidya struct {
	TriggerDate string `json:"triggerDate"` // day signal was triggered
	Low2Date    string `json:"low2Date"`    // current low
	Low1Date    string `json:"low1Date"`    // previous low
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

func calcBarCloseSMA(bars []marketdata.Bar) float64 {
	// Extract an array of just close values.
	closeValues := []float64{}
	for _, bar := range bars {
		closeValues = append(closeValues, bar.Close)
	}

	return calcSMA(closeValues)
}

func calcSMA(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

func calcEMA(price, lastEMA float64, period int, smoothing float64) float64 {
	multiplier := smoothing / float64(period+1)
	return price*multiplier + lastEMA*(1-multiplier)
}

/*
gainOrLoss bool: true -> calculate avg. gain, false -> calculate avg. loss
*/
func calcFirstAvgGainLoss(bars []marketdata.Bar, gainOrLoss bool) float64 {
	var gainLossSum float64 = 0

	// sum gain or loss
	for i, bar := range bars {
		// skip first bar
		if i == 0 {
			continue
		}

		if gainOrLoss {
			// calculate gain
			gainLossSum += math.Max(bar.Close-bars[i-1].Close, 0)
		} else {
			// sum loss
			gainLossSum += math.Max(bars[i-1].Close-bar.Close, 0)
		}
	}

	return gainLossSum / float64(len(bars)-1)
}

/*
gainOrLoss bool: true -> calculate avg. gain, false -> calculate avg. loss
*/
func calcAvgGainLoss(period int, prevGainLoss float64, last, current float64, gainOrLoss bool) float64 {
	var gainLoss float64
	if gainOrLoss {
		// calculate gain
		gainLoss = math.Max(current-last, 0)
	} else {
		// calculate loss
		gainLoss = math.Max(last-current, 0)
	}

	return (prevGainLoss*float64(period-1) + gainLoss) / float64(period)
}
