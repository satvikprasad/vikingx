package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/satvikprasad/vikingx/okx"
	"net/http"
	"os"
)

type Context struct {
	a *okx.OkApi
	c *gin.Context
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

func main() {
	godotenv.Load(".env")

	a := okx.NewOkApi(true, ".env")
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))
	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	r.Static("/css", "public/css")
	r.Static("/js", "public/js")

	r.POST("/webhook", makeAPIFunc(a, handleWebhook))
	r.GET("/api/balance", makeAPIFunc(a, handleBalance))
	r.GET("/api/bidask/:ticker", makeAPIFunc(a, handleBidAsk))
	r.GET("/api/instruments/:instType", makeAPIFunc(a, handleInstruments))

	r.Run(":" + os.Getenv("PORT"))
}

type apiFunc func(c *Context) error

func makeAPIFunc(a *okx.OkApi, fn apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			a: a,
			c: c,
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
