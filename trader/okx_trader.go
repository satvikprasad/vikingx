package trader

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

type OkxCredentials struct {
	Passphrase string
	AccessKey  string
	SecretKey  string
}

type OkxRequestParams struct {
	Method      string
	RequestPath string
	Body        string
}

type OkxTrader struct {
	OkxCredentials

	Demo bool
}

func NewOkxTrader(demo bool) Trader {
	godotenv.Load()

	return &OkxTrader{
		OkxCredentials: OkxCredentials{
			Passphrase: os.Getenv("PASSPHRASE"),
			AccessKey:  os.Getenv("ACCESS_KEY"),
			SecretKey:  os.Getenv("SECRET_KEY"),
		},
		Demo: demo,
	}
}

// TODO(satvik): Make this return ct size factored in
func (a *OkxTrader) Positions() ([]Position, error) {
	p := OkxRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/account/positions",
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return nil, err
	}

	positionRes := &OkxPositionsResponse{}
	if err := json.Unmarshal([]byte(res), &positionRes); err != nil {
		return nil, err
	}

	if positionRes.Code != "0" {
		return nil, fmt.Errorf("Error getting positions: %s",
			positionRes.Msg)
	}

	positions := []Position{}
	for _, o := range positionRes.Data {
		cx := 1.0
		if o.InstType == "SWAP" {
			ctSize, err := a.TickerCtSize(o.InstID)
			if err != nil {
				return nil, err
			}
			cx = ctSize
		}

		positionSize, err := strconv.ParseFloat(o.Pos, 64)
		if err != nil {
			return nil, err
		}

		positions = append(positions, Position{
			Size:   positionSize * cx,
			Symbol: o.InstID,
			Type:   o.InstType,
		})
	}

	return positions, nil
}

func (a *OkxTrader) Tickers(instType string) ([]Ticker, error) {
	p := OkxRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/market/tickers?instType=" + instType,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return nil, err
	}

	tickerRes := &OkxTickersResponse{}
	if err := json.Unmarshal([]byte(res), &tickerRes); err != nil {
		return nil, err
	}

	if tickerRes.Code != "0" {
		return nil, fmt.Errorf("Error getting tickers: %s",
			tickerRes.Msg)
	}

	tickers := []Ticker{}
	for _, t := range tickerRes.Data {
		bidPx, err := strconv.ParseFloat(t.BidPx, 64)
		if err != nil {
			continue
		}

		askPx, err := strconv.ParseFloat(t.AskPx, 64)
		if err != nil {
			continue
		}

		vol24H, err := strconv.ParseFloat(t.Vol24H, 64)
		if err != nil {
			continue
		}

		lastTradedPx, err := strconv.ParseFloat(t.Last, 64)
		if err != nil {
			continue
		}

		tickers = append(tickers, Ticker{
			Symbol:          t.InstID,
			BidPrice:        bidPx,
			AskPrice:        askPx,
			Vol24H:          vol24H,
			LastTradedPrice: lastTradedPx,
		})
	}

	return tickers, nil
}

func (a *OkxTrader) Instruments(instType string) ([]Instrument, error) {
	p := OkxRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/public/instruments?instType=" + instType,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return nil, err
	}

	instrumentsRes := &OkxInstrumentsResponse{}
	if err := json.Unmarshal([]byte(res), &instrumentsRes); err != nil {
		return nil, err
	}

	if instrumentsRes.Code != "0" {
		return nil, fmt.Errorf("Error getting instruments: %s",
			instrumentsRes.Msg)
	}

	instruments := []Instrument{}
	for _, i := range instrumentsRes.Data {
		instruments = append(instruments, Instrument{
			Symbol:    i.InstID,
			BaseCcy:   i.BaseCcy,
			QuoteCcy:  i.QuoteCcy,
			CtValCcy:  i.CtValCcy,
			CtVal:     i.CtVal,
			SettleCcy: i.SettleCcy,
		})
	}

	return instruments, nil
}

func (a *OkxTrader) LimitSwapPrice(symbol string) (float64, float64, error) {
	p := OkxRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/public/price-limit?instId=" + symbol,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return -1, -1, err
	}

	limitPricesRes := &OkxLimitPricesResponse{}
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

