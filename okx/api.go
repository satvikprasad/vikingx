package okx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type OkCredentials struct {
	Passphrase string
	AccessKey  string
	SecretKey  string
}

type OkRequestParams struct {
	Method      string
	RequestPath string
	Body        string
}

type OkApi struct {
	OkCredentials

	Demo bool
}

func NewOkApi(demo bool, envPath string) *OkApi {
	godotenv.Load(envPath)

	return &OkApi{
		OkCredentials: OkCredentials{
			Passphrase: os.Getenv("PASSPHRASE"),
			AccessKey:  os.Getenv("ACCESS_KEY"),
			SecretKey:  os.Getenv("SECRET_KEY"),
		},
		Demo: demo,
	}
}

func (a *OkApi) GetTickers(instType string) ([]OkTicker, error) {
	p := OkRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/market/tickers?instType=" + instType,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return nil, err
	}

	tickerRes := &OkTickersResponse{}
	if err := json.Unmarshal([]byte(res), &tickerRes); err != nil {
		return nil, err
	}

	if tickerRes.Code != "0" {
		return nil, fmt.Errorf("Error getting tickers: %s",
			tickerRes.Msg)
	}

	return tickerRes.Data, nil
}

func (a *OkApi) GetInstruments(instType string) ([]OkInstruments, error) {
	p := OkRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/public/instruments?instType=" + instType,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return nil, err
	}

	instrumentsRes := &OkInstrumentsResponse{}
	if err := json.Unmarshal([]byte(res), &instrumentsRes); err != nil {
		return nil, err
	}

	if instrumentsRes.Code != "0" {
		return nil, fmt.Errorf("Error getting instruments: %s",
			instrumentsRes.Msg)
	}

	return instrumentsRes.Data, nil
}

func (a *OkApi) GetLimitSwapPrice(symbol string) (buy float64,
	sell float64, err error) {
	p := OkRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/public/price-limit?instId=" + symbol,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return -1, -1, err
	}

	limitPricesRes := &OkLimitPricesResponse{}
	if err := json.Unmarshal([]byte(res), &limitPricesRes); err != nil {
		return -1, -1, err
	}

	if limitPricesRes.Code != "0" {
		return -1, -1, fmt.Errorf("Error getting limit prices: %s",
			limitPricesRes.Msg)
	}

	buyLmt, err := strconv.ParseFloat(limitPricesRes.Data[0].BuyLmt, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("Could not parse limit price: %s",
			limitPricesRes.Data[0].BuyLmt)
	}
	sellLmt, err := strconv.ParseFloat(limitPricesRes.Data[0].SellLmt, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("Could not parse limit price: %s",
			limitPricesRes.Data[0].SellLmt)
	}

	return buyLmt, sellLmt, nil
}

func (a *OkApi) GetCandleSticks(symbol string,
	time string) ([]OkCandlestick, error) {
	body := fmt.Sprintf(`{
        bar: %s
    }`, time)

	p := OkRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/market/index-candles?instId=" + symbol,
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return nil, err
	}

	candles := &OkCandlestickResponse{}
	if err := json.Unmarshal([]byte(res), &candles); err != nil {
		return nil, err
	}

	if candles.Code != "0" {
		return nil, fmt.Errorf("Error getting candles: %s", candles.Msg)
	}

	c := []OkCandlestick{}

	for i, v := range candles.Data {
		timestamp, err := unixMsToTime(v[0])
		if err != nil {
			return nil, fmt.Errorf("Error converting timestamp %s on bar %d",
				v[0], i)
		}

		open, err := strconv.ParseFloat(v[1], 64)
		high, err := strconv.ParseFloat(v[2], 64)
		low, err := strconv.ParseFloat(v[3], 64)
		close, err := strconv.ParseFloat(v[4], 64)

		c = append(c, OkCandlestick{
			Timestamp: timestamp,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
		})
	}

	return c, nil
}

// @TODO(satvikprasad): extract tradeMode into an enum
func (a *OkApi) LimitOrder(symbol string, tradeMode string,
	side string, price float64, size float64) error {
	body := fmt.Sprintf(`{
        "instId": "%s",
        "tdMode": "%s",
        "side": "%s",
        "ordType": "limit",
        "px": "%.3f",
        "sz": "%.3f"
    }`, symbol, tradeMode, side, price, size)

	p := OkRequestParams{
		Method:      "POST",
		RequestPath: "/api/v5/trade/order",
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return err
	}

	orderRes := &OkOrderResponse{}
	if err := json.Unmarshal([]byte(res), &orderRes); err != nil {
		return err
	}

	if orderRes.Code != "0" {
		return fmt.Errorf("Error placing limit order: %s: %s",
			orderRes.Data[0].SMsg, orderRes.Msg)
	}

	return nil
}

