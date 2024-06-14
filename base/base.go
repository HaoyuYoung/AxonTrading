package base

import (
	"errors"
)

// 交易所名称
var (
	KUCOIN    = "KUCOIN"
	BINANCE   = "BINANCE"
	MEXC      = "MEXC"
	MEXCV3    = "MEXCV3"
	GATEIO    = "GATEIO"
	BITMART   = "BITMART"
	COINSTORE = "CoinStore"
	PROBIT    = "Probit"
	BYBIT     = "Bybit"
	OKEX      = "OKEX"
	BITGET    = "Bitget"
	XT        = "XT"
	BITRUE    = "BITRUE"
	COINLIST  = "COINLIST"
	BTSE      = "BTSE"
	LBANK     = "LBANK"
)

// 枚举类型
var (
	//统一枚举类型

	BID = "bid" //买
	ASK = "ask" //卖

	LONG  = "long"  //多头
	SHORT = "short" //空头

	ISOLATED = "isolated" //逐仓
	CROSSED  = "crossed"  //全仓

	ADDMARGIN    = 1 //增加保证金
	REMOVEMARGIN = 2 //减少保证金

	OPEN      = "open"      //完全未成交
	FILLED    = "filled"    //完全成交
	CANCELED  = "canceled"  //取消
	PARTIALLY = "partially" //部分成交

	LIMIT       = "limit"
	LIMITHIDDEN = "limit_hidden"
	MAKER       = "maker"
	TAKER       = "taker"
	MARKET      = "market"

	STOP             = "stop"               //止损
	STOPMARKET       = "stop_market"        //市价止损
	TAKEPROFIT       = "take_profit"        //止盈
	TAKEPROFITMARKET = "take_profit_market" //市价止盈

	// UnifiedBuy  UnifiedSell  统一的 Side 方向
	UnifiedBuy  = "buy"
	UnifiedSell = "sell"

	ICEBERG = "iceBerg"
)

// 订单状态
var (
	Pending  = 0
	Filled   = 1
	Canceled = 2
)

var (
	ErrResponse = errors.New("response error , code is not 200")
)
