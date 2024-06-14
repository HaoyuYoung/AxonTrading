package binance

import (
	"AxonTrading/base"
	"AxonTrading/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Client       *binance.Client
	FutureClient *futures.Client
}

func (c *Client) GetFutureTradingFee(symbol string) (models.TradingFee, error) {
	result, err := c.FutureClient.NewCommissionRateService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return models.TradingFee{}, err
	}
	tradingFee := models.TradingFee{
		Symbol:                result.Symbol,
		TakerFeeFromApi:       result.TakerCommissionRate,
		MakerFeeFromApi:       result.MakerCommissionRate,
		TakerFeeFromRealOrder: "",
		MakerFeeFromRealOrder: "",
		IfDiscount:            false,
	}
	return tradingFee, err
}

func (c *Client) GetPositionRisk(symbol string) ([]models.PositionInfo, error) {
	result, err := c.FutureClient.NewGetPositionRiskService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}
	var positionInfo []models.PositionInfo
	var marginType, positionSide string
	for _, r := range result {
		if r.MarginType == "ISOLATED" {
			marginType = base.ISOLATED
		} else {
			marginType = base.CROSSED
		}
		if r.PositionSide == "LONG" {
			positionSide = base.LONG
		} else if r.PositionSide == "SHORT" {
			positionSide = base.SHORT
		} else {
			positionSide = "BOTH"
		}

		p := models.PositionInfo{
			Symbol:           r.Symbol,
			PositionAmt:      r.PositionAmt,
			EntryPrice:       r.EntryPrice,
			MarkPrice:        r.MarkPrice,
			UnRealizedProfit: r.UnRealizedProfit,
			LiquidationPrice: r.LiquidationPrice,
			Leverage:         r.Leverage,
			MaxNotionalValue: r.MaxNotionalValue,
			MarginType:       marginType,
			IsolatedMargin:   r.IsolatedMargin,
			IsAutoAddMargin:  r.IsAutoAddMargin,
			PositionSide:     positionSide,
			Notional:         r.Notional,
			IsolatedWallet:   r.IsolatedWallet,
			UpdateTime:       0,
		}
		positionInfo = append(positionInfo, p)

	}
	return positionInfo, err
}

func (c *Client) ChangePositionMargin(symbol, positionSide, amount string, typ int) (bool, error) {
	var pSide string
	dual, err := c.CheckDual()
	if err != nil {
		return false, err
	}

	if dual == true && positionSide == base.LONG {
		pSide = "LONG"
	} else if dual == true && positionSide == base.SHORT {
		pSide = "SHORT"
	} else if dual == false {
		pSide = "BOTH"
	}

	err = c.FutureClient.NewUpdatePositionMarginService().Symbol(symbol).PositionSide(futures.PositionSideType(pSide)).Amount(amount).Type(typ).Do(context.Background())
	if err != nil {
		return false, err
	}
	return true, err
}

func (c *Client) ChangeMarginType(symbol, typ string) error {
	var marginType string
	if typ == base.ISOLATED {
		marginType = "ISOLATED"
	} else if typ == base.CROSSED {
		marginType = "CROSSED"
	}
	err := c.FutureClient.NewChangeMarginTypeService().Symbol(symbol).MarginType(futures.MarginType(marginType)).Do(context.Background())
	if err != nil {
		return err
	}
	return err
}

func (c *Client) ChangeLeverage(symbol string, leverage int) (string, error) {
	reslut, err := c.FutureClient.NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(context.Background())
	if err != nil {
		return strconv.Itoa(reslut.Leverage) + " " + reslut.Symbol, err
	}
	return strconv.Itoa(reslut.Leverage) + " " + reslut.Symbol, err
}