func (a *OkApi) MarketOrderSwap(symbol string,
	side string, size float64) error {
	body := fmt.Sprintf(`{
        "instId": "%s",
        "tdMode": "isolated",
        "side": "%s",
        "ordType": "market",
        "sz": "%.3f"
    }`, symbol, side, size)

	p := OkRequestParams{
		Method:      "POST",
		RequestPath: "/api/v5/trade/order",
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return err
	}

	orderRes := &OkOrderResponse{}
	if err := json.Unmarshal([]byte(res), &orderRes); err != nil {
		return err
	}

	if orderRes.Code != "0" {
		return fmt.Errorf("Error placing market order: %s: %s",
			orderRes.Data[0].SMsg, orderRes.Msg)
	}

	fmt.Println(orderRes.Data[0].SMsg)

	return nil
}

func (a *OkApi) MarketOrder(symbol string, side string, size float64) error {
	body := fmt.Sprintf(`{
        "instId": "%s",
        "tdMode": "cash",
        "side": "%s",
        "ordType": "market",
        "sz": "%.3f"
    }`, symbol, side, size)

	p := OkRequestParams{
		Method:      "POST",
		RequestPath: "/api/v5/trade/order",
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return err
	}

	orderRes := &OkOrderResponse{}
	if err := json.Unmarshal([]byte(res), &orderRes); err != nil {
		return err
	}

	if orderRes.Code != "0" {
		return fmt.Errorf("Error placing market order: %s: %s",
			orderRes.Data[0].SMsg, orderRes.Msg)
	}

	fmt.Println(orderRes.Data[0].SMsg)

	return nil
}

func (a *OkApi) SetLeverage(symbol string, leverage int) error {
	body := fmt.Sprintf(`{
        "instId":"%s",
        "lever":"%d",
        "mgnMode":"isolated"
    }`, symbol, leverage)

	p := OkRequestParams{
		Method:      "POST",
		RequestPath: "/api/v5/account/set-leverage",
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return err
	}

	leverageRes := &OkDefaultResponse{}
	if err := json.Unmarshal([]byte(res), &leverageRes); err != nil {
		return err
	}

	if leverageRes.Code != "0" {
		return fmt.Errorf("Error setting leverage: %s", leverageRes.Msg)
	}

	return nil
}

func (a *OkApi) GetBalance(symbol string) (float64, error) {
	p := OkRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/account/balance?ccy=" + symbol,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return -1, err
	}

	balanceRes := &OkBalanceResponse{}
	if err := json.Unmarshal([]byte(res), &balanceRes); err != nil {
		return -1, err
	}

	balance, err := strconv.ParseFloat(
		balanceRes.Data[0].Details[0].AvailBal, 32)
	if err != nil {
		return -1, err
	}

	return balance, nil
}

func (a *OkApi) SendRequest(p OkRequestParams) (string, error) {
	timestamp := formatUTCTimestamp(time.Now().UTC())
	keyHash := calculateHash(timestamp, p.Method,
		p.RequestPath, p.Body, a.SecretKey)

	client := &http.Client{}
	req, err := http.NewRequest(p.Method, "https://www.okx.com"+p.RequestPath,
		bytes.NewBuffer([]byte(p.Body)))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", a.AccessKey)
	req.Header.Add("OK-ACCESS-SIGN", keyHash)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", a.Passphrase)

	if a.Demo {
		req.Header.Add("x-simulated-trading", "1")
	}

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

func (a *OkApi) GetTickerCtSize(ticker string) (float64, error) {
	inst, err := a.GetInstruments("SWAP")
	if err != nil {
		return 0, err
	}

	for _, instrument := range inst {
		if strings.Contains(ticker, instrument.CtValCcy) &&
			strings.Contains(ticker, instrument.SettleCcy) &&
			strings.Contains(instrument.InstID, "USDT") {
			ctVal, err := strconv.ParseFloat(instrument.CtVal, 64)
			if err != nil {
				return 0, err
			}

			return ctVal, nil
		}
	}

	return 0, fmt.Errorf("Could not get ticker %s", ticker)
}

func (a *OkApi) ConvertTickerName(instType string,
	ticker string) (string, error) {
	inst, err := a.GetInstruments(instType)
	if err != nil {
		return "", err
	}

	switch instType {
	case "SWAP":
		for _, instrument := range inst {
			if strings.Contains(ticker, instrument.CtValCcy) &&
				strings.Contains(ticker, instrument.SettleCcy) &&
				strings.Contains(instrument.InstID, "USDT") {
				return instrument.InstID, nil
			}
		}
	case "SPOT":
		for _, instrument := range inst {
			if strings.Contains(ticker, instrument.BaseCcy) &&
				strings.Contains(ticker, instrument.QuoteCcy) &&
				strings.Contains(instrument.InstID, "USDT") {
				return instrument.InstID, nil
			}
		}
	}

	return "", fmt.Errorf("Could not get ticker %s", ticker)
}

func calculateHash(timestamp string, method string, requestPath string,
	body string, secretKey string) string {
	key := timestamp + method + requestPath + body

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(key))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func formatUTCTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}

func unixMsToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}
