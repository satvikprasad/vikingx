package routes

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/satvikprasad/vikingx/models"
	"github.com/satvikprasad/vikingx/server"
	"github.com/satvikprasad/vikingx/trader"
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

	writeJSON(c.Context, http.StatusOK, map[string]string{"balance": fmt.Sprintf("%f", bal)})
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
	symbol := c.Context.Param("symbol")

	candles, err := trader.GetCandles(symbol, "D", time.Date(2019, time.January, 0, 0, 0, 0, 0, time.Now().Location()))
	if err != nil {
		return err
	}

	writeJSON(c.Context, http.StatusOK, candles)
	return nil
}

func HandleWebhook(c *server.Context) error {
	r := WebhookRequest{}

	if err := c.Context.BindJSON(&r); err != nil {
		return fmt.Errorf("Error binding to json body")
	}

	if r.Passphrase != os.Getenv("WEBHOOK_PHRASE") {
		return fmt.Errorf("Error")
	}

	ticker, err := convertTickerName(r.Ticker)
	if err != nil {
		return err
	}

	sz, err := c.Trader.ContractSize(r.Ticker)
	if err != nil {
		sz = 1.0
	}

	if err := c.Trader.MarketOrder(ticker, r.Strategy.OrderAction,
		float64(r.Strategy.OrderContracts)/sz); err != nil {
		return fmt.Errorf("Error placing order: %s", err)
	}

	trade := models.Trade{
		Ticker:                 r.Ticker,
		Side:                   r.Strategy.OrderAction,
		Size:                   r.Strategy.OrderContracts,
		Price:                  r.Strategy.OrderPrice,
		MarketPositionSize:     r.Strategy.MarketPositionSize,
		PrevMarketPositionSize: r.Strategy.PrevMarketPositionSize,
	}

	c.Database.CreateTrade(&trade)

	writeJSON(c.Context, http.StatusOK, trade)
	return nil
}

func convertTickerName(ticker string) (string, error) {
	reg, err := regexp.Compile("(.{2,4})(USDT|USD)(.P)*")
	if err != nil {
		return "", err
	}

	matches := reg.FindAllStringSubmatch(ticker, -1)
	if matches == nil {
		return "", fmt.Errorf("o ticker matches")
	}

	if len(matches[0]) < 4 {
		return "", fmt.Errorf("ticker not formatted properly")
	}

	base := matches[0][1]
	quote := matches[0][2]
	perp := matches[0][3]

	if perp == "" {
		return fmt.Sprintf("%s-%s", base, quote), nil
	}

	return fmt.Sprintf("%s-%s-SWAP", base, quote), nil
}
