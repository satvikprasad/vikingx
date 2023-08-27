package okx

import "time"

type OkBalanceResponse struct {
	Code string `json:"code"`
	Data []struct {
		AdjEq      string `json:"adjEq"`
		BorrowFroz string `json:"borrowFroz"`
		Details    []struct {
			AvailBal      string `json:"availBal"`
			AvailEq       string `json:"availEq"`
			BorrowFroz    string `json:"borrowFroz"`
			CashBal       string `json:"cashBal"`
			Ccy           string `json:"ccy"`
			CrossLiab     string `json:"crossLiab"`
			DisEq         string `json:"disEq"`
			Eq            string `json:"eq"`
			EqUsd         string `json:"eqUsd"`
			FixedBal      string `json:"fixedBal"`
			FrozenBal     string `json:"frozenBal"`
			Interest      string `json:"interest"`
			IsoEq         string `json:"isoEq"`
			IsoLiab       string `json:"isoLiab"`
			IsoUpl        string `json:"isoUpl"`
			Liab          string `json:"liab"`
			MaxLoan       string `json:"maxLoan"`
			MgnRatio      string `json:"mgnRatio"`
			NotionalLever string `json:"notionalLever"`

			OrdFrozen    string `json:"ordFrozen"`
			SpotInUseAmt string `json:"spotInUseAmt"`
			StgyEq       string `json:"stgyEq"`
			Twap         string `json:"twap"`
			UTime        string `json:"uTime"`

			Upl     string `json:"upl"`
			UplLiab string `json:"uplLiab"`
		} `json:"details"`
		Imr         string `json:"imr"`
		IsoEq       string `json:"isoEq"`
		MgnRatio    string `json:"mgnRatio"`
		Mmr         string `json:"mmr"`
		NotionalUsd string `json:"notionalUsd"`
		OrdFroz     string `json:"ordFroz"`
		TotalEq     string `json:"totalEq"`

		UTime string `json:"uTime"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type OkDefaultResponse struct {
	Code string `json:"code"`
	Data []any  `json:"data"`

	Msg string `json:"msg"`
}

type OkOrderResponse struct {
	Code string `json:"code"`
	Data []struct {
		ClOrdID string `json:"clOrdId"`
		OrdID   string `json:"ordId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
		Tag     string `json:"tag"`
	} `json:"data"`

	InTime string `json:"inTime"`

	Msg     string `json:"msg"`
	OutTime string `json:"outTime"`
}

type OkLimitPricesResponse struct {
	Code string `json:"code"`
	Data []struct {
		BuyLmt   string `json:"buyLmt"`
		InstID   string `json:"instId"`
		InstType string `json:"instType"`
		SellLmt  string `json:"sellLmt"`
		Ts       string `json:"ts"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type OkCandlestickResponse struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data [][]string `json:"data"`
}

type OkInstrumentsResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []OkInstruments `json:"data"`
}

type OkInstruments struct {
	InstType string `json:"instType"`

	InstID     string `json:"instId"`
	InstFamily string `json:"instFamily"`
	Uly        string `json:"uly"`
	Category   string `json:"category"`

	BaseCcy   string `json:"baseCcy"`
	QuoteCcy  string `json:"quoteCcy"`
	SettleCcy string `json:"settleCcy"`

	CtVal    string `json:"ctVal"`
	CtMult   string `json:"ctMult"`
	CtValCcy string `json:"ctValCcy"`
	OptType  string `json:"optType"`
	Stk      string `json:"stk"`
	ListTime string `json:"listTime"`
	ExpTime  string `json:"expTime"`
	Lever    string `json:"lever"`
	TickSz   string `json:"tickSz"`
	LotSz    string `json:"lotSz"`
	MinSz    string `json:"minSz"`
	CtType   string `json:"ctType"`
	Alias    string `json:"alias"`

	State        string `json:"state"`
	MaxLmtSz     string `json:"maxLmtSz"`
	MaxMktSz     string `json:"maxMktSz"`
	MaxTwapSz    string `json:"maxTwapSz"`
	MaxIcebergSz string `json:"maxIcebergSz"`

	MaxTriggerSz string `json:"maxTriggerSz"`
	MaxStopSz    string `json:"maxStopSz"`
}

type OkCandlestick struct {
	Timestamp              time.Time
	Open, High, Low, Close float64
}

type OkTicker struct {
	InstType  string `json:"instType"`
	InstID    string `json:"instId"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24H   string `json:"open24h"`
	High24H   string `json:"high24h"`
	Low24H    string `json:"low24h"`
	VolCcy24H string `json:"volCcy24h"`
	Vol24H    string `json:"vol24h"`
	SodUtc0   string `json:"sodUtc0"`
	SodUtc8   string `json:"sodUtc8"`
	Ts        string `json:"ts"`
}

type OkTickersResponse struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data []OkTicker `json:"data"`
}
