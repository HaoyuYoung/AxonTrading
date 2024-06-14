package models

import (
	"AxonTrading/base"

	"sync"
)

type WsData struct {
	Time int64        `json:"time"`
	Bids []PriceLevel `json:"bids"`
	Asks []PriceLevel `json:"asks"`
}

type PositionInfo struct {
	Symbol           string `json:"symbol"`
	PositionAmt      string `json:"positionAmt"`
	EntryPrice       string `json:"entryPrice"`
	MarkPrice        string `json:"markPrice"`
	UnRealizedProfit string `json:"unRealizedProfit"`
	LiquidationPrice string `json:"liquidationPrice"`
	Leverage         string `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"`
	MarginType       string `json:"marginType"`
	IsolatedMargin   string `json:"isolatedMargin"`
	IsAutoAddMargin  string `json:"isAutoAddMargin"`
	PositionSide     string `json:"positionSide"`
	Notional         string `json:"notional"`
	IsolatedWallet   string `json:"isolatedWallet"`
	UpdateTime       int64  `json:"updateTime"`
}

type FutureBalance struct {
	Asset            string `json:"asset"`
	TotalBalance     string `json:"totalBalance"`
	CrossBalance     string `json:"crossBalance"`
	AvailableBalance string `json:"availableBalance"`
}
type FutureOrderInfo struct {
	AvgPrice      string `json:"avgPrice"`
	CumQuote      string `json:"cumQuote"`
	ExecutedQty   string `json:"executedQty"`
	OrderId       int    `json:"orderId"`
	OrigQty       string `json:"origQty"`
	OrigType      string `json:"origType"`
	Price         string `json:"price"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	ClosePosition bool   `json:"closePosition"`
	Symbol        string `json:"symbol"`
	Time          int64  `json:"time"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	UpdateTime    int64  `json:"updateTime"`
	PriceProtect  bool   `json:"priceProtect"`
}
type FundingRate struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`  // 标记价格
	IndexPrice           string `json:"indexPrice"` // 指数价格
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"` // 当前 Funding rate
	NextFundingTime      int64  `json:"nextFundingTime"` // next 更新时间 时间cuo
	InterestRate         string `json:"interestRate"`    // 0 0
	Time                 int64  `json:"time"`
}
type Candle struct {
	Symbol  string   `json:"symbol"`
	Candles []string `json:"candles"`
	Time    int64    `json:"time"`
}
type CandleInfo struct {
	Opentime string `json:"opentime"`
	Open     string `json:"open"`
	Close    string `json:"close"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Volume   string `json:"volume"`
	Turnover string `json:"turnover"`
}

type PrvWsTradeOrderInfo struct {
	Symbol     string `json:"symbol"`
	OrderType  string `json:"orderType"`
	Type       string `json:"type"`
	OrderId    string `json:"orderId"`
	OrderTime  int64  `json:"orderTime"`
	Size       string `json:"size"`
	FilledSize string `json:"filledSize"`
	Price      string `json:"price"`
	ClientOid  string `json:"clientOid"`
	RemainSize string `json:"remainSize"`
	Status     string `json:"status"`
	Ts         int64  `json:"ts"`
	Side       string `json:"side"`
}

type PriceLevel struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type OrderInfo struct {
	OrderID string `json:"order_id"`
	Symbol  string `json:"symbol"`
	Side    string `json:"side"`
	Price   string `json:"price"`
	// Size     string `json:"size"`
	Quantity string `json:"quantity"`
	Type     string `json:"type"`
	Filled   string `json:"filled"`
	USDT     string `json:"usdt"`
	Status   string `json:"status"`
	Time     int64  `json:"time"`
}

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type MexcOrderInfo struct {
	Symbol              string `json:"symbol"`
	OrderId             string `json:"orderId"`
	OrderListId         int    `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          string `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

type CSInfo struct {
	Price    string `json:"price"`
	Quantity string `json:"origQty"`
	Type     int    `json:"type"`
}

type OrderList struct {
	Side  string `json:"side"`
	Price string `json:"price"`
	Size  string `json:"size"`
	Type  string `json:"type"`
}

type TradingFee struct {
	Symbol                string `json:"symbol"`
	TakerFeeFromApi       string `json:"taker_fee_from_api"`
	MakerFeeFromApi       string `json:"maker_fee_from_api"`
	TakerFeeFromRealOrder string `json:"taker_fee_from_real_order"`
	MakerFeeFromRealOrder string `json:"maker_fee_from_real_order"`
	IfDiscount            bool   `json:"if_discount"`
}
type PairInfo struct {
	MinBaseAmount   string `json:"min_base_amount"`
	MinQuoteAmount  string `json:"min_quote_amount"`
	AmountPrecision int    `json:"amount_precision"`
	Precision       int    `json:"precision"`
}

// SideAdaptor Uniformity Side 统一 side 适配器
func SideAdaptor(name, side string) string {
	switch name {
	case base.BINANCE:
		if side == base.UnifiedBuy {
			return "bid"
		} else {
			return "ask"
		}
	case base.BITMART:
		if side == base.UnifiedBuy {
			return "buy"
		} else {
			return "sell"
		}
	case base.COINSTORE:
		if side == base.UnifiedBuy {
			return "bid"
		} else {
			return "ask"
		}
	case base.GATEIO:
		if side == base.UnifiedBuy {
			return "bid"
		} else {
			return "ask"
		}
	case base.KUCOIN:
		if side == base.UnifiedBuy {
			return "bid"
		} else {
			return "ask"
		}
	case base.MEXC:
		if side == base.UnifiedBuy {
			return "bid"
		} else {
			return "ask"
		}
	default:
		return ""
	}
}

type SymbolTicker struct {
	Symbol     string `json:"symbol"`
	ChangeRate string `json:"change_rate"`
	Volume     string `json:"volume"`
	LastPrice  string `json:"last_price"`
}

type PriceRecord struct {
	Price      []float64 `json:"price"`
	Volome     string    `json:"volome"`
	ChangeRate string    `json:"change_rate"`
	Prec       int       `json:"prec"`
}

type PR struct {
	sync.RWMutex
	M map[string]PriceRecord `json:"m"`
}
