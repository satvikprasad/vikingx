package routes

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satvikprasad/vikingx/models"
	"github.com/satvikprasad/vikingx/server"
)

func RenderPositionsList(c *server.Context) error {
	positions, err := c.Trader.Positions()
	if err != nil {
		return err
	}

	c.Context.HTML(http.StatusOK, "home/positions.tmpl", positions)
	return nil
}

func RenderPlaceMarketOrder(c *server.Context) error {
	ticker := c.Context.Request.PostFormValue("ticker")
	side := c.Context.Request.PostFormValue("side")

	size, err := strconv.ParseFloat(c.Context.Request.PostFormValue("size"), 64)
	if err != nil {
		return err
	}

	sz, err := c.Trader.TickerCtSize(ticker)
	if err != nil {
		sz = 1.0
	}

	if err := c.Trader.MarketOrder(ticker, side, float64(size)/sz); err != nil {
		return err
	}

	market, err := c.Trader.MarkPrice(ticker)
	if err != nil {
		return err
	}

	trade := models.Trade{
		Ticker: ticker,
		Side:   side,
		Size:   size,
		Price:  market,
	}
	c.Database.CreateTrade(&trade)

	positions, err := c.Trader.Positions()
	if err != nil {
		return err
	}

	c.Context.HTML(http.StatusOK, "home/positions.tmpl", positions)
	return nil
}

func RenderInstruments(c *server.Context) error {
	tickers, err := c.Trader.Tickers("SWAP")
	if err != nil {
		return err
	}

	sort.Slice(tickers, func(a, b int) bool {
		return tickers[a].Vol24H > tickers[b].Vol24H
	})

	tickers = tickers[0:10]

	c.Context.HTML(http.StatusOK, "home/instruments.tmpl", tickers)
	return nil
}

func RenderTrades(c *server.Context) error {
	trades := c.Database.Trades()

	sort.Slice(trades, func(a, b int) bool {
		return trades[a].CreatedAt.Unix() > trades[b].CreatedAt.Unix()
	})

	c.Context.HTML(http.StatusOK, "home/trades.tmpl", trades)
	return nil
}

func RenderHome(c *server.Context) error {
	c.Context.HTML(http.StatusOK, "home/index.tmpl", nil)
	return nil
}

func writeJSON(c *gin.Context, code int, v any) {
	c.Header("Content-Type", "application/json")
	c.JSON(code, v)
}
