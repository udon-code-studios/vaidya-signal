package utils

import (
	"math"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

func CalcBarCloseSMA(bars []marketdata.Bar) float64 {
	// Extract an array of just close values.
	closeValues := []float64{}
	for _, bar := range bars {
		closeValues = append(closeValues, bar.Close)
	}

	return CalcSMA(closeValues)
}

func CalcSMA(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

func CalcEMA(price, lastEMA float64, period int, smoothing float64) float64 {
	multiplier := smoothing / float64(period+1)
	return price*multiplier + lastEMA*(1-multiplier)
}

/*
gainOrLoss bool: true -> calculate avg. gain, false -> calculate avg. loss
*/
func CalcFirstAvgGainLoss(bars []marketdata.Bar, gainOrLoss bool) float64 {
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
func CalcAvgGainLoss(period int, prevGainLoss float64, last, current float64, gainOrLoss bool) float64 {
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
