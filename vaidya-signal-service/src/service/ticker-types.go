package service

import "time"

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
	Ticker    string    `db:"ticker"`
	FirstDate time.Time `db:"first_date"`
	LastDate  time.Time `db:"last_date"`
}

type VaidyaSignalsTable struct {
	TriggerDate time.Time `db:"trigger_date" json:"trigger_date"` // day signal was triggered
	Low2Date    time.Time `db:"low_2_date" json:"low_2_date"`     // current low
	Low1Date    time.Time `db:"low_1_date" json:"low_1_date"`     // previous low
}

type WatchlistTable struct {
	Ticker      string    `db:"ticker"`
	LastTrigger time.Time `db:"last_trigger"`
}
