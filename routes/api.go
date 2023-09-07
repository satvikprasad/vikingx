package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/satvikprasad/vikingx/models"
	"github.com/satvikprasad/vikingx/server"
)

type WebhookRequest struct {
	Passphrase string    `json:"passphrase"`
	Time       time.Time `json:"time"`
	Exchange   string    `json:"exchange"`
	Ticker     string    `json:"ticker"`
	Bar        struct {
		Time time.Time `json:"time"`
		Open int       `json:"open"`
		High float64   `json:"high"`
		Low  float64   `json:"low"`

		Close float64 `json:"close"`

		Volume float64 `json:"volume"`
	} `json:"bar"`
	Strategy struct {
		PositionSize           int     `json:"position_size"`
		OrderAction            string  `json:"order_action"`
		OrderContracts         float64 `json:"order_contracts"`
		OrderPrice             float64 `json:"order_price"`
		OrderID                string  `json:"order_id"`
		MarketPosition         string  `json:"market_position"`
		MarketPositionSize     float64 `json:"market_position_size"`
		PrevMarketPosition     string  `json:"prev_market_position"`
		PrevMarketPositionSize float64 `json:"prev_market_position_size"`
	} `json:"strategy"`
}

func CreateTrade(c *server.Context) error {
	trade := new(models.Trade)
	if err := c.Context.BindJSON(trade); err != nil {
		return err
	}

	c.Database.CreateTrade(trade)

	writeJSON(c.Context, http.StatusCreated, trade)
	return nil
}

func Trades(c *server.Context) error {
	trades := c.Database.Trades()

	writeJSON(c.Context, http.StatusOK, trades)
	return nil
}

func Balance(c *server.Context) error {
	bal, err := c.Trader.Balance("USDT")
	if err != nil {
		return err
	}

	writeJSON(c.Context, http.StatusOK, []string{fmt.Sprintf("%f", bal)})
	return nil
}

func BidAsk(c *server.Context) error {
	bid, ask, err := c.Trader.LimitSwapPrice(c.Context.Param("ticker"))
	if err != nil {
		return err
	}

	writeJSON(c.Context, http.StatusOK, []float64{
		bid,
		ask,
	})
	return nil
}

func Instruments(c *server.Context) error {
	instruments, err := c.Trader.Instruments(c.Context.Param("instType"))
	if err != nil {
		return err
	}

	writeJSON(c.Context, http.StatusOK, instruments)
	return nil
}

func Candles(c *server.Context) error {
	candles, err := c.Trader.CandleSticks("ETH-USDT", "1D")
	if err != nil {
		return err
	}

	writeJSON(c.Context, http.StatusOK, candles)
	return nil
}

/**
func HandleWebhook(c *server.Context) error {
	t := WebhookRequest{}

	if err := c.Context.BindJSON(&t); err != nil {
		return fmt.Errorf("Error binding to json body")
	}

	if t.Passphrase != os.Getenv("WEBHOOK_PHRASE") {
		return fmt.Errorf("Error")
	}

	switch strings.Contains(t.Ticker, ".P") || strings.Contains(t.Ticker, "SWAP") {
	case true:
		ticker, err := c.Trader.ConvertTickerName("SWAP", t.Ticker)
		if err != nil {
			return err
		}

		sz, err := c.Trader.TickerCtSize(t.Ticker)
		if err != nil {
			return err
		}

		if err := c.Trader.SetLeverage(ticker, 20); err != nil {
			fmt.Printf("Could not set leverage: %s\n", err)
		}

		if err := c.Trader.MarketOrderSwap(ticker, t.Strategy.OrderAction,
			float64(t.Strategy.OrderContracts)/sz); err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error placing order: %s", err)
		}
	case false:
		ticker, err := c.Trader.ConvertTickerName("SPOT", t.Ticker)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error decoding ticker info: %s", err)
		}

		if err := c.Trader.MarketOrder(ticker, t.Strategy.OrderAction,
			float64(t.Strategy.OrderContracts)); err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error placing order: %s", err)
		}
	default:
		return fmt.Errorf("Could not decode ticker")
	}

	trade := models.Trade{
		Ticker:                 t.Ticker,
		Side:                   t.Strategy.OrderAction,
		Size:                   t.Strategy.OrderContracts,
		Price:                  t.Strategy.OrderPrice,
		MarketPositionSize:     t.Strategy.MarketPositionSize,
		PrevMarketPositionSize: t.Strategy.PrevMarketPositionSize,
	}

	c.Database.CreateTrade(&trade)

	writeJSON(c.Context, http.StatusOK, trade)
	return nil
}
**/
