package db

import "github.com/satvikprasad/vikingx/models"

type Database interface {
	CreateTrade(t *models.Trade)
	Trades() []*models.Trade
}
