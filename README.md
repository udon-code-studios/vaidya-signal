# Vaidya Signal

## Vision

1. A program which at the end of each business day can scan a list of tickers to check for triggers of the Vaidya Signal. It should show all instances when the Vaidya Signal triggered in the past 5 years (using a bar size of 1 day).
2. Enable the script to be run for weekly and monthly bars.
3. A GUI and service to visualize when the signal triggered.

## The Vaidya Signal (Bottom Divergence)

The Vaidya Signal is triggered when these three conditions are met:
- the current low* (L2) is lower than the previous low (L1)
- MACD and RSI are both higher at L2 than they were at L1
- volume at the L2 is lower than it was at L1

\*a low is defined as having 3 days before and after that are higher than the low (using close values)

## TODO

- [x] create a program which can loop over a list of tickers
- [x] collect day bars each ticker for the past 6 years
- [x] calculate MACD and RSI for each day (using two years as ramp-up)
- [x] find all instances of the Vaidya Signal triggering over the past 5 years
- [x] output instances to a JSON file