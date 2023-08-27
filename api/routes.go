package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/satvikprasad/vikingx/db"
	"github.com/satvikprasad/vikingx/models"
	"github.com/satvikprasad/vikingx/okx"
)

type Context struct {
	a  *okx.OkApi
	db *db.DbInstance
	c  *gin.Context
}

func ListenAndServe(db *db.DbInstance, a *okx.OkApi, port string) {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))
	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	r.Static("/css", "public/css")
	r.Static("/js", "public/js")

	r.POST("/api/webhook", makeAPIFunc(a, db, handleWebhook))
	r.POST("/api/createTrade", makeAPIFunc(a, db, handleCreateTrade))
	r.GET("/api/trades", makeAPIFunc(a, db, handleTrades))
	r.GET("/api/balance", makeAPIFunc(a, db, handleBalance))
	r.GET("/api/bidask/:ticker", makeAPIFunc(a, db, handleBidAsk))
	r.GET("/api/instruments/:instType", makeAPIFunc(a, db, handleInstruments))

	r.Run(":" + port)
}

func handleCreateTrade(c *Context) error {
	trade := new(models.Trade)
	if err := c.c.BindJSON(trade); err != nil {
		return err
	}

	c.db.Db.Create(&trade)

	writeJSON(c.c, http.StatusCreated, trade)
	return nil
}

func handleTrades(c *Context) error {
	trades := []models.Trade{}

	c.db.Db.Find(&trades)

	writeJSON(c.c, http.StatusOK, trades)
	return nil
}

func handleBalance(c *Context) error {
	bal, err := c.a.GetBalance("USDT")
	if err != nil {
		return err
	}

	writeJSON(c.c, http.StatusOK, []string{fmt.Sprintf("%f", bal)})
	return nil
}

func handleBidAsk(c *Context) error {
	bid, ask, err := c.a.GetLimitSwapPrice(c.c.Param("ticker"))
	if err != nil {
		return err
	}

	writeJSON(c.c, http.StatusOK, []float64{
		bid,
		ask,
	})
	return nil
}

func handleInstruments(c *Context) error {
	instruments, err := c.a.GetInstruments(c.c.Param("instType"))
	if err != nil {
		return err
	}

	writeJSON(c.c, http.StatusOK, instruments)
	return nil
}

func handleWebhook(c *Context) error {
	t := WebhookRequest{}

	if err := c.c.BindJSON(&t); err != nil {
		return fmt.Errorf("Error binding to json body")
	}

	if t.Passphrase != os.Getenv("WEBHOOK_PHRASE") {
		return fmt.Errorf("Error")
	}

	switch strings.Contains(t.Ticker, ".P") || strings.Contains(t.Ticker, "SWAP") {
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

	c.db.Db.Create(&trade)

	writeJSON(c.c, http.StatusOK, trade)
	return nil
}

type apiFunc func(c *Context) error

func makeAPIFunc(a *okx.OkApi, db *db.DbInstance, fn apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			a:  a,
			db: db,
			c:  c,
		}
		if err := fn(ctx); err != nil {
			writeJSON(c, http.StatusInternalServerError,
				map[string]string{"error": err.Error()})
		}
	}
}

func writeJSON(c *gin.Context, code int, v any) {
	c.Header("Content-Type", "application/json")
	c.JSON(code, v)
}
