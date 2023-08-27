package models

import (
	"gorm.io/gorm"
)

type Trade struct {
	gorm.Model
	Ticker                 string  `json:"ticker"`
	Side                   string  `json:"side"`
	Size                   float64 `json:"size"`
	Price                  float64 `json:"price"`
	MarketPositionSize     float64 `json:"marketPositionSize"`
	PrevMarketPositionSize float64 `json:"prevMarketPositionSize"`
}
