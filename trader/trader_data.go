package trader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TraderDataCandles struct {
	C []float64 `json:"c"`
	H []float64 `json:"h"`
	L []float64 `json:"l"`
	O []float64 `json:"o"`
	S string    `json:"s"`
	T []int     `json:"t"`
	V []float64 `json:"v"`
}

type TraderDataRequest struct {
	Method      string
	RequestPath string
	Body        string
}

// TODO(satvik): Take resolution as input
func GetCandles(symbol string, resolution string, from time.Time) ([]Candlestick, error) {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	fromTime := strconv.FormatInt(from.Unix(), 10)
	fmt.Println(fromTime)

	reqPath := "/crypto/candle?symbol=BINANCE:" + symbol + "&resolution=" + resolution + "&from=" + fromTime + "&to=" + currentTime

	fmt.Println(reqPath)

	p := TraderDataRequest{
		Method:      "GET",
		RequestPath: reqPath,
	}

	res, err := sendFinnhubRequest(p)
	if err != nil {
		return nil, err
	}

	candlesRes := &TraderDataCandles{}
	if err := json.Unmarshal([]byte(res), &candlesRes); err != nil {
		return nil, err
	}

	if candlesRes.S == "no_data" {
		return nil, fmt.Errorf("no data for ticker %s", symbol)
	}

	candles := []Candlestick{}
	for i := 0; i < len(candlesRes.C); i++ {
		time := time.Unix(int64(candlesRes.T[i]), 0)

		candles = append(candles, Candlestick{
			Open:      candlesRes.O[i],
			Close:     candlesRes.C[i],
			High:      candlesRes.H[i],
			Low:       candlesRes.L[i],
			Timestamp: time,
		})
	}

	return candles, nil
}

func sendFinnhubRequest(p TraderDataRequest) (string, error) {
	path := "https://finnhub.io/api/v1" + p.RequestPath

	client := &http.Client{}
	req, err := http.NewRequest(p.Method, path, bytes.NewBuffer([]byte(p.Body)))
	if err != nil {
		return "", err
	}

	req.Header.Add("X-Finnhub-Token", os.Getenv("FINNHUB_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