func (a *OkxTrader) Candlesticks(symbol string,
	time string) ([]Candlestick, error) {
	p := OkxRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/market/candles?instId=" + symbol + "&bar=" + time + "&limit=300",
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return nil, err
	}

	candles := &OkxCandlestickResponse{}
	if err := json.Unmarshal([]byte(res), &candles); err != nil {
		return nil, err
	}

	if candles.Code != "0" {
		return nil, fmt.Errorf("Error getting candles: %s", candles.Msg)
	}

	c := []Candlestick{}

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

		c = append(c, Candlestick{
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
func (a *OkxTrader) LimitOrder(symbol string, tradeMode string,
	side string, price float64, size float64) error {
	body := fmt.Sprintf(`{
        "instId": "%s",
        "tdMode": "%s",
        "side": "%s",
        "ordType": "limit",
        "px": "%.3f",
        "sz": "%.3f"
    }`, symbol, tradeMode, side, price, size)

	p := OkxRequestParams{
		Method:      "POST",
		RequestPath: "/api/v5/trade/order",
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return err
	}

	orderRes := &OkxOrderResponse{}
	if err := json.Unmarshal([]byte(res), &orderRes); err != nil {
		return err
	}

	if orderRes.Code != "0" {
		return fmt.Errorf("Error placing limit order: %s: %s",
			orderRes.Data[0].SMsg, orderRes.Msg)
	}

	return nil
}

func (a *OkxTrader) MarketOrder(symbol string, side string, size float64) error {
	tdMode := "cash"
	if strings.Contains(symbol, "SWAP") {
		tdMode = "isolated"
	}

	body := fmt.Sprintf(`{
        "instId": "%s",
        "tdMode": "%s",
        "side": "%s",
        "ordType": "market",
        "sz": "%.3f"
    }`, symbol, tdMode, side, size)

	fmt.Println(body)

	p := OkxRequestParams{
		Method:      "POST",
		RequestPath: "/api/v5/trade/order",
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return err
	}

	orderRes := &OkxOrderResponse{}
	if err := json.Unmarshal([]byte(res), &orderRes); err != nil {
		return err
	}

	fmt.Println(orderRes)

	if orderRes.Code != "0" {
		return fmt.Errorf("Error placing market order: %s: %s",
			orderRes.Data[0].SMsg, orderRes.Msg)
	}

	fmt.Println(orderRes.Data[0].SMsg)

	return nil
}

func (a *OkxTrader) SetLeverage(symbol string, leverage int) error {
	body := fmt.Sprintf(`{
        "instId":"%s",
        "lever":"%d",
        "mgnMode":"isolated"
    }`, symbol, leverage)

	p := OkxRequestParams{
		Method:      "POST",
		RequestPath: "/api/v5/account/set-leverage",
		Body:        body,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return err
	}

	leverageRes := &OkxDefaultResponse{}
	if err := json.Unmarshal([]byte(res), &leverageRes); err != nil {
		return err
	}

	if leverageRes.Code != "0" {
		return fmt.Errorf("Error setting leverage: %s", leverageRes.Msg)
	}

	return nil
}

func (a *OkxTrader) Balance(symbol string) (float64, error) {
	p := OkxRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/account/balance?ccy=" + symbol,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return -1, err
	}

	balanceRes := &OkxBalanceResponse{}
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

func (a *OkxTrader) MarkPrice(instId string) (float64, error) {
	p := OkxRequestParams{
		Method:      "GET",
		RequestPath: "/api/v5/public/mark-price?instId=" + instId,
	}

	res, err := a.SendRequest(p)
	if err != nil {
		return -1, err
	}

	markPriceRes := &OkxMarkPriceResponse{}
	if err := json.Unmarshal([]byte(res), &markPriceRes); err != nil {
		return -1, err
	}
	fmt.Printf("%+v\n", markPriceRes)

	mark, err := strconv.ParseFloat(markPriceRes.Data[0].MarkPx, 64)
	if err != nil {
		return -1, err
	}

	return mark, nil
}

func (a *OkxTrader) SendRequest(p OkxRequestParams) (string, error) {
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

	if a.Demo == true {
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

func (a *OkxTrader) TickerCtSize(ticker string) (float64, error) {
	if !strings.Contains(ticker, "SWAP") {
		return -1, fmt.Errorf("ct size requested for ticker %s is not perpetual", ticker)
	}

	inst, err := a.Instruments("SWAP")
	if err != nil {
		return 0, err
	}

	for _, instrument := range inst {
		if strings.Contains(ticker, instrument.CtValCcy) &&
			strings.Contains(ticker, instrument.SettleCcy) &&
			strings.Contains(instrument.Symbol, "USDT") {
			ctVal, err := strconv.ParseFloat(instrument.CtVal, 64)
			if err != nil {
				return 0, err
			}

			return ctVal, nil
		}
	}

	return 0, fmt.Errorf("Could not get ticker %s", ticker)
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