func (c *Client) GetFutureOpenOrders(symbol string) ([]models.FutureOrderInfo, error) {
	reslut, err := c.FutureClient.NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}
	var opens []models.FutureOrderInfo
	var orderSide, pSide, orderState, orderType string
	for _, r := range reslut {
		if r.Side == "BUY" {
			orderSide = base.BID
		} else {
			orderSide = base.ASK
		}
		if r.PositionSide == "LONG" {
			pSide = base.LONG
		} else {
			pSide = base.SHORT
		}

		if r.Status == "NEW" {
			orderState = base.OPEN
		} else if r.Status == "CANCELED" || r.Status == "EXPIRED" {
			orderState = base.CANCELED
		} else if r.Status == "FILLED" {
			orderState = base.FILLED
		} else if r.Status == "PARTIALLY_FILLED" {
			orderState = base.PARTIALLY
		}

		if r.Type == "LIMIT" {
			orderType = base.LIMIT
		} else if r.Type == "MARKET" {
			orderType = base.MARKET
		} else if r.Type == "STOP" {
			orderType = base.STOP
		} else if r.Type == "STOP_MARKET" {
			orderType = base.STOPMARKET
		} else if r.Type == "TAKE_PROFIT" {
			orderType = base.TAKEPROFIT
		} else if r.Type == "TAKE_PROFIT_MARKET" {
			orderType = base.TAKEPROFITMARKET
		}

		orderInfo := models.FutureOrderInfo{
			AvgPrice:      r.AvgPrice,
			CumQuote:      r.CumQuote,
			ExecutedQty:   r.ExecutedQuantity,
			OrderId:       int(r.OrderID),
			OrigQty:       r.OrigQuantity,
			OrigType:      string(r.OrigType),
			Price:         r.Price,
			ReduceOnly:    r.ReduceOnly,
			Side:          orderSide,
			PositionSide:  pSide,
			Status:        orderState,
			StopPrice:     r.StopPrice,
			ClosePosition: r.ClosePosition,
			Symbol:        r.Symbol,
			Time:          r.Time,
			TimeInForce:   string(r.TimeInForce),
			Type:          orderType,
			UpdateTime:    r.UpdateTime,
			PriceProtect:  r.PriceProtect,
		}
		opens = append(opens, orderInfo)
	}
	return opens, err
}

func (c *Client) CancelFutureOrder(symbol, orderID string) (bool, error) {
	id, _ := strconv.ParseInt(orderID, 10, 64)

	_, err := c.FutureClient.NewCancelOrderService().Symbol(symbol).OrderID(id).Do(context.Background())
	if err != nil {
		return false, err
	}
	return true, err
}

func (c *Client) CancelFutureOrders(symbol string) error {
	err := c.FutureClient.NewCancelAllOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return err
	}
	return err
}

func (c *Client) GetFutureOrder(symbol, orderID string) (models.FutureOrderInfo, error) {
	id, _ := strconv.ParseInt(orderID, 10, 64)
	result, err := c.FutureClient.NewGetOrderService().Symbol(symbol).OrderID(id).Do(context.Background())
	if err != nil {
		return models.FutureOrderInfo{}, err
	}
	var orderSide, pSide, orderState, orderType string
	if result.Side == "BUY" {
		orderSide = base.BID
	} else {
		orderSide = base.ASK
	}
	if result.PositionSide == "LONG" {
		pSide = base.LONG
	} else {
		pSide = base.SHORT
	}

	if result.Status == "NEW" {
		orderState = base.OPEN
	} else if result.Status == "CANCELED" || result.Status == "EXPIRED" {
		orderState = base.CANCELED
	} else if result.Status == "FILLED" {
		orderState = base.FILLED
	} else if result.Status == "PARTIALLY_FILLED" {
		orderState = base.PARTIALLY
	}

	if result.Type == "LIMIT" {
		orderType = base.LIMIT
	} else if result.Type == "MARKET" {
		orderType = base.MARKET
	} else if result.Type == "STOP" {
		orderType = base.STOP
	} else if result.Type == "STOP_MARKET" {
		orderType = base.STOPMARKET
	} else if result.Type == "TAKE_PROFIT" {
		orderType = base.TAKEPROFIT
	} else if result.Type == "TAKE_PROFIT_MARKET" {
		orderType = base.TAKEPROFITMARKET
	}
	orderInfo := models.FutureOrderInfo{
		AvgPrice:      result.AvgPrice,
		CumQuote:      result.CumQuote,
		ExecutedQty:   result.ExecutedQuantity,
		OrderId:       int(id),
		OrigQty:       result.OrigQuantity,
		OrigType:      string(result.OrigType),
		Price:         result.Price,
		ReduceOnly:    result.ReduceOnly,
		Side:          orderSide,
		PositionSide:  pSide,
		Status:        orderState,
		StopPrice:     result.StopPrice,
		ClosePosition: result.ClosePosition,
		Symbol:        result.Symbol,
		Time:          result.Time,
		TimeInForce:   string(result.TimeInForce),
		Type:          orderType,
		UpdateTime:    result.UpdateTime,
		PriceProtect:  result.PriceProtect,
	}
	return orderInfo, err
}

