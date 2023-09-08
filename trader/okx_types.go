package trader

type OkxBalanceResponse struct {
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

type OkxDefaultResponse struct {
	Code string `json:"code"`
	Data []any  `json:"data"`

	Msg string `json:"msg"`
}

type OkxOrderResponse struct {
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

type OkxLimitPricesResponse struct {
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

type OkxCandlestickResponse struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data [][]string `json:"data"`
}

type OkxInstrumentsResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data []OkxInstrument `json:"data"`
}

type OkxInstrument struct {
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

type OkxTicker struct {
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

type OkxTickersResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data []OkxTicker `json:"data"`
}

type OkxMarkPriceResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType string `json:"instType"`
		InstID   string `json:"instId"`
		MarkPx   string `json:"markPx"`
		Ts       string `json:"ts"`
	} `json:"data"`
}

type OkxPosition struct {
	Adl            string `json:"adl"`
	AvailPos       string `json:"availPos"`
	AvgPx          string `json:"avgPx"`
	CTime          string `json:"cTime"`
	Ccy            string `json:"ccy"`
	DeltaBS        string `json:"deltaBS"`
	DeltaPA        string `json:"deltaPA"`
	GammaBS        string `json:"gammaBS"`
	GammaPA        string `json:"gammaPA"`
	Imr            string `json:"imr"`
	InstID         string `json:"instId"`
	InstType       string `json:"instType"`
	Interest       string `json:"interest"`
	IdxPx          string `json:"idxPx"`
	UsdPx          string `json:"usdPx"`
	Last           string `json:"last"`
	Lever          string `json:"lever"`
	Liab           string `json:"liab"`
	LiabCcy        string `json:"liabCcy"`
	LiqPx          string `json:"liqPx"`
	MarkPx         string `json:"markPx"`
	Margin         string `json:"margin"`
	MgnMode        string `json:"mgnMode"`
	MgnRatio       string `json:"mgnRatio"`
	Mmr            string `json:"mmr"`
	NotionalUsd    string `json:"notionalUsd"`
	OptVal         string `json:"optVal"`
	PTime          string `json:"pTime"`
	Pos            string `json:"pos"`
	BaseBorrowed   string `json:"baseBorrowed"`
	BaseInterest   string `json:"baseInterest"`
	QuoteBorrowed  string `json:"quoteBorrowed"`
	QuoteInterest  string `json:"quoteInterest"`
	PosCcy         string `json:"posCcy"`
	PosID          string `json:"posId"`
	PosSide        string `json:"posSide"`
	SpotInUseAmt   string `json:"spotInUseAmt"`
	SpotInUseCcy   string `json:"spotInUseCcy"`
	BizRefID       string `json:"bizRefId"`
	BizRefType     string `json:"bizRefType"`
	ThetaBS        string `json:"thetaBS"`
	ThetaPA        string `json:"thetaPA"`
	TradeID        string `json:"tradeId"`
	UTime          string `json:"uTime"`
	Upl            string `json:"upl"`
	UplLastPx      string `json:"uplLastPx"`
	UplRatio       string `json:"uplRatio"`
	UplRatioLastPx string `json:"uplRatioLastPx"`
	VegaBS         string `json:"vegaBS"`
	VegaPA         string `json:"vegaPA"`
	CloseOrderAlgo []struct {
		AlgoID          string `json:"algoId"`
		SlTriggerPx     string `json:"slTriggerPx"`
		SlTriggerPxType string `json:"slTriggerPxType"`
		TpTriggerPx     string `json:"tpTriggerPx"`
		TpTriggerPxType string `json:"tpTriggerPxType"`
		CloseFraction   string `json:"closeFraction"`
	} `json:"closeOrderAlgo"`
}

type OkxPositionsResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []OkxPosition `json:"data"`
}
