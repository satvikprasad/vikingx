package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/satvikprasad/vikingx/db"
	"github.com/satvikprasad/vikingx/routes"
	"github.com/satvikprasad/vikingx/server"
	"github.com/satvikprasad/vikingx/trader"
)

func main() {
	godotenv.Load()

	db, err := db.NewDB()
	if err != nil {
		fmt.Printf("Error creating database: %s", err)
	}
	a := trader.NewOkxTrader(true)

	s := server.CreateServer(db, a, os.Getenv("PORT"))

	s.GET("/", routes.RenderHome)
	s.GET("/positions", routes.RenderPositionsList)
	s.GET("/trades", routes.RenderTrades)
	s.GET("/instruments", routes.RenderInstruments)
	s.POST("/market-order", routes.RenderPlaceMarketOrder)

	s.POST("/api/webhook", routes.HandleWebhook)
	s.POST("/api/create-trade", routes.CreateTrade)
	s.GET("/api/trades", routes.Trades)
	s.GET("/api/balance", routes.Balance)
	s.GET("/api/candles", routes.Candles)
	s.GET("/api/bidask/:ticker", routes.BidAsk)
	s.GET("/api/instruments/:instType", routes.Instruments)

	s.Listen()
}