func (c *Client) NewFutureOrder(symbol, side, positionSide, typ, size, price, stopPrice, positionType string, closePosition, priceProtect bool) (string, error) {
	err := c.ChangeMarginType(symbol, positionType)
	if err != nil {
		return "", err
	}
	var orderSide, pSide, orderType string
	if side == base.BID {
		orderSide = "BUY"
	} else if side == base.ASK {
		orderSide = "SELL"
	}
	dual, err := c.CheckDual()
	if err != nil {
		return "", err
	}

	if dual == true && positionSide == base.LONG {
		pSide = "LONG"
	} else if dual == true && positionSide == base.SHORT {
		pSide = "SHORT"
	} else if dual == false {
		pSide = "BOTH"
	}

	if typ == base.LIMIT {
		orderType = "LIMIT"
	} else if typ == base.MARKET {
		orderType = "MARKET"
	} else if typ == base.STOP {
		orderType = "STOP"
	} else if typ == base.STOPMARKET {
		orderType = "STOP_MARKET"
	} else if typ == base.TAKEPROFIT {
		orderType = "TAKE_PROFIT"
	} else if typ == base.TAKEPROFITMARKET {
		orderType = "TAKE_PROFIT_MARKET"
	}
	if typ == base.LIMIT {
		result, err := c.FutureClient.NewCreateOrderService().
			Symbol(symbol).
			Side(futures.SideType(orderSide)).
			PositionSide(futures.PositionSideType(pSide)).
			Type(futures.OrderType(orderType)).TimeInForce("GTC").
			Price(price).Quantity(size).
			ClosePosition(closePosition).
			PriceProtect(priceProtect).
			Do(context.Background())
		if err != nil {
			return "", err
		}
		return fmt.Sprint(result.OrderID), err

	} else if typ == base.MARKET {
		result, err := c.FutureClient.NewCreateOrderService().
			Symbol(symbol).
			Side(futures.SideType(orderSide)).
			PositionSide(futures.PositionSideType(pSide)).
			Type(futures.OrderType(orderType)).
			Quantity(size).
			ClosePosition(closePosition).
			PriceProtect(priceProtect).
			Do(context.Background())
		if err != nil {
			return "", err
		}
		return fmt.Sprint(result.OrderID), err

	} else if typ == base.STOP || typ == base.TAKEPROFIT {
		result, err := c.FutureClient.NewCreateOrderService().
			Symbol(symbol).
			Side(futures.SideType(orderSide)).
			PositionSide(futures.PositionSideType(pSide)).
			Type(futures.OrderType(orderType)).
			Price(price).Quantity(size).
			StopPrice(stopPrice).
			ClosePosition(closePosition).
			PriceProtect(priceProtect).
			Do(context.Background())
		if err != nil {
			return "", err
		}
		return fmt.Sprint(result.OrderID), err

	} else {
		result, err := c.FutureClient.NewCreateOrderService().
			Symbol(symbol).
			Side(futures.SideType(orderSide)).
			PositionSide(futures.PositionSideType(pSide)).
			Type(futures.OrderType(orderType)).
			Quantity(size).
			StopPrice(stopPrice).
			ClosePosition(closePosition).
			PriceProtect(priceProtect).
			Do(context.Background())
		if err != nil {
			return "", err
		}
		return fmt.Sprint(result.OrderID), err

	}

}

