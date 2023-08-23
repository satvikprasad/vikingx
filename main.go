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

type WebhookContext struct {
	a *okx.OkApi
	c *gin.Context
}

func main() {
	godotenv.Load(".env")

	a := okx.NewOkApi(true, ".env")
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	bal, err := a.GetBalance("USDT")
	if err != nil {
		fmt.Println(err)
	}

	r.Use(cors.New(config))
	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	r.Static("/css", "public/css")
	r.Static("/js", "public/js")

	r.POST("/webhook", makeWebhookHttpAPIFunc(a, handleWebhook))
	r.GET("/api/balance", func(c *gin.Context) {
		c.JSON(200, []string{
			fmt.Sprintf("%f", bal),
		})
	})
	r.GET("/api/ticker", func(c *gin.Context) {
		buy, sell, err := a.GetLimitSwapPrice("BTC-USDT-SWAP")
		if err != nil {
			writeJSON(c, http.StatusInternalServerError, err)
		}

		writeJSON(c, http.StatusOK, []float64{
			buy,
			sell,
		})
	})

	r.Run(":" + os.Getenv("PORT"))
}

type webhookFunc func(c *WebhookContext) error

func makeWebhookHttpAPIFunc(a *okx.OkApi, fn webhookFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &WebhookContext{
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
