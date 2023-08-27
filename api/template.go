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

func createTemplateRoutes(r *gin.Engine, a *okx.OkApi, db *db.DbInstance) {
	r.GET("/", makeTemplateAPIFunc(a, db, handleHomeTemplate))

	r.GET("/trades", makeTemplateAPIFunc(a, db, handleTradesTemplate))
	r.GET("/instruments", makeTemplateAPIFunc(a, db, handleInstrumentsTemplate))

	r.POST("/market-order", makeTemplateAPIFunc(a, db, handlePlaceMarketOrderTemplate))
}

func handlePlaceMarketOrderTemplate(c *Context) error {
	ticker := c.c.Request.PostFormValue("ticker")
	side := c.c.Request.PostFormValue("side")

	size, err := strconv.ParseFloat(c.c.Request.PostFormValue("size"), 64)
	if err != nil {
		return err
	}

	sz, err := c.a.GetTickerCtSize(ticker)
	if err != nil {
		return err
	}

	if err := c.a.MarketOrderSwap(ticker, side, float64(size)/sz); err != nil {
		return err
	}

	trades := []models.Trade{}
	c.db.Db.Find(&trades)

	trade := models.Trade{
		Ticker: ticker,
		Side:   side,
		Size:   1,
	}
	c.db.Db.Create(&trade)

	trades = append(trades, trade)

	sort.Slice(trades, func(a, b int) bool {
		return trades[a].CreatedAt.Unix() > trades[b].CreatedAt.Unix()
	})

	writeHTML(c.c, http.StatusOK, "home/tradelist.tmpl", trades)
	return nil
}

func handleInstrumentsTemplate(c *Context) error {
	tickers, err := c.a.GetTickers("SWAP")
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
	trades := []models.Trade{}

	c.db.Db.Find(&trades)

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

func makeTemplateAPIFunc(a *okx.OkApi, db *db.DbInstance, fn apiFunc) gin.HandlerFunc {
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
