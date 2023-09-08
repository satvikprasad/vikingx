package trader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandles(t *testing.T) {
	d := TraderData{}

	candles, err := d.GetCandles()
	assert.Nil(t, err)
	fmt.Println(candles)
}
