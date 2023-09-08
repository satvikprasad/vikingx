package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertperpTickerName(t *testing.T) {
	_, err := convertTickerName("BTCUX:SDT")
	assert.NotNil(t, err)

	perpTicker, err := convertTickerName("BTCUSDT.P")
	assert.Nil(t, err)
	assert.Equal(t, perpTicker, "BTC-USDT-SWAP")

	ticker, err := convertTickerName("ETHUSDT")
	assert.Nil(t, err)
	assert.Equal(t, ticker, "ETH-USDT")
}
