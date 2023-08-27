package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satvikprasad/vikingx/db"
	"github.com/satvikprasad/vikingx/models"
	"github.com/satvikprasad/vikingx/okx"
)

type Position struct {
	Size   float64
	Type   string
	Symbol string
}

func createTemplateRoutes(r *gin.Engine, a *okx.OkxApi, db db.Database) {
	r.GET("/", makeTemplateAPIFunc(a, db, handleHomeTemplate))

	r.GET("/positions-list", makeTemplateAPIFunc(a, db,
		handlePositionsListTemplate))
	r.GET("/trades-list", makeTemplateAPIFunc(a, db,
		handleTradesTemplate))
	r.GET("/instruments-list", makeTemplateAPIFunc(a, db,
		handleInstrumentsTemplate))

	r.POST("/market-order", makeTemplateAPIFunc(a, db,
		handlePlaceMarketOrderTemplate))
}

func handlePositionsListTemplate(c *Context) error {
	okxPositions, err := c.a.Positions()
	if err != nil {
		return err
	}

	positions := []Position{}
	for _, o := range okxPositions {
		cx := 1.0
		if o.InstType == "SWAP" {
			ctSize, err := c.a.TickerCtSize(o.InstID)
			if err != nil {
				return err
			}
			cx = ctSize
		}

		positionSize, err := strconv.ParseFloat(o.Pos, 64)
		if err != nil {
			return err
		}

		positions = append(positions, Position{
			Size:   positionSize * cx,
			Symbol: o.InstID,
			Type:   o.InstType,
		})
	}

	writeHTML(c.c, http.StatusOK, "home/positionslist.tmpl", positions)
	return nil
}

func handlePlaceMarketOrderTemplate(c *Context) error {
	ticker := c.c.Request.PostFormValue("ticker")
	side := c.c.Request.PostFormValue("side")

	size, err := strconv.ParseFloat(c.c.Request.PostFormValue("size"), 64)
	if err != nil {
		return err
	}

	sz, err := c.a.TickerCtSize(ticker)
	if err != nil {
		return err
	}

	if err := c.a.MarketOrderSwap(ticker, side, float64(size)/sz); err != nil {
		return err
	}

	market, err := c.a.MarkPrice(ticker)
	if err != nil {
		return err
	}

	trades := c.db.Trades()

	fmt.Println(market)

	trade := models.Trade{
		Ticker: ticker,
		Side:   side,
		Size:   size,
		Price:  market,
	}
	c.db.CreateTrade(&trade)

	trades = append(trades, &trade)

	sort.Slice(trades, func(a, b int) bool {
		return trades[a].CreatedAt.Unix() > trades[b].CreatedAt.Unix()
	})

	writeHTML(c.c, http.StatusOK, "home/tradelist.tmpl", trades)
	return nil
}

func handleInstrumentsTemplate(c *Context) error {
	tickers, err := c.a.Tickers("SWAP")
	if err != nil {
		return err
	}

	sort.Slice(tickers, func(a, b int) bool {
		return tickers[a].VolCcy24H > tickers[b].VolCcy24H
	})

	tickers = tickers[0:10]

	writeHTML(c.c, http.StatusOK, "home/instrumentslist.tmpl", tickers)
	return nil
}

func handleTradesTemplate(c *Context) error {
	trades := c.db.Trades()

	sort.Slice(trades, func(a, b int) bool {
		return trades[a].CreatedAt.Unix() > trades[b].CreatedAt.Unix()
	})

	writeHTML(c.c, http.StatusOK, "home/tradelist.tmpl", trades)
	return nil
}

func handleHomeTemplate(c *Context) error {
	writeHTML(c.c, http.StatusOK, "home/index.tmpl", nil)
	return nil
}

func makeTemplateAPIFunc(a *okx.OkxApi, db db.Database, fn apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			writeHTML(c, http.StatusInternalServerError,
				"error/index.tmpl",
				map[string]string{"Error": "Could not initialise database"})
			return
		}

		ctx := &Context{
			a:  a,
			db: db,
			c:  c,
		}

		if err := fn(ctx); err != nil {
			fmt.Println(err.Error())
			writeHTML(ctx.c, http.StatusOK,
				"error/index.tmpl",
				map[string]string{"Error": err.Error()})
		}
	}
}

func writeHTML(c *gin.Context, code int, tmpl string, d any) {
	c.HTML(code, tmpl, d)
}
