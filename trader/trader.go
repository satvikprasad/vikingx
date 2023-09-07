package trader

type Trader interface {
	Positions() ([]OkxPosition, error)
	Tickers(string) ([]OkxTicker, error)
	Instruments(string) ([]OkxInstruments, error)
	LimitSwapPrice(string) (float64, float64, error)
	CandleSticks(string, string) ([]OkxCandlestick, error)
	LimitOrder(string, string, string, float64, float64) error
	MarketOrderSwap(string, string, float64) error
	MarketOrder(string, string, float64) error
	SetLeverage(string, int) error
	Balance(string) (float64, error)
	MarkPrice(string) (float64, error)
	TickerCtSize(string) (float64, error)
}
