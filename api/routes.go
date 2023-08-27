package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
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
	r.LoadHTMLGlob("templates/**/*")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

	createWebhookRoutes(r, a, db)
	createTemplateRoutes(r, a, db)
	createApiRoutes(r, a, db)

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

func createApiRoutes(r *gin.Engine, a *okx.OkApi, db *db.DbInstance) {
	r.POST("/api/create-trade", makeAPIFunc(a, db, handleCreateTrade))

	r.GET("/api/trades", makeAPIFunc(a, db, handleTrades))
	r.GET("/api/balance", makeAPIFunc(a, db, handleBalance))

	r.GET("/api/bidask/:ticker", makeAPIFunc(a, db, handleBidAsk))
	r.GET("/api/instruments/:instType", makeAPIFunc(a, db, handleInstruments))
}

type apiFunc func(c *Context) error

func makeAPIFunc(a *okx.OkApi, db *db.DbInstance, fn apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			writeJSON(c, http.StatusInternalServerError,
				map[string]string{"error": "Could not initialise database"})
			return
		}

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
