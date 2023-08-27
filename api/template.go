package api

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/satvikprasad/vikingx/db"
	"github.com/satvikprasad/vikingx/models"
	"github.com/satvikprasad/vikingx/okx"
)

func createTemplateRoutes(r *gin.Engine, a *okx.OkApi, db *db.DbInstance) {
	r.GET("/", makeAPIFunc(a, db, handleHomeTemplate))

	r.GET("/trades", makeAPIFunc(a, db, handleTradesTemplate))
	r.GET("/instruments", makeAPIFunc(a, db, handleInstrumentsTemplate))
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

	fmt.Println("lskdjf")

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