func (c *Client) Dual(dualSize bool) (bool, error) {
	err := c.FutureClient.
		NewChangePositionModeService().
		DualSide(dualSize).
		Do(context.Background())

	if err != nil {
		return false, err
	}
	return true, err
}

func (c *Client) CheckDual() (bool, error) {
	dual, err := c.FutureClient.
		NewGetPositionModeService().
		Do(context.Background())

	if err != nil {
		return false, err
	}
	return dual.DualSidePosition, err
}

func (c *Client) NewFuture(params []byte) error {
	if c.FutureClient == nil {

		sj, err := simplejson.NewJson(params)
		if err != nil {
			return err
		}
		_ = sj.Get("url").MustString()
		apiKey := sj.Get("apiKey").MustString()
		secretKey := sj.Get("secretKey").MustString()
		_ = sj.Get("password").MustString()

		// init client by config
		binance.UseTestnet = false
		c.FutureClient = futures.NewClient(apiKey, secretKey)

		return nil
	}

	return errors.New("binance client has not been initialized")
}

func (c *Client) GetFutureBalance() (models.FutureBalance, error) {
	account, err := c.FutureClient.
		NewGetBalanceService().
		Do(context.Background())

	if err != nil {
		return models.FutureBalance{}, err
	}
	var balance models.FutureBalance
	for _, b := range account {
		if b.Asset == "USDT" {
			balance = models.FutureBalance{
				Asset:            b.Asset,
				TotalBalance:     b.Balance,
				CrossBalance:     b.CrossWalletBalance,
				AvailableBalance: b.AvailableBalance,
			}
		}
	}

	return balance, nil
}

func (c *Client) FutureDepth(symbol, limit string) (models.WsData, error) {
	parseInt, err := strconv.Atoi(limit)
	if err != nil {
		return models.WsData{}, err
	}
	var info struct {
		LastUpdateId int64         `json:"lastUpdateId"`
		Bids         []binance.Bid `json:"bids"`
		Asks         []binance.Ask `json:"asks"`
	}

	res, err := c.FutureClient.NewDepthService().
		Symbol(symbol).
		Limit(parseInt).
		Do(context.Background())

	resByre, err := json.Marshal(res)
	if err != nil {
		return models.WsData{}, err
	}

	err = json.Unmarshal(resByre, &info)
	if err != nil {
		return models.WsData{}, err
	}
	var rawD models.WsData
	var bids []models.PriceLevel
	var asks []models.PriceLevel

	for _, bid := range info.Bids {
		b := models.PriceLevel{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}

		bids = append(bids, b)
	}

	for _, ask := range info.Asks {

		a := models.PriceLevel{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}

		asks = append(asks, a)
	}

	rawD.Time = info.LastUpdateId
	rawD.Bids = bids
	rawD.Asks = asks

	return rawD, nil
}

func (c *Client) GetFutureMarketPrice(symbol string) (string, error) {
	prices, err := c.FutureClient.NewListPricesService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return "", err
	}

	return prices[0].Price, nil
}

func (c *Client) GetMarkPriceAndFundingRate(symbol string) (models.FundingRate, error) {
	FR, err := c.FutureClient.NewPremiumIndexService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return models.FundingRate{}, err
	}
	result := models.FundingRate{
		Symbol:               FR[0].Symbol,
		MarkPrice:            FR[0].MarkPrice,
		IndexPrice:           "",
		EstimatedSettlePrice: "",
		LastFundingRate:      FR[0].LastFundingRate,
		NextFundingTime:      FR[0].NextFundingTime,
		InterestRate:         "",
		Time:                 FR[0].Time,
	}
	return result, err
}

func (c *Client) GetDepositAddress(token, chain string) (string, error) {
	address, err := c.Client.NewGetDepositAddressService().
		Coin(token).Network(chain).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return address.Address, err
}

func (c *Client) Withdraw(token, chain, to, amount string) (string, error) {

	//TODO implement me
	panic("implement me")
}

func (c *Client) GetAllTickers() ([]models.SymbolTicker, error) {

	//TODO implement me
	panic("implement me")
}

