package trader

import "time"

type Position struct {
	Size   float64
	Type   string
	Symbol string
}

type Ticker struct {
	Symbol          string
	BidPrice        float64
	AskPrice        float64
	Vol24H          float64
	LastTradedPrice float64
}

type Instrument struct {
	Symbol    string
	BaseCcy   string
	QuoteCcy  string
	CtValCcy  string
	SettleCcy string
	CtVal     string
}

type Candlestick struct {
	Timestamp              time.Time
	Open, High, Low, Close float64
}

type Trader interface {
	Positions() ([]Position, error)
	Tickers(string) ([]Ticker, error)
	Instruments(string) ([]Instrument, error)
	LimitSwapPrice(string) (float64, float64, error)
	Candlesticks(string, string) ([]Candlestick, error)
	LimitOrder(string, string, string, float64, float64) error
	MarketOrder(string, string, float64) error
	SetLeverage(string, int) error
	Balance(string) (float64, error)
	MarkPrice(string) (float64, error)
	TickerCtSize(string) (float64, error)
}
