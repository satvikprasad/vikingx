package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
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
		OrderContracts         int     `json:"order_contracts"`
		OrderPrice             float64 `json:"order_price"`
		OrderID                string  `json:"order_id"`
		MarketPosition         string  `json:"market_position"`
		MarketPositionSize     int     `json:"market_position_size"`
		PrevMarketPosition     string  `json:"prev_market_position"`
		PrevMarketPositionSize int     `json:"prev_market_position_size"`
	} `json:"strategy"`
}

func handleWebhook(c *Context) error {
	godotenv.Load(".env")

	t := WebhookRequest{}

	if err := c.c.BindJSON(&t); err != nil {
		return fmt.Errorf("Error binding to json body")
	}

	if t.Passphrase != os.Getenv("WEBHOOK_PHRASE") {
		return fmt.Errorf("Error")
	}

	switch strings.Contains(t.Ticker, ".P") {
	case true:
		ticker, err := c.a.ConvertTickerName("SWAP", t.Ticker)
		if err != nil {
			return err
		}

		sz, err := c.a.GetTickerCtSize(t.Ticker)
		if err != nil {
			return err
		}

		if err := c.a.SetLeverage(ticker, 20); err != nil {
			fmt.Printf("Could not set leverage: %s\n", err)
		}

		if err := c.a.MarketOrderSwap(ticker, t.Strategy.OrderAction,
			float64(t.Strategy.OrderContracts)/sz); err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error placing order: %s", err)
		}

		writeJSON(c.c, http.StatusOK, "Order placed")
	case false:
		ticker, err := c.a.ConvertTickerName("SPOT", t.Ticker)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error decoding ticker info: %s", err)
		}

		if err := c.a.MarketOrder(ticker, t.Strategy.OrderAction,
			float64(t.Strategy.OrderContracts)); err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error placing order: %s", err)
		}

		writeJSON(c.c, http.StatusOK, "Order placed")
	default:
		return fmt.Errorf("Could not decode ticker")
	}

	return nil
}