func (c *Client) LimitHiddenOrder(symbol, side, price, size string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) LimitHiddenOrders(symbol string, ol []models.OrderList) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetFeeFromFilled(symbol, id string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) New(params []byte) error {
	if c.Client == nil {

		sj, err := simplejson.NewJson(params)
		if err != nil {
			return err
		}
		_ = sj.Get("url").MustString()
		apiKey := sj.Get("apiKey").MustString()
		secretKey := sj.Get("secretKey").MustString()
		_ = sj.Get("password").MustString()

		// init client by config
		binance.UseTestnet = false
		c.Client = binance.NewClient(apiKey, secretKey)

		return nil
	}

	return errors.New("binance client has not been initialized")
}

func (c *Client) NewWithParams(apiKey, secretKey string) error {
	if c.Client == nil {
		// init client by config
		binance.UseTestnet = false
		c.Client = binance.NewClient(apiKey, secretKey)
		return nil
	}

	return errors.New("binance client has not been initialized")
}

func (c *Client) IceBergOrder(symbol, side, typ, price, size, ice string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetAccountBalance(currency string) ([]string, error) {

	account, err := c.Client.
		NewGetAccountService().
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	var res []string

	for _, balance := range account.Balances {
		if balance.Asset == currency {
			a, _ := strconv.ParseFloat(balance.Free, 64)
			b, _ := strconv.ParseFloat(balance.Locked, 64)
			d := strconv.FormatFloat(a+b, 'f', 5, 64)
			res = append(res, balance.Free, balance.Locked, d)
			return res, nil
		}
	}
	return nil, nil
}

func (c *Client) MarketOrder(symbol, side, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.Client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeMarket).
		//	TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(size).
		//	Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(order.OrderID, 10), nil
}

func (c *Client) LimitOrder(symbol, side, price, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.Client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(size).
		Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(order.OrderID, 10), nil
}

func (c *Client) TakerOrder(symbol, side, price, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.Client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeIOC).
		Quantity(size).
		Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(order.OrderID, 10), nil
}

func (c *Client) CancelOrder(symbol, id string) (bool, error) {
	parseInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return false, err
	}

	resp, err := c.Client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(parseInt).
		Do(context.Background())
	if err != nil {
		return false, err
	}

	if resp.Status != "CANCELED" {
		return false, err
	}
	return true, nil
}

func (c *Client) CancelOrders(symbol string) error {
	_, err := c.Client.NewCancelOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return err
	}

	return err
}

func (c *Client) GetOrder(symbol, id string) (models.OrderInfo, error) {

	oid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return models.OrderInfo{}, err
	}

	order, err := c.Client.NewGetOrderService().
		Symbol(symbol).
		OrderID(oid).
		Do(context.Background())
	if err != nil {
		return models.OrderInfo{}, err
	}

	var side string

	if order.Side == "BUY" {
		side = base.BID
	} else if order.Side == "SELL" {
		side = base.ASK
	}
	var status string
	if order.Status == "NEW" {
		status = base.OPEN
	} else if order.Status == "CANCELED" || order.Status == "EXPIRED" {
		status = base.CANCELED
	} else if order.Status == "FILLED" {
		status = base.FILLED
	} else if order.Status == "PARTIALLY_FILLED" {
		status = base.PARTIALLY
	}
	var typ string
	if order.Type == "LIMIT" && order.TimeInForce == "GTC" {
		typ = base.LIMIT
	} else if order.Type == "LIMIT" && order.TimeInForce == "IOC" {
		typ = base.TAKER
	} else if order.Type == "LIMIT_MAKER" {
		typ = base.MAKER
	} else if order.Type == "MARKET" {
		typ = base.MARKET
	}

	o := models.OrderInfo{
		OrderID:  strconv.FormatInt(order.OrderID, 10),
		Symbol:   order.Symbol,
		Side:     side,
		Price:    order.Price,
		Quantity: order.OrigQuantity,
		Status:   status,
		Type:     typ,
		USDT:     order.CummulativeQuoteQuantity,
		Filled:   order.ExecutedQuantity,
		Time:     order.Time,
	}

	return o, nil
}

