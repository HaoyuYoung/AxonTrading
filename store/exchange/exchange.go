package exchange

import (
	"AxonTrading/base"
	"AxonTrading/exchanges/binance"
	"AxonTrading/exchanges/okx"
	"AxonTrading/models"
)

type Exchange interface {
	New(params []byte) error
	GetAccountBalance(currency string) ([]string, error)
	MarketOrder(symbol, side, size string) (string, error)
	LimitOrder(symbol, side, price, size string) (string, error)
	LimitHiddenOrder(symbol, side, price, size string) (string, error)
	LimitOrders(symbol string, ol []models.OrderList) ([]string, error)
	MakerOrder(symbol, side, price, size string) (string, error)
	MakerOrders(symbol string, ol []models.OrderList) ([]string, error)
	TakerOrder(symbol, side, price, size string) (string, error)
	TakerOrders(symbol string, ol []models.OrderList) ([]string, error)
	CancelOrder(symbol, id string) (bool, error)
	CancelOrders(symbol string) error
	GetOrder(symbol, id string) (models.OrderInfo, error)
	GetOpenOrders(symbol string) ([]models.OrderInfo, error)
	GetOpenOrdersWithSide(symbol, side string) ([]models.OrderInfo, error)
	GetOpenSplitOrders(symbol string) ([]models.OrderInfo, []models.OrderInfo, error)
	GetMarketPrice(symbol string) (string, error)
	Depth(symbol, limit string) (models.WsData, error)
	GetTradingFee(symbol string) (models.TradingFee, error)
	GetPairInfo(symbol string) (models.PairInfo, error)
	GetFeeFromFilled(symbol, id string) (string, string, error)
	IceBergOrder(symbol, side, typ, price, size, ice string) (string, error)
	LimitHiddenOrders(symbol string, ol []models.OrderList) ([]string, error)
	GetDepositAddress(token, chain string) (string, error)
	Withdraw(token, chain, to, amount string) (string, error)

	// TODO 期货

	// NewFuture 新建 client
	NewFuture(params []byte) error
	// GetFutureBalance 获取期货账户底仓（U本位）
	GetFutureBalance() (models.FutureBalance, error)
	// FutureDepth 期货深度，symbol 币对 limit 深度档位
	FutureDepth(symbol, limit string) (models.WsData, error)
	// GetFutureMarketPrice 获取期货市场价格 symbol 币对
	GetFutureMarketPrice(symbol string) (string, error)
	// GetMarkPriceAndFundingRate 获取标记价格 Fund rate
	GetMarkPriceAndFundingRate(symbol string) (models.FundingRate, error)
	// Dual 改变持仓方向 true 双向 false 单向
	Dual(dualSize bool) (bool, error)
	// CheckDual 检查当前是否为双向持仓（true）
	CheckDual() (bool, error)
	NewFutureOrder(symbol, side, positionSide, typ, size, price, stopPrice, positionType string, closePosition, priceProtect bool) (string, error)
	GetFutureOrder(symbol, orderID string) (models.FutureOrderInfo, error)
	// CancelFutureOrder 取消挂单
	CancelFutureOrder(symbol, orderID string) (bool, error)
	// CancelFutureOrders 取消全部挂单
	CancelFutureOrders(symbol string) error
	// GetFutureOpenOrders 获取 open 状态的挂单
	GetFutureOpenOrders(symbol string) ([]models.FutureOrderInfo, error)
	// ChangeLeverage 改变杠杆倍数
	ChangeLeverage(symbol string, leverage int) (string, error)
	// ChangeMarginType 改变仓位
	ChangeMarginType(symbol, typ string) error
	// ChangePositionMargin 改变逐仓保证金
	ChangePositionMargin(symbol, positionSide, amount string, typ int) (bool, error)
	// GetPositionRisk 获取当前仓位
	GetPositionRisk(symbol string) ([]models.PositionInfo, error)
	// GetFutureTradingFee 获取手续费
	GetFutureTradingFee(symbol string) (models.TradingFee, error)

	// TODO
}

type ExchangeFactory struct {
}

func (e ExchangeFactory) CreateClient(exchange string) Exchange {

	switch exchange {
	case base.BINANCE:
		return &binance.Client{}
	case base.OKEX:
		return &okx.Client{}

	default:
		return nil
	}

}