func (c *Client) GetOpenOrders(symbol string) ([]models.OrderInfo, error) {
	orders, err := c.Client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, nil
	}

	var orderInfos []models.OrderInfo
	var typ string

	for _, order := range orders {

		var side string

		if order.Side == "BUY" {
			side = base.BID
		} else if order.Side == "SELL" {
			side = base.ASK
		}
		if order.Type == "LIMIT" && order.TimeInForce == "GTC" {
			typ = base.LIMIT
		} else if order.Type == "LIMIT" && order.TimeInForce == "IOC" {
			typ = base.TAKER
		} else if order.Type == "LIMIT_MAKER" {
			typ = base.MAKER
		} else if order.Type == "MARKET" {
			typ = base.MARKET
		}

		o := models.OrderInfo{
			OrderID:  strconv.FormatInt(order.OrderID, 10),
			Symbol:   order.Symbol,
			Side:     side,
			Price:    order.Price,
			Quantity: order.OrigQuantity,
			Filled:   order.ExecutedQuantity,
			Type:     typ,
			Time:     order.Time,
		}

		orderInfos = append(orderInfos, o)

	}

	return orderInfos, nil
}

func (c *Client) GetOpenOrdersWithSide(symbol, side string) ([]models.OrderInfo, error) {
	orders, err := c.GetOpenOrders(symbol)
	if err != nil {
		return nil, err
	}
	var info []models.OrderInfo
	for _, o := range orders {
		if o.Side == side {
			info = append(info, o)
		}
	}
	return info, err
}

func (c *Client) GetOpenSplitOrders(symbol string) ([]models.OrderInfo, []models.OrderInfo, error) {
	orders, err := c.Client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, nil, err
	}

	if len(orders) == 0 {
		return nil, nil, nil
	}

	var buys []models.OrderInfo
	var sells []models.OrderInfo
	var typ string

	for _, order := range orders {

		var side string
		if order.Type == "LIMIT" && order.TimeInForce == "GTC" {
			typ = base.LIMIT
		} else if order.Type == "LIMIT" && order.TimeInForce == "IOC" {
			typ = base.TAKER
		} else if order.Type == "LIMIT_MAKER" {
			typ = base.MAKER
		} else if order.Type == "MARKET" {
			typ = base.MARKET
		}

		if order.Side == "BUY" {
			side = base.BID

			o := models.OrderInfo{
				OrderID:  strconv.FormatInt(order.OrderID, 10),
				Symbol:   order.Symbol,
				Side:     side,
				Price:    order.Price,
				Quantity: order.OrigQuantity,
				Filled:   order.ExecutedQuantity,
				Type:     typ,
				Time:     order.Time,
			}

			buys = append(buys, o)

		} else if order.Side == "SELL" {
			side = base.ASK

			o := models.OrderInfo{
				OrderID:  strconv.FormatInt(order.OrderID, 10),
				Symbol:   order.Symbol,
				Side:     side,
				Price:    order.Price,
				Quantity: order.OrigQuantity,
				Type:     typ,
				Time:     order.Time,
			}

			sells = append(sells, o)
		}

	}

	return buys, sells, nil
}

func (c *Client) GetMarketPrice(symbol string) (string, error) {
	prices, err := c.Client.NewListPricesService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return "", err
	}

	return prices[0].Price, nil
}

func (c *Client) Depth(symbol, limit string) (models.WsData, error) {
	parseInt, err := strconv.Atoi(limit)
	if err != nil {
		return models.WsData{}, err
	}
	var info struct {
		LastUpdateId int64         `json:"lastUpdateId"`
		Bids         []binance.Bid `json:"bids"`
		Asks         []binance.Ask `json:"asks"`
	}

	res, err := c.Client.NewDepthService().
		Symbol(symbol).
		Limit(parseInt).
		Do(context.Background())

	resByre, err := json.Marshal(res)
	if err != nil {
		return models.WsData{}, err
	}

	err = json.Unmarshal(resByre, &info)
	if err != nil {
		return models.WsData{}, err
	}
	var rawD models.WsData
	var bids []models.PriceLevel
	var asks []models.PriceLevel

	for _, bid := range info.Bids {
		b := models.PriceLevel{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}

		bids = append(bids, b)
	}

	for _, ask := range info.Asks {

		a := models.PriceLevel{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}

		asks = append(asks, a)
	}

	rawD.Time = info.LastUpdateId
	rawD.Bids = bids
	rawD.Asks = asks

	return rawD, nil
}

func (c *Client) LimitOrders(symbol string, ol []models.OrderList) ([]string, error) {
	l := len(ol)
	IdList := make([]string, 0, l)
	var err error
	var id string

	for i := 0; i < l; i++ {

		side := ol[i].Side

		price := ol[i].Price
		size := ol[i].Size
		id, err = c.LimitOrder(symbol, side, price, size)
		if err != nil {
			return nil, err
		}

		IdList = append(IdList, id)
		time.Sleep(100 * time.Millisecond)
	}
	return IdList, err
}

func (c *Client) TakerOrders(symbol string, ol []models.OrderList) ([]string, error) {
	l := len(ol)
	IdList := make([]string, 0, l)
	var err error
	var id string

	for i := 0; i < l; i++ {

		side := ol[i].Side
		price := ol[i].Price
		size := ol[i].Size
		id, err = c.TakerOrder(side, symbol, price, size)
		if err != nil {
			return nil, err
		}
		IdList = append(IdList, id)
		time.Sleep(100 * time.Millisecond)
	}
	return IdList, err
}

func (c *Client) MakerOrder(symbol, side, price, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.Client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeLimitMaker).
		Quantity(size).
		Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(order.OrderID, 10), err
}

func (c *Client) MakerOrders(symbol string, ol []models.OrderList) ([]string, error) {
	l := len(ol)
	IdList := make([]string, 0, l)
	var err error
	var id string

	for i := 0; i < l; i++ {

		side := ol[i].Side
		price := ol[i].Price
		size := ol[i].Size
		id, err = c.MakerOrder(side, symbol, price, size)
		if err != nil {
			return nil, err
		}
		IdList = append(IdList, id)
		time.Sleep(100 * time.Millisecond)
	}
	return IdList, err
}

func (c *Client) GetPairInfo(symbol string) (models.PairInfo, error) {
	pair, err := c.Client.
		NewExchangeInfoService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return models.PairInfo{}, err
	}

	depth, err := c.Depth(symbol, "5")
	if err != nil {
		return models.PairInfo{}, err
	}
	q := depth.Bids[0].Quantity
	l := len(q)
	for i := l; i >= 0; i-- {
		if q[len(q)-1] == '0' {
			q = q[0 : len(q)-1]
		} else {
			break
		}
	}

	p := depth.Bids[0].Price
	l = len(p)
	for i := l; i >= 0; i-- {
		if p[len(p)-1] == '0' {
			p = p[0 : len(p)-1]
		} else {
			break
		}
	}

	info := models.PairInfo{
		MinBaseAmount:   strconv.FormatFloat(float64(1/10^pair.Symbols[0].BaseAssetPrecision), 'f', 7, 64),
		MinQuoteAmount:  strconv.FormatFloat(float64(1/10^pair.Symbols[0].QuoteAssetPrecision), 'f', 7, 64),
		AmountPrecision: len(strings.Split(q, ".")[1]),
		Precision:       len(strings.Split(p, ".")[1]),
	}
	return info, err

}

func (c *Client) GetTradingFee(symbol string) (models.TradingFee, error) {
	account, err := c.Client.NewTradeFeeService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return models.TradingFee{}, err
	}

	info := models.TradingFee{
		Symbol:                symbol,
		TakerFeeFromApi:       account[0].TakerCommission,
		MakerFeeFromApi:       account[0].MakerCommission,
		TakerFeeFromRealOrder: "",
		MakerFeeFromRealOrder: "",
		IfDiscount:            false,
	}
	return info, err
}
