package okx

import (
	"AxonTrading/base"
	"AxonTrading/models"
	"AxonTrading/tools"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	BaseUrl   string
	AccessKey string
	SecretKey string
	Password  string
	Client    *http.Client
}

type convertCoin struct {
	InstId string `json:"instId"`
	Px     string `json:"px"`
	Sz     string `json:"sz"`
	Type   string `json:"type"`
	Unit   string `json:"unit"`
}

func (c *Client) convertContractCoin(typ, symbol, size string) (convertCoin, error) {
	param := map[string]string{"type": typ, "instId": symbol, "sz": size, "coin": "coin"}
	url := "/api/v5/public/convert-contract-coin"
	resp, err := c.do(http.MethodGet, url, false, param)
	if err != nil {
		return convertCoin{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return convertCoin{}, errors.New("http request error" + string(rune(resp.StatusCode)))
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return convertCoin{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstId string `json:"instId"`
			Px     string `json:"px"`
			Sz     string `json:"sz"`
			Type   string `json:"type"`
			Unit   string `json:"unit"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return convertCoin{}, err
	}
	if bodyMarshal.Code != "0" {
		return convertCoin{}, errors.New(bodyMarshal.Msg)
	}
	rstData := bodyMarshal.Data[0]
	return convertCoin{InstId: rstData.InstId, Px: rstData.Px, Sz: rstData.Sz, Type: rstData.Type, Unit: rstData.Unit}, nil
}

func (c *Client) NewFuture(params []byte) error {
	c.Client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 2,
		},
	}

	sj, err := simplejson.NewJson(params)
	if err != nil {
		return err
	}
	baseUrl := sj.Get("url").MustString()
	apiKey := sj.Get("apiKey").MustString()
	secretKey := sj.Get("secretKey").MustString()
	password := sj.Get("password").MustString()

	c.BaseUrl = baseUrl
	c.AccessKey = apiKey
	c.SecretKey = secretKey
	c.Password = password

	return nil
}

func (c *Client) ChangeMarginType(symbol, typ string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) ChangePositionMargin(symbol, positionSide, amount string, typ int) (bool, error) {
	url := "/api/v5/account/position/margin-balance"
	param := map[string]string{"instId": symbol + "-SWAP", "posSide": "", "amt": amount, "type": ""}
	if typ == base.ADDMARGIN {
		param["type"] = "add"
	} else if typ == base.REMOVEMARGIN {
		param["type"] = "reduce"
	}
	if positionSide == base.LONG {
		param["posSide"] = "long"
	} else if positionSide == base.SHORT {
		param["posSide"] = "short"
	}
	paramByte, _ := json.Marshal(param)
	resp, err := c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Amt      string `json:"amt"`
			Ccy      string `json:"ccy"`
			InstId   string `json:"instId"`
			Leverage string `json:"leverage"`
			PosSide  string `json:"posSide"`
			Type     string `json:"type"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return false, err
	}
	if bodyMarshal.Code != "0" {
		return false, errors.New(bodyMarshal.Msg)
	}
	return true, nil
}

func HttpErr(code int) error {
	return errors.New(fmt.Sprintf("http request err code:%v", code))
}

func (c *Client) GetFutureBalance() (models.FutureBalance, error) {
	url := "/api/v5/account/balance"
	//asset := "USDT"
	param := map[string]string{}
	resp, err := c.do(http.MethodGet, url, true, param)
	defer resp.Body.Close()
	if err != nil {
		return models.FutureBalance{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.FutureBalance{}, HttpErr(resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.FutureBalance{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			AdjEq   string `json:"adjEq"`
			Details []struct {
				AvailBal      string `json:"availBal"`
				AvailEq       string `json:"availEq"`
				CashBal       string `json:"cashBal"`
				Ccy           string `json:"ccy"`
				CrossLiab     string `json:"crossLiab"`
				DisEq         string `json:"disEq"`
				Eq            string `json:"eq"`
				EqUsd         string `json:"eqUsd"`
				FrozenBal     string `json:"frozenBal"`
				Interest      string `json:"interest"`
				IsoEq         string `json:"isoEq"`
				IsoLiab       string `json:"isoLiab"`
				IsoUpl        string `json:"isoUpl"`
				Liab          string `json:"liab"`
				MaxLoan       string `json:"maxLoan"`
				MgnRatio      string `json:"mgnRatio"`
				NotionalLever string `json:"notionalLever"`
				OrdFrozen     string `json:"ordFrozen"`
				Twap          string `json:"twap"`
				UTime         string `json:"uTime"`
				Upl           string `json:"upl"`
				UplLiab       string `json:"uplLiab"`
				StgyEq        string `json:"stgyEq"`
				SpotInUseAmt  string `json:"spotInUseAmt"`
			} `json:"details"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &bodyMarshal)
	if err != nil {
		return models.FutureBalance{}, err
	}
	if bodyMarshal.Code != "0" {
		return models.FutureBalance{}, errors.New(bodyMarshal.Msg)
	}
	for _, v := range bodyMarshal.Data[0].Details {
		if v.Ccy == "USDT" {
			return models.FutureBalance{Asset: v.Ccy, TotalBalance: v.CashBal, CrossBalance: "", AvailableBalance: v.AvailBal}, nil
		}
	}
	return models.FutureBalance{}, errors.New("USDT NOT FOUND")
}

// FutureDepth
// Example: c.FutureDepth("BTC-USDT", "5")
func (c *Client) FutureDepth(symbol, limit string) (models.WsData, error) {
	url := "/api/v5/market/books"
	param := map[string]string{"instId": symbol + "-SWAP", "sz": limit}
	resp, err := c.do(http.MethodGet, url, false, param)
	defer resp.Body.Close()
	if err != nil {
		return models.WsData{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.WsData{}, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.WsData{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Asks [][]string `json:"asks"`
			Bids [][]string `json:"bids"`
			Ts   string     `json:"ts"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return models.WsData{}, err
	}
	if bodyMarshal.Code != "0" {
		return models.WsData{}, errors.New(bodyMarshal.Msg)
	}
	var rst models.WsData
	rst.Time, _ = strconv.ParseInt(bodyMarshal.Data[0].Ts, 10, 64)
	for _, v := range bodyMarshal.Data[0].Asks {
		rst.Asks = append(rst.Asks, models.PriceLevel{Price: v[0], Quantity: v[3]})
	}
	for _, v := range bodyMarshal.Data[0].Bids {
		rst.Bids = append(rst.Bids, models.PriceLevel{Price: v[0], Quantity: v[3]})
	}
	return rst, nil
}

func (c *Client) GetFutureMarketPrice(symbol string) (string, error) {
	url := "/api/v5/market/ticker"
	param := map[string]string{"instId": symbol + "-SWAP"}
	resp, err := c.do(http.MethodGet, url, false, param)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Last string `json:"last"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return "", err
	}
	if bodyMarshal.Code != "0" {
		return "", errors.New(bodyMarshal.Msg)
	}
	return bodyMarshal.Data[0].Last, nil
}

func (c *Client) GetFundingRate(symbol string) (models.FundingRate, error) {
	url := "/api/v5/public/funding-rate"
	param := map[string]string{"instId": symbol + "-SWAP"}
	resp, err := c.do(http.MethodGet, url, false, param)
	defer resp.Body.Close()
	if err != nil {
		return models.FundingRate{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.FundingRate{}, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.FundingRate{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			FundingRate     string `json:"fundingRate"`
			FundingTime     string `json:"fundingTime"`
			InstId          string `json:"instId"`
			InstType        string `json:"instType"`
			NextFundingRate string `json:"nextFundingRate"`
			NextFundingTime string `json:"nextFundingTime"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return models.FundingRate{}, err
	}
	if bodyMarshal.Code != "0" {
		return models.FundingRate{}, errors.New(bodyMarshal.Msg)
	}
	var rst models.FundingRate
	NextFundingTime, err := strconv.ParseInt(bodyMarshal.Data[0].NextFundingTime, 10, 64)
	rst = models.FundingRate{LastFundingRate: bodyMarshal.Data[0].FundingRate, NextFundingTime: NextFundingTime}
	return rst, nil
}

func (c *Client) GetMarkPriceAndFundingRate(symbol string) (models.FundingRate, error) {
	url := "/api/v5/public/mark-price"
	param := map[string]string{"instId": symbol + "-SWAP", "instType": "SWAP"}
	resp, err := c.do(http.MethodGet, url, false, param)
	defer resp.Body.Close()
	if err != nil {
		return models.FundingRate{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.FundingRate{}, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.FundingRate{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			MarkPrice string `json:"markPx"`
			InstId    string `json:"instId"`
			InstType  string `json:"instType"`
			Ts        string `json:"ts"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return models.FundingRate{}, err
	}
	if bodyMarshal.Code != "0" {
		return models.FundingRate{}, errors.New(bodyMarshal.Msg)
	}
	MarkPrice := bodyMarshal.Data[0].MarkPrice
	Symbol := bodyMarshal.Data[0].InstId
	t, err := strconv.ParseInt(bodyMarshal.Data[0].Ts, 10, 64)
	if err != nil {
		return models.FundingRate{}, err
	}
	var rst models.FundingRate
	partFundingData, err := c.GetFundingRate(symbol)
	if err != nil {
		return models.FundingRate{}, err
	}
	rst = models.FundingRate{MarkPrice: MarkPrice, Symbol: Symbol, LastFundingRate: partFundingData.LastFundingRate, NextFundingTime: partFundingData.NextFundingTime, Time: t}
	return rst, nil
}

func (c *Client) Dual(dualSize bool) (bool, error) {
	url := "/api/v5/account/set-position-mode"
	var posMode string
	if dualSize {
		posMode = "long_short_mode"
	} else {
		posMode = "net_mode"
	}
	param := map[string]string{"posMode": posMode}
	paramByte, _ := json.Marshal(param)
	resp, err := c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			PosMode string `json:"posMode"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return false, err
	}
	if bodyMarshal.Code != "0" {
		return false, errors.New(bodyMarshal.Msg)
	}
	return bodyMarshal.Data[0].PosMode == posMode, nil
}

func (c *Client) CheckDual() (bool, error) {
	url := "/api/v5/account/config"
	param := map[string]string{}
	resp, err := c.do(http.MethodGet, url, true, param)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			PosMode string `json:"posMode"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return false, err
	}
	if bodyMarshal.Code != "0" {
		return false, errors.New(bodyMarshal.Msg)
	}
	return bodyMarshal.Data[0].PosMode == "long_short_mode", nil
}

// NewFutureOrder 下单
func (c *Client) NewFutureOrder(symbol, side, positionSide, typ, size, price, stopPrice, positionType string, closePosition, priceProtect bool) (string, error) {
	coin, err := c.convertContractCoin("1", symbol+"-SWAP", size)
	if err != nil {
		return "", err
	}

	Newsize := coin.Sz
	url := "/api/v5/trade/order"
	NewsizeFloat, _ := strconv.ParseFloat(Newsize, 64)
	Newsize = strconv.FormatFloat(NewsizeFloat*10, 'f', 8, 64)

	var orderSide, pSide, orderType, pType string
	if positionType == base.ISOLATED {
		pType = "isolated"
	} else if positionType == base.CROSSED {
		pType = "cross"
	}
	if side == base.BID {
		orderSide = "buy"
	} else if side == base.ASK {
		orderSide = "sell"
	}
	dual, err := c.CheckDual()
	if err != nil {
		return "", err
	}

	if dual == true && positionSide == base.LONG {
		pSide = "long"
	} else if dual == true && positionSide == base.SHORT {
		pSide = "short"
	} else if dual == false {
		pSide = "net"
	}
	param := map[string]string{"instId": symbol + "-SWAP", "side": orderSide, "posSide": pSide, "sz": Newsize, "tdMode": pType}

	if typ == base.LIMIT {
		orderType = "limit"
		param["ordType"] = orderType
		param["px"] = price

	} else if typ == base.MARKET {
		orderType = "market"
		param["ordType"] = orderType
	} else if typ == base.STOP {
		orderType = "limit"
		param["ordType"] = orderType
		param["slOrdPx"] = price
		param["slTriggerPx"] = stopPrice
		param["px"] = price

	} else if typ == base.STOPMARKET {
		orderType = "limit"
		param["ordType"] = orderType
		param["slOrdPx"] = "-1"
		param["slTriggerPx"] = stopPrice
		param["px"] = price

	} else if typ == base.TAKEPROFIT {
		orderType = "limit"
		param["ordType"] = orderType
		param["tpOrdPx"] = price
		param["tpTriggerPx"] = stopPrice
		param["px"] = price

	} else if typ == base.TAKEPROFITMARKET {
		orderType = "limit"
		param["ordType"] = orderType
		param["tpOrdPx"] = "-1"
		param["tpTriggerPx"] = stopPrice
		param["px"] = price

	}

	paramByte, _ := json.Marshal(param)
	fmt.Println(param)
	resp, err := c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	if err != nil {
		return "", err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			ClOrdId string `json:"clOrdId"`
			OrdId   string `json:"ordId"`
			Tag     string `json:"tag"`
			SCode   string `json:"sCode"`
			SMsg    string `json:"sMsg"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return "", err
	}
	if bodyMarshal.Code != "0" {
		return "", errors.New(bodyMarshal.Msg)
	}
	return bodyMarshal.Data[0].OrdId, nil
}

func (c *Client) GetFutureOrder(symbol, orderID string) (models.FutureOrderInfo, error) {
	url := "/api/v5/trade/order"
	// TODO: instId 确认
	param := map[string]string{"instId": symbol + "-SWAP", "ordId": orderID}
	resp, err := c.do(http.MethodGet, url, true, param)
	defer resp.Body.Close()
	if err != nil {
		return models.FutureOrderInfo{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.FutureOrderInfo{}, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.FutureOrderInfo{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstType           string `json:"instType"`
			InstId             string `json:"instId"`
			Ccy                string `json:"ccy"`
			OrdId              string `json:"ordId"`
			ClOrdId            string `json:"clOrdId"`
			Tag                string `json:"tag"`
			Px                 string `json:"px"`
			Sz                 string `json:"sz"`
			Pnl                string `json:"pnl"`
			OrdType            string `json:"ordType"`
			Side               string `json:"side"`
			PosSide            string `json:"posSide"`
			TdMode             string `json:"tdMode"`
			AccFillSz          string `json:"accFillSz"`
			FillPx             string `json:"fillPx"`
			TradeId            string `json:"tradeId"`
			FillSz             string `json:"fillSz"`
			FillTime           string `json:"fillTime"`
			Source             string `json:"source"`
			State              string `json:"state"`
			AvgPx              string `json:"avgPx"`
			Lever              string `json:"lever"`
			AttachAlgoClOrdId  string `json:"attachAlgoClOrdId"`
			TpTriggerPx        string `json:"tpTriggerPx"`
			TpTriggerPxType    string `json:"tpTriggerPxType"`
			TpOrdPx            string `json:"tpOrdPx"`
			SlTriggerPx        string `json:"slTriggerPx"`
			SlTriggerPxType    string `json:"slTriggerPxType"`
			SlOrdPx            string `json:"slOrdPx"`
			StpId              string `json:"stpId"`
			StpMode            string `json:"stpMode"`
			FeeCcy             string `json:"feeCcy"`
			Fee                string `json:"fee"`
			RebateCcy          string `json:"rebateCcy"`
			Rebate             string `json:"rebate"`
			TgtCcy             string `json:"tgtCcy"`
			Category           string `json:"category"`
			ReduceOnly         string `json:"reduceOnly"`
			CancelSource       string `json:"cancelSource"`
			CancelSourceReason string `json:"cancelSourceReason"`
			QuickMgnType       string `json:"quickMgnType"`
			AlgoClOrdId        string `json:"algoClOrdId"`
			AlgoId             string `json:"algoId"`
			UTime              string `json:"uTime"`
			CTime              string `json:"cTime"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return models.FutureOrderInfo{}, err
	}
	if bodyMarshal.Code != "0" {
		return models.FutureOrderInfo{}, errors.New(bodyMarshal.Msg)
	}
	bodyMarshalData := bodyMarshal.Data[0]
	var rst models.FutureOrderInfo
	ordID, err := strconv.Atoi(bodyMarshalData.OrdId)
	if err != nil {
		return models.FutureOrderInfo{}, errors.New(bodyMarshal.Msg)
	}
	uTime, err := strconv.ParseInt(bodyMarshalData.UTime, 10, 64)
	if err != nil {
		return models.FutureOrderInfo{}, errors.New(bodyMarshal.Msg)
	}
	fillTime, err := strconv.ParseInt(bodyMarshalData.FillTime, 10, 64)

	var ReduceOnly bool
	if bodyMarshalData.ReduceOnly == "true" {
		ReduceOnly = true
	} else {
		ReduceOnly = false
	}
	var side, posSide, state, orderType string
	if bodyMarshalData.Side == "buy" {
		side = base.BID
	} else if bodyMarshalData.Side == "sell" {
		side = base.ASK
	}

	if bodyMarshalData.PosSide == "long" {
		posSide = base.LONG
	} else if bodyMarshalData.PosSide == "short" {
		posSide = base.SHORT
	}

	if bodyMarshalData.State == "filled" {
		state = base.FILLED
	} else if bodyMarshalData.State == "canceled" {
		state = base.CANCELED
	} else if bodyMarshalData.State == "partially_filled" {
		state = base.PARTIALLY
	} else if bodyMarshalData.State == "live" {
		state = base.OPEN
	}

	if bodyMarshalData.OrdType == "market" {
		orderType = base.MARKET
	} else if bodyMarshalData.OrdType == "limit" {
		orderType = base.LIMIT
	} else if bodyMarshalData.OrdType == "post_only" {
		orderType = base.MAKER
	} else if bodyMarshalData.OrdType == "limit" && bodyMarshalData.SlTriggerPx != "" && bodyMarshalData.SlTriggerPx != "-1" {
		orderType = base.STOP

	} else if bodyMarshalData.OrdType == "limit" && bodyMarshalData.SlTriggerPx == "-1" {
		orderType = base.STOPMARKET

	} else if bodyMarshalData.OrdType == "limit" && bodyMarshalData.TpTriggerPx != "" && bodyMarshalData.TpTriggerPx != "-1" {
		orderType = base.TAKEPROFIT

	} else if bodyMarshalData.OrdType == "limit" && bodyMarshalData.TpTriggerPx == "-1" {
		orderType = base.TAKEPROFITMARKET

	}

	rst = models.FutureOrderInfo{
		AvgPrice: bodyMarshalData.AvgPx, OrderId: ordID, Status: state,
		UpdateTime: uTime, Type: orderType, Side: side,
		Symbol: bodyMarshalData.InstId, Price: bodyMarshalData.Px, Time: fillTime,
		PositionSide: posSide, ReduceOnly: ReduceOnly, StopPrice: bodyMarshalData.TpTriggerPx,
		ClosePosition: false, PriceProtect: false,
	}
	return rst, nil
}

func (c *Client) CancelFutureOrder(symbol, orderID string) (bool, error) {
	url := "/api/v5/trade/cancel-order"
	param := map[string]string{"instId": symbol + "-SWAP", "ordId": orderID}
	paramByte, err := json.Marshal(param)
	resp, err := c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			ClOrdId string `json:"clOrdId"`
			OrdId   string `json:"ordId"`
			SCode   string `json:"sCode"`
			SMsg    string `json:"sMsg"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return false, err
	}
	if bodyMarshal.Code != "0" {
		return false, errors.New(bodyMarshal.Msg)
	}
	if bodyMarshal.Data[0].SCode == "0" {
		return true, nil
	} else {
		return false, errors.New(bodyMarshal.Data[0].SMsg)
	}
}

func (c *Client) CancelFutureOrders(symbol string) error {
	cancelOrders, err := c.GetFutureOpenOrders(symbol)
	if err != nil {
		return err
	}
	fmt.Println(cancelOrders)
	param := make([]map[string]string, 0, len(cancelOrders))
	for _, v := range cancelOrders {
		param = append(param, map[string]string{"instId": symbol + "-SWAP", "ordId": strconv.Itoa(v.OrderId)})
	}
	url := "/api/v5/trade/cancel-batch-orders"
	paramByte, err := json.Marshal(param)
	fmt.Println(param)
	resp, err := c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respBody))
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			ClOrdId string `json:"clOrdId"`
			OrdId   string `json:"ordId"`
			SCode   string `json:"sCode"`
			SMsg    string `json:"sMsg"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return err
	}
	if bodyMarshal.Code != "0" {
		return errors.New(bodyMarshal.Msg)
	}
	if bodyMarshal.Data[0].SCode == "0" {
		return nil
	} else {
		return errors.New(bodyMarshal.Data[0].SMsg)
	}
}

func (c *Client) GetFutureOpenOrders(symbol string) ([]models.FutureOrderInfo, error) {
	url := "/api/v5/trade/orders-pending"
	param := map[string]string{"instId": symbol + "-SWAP", "instType": "SWAP"}
	resp, err := c.do(http.MethodGet, url, true, param)
	defer resp.Body.Close()
	if err != nil {
		return []models.FutureOrderInfo{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []models.FutureOrderInfo{}, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.FutureOrderInfo{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			AccFillSz         string `json:"accFillSz"`
			AvgPx             string `json:"avgPx"`
			CTime             string `json:"cTime"`
			Category          string `json:"category"`
			Ccy               string `json:"ccy"`
			ClOrdId           string `json:"clOrdId"`
			Fee               string `json:"fee"`
			FeeCcy            string `json:"feeCcy"`
			FillPx            string `json:"fillPx"`
			FillSz            string `json:"fillSz"`
			FillTime          string `json:"fillTime"`
			InstId            string `json:"instId"`
			InstType          string `json:"instType"`
			Lever             string `json:"lever"`
			OrdId             string `json:"ordId"`
			OrdType           string `json:"ordType"`
			Pnl               string `json:"pnl"`
			PosSide           string `json:"posSide"`
			Px                string `json:"px"`
			Rebate            string `json:"rebate"`
			RebateCcy         string `json:"rebateCcy"`
			Side              string `json:"side"`
			AttachAlgoClOrdId string `json:"attachAlgoClOrdId"`
			SlOrdPx           string `json:"slOrdPx"`
			SlTriggerPx       string `json:"slTriggerPx"`
			SlTriggerPxType   string `json:"slTriggerPxType"`
			Source            string `json:"source"`
			State             string `json:"state"`
			StpId             string `json:"stpId"`
			StpMode           string `json:"stpMode"`
			Sz                string `json:"sz"`
			Tag               string `json:"tag"`
			TdMode            string `json:"tdMode"`
			TgtCcy            string `json:"tgtCcy"`
			TpOrdPx           string `json:"tpOrdPx"`
			TpTriggerPx       string `json:"tpTriggerPx"`
			TpTriggerPxType   string `json:"tpTriggerPxType"`
			TradeId           string `json:"tradeId"`
			ReduceOnly        string `json:"reduceOnly"`
			QuickMgnType      string `json:"quickMgnType"`
			AlgoClOrdId       string `json:"algoClOrdId"`
			AlgoId            string `json:"algoId"`
			UTime             string `json:"uTime"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return []models.FutureOrderInfo{}, err
	}
	if bodyMarshal.Code != "0" {
		return []models.FutureOrderInfo{}, errors.New(bodyMarshal.Msg)
	}

	var rst = make([]models.FutureOrderInfo, 0, 10)
	var ReduceOnly bool
	var side, posSide, state, orderType string
	for _, v := range bodyMarshal.Data {
		ordID, err := strconv.Atoi(v.OrdId)
		if err != nil {
			continue
		}
		uTime, err := strconv.ParseInt(v.UTime, 10, 64)
		if err != nil {
			continue
		}
		if v.ReduceOnly == "true" {
			ReduceOnly = true
		} else {
			ReduceOnly = false
		}
		fillTime, err := strconv.ParseInt(v.FillTime, 10, 64)
		if v.Side == "buy" {
			side = base.UnifiedBuy
		} else if v.Side == "sell" {
			side = base.UnifiedSell
		}
		if v.PosSide == "long" {
			posSide = base.LONG
		} else if v.PosSide == "short" {
			posSide = base.SHORT
		}

		if v.State == "filled" {
			state = base.FILLED
		} else if v.State == "canceled" {
			state = base.CANCELED
		} else if v.State == "partially_filled" {
			state = base.PARTIALLY
		} else if v.State == "live" {
			state = base.OPEN
		}

		if v.OrdType == "market" {
			orderType = base.MARKET
		} else if v.OrdType == "limit" {
			orderType = base.LIMIT
		} else if v.OrdType == "post_only" {
			orderType = base.MAKER
		} else if v.OrdType == "limit" && v.SlTriggerPx != "" && v.SlTriggerPx != "-1" {
			orderType = base.STOP

		} else if v.OrdType == "limit" && v.SlTriggerPx == "-1" {
			orderType = base.STOPMARKET

		} else if v.OrdType == "limit" && v.TpTriggerPx != "" && v.TpTriggerPx != "-1" {
			orderType = base.TAKEPROFIT

		} else if v.OrdType == "limit" && v.TpTriggerPx == "-1" {
			orderType = base.TAKEPROFITMARKET

		}

		// TODO:订单状态常量确认 OrdType,State
		rst = append(rst, models.FutureOrderInfo{
			AvgPrice: v.AvgPx, OrderId: ordID, Status: state, UpdateTime: uTime,
			Type: orderType, Side: side, Symbol: v.InstId, Price: v.Px, Time: fillTime,
			PositionSide: posSide, ReduceOnly: ReduceOnly, StopPrice: v.TpTriggerPx,
			ClosePosition: false, PriceProtect: false,
		})
	}
	return rst, nil
}

func (c *Client) ChangeLeverage(symbol string, leverage int) (string, error) {
	url := "/api/v5/account/set-leverage"
	param := map[string]string{"instId": symbol + "-SWAP", "lever": strconv.Itoa(leverage), "mgnMode": "cross"}
	paramByte, _ := json.Marshal(param)
	resp, err := c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Lever   string `json:"lever"`
			MgnMode string `json:"mgnMode"`
			InstId  string `json:"instId"`
			PosSide string `json:"posSide"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return "", err
	}
	if bodyMarshal.Code != "0" {
		return "", errors.New(bodyMarshal.Msg)
	}

	param = map[string]string{"instId": symbol + "-SWAP", "lever": strconv.Itoa(leverage), "mgnMode": "isolated", "posSide": "long"}
	paramByte, _ = json.Marshal(param)
	resp, err = c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", HttpErr(resp.StatusCode)
	}
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return "", err
	}
	if bodyMarshal.Code != "0" {
		return "", errors.New(bodyMarshal.Msg)
	}

	param = map[string]string{"instId": symbol + "-SWAP", "lever": strconv.Itoa(leverage), "mgnMode": "isolated", "posSide": "short"}
	paramByte, _ = json.Marshal(param)
	resp, err = c.doPost(http.MethodPost, url, true, paramByte)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", HttpErr(resp.StatusCode)
	}
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return "", err
	}
	if bodyMarshal.Code != "0" {
		return "", errors.New(bodyMarshal.Msg)
	}
	return bodyMarshal.Data[0].Lever, nil
}

func (c *Client) GetPositionRisk(symbol string) ([]models.PositionInfo, error) {
	url := "/api/v5/account/positions"
	param := map[string]string{"instId": symbol + "-SWAP"}
	//param := map[string]string{}
	resp, err := c.do(http.MethodGet, url, true, param)
	defer resp.Body.Close()
	if err != nil {
		return []models.PositionInfo{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []models.PositionInfo{}, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.PositionInfo{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
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
			InstId         string `json:"instId"`
			InstType       string `json:"instType"`
			Interest       string `json:"interest"`
			IdxPx          string `json:"idxPx"`
			Last           string `json:"last"`
			UsdPx          string `json:"usdPx"`
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
			PosCcy         string `json:"posCcy"`
			PosId          string `json:"posId"`
			PosSide        string `json:"posSide"`
			SpotInUseAmt   string `json:"spotInUseAmt"`
			SpotInUseCcy   string `json:"spotInUseCcy"`
			ThetaBS        string `json:"thetaBS"`
			ThetaPA        string `json:"thetaPA"`
			TradeId        string `json:"tradeId"`
			BizRefId       string `json:"bizRefId"`
			BizRefType     string `json:"bizRefType"`
			QuoteBal       string `json:"quoteBal"`
			BaseBal        string `json:"baseBal"`
			BaseBorrowed   string `json:"baseBorrowed"`
			BaseInterest   string `json:"baseInterest"`
			QuoteBorrowed  string `json:"quoteBorrowed"`
			QuoteInterest  string `json:"quoteInterest"`
			UTime          string `json:"uTime"`
			Upl            string `json:"upl"`
			UplLastPx      string `json:"uplLastPx"`
			UplRatio       string `json:"uplRatio"`
			UplRatioLastPx string `json:"uplRatioLastPx"`
			VegaBS         string `json:"vegaBS"`
			VegaPA         string `json:"vegaPA"`
			CloseOrderAlgo []struct {
				AlgoId          string `json:"algoId"`
				SlTriggerPx     string `json:"slTriggerPx"`
				SlTriggerPxType string `json:"slTriggerPxType"`
				TpTriggerPx     string `json:"tpTriggerPx"`
				TpTriggerPxType string `json:"tpTriggerPxType"`
				CloseFraction   string `json:"closeFraction"`
			} `json:"closeOrderAlgo"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return []models.PositionInfo{}, err
	}
	if bodyMarshal.Code != "0" {
		return []models.PositionInfo{}, errors.New(bodyMarshal.Msg)
	}
	var rst = make([]models.PositionInfo, 0, 20)
	var mgnMode, posSide string
	for _, v := range bodyMarshal.Data {
		uTime, _ := strconv.ParseInt(v.UTime, 10, 64)
		if v.MgnMode == "cross" {
			mgnMode = base.CROSSED
		} else if v.MgnMode == "isolated" {
			mgnMode = base.ISOLATED
		}

		if v.PosSide == "long" {
			posSide = base.LONG
		} else if v.PosSide == "short" {
			posSide = base.SHORT
		}
		rst = append(rst, models.PositionInfo{
			Symbol: v.InstId, MarkPrice: v.MarkPx, MarginType: mgnMode,
			Leverage: v.Lever, UnRealizedProfit: v.Upl, UpdateTime: uTime,
			LiquidationPrice: v.LiqPx, IsolatedMargin: v.Imr, EntryPrice: v.AvgPx,
			PositionAmt: v.Pos, PositionSide: posSide, Notional: "", IsAutoAddMargin: "",
			IsolatedWallet: "", MaxNotionalValue: "",
		})
	}
	return rst, nil
}

func (c *Client) GetFutureTradingFee(symbol string) (models.TradingFee, error) {
	url := "/api/v5/account/trade-fee"
	param := map[string]string{"instType": "SWAP"}
	resp, err := c.do(http.MethodGet, url, true, param)
	defer resp.Body.Close()
	if err != nil {
		return models.TradingFee{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.TradingFee{}, HttpErr(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.TradingFee{}, err
	}
	var bodyMarshal struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Delivery  string `json:"delivery"`
			Exercise  string `json:"exercise"`
			InstType  string `json:"instType"`
			Level     string `json:"level"`
			Maker     string `json:"maker"`
			MakerU    string `json:"makerU"`
			MakerUSDC string `json:"makerUSDC"`
			Taker     string `json:"taker"`
			TakerU    string `json:"takerU"`
			TakerUSDC string `json:"takerUSDC"`
			Ts        string `json:"ts"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &bodyMarshal)
	if err != nil {
		return models.TradingFee{}, err
	}
	if bodyMarshal.Code != "0" {
		return models.TradingFee{}, errors.New(bodyMarshal.Msg)
	}
	return models.TradingFee{Symbol: symbol, TakerFeeFromApi: bodyMarshal.Data[0].TakerU, MakerFeeFromApi: bodyMarshal.Data[0].MakerU}, nil
}

func (c *Client) GetDepositAddress(token, chain string) (string, error) {
	p := "/api/v5/asset/deposit-address"
	m := make(map[string]string)

	m["ccy"] = token

	res, err := c.do(http.MethodGet, p, true, m)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Data []struct {
			Chain    string `json:"chain"`
			CtAddr   string `json:"ctAddr"`
			Ccy      string `json:"ccy"`
			To       string `json:"to"`
			Addr     string `json:"addr"`
			Selected bool   `json:"selected"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	err = d.Decode(&response)
	if response.Code != "0" {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}
	for _, add := range response.Data {
		if add.Chain == token+"-"+chain && add.To == "18" {
			return add.Addr, err
		}
	}
	return "", err
}

func (c *Client) Withdraw(token, chain, to, amount string) (string, error) {
	c.InnerTrans(token, amount, "18", "6")
	fee, err := c.CurrencyInfo(token, chain)
	if err != nil {
		return "", err
	}

	p := "/api/v5/asset/withdrawal"
	m := make(map[string]string)

	m["ccy"] = token
	m["amt"] = amount
	m["dest"] = "4"
	m["toAddr"] = to
	m["fee"] = fee
	m["chain"] = token + "-" + chain

	res, err := c.do(http.MethodPost, p, true, m)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Amt      string `json:"amt"`
			WdId     string `json:"wdId"`
			Ccy      string `json:"ccy"`
			ClientId string `json:"clientId"`
			Chain    string `json:"chain"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	return response.Data[0].WdId, err
}

func (c *Client) CurrencyInfo(token, chain string) (string, error) {
	p := "/api/v5/asset/currencies"
	m := make(map[string]string)

	m["ccy"] = token

	res, err := c.do(http.MethodGet, p, true, m)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			CanDep               bool   `json:"canDep"`
			CanInternal          bool   `json:"canInternal"`
			CanWd                bool   `json:"canWd"`
			Ccy                  string `json:"ccy"`
			Chain                string `json:"chain"`
			DepQuotaFixed        string `json:"depQuotaFixed"`
			DepQuoteDailyLayer2  string `json:"depQuoteDailyLayer2"`
			LogoLink             string `json:"logoLink"`
			MainNet              bool   `json:"mainNet"`
			MaxFee               string `json:"maxFee"`
			MaxWd                string `json:"maxWd"`
			MinDep               string `json:"minDep"`
			MinDepArrivalConfirm string `json:"minDepArrivalConfirm"`
			MinFee               string `json:"minFee"`
			MinWd                string `json:"minWd"`
			MinWdUnlockConfirm   string `json:"minWdUnlockConfirm"`
			Name                 string `json:"name"`
			NeedTag              bool   `json:"needTag"`
			UsedDepQuotaFixed    string `json:"usedDepQuotaFixed"`
			UsedWdQuota          string `json:"usedWdQuota"`
			WdQuota              string `json:"wdQuota"`
			WdTickSz             string `json:"wdTickSz"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}
	for _, info := range response.Data {
		if info.Chain == token+"-"+chain {
			return info.MinFee, err
		}
	}

	return "", err
}

func (c *Client) InnerTrans(token, amount, from, to string) (string, error) {
	p := "/api/v5/asset/transfer"
	m := make(map[string]string)

	m["ccy"] = token
	m["amt"] = amount
	m["from"] = from
	m["to"] = to

	res, err := c.do(http.MethodPost, p, true, m)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	return "", err
}

func (c *Client) New(params []byte) error {
	c.Client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 2,
		},
	}

	sj, err := simplejson.NewJson(params)
	if err != nil {
		return err
	}
	baseUrl := sj.Get("url").MustString()
	apiKey := sj.Get("apiKey").MustString()
	secretKey := sj.Get("secretKey").MustString()
	password := sj.Get("password").MustString()

	c.BaseUrl = baseUrl
	c.AccessKey = apiKey
	c.SecretKey = secretKey
	c.Password = password

	return nil
}

// Do the http request to the server
func (c *Client) do(method, path string, private bool, params ...map[string]string) (*http.Response, error) {
	u := fmt.Sprintf("%s%s", c.BaseUrl, path)
	var (
		r    *http.Request
		err  error
		j    []byte
		body string
	)
	if method == http.MethodGet {
		r, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}

		if len(params) > 0 {
			q := r.URL.Query()
			for k, v := range params[0] {
				q.Add(k, strings.ReplaceAll(v, "\"", ""))
			}
			r.URL.RawQuery = q.Encode()
			if len(params[0]) > 0 {
				path += "?" + r.URL.RawQuery
			}
		}
	} else {
		j, err = json.Marshal(params[0])
		if err != nil {
			return nil, err
		}
		body = string(j)
		if body == "{}" {
			body = ""
		}
		r, err = http.NewRequest(method, u, bytes.NewBuffer(j))
		if err != nil {
			return nil, err
		}
		r.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		return nil, err
	}
	if private {
		timestamp, sign := c.sign(method, path, body)
		r.Header.Add("OK-ACCESS-KEY", c.AccessKey)
		r.Header.Add("OK-ACCESS-PASSPHRASE", c.Password)
		r.Header.Add("OK-ACCESS-SIGN", sign)
		r.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		//r.Header.Add("x-simulated-trading", "1")
	}
	return c.Client.Do(r)
}

func (c *Client) doPost(method, path string, private bool, params []byte) (*http.Response, error) {
	u := fmt.Sprintf("%s%s", c.BaseUrl, path)
	var (
		r    *http.Request
		err  error
		body string
	)
	body = string(params)
	fmt.Println(body)
	if body == "{}" {
		body = ""
	}
	r, err = http.NewRequest(method, u, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}
	if private {
		timestamp, sign := c.sign(method, path, body)
		r.Header.Add("OK-ACCESS-KEY", c.AccessKey)
		r.Header.Add("OK-ACCESS-PASSPHRASE", c.Password)
		r.Header.Add("OK-ACCESS-SIGN", sign)
		r.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		//r.Header.Add("x-simulated-trading", "1")

	}
	return c.Client.Do(r)
}

func (c *Client) sign(method, path, body string) (string, string) {
	format := "2006-01-02T15:04:05.999Z07:00"
	t := time.Now().UTC().Format(format)
	ts := fmt.Sprint(t)
	s := ts + method + path + body
	p := []byte(s)
	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write(p)
	return ts, base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c *Client) GetAccountBalance(currency string) ([]string, error) {
	p := "/api/v5/account/balance"
	m := make(map[string]string)

	if len(currency) > 0 {
		m["ccy"] = currency
	}
	res, err := c.do(http.MethodGet, p, true, m)
	if err != nil {
		return []string{"0", "0", "0"}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return []string{"0", "0", "0"}, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Data []struct {
			Details []struct {
				AvailBal  string `json:"availBal"`
				CashBal   string `json:"cashBal"`
				Ccy       string `json:"ccy"`
				FrozenBal string `json:"frozenBal"`
			} `json:"details"`
			UTime string `json:"uTime"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
	err = d.Decode(&response)
	if response.Code != "0" {
		data, _ := ioutil.ReadAll(res.Body)
		return []string{"0", "0", "0"}, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	var bs []string
	if len(response.Data) > 0 && len(response.Data[0].Details) > 0 {
		bs = append(bs, response.Data[0].Details[0].AvailBal, response.Data[0].Details[0].FrozenBal, response.Data[0].Details[0].CashBal)
		return bs, nil
	} else {
		return []string{"0", "0", "0"}, fmt.Errorf("not found the currency")
	}
}

type (
	PlaceOrder struct {
		ID         string `json:"-"`
		InstID     string `json:"instId"`
		Ccy        string `json:"ccy,omitempty"`
		ClOrdID    string `json:"clOrdId,omitempty"`
		Tag        string `json:"tag,omitempty"`
		ReduceOnly bool   `json:"reduceOnly,omitempty"`
		Sz         string `json:"sz"`
		Px         string `json:"px,omitempty"`
		TdMode     string `json:"tdMode"`
		Side       string `json:"side"`
		PosSide    string `json:"posSide,omitempty"`
		OrdType    string `json:"ordType"`
		TgtCcy     string `json:"tgtCcy,omitempty"`
	}
	PlaceOrderResp struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			ClOrdId string `json:"clOrdId"`
			OrdId   string `json:"ordId"`
			Tag     string `json:"tag"`
			SCode   string `json:"sCode"`
			SMsg    string `json:"sMsg"`
		} `json:"data"`
	}
	CancelR struct {
		ID     string `json:"-"`
		InstID string `json:"instId"`
		OrdID  string `json:"ordId"`
	}
)

func (c *Client) placeOrder(req []PlaceOrder) (response PlaceOrderResp, err error) {
	p := "/api/v5/trade/order"
	var tmp interface{}
	tmp = req[0]
	m, err := json.Marshal(tmp)
	if len(req) > 1 {
		tmp = req
		p = "/api/v5/trade/batch-orders"
	}
	// m := tool.S2M(tmp)

	res, err := c.doPost(http.MethodPost, p, true, m)
	// res, err := c.do(http.MethodPost, p, true, m)
	if err != nil {
		return
	}
	defer res.Body.Close()
	d := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		// data, _ := ioutil.ReadAll(res.Body)
		err = fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, res)
		return
	}

	err = d.Decode(&response)

	if response.Code != "0" {
		err = fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response.Data)
		return
	}

	return
}

func (c *Client) setSide(side string) string {
	if side == base.BID {
		return "buy"
	}

	if side == base.ASK {
		return "sell"
	}

	return ""
}

func (c *Client) MarketOrder(symbol, side, size string) (string, error) {
	o := PlaceOrder{
		InstID:  symbol,
		TdMode:  "cash",
		Side:    c.setSide(side),
		OrdType: "market",
		Sz:      size,
		TgtCcy:  "base_ccy",
	}
	placeOrderResp, err := c.placeOrder([]PlaceOrder{o})
	if err != nil {
		return "", err
	}
	if len(placeOrderResp.Data) > 0 {
		return placeOrderResp.Data[0].OrdId, nil
	}
	return "", fmt.Errorf("not get the order")
}

func (c *Client) LimitOrder(symbol, side, price, size string) (string, error) {
	o := PlaceOrder{
		InstID:  symbol,
		TdMode:  "cash",
		Side:    c.setSide(side),
		OrdType: "limit",
		Sz:      size,
		Px:      price,
	}
	placeOrderResp, err := c.placeOrder([]PlaceOrder{o})
	if err != nil {
		return "", err
	}
	if len(placeOrderResp.Data) > 0 {
		return placeOrderResp.Data[0].OrdId, nil
	}
	return "", fmt.Errorf("not get the order")
}

func (c *Client) LimitHiddenOrder(symbol, side, price, size string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) LimitOrders(symbol string, ol []models.OrderList) ([]string, error) {

	var os []PlaceOrder
	for _, order := range ol {
		o := PlaceOrder{
			InstID:  symbol,
			TdMode:  "cash",
			Side:    c.setSide(order.Side),
			OrdType: "limit",
			Sz:      order.Size,
			Px:      order.Price,
		}
		os = append(os, o)
	}
	placeOrderResp, err := c.placeOrder(os)
	if err != nil {
		fmt.Println(placeOrderResp)
		return nil, err
	}
	if len(placeOrderResp.Data) > 0 {
		var orderIds []string
		for _, datum := range placeOrderResp.Data {
			orderIds = append(orderIds, datum.OrdId)
		}
		return orderIds, nil
	}
	return nil, fmt.Errorf("not get the order")
}

func (c *Client) MakerOrder(symbol, side, price, size string) (string, error) {
	o := PlaceOrder{
		InstID:  symbol,
		TdMode:  "cash",
		Side:    c.setSide(side),
		OrdType: "post_only",
		Sz:      size,
		Px:      price,
	}
	placeOrderResp, err := c.placeOrder([]PlaceOrder{o})
	if err != nil {
		return "", err
	}
	if len(placeOrderResp.Data) > 0 {
		return placeOrderResp.Data[0].OrdId, nil
	}
	return "", fmt.Errorf("not get the order")
}

func (c *Client) MakerOrders(symbol string, ol []models.OrderList) ([]string, error) {
	var os []PlaceOrder
	for _, order := range ol {
		o := PlaceOrder{
			InstID:  symbol,
			TdMode:  "cash",
			Side:    c.setSide(order.Side),
			OrdType: "post_only",
			Sz:      order.Size,
			Px:      order.Price,
		}
		os = append(os, o)
	}
	placeOrderResp, err := c.placeOrder(os)
	if err != nil {
		return nil, err
	}
	if len(placeOrderResp.Data) > 0 {
		var orderIds []string
		for _, datum := range placeOrderResp.Data {
			orderIds = append(orderIds, datum.OrdId)
		}
		return orderIds, nil
	}
	return nil, fmt.Errorf("not get the order")
}

func (c *Client) TakerOrder(symbol, side, price, size string) (string, error) {
	o := PlaceOrder{
		InstID:  symbol,
		TdMode:  "cash",
		Side:    c.setSide(side),
		OrdType: "ioc",
		Sz:      size,
		Px:      price,
	}
	placeOrderResp, err := c.placeOrder([]PlaceOrder{o})
	if err != nil {
		return "", err
	}
	if len(placeOrderResp.Data) > 0 {
		return placeOrderResp.Data[0].OrdId, nil
	}
	return "", fmt.Errorf("not get the order")
}

func (c *Client) TakerOrders(symbol string, ol []models.OrderList) ([]string, error) {
	var os []PlaceOrder
	for _, order := range ol {
		o := PlaceOrder{
			InstID:  symbol,
			TdMode:  "cash",
			Side:    c.setSide(order.Side),
			OrdType: "ioc",
			Sz:      order.Size,
			Px:      order.Price,
		}
		os = append(os, o)
	}
	placeOrderResp, err := c.placeOrder(os)
	if err != nil {
		return nil, err
	}
	if len(placeOrderResp.Data) > 0 {
		var orderIds []string
		for _, datum := range placeOrderResp.Data {
			orderIds = append(orderIds, datum.OrdId)
		}
		return orderIds, nil
	}
	return nil, fmt.Errorf("not get the order")
}

func (c *Client) CancelOrder(symbol, id string) (bool, error) {

	p := "/api/v5/trade/cancel-order"
	m := make(map[string]string)

	m["ordId"] = id
	m["instId"] = symbol

	res, err := c.do(http.MethodPost, p, true, m)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return false, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			ClOrdId string `json:"clOrdId"`
			OrdId   string `json:"ordId"`
			SCode   string `json:"sCode"`
			SMsg    string `json:"sMsg"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" || len(response.Data) == 0 {
		return false, fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
	}
	if response.Data[0].SCode != "0" {
		return false, nil
	}
	return true, nil
}

func (c *Client) CancelOrders(symbol string) error {

	openOrders, err := c.GetOpenOrders(symbol)
	if err != nil {
		return err
	}
	newOpenOrders := arrayInGroupsOf(openOrders, 20)
	for _, orders := range newOpenOrders {
		var cancelRs []CancelR
		for _, o := range orders {
			cancelRs = append(cancelRs, CancelR{
				InstID: symbol,
				OrdID:  o.OrderID,
			})
		}

		p := "/api/v5/trade/cancel-batch-orders"
		marshal, err := json.Marshal(cancelRs)
		if err != nil {
			return err
		}
		res, err := c.doPost(http.MethodPost, p, true, marshal)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			data, _ := ioutil.ReadAll(res.Body)
			return fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
		}

		d := json.NewDecoder(res.Body)

		var response struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
			Data []struct {
				ClOrdId string `json:"clOrdId"`
				OrdId   string `json:"ordId"`
				SCode   string `json:"sCode"`
				SMsg    string `json:"sMsg"`
			} `json:"data"`
		}

		err = d.Decode(&response)
		if response.Code != "0" {
			return fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
		}
	}
	return nil
}

func (c *Client) GetOrder(symbol, id string) (models.OrderInfo, error) {
	p := "/api/v5/trade/order"
	m := make(map[string]string)
	m["instId"] = symbol
	m["ordId"] = id

	orderInfo := models.OrderInfo{}

	res, err := c.do(http.MethodGet, p, true, m)
	if err != nil {
		return orderInfo, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return orderInfo, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstType           string `json:"instType"`
			InstId             string `json:"instId"`
			Ccy                string `json:"ccy"`
			OrdId              string `json:"ordId"`
			ClOrdId            string `json:"clOrdId"`
			Tag                string `json:"tag"`
			Px                 string `json:"px"`
			Sz                 string `json:"sz"`
			Pnl                string `json:"pnl"`
			OrdType            string `json:"ordType"`
			Side               string `json:"side"`
			PosSide            string `json:"posSide"`
			TdMode             string `json:"tdMode"`
			AccFillSz          string `json:"accFillSz"`
			FillPx             string `json:"fillPx"`
			TradeId            string `json:"tradeId"`
			FillSz             string `json:"fillSz"`
			FillTime           string `json:"fillTime"`
			Source             string `json:"source"`
			State              string `json:"state"`
			AvgPx              string `json:"avgPx"`
			Lever              string `json:"lever"`
			TpTriggerPx        string `json:"tpTriggerPx"`
			TpTriggerPxType    string `json:"tpTriggerPxType"`
			TpOrdPx            string `json:"tpOrdPx"`
			SlTriggerPx        string `json:"slTriggerPx"`
			SlTriggerPxType    string `json:"slTriggerPxType"`
			SlOrdPx            string `json:"slOrdPx"`
			FeeCcy             string `json:"feeCcy"`
			Fee                string `json:"fee"`
			RebateCcy          string `json:"rebateCcy"`
			Rebate             string `json:"rebate"`
			TgtCcy             string `json:"tgtCcy"`
			Category           string `json:"category"`
			ReduceOnly         string `json:"reduceOnly"`
			CancelSource       string `json:"cancelSource"`
			CancelSourceReason string `json:"cancelSourceReason"`
			QuickMgnType       string `json:"quickMgnType"`
			AlgoClOrdId        string `json:"algoClOrdId"`
			AlgoId             string `json:"algoId"`
			UTime              string `json:"uTime"`
			CTime              string `json:"cTime"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" || len(response.Data) == 0 {
		return orderInfo, fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
	}

	o := response.Data[0]
	orderInfo.Symbol = o.InstId
	orderInfo.Side = o.Side
	orderInfo.OrderID = o.OrdId
	orderInfo.Price = o.Px
	orderInfo.Filled = o.AccFillSz
	orderInfo.Quantity = o.Sz
	orderInfo.Status = o.State
	orderInfo.Type = o.InstType

	ctime, _ := strconv.Atoi(o.CTime)
	orderInfo.Time = int64(ctime)

	return orderInfo, nil
}

func (c *Client) GetOpenOrders(symbol string) ([]models.OrderInfo, error) {
	p := "/api/v5/trade/orders-pending"
	m := make(map[string]string)
	m["instType"] = "SPOT"
	m["instId"] = symbol

	var orders []models.OrderInfo

	res, err := c.do(http.MethodGet, p, true, m)
	if err != nil {
		return orders, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return orders, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			AccFillSz       string `json:"accFillSz"`
			AvgPx           string `json:"avgPx"`
			CTime           string `json:"cTime"`
			Category        string `json:"category"`
			Ccy             string `json:"ccy"`
			ClOrdId         string `json:"clOrdId"`
			Fee             string `json:"fee"`
			FeeCcy          string `json:"feeCcy"`
			FillPx          string `json:"fillPx"`
			FillSz          string `json:"fillSz"`
			FillTime        string `json:"fillTime"`
			InstId          string `json:"instId"`
			InstType        string `json:"instType"`
			Lever           string `json:"lever"`
			OrdId           string `json:"ordId"`
			OrdType         string `json:"ordType"`
			Pnl             string `json:"pnl"`
			PosSide         string `json:"posSide"`
			Px              string `json:"px"`
			Rebate          string `json:"rebate"`
			RebateCcy       string `json:"rebateCcy"`
			Side            string `json:"side"`
			SlOrdPx         string `json:"slOrdPx"`
			SlTriggerPx     string `json:"slTriggerPx"`
			SlTriggerPxType string `json:"slTriggerPxType"`
			Source          string `json:"source"`
			State           string `json:"state"`
			Sz              string `json:"sz"`
			Tag             string `json:"tag"`
			TdMode          string `json:"tdMode"`
			TgtCcy          string `json:"tgtCcy"`
			TpOrdPx         string `json:"tpOrdPx"`
			TpTriggerPx     string `json:"tpTriggerPx"`
			TpTriggerPxType string `json:"tpTriggerPxType"`
			TradeId         string `json:"tradeId"`
			ReduceOnly      string `json:"reduceOnly"`
			QuickMgnType    string `json:"quickMgnType"`
			AlgoClOrdId     string `json:"algoClOrdId"`
			AlgoId          string `json:"algoId"`
			UTime           string `json:"uTime"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" {
		return orders, fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
	}

	for _, o := range response.Data {
		var orderInfo models.OrderInfo
		orderInfo.Symbol = o.InstId
		orderInfo.Side = o.Side
		orderInfo.OrderID = o.OrdId
		orderInfo.Price = o.Px
		orderInfo.Filled = o.AccFillSz
		orderInfo.Quantity = o.Sz
		orderInfo.Status = o.State
		orderInfo.Type = o.InstType

		ctime, _ := strconv.Atoi(o.CTime)
		orderInfo.Time = int64(ctime)

		orders = append(orders, orderInfo)
	}

	return orders, nil
}

func (c *Client) GetOpenOrdersWithSide(symbol, side string) ([]models.OrderInfo, error) {
	var orders []models.OrderInfo
	openOrders, err := c.GetOpenOrders(symbol)
	if err != nil {
		return nil, err
	}
	for _, order := range openOrders {
		if order.Side == c.setSide(side) {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (c *Client) GetOpenSplitOrders(symbol string) ([]models.OrderInfo, []models.OrderInfo, error) {
	var bidOrders []models.OrderInfo
	var askOrders []models.OrderInfo
	openOrders, err := c.GetOpenOrders(symbol)
	if err != nil {
		return nil, nil, err
	}
	for _, order := range openOrders {
		if order.Side == c.setSide(base.BID) {
			bidOrders = append(bidOrders, order)
		}

		if order.Side == c.setSide(base.ASK) {
			askOrders = append(askOrders, order)
		}
	}
	return bidOrders, askOrders, nil
}

func (c *Client) GetMarketPrice(symbol string) (string, error) {
	p := "/api/v5/market/ticker"
	m := make(map[string]string)
	m["instId"] = symbol
	res, err := c.do(http.MethodGet, p, false, m)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstType  string `json:"instType"`
			InstId    string `json:"instId"`
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
			Ts        string `json:"ts"`
			SodUtc0   string `json:"sodUtc0"`
			SodUtc8   string `json:"sodUtc8"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" || len(response.Data) == 0 {
		return "", fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
	}
	return response.Data[0].Last, nil
}

func (c *Client) Depth(symbol, limit string) (models.WsData, error) {

	var ws models.WsData

	p := "/api/v5/market/books"
	m := make(map[string]string)
	m["sz"] = limit
	m["instId"] = symbol
	res, err := c.do(http.MethodGet, p, false, m)
	if err != nil {
		return ws, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return ws, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Asks [][]string `json:"asks"`
			Bids [][]string `json:"bids"`
			Ts   string     `json:"ts"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" || len(response.Data) == 0 || len(response.Data[0].Asks) == 0 || len(response.Data[0].Bids) == 0 {
		return ws, fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
	}

	var bids []models.PriceLevel
	var asks []models.PriceLevel

	for _, ask := range response.Data[0].Asks {
		asks = append(asks, models.PriceLevel{
			Price:    ask[0],
			Quantity: ask[1],
		})
	}

	for _, bid := range response.Data[0].Bids {
		bids = append(bids, models.PriceLevel{
			Price:    bid[0],
			Quantity: bid[1],
		})
	}

	ws.Asks = asks
	ws.Bids = bids
	ws.Time = time.Now().Unix()

	return ws, nil
}

func (c *Client) GetTradingFee(symbol string) (models.TradingFee, error) {
	p := "/api/v5/account/trade-fee"
	var fee models.TradingFee
	m := make(map[string]string)
	m["instType"] = "SPOT"
	m["instId"] = symbol
	res, err := c.do(http.MethodGet, p, true, m)
	if err != nil {
		return fee, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return fee, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Category  string `json:"category"`
			Delivery  string `json:"delivery"`
			Exercise  string `json:"exercise"`
			InstType  string `json:"instType"`
			Level     string `json:"level"`
			Maker     string `json:"maker"`
			MakerU    string `json:"makerU"`
			MakerUSDC string `json:"makerUSDC"`
			Taker     string `json:"taker"`
			TakerU    string `json:"takerU"`
			TakerUSDC string `json:"takerUSDC"`
			Ts        string `json:"ts"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" || len(response.Data) == 0 {
		return fee, fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
	}

	fee.Symbol = symbol
	fee.MakerFeeFromApi = response.Data[0].Maker
	fee.TakerFeeFromApi = response.Data[0].Taker

	return fee, nil
}

func (c *Client) GetPairInfo(symbol string) (models.PairInfo, error) {

	var pairInfo models.PairInfo

	p := "/api/v5/public/instruments"
	m := make(map[string]string)
	m["instType"] = "SPOT"
	m["instId"] = symbol
	res, err := c.do(http.MethodGet, p, false, m)
	if err != nil {
		return pairInfo, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return pairInfo, fmt.Errorf("response status code is not OK, response code is %d, body:%s", res.StatusCode, string(data))
	}

	d := json.NewDecoder(res.Body)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Alias        string `json:"alias"`
			BaseCcy      string `json:"baseCcy"`
			Category     string `json:"category"`
			CtMult       string `json:"ctMult"`
			CtType       string `json:"ctType"`
			CtVal        string `json:"ctVal"`
			CtValCcy     string `json:"ctValCcy"`
			ExpTime      string `json:"expTime"`
			InstFamily   string `json:"instFamily"`
			InstId       string `json:"instId"`
			InstType     string `json:"instType"`
			Lever        string `json:"lever"`
			ListTime     string `json:"listTime"`
			LotSz        string `json:"lotSz"`
			MaxIcebergSz string `json:"maxIcebergSz"`
			MaxLmtSz     string `json:"maxLmtSz"`
			MaxMktSz     string `json:"maxMktSz"`
			MaxStopSz    string `json:"maxStopSz"`
			MaxTriggerSz string `json:"maxTriggerSz"`
			MaxTwapSz    string `json:"maxTwapSz"`
			MinSz        string `json:"minSz"`
			OptType      string `json:"optType"`
			QuoteCcy     string `json:"quoteCcy"`
			SettleCcy    string `json:"settleCcy"`
			State        string `json:"state"`
			Stk          string `json:"stk"`
			TickSz       string `json:"tickSz"`
			Uly          string `json:"uly"`
		} `json:"data"`
	}

	err = d.Decode(&response)
	if response.Code != "0" || len(response.Data) == 0 {
		return pairInfo, fmt.Errorf("response status code is not OK, response code is %d, body:%+v", res.StatusCode, response)
	}

	pairInfo.Precision = tools.GetDecimalPlacesStr(response.Data[0].TickSz)
	pairInfo.AmountPrecision = tools.GetDecimalPlacesStr(response.Data[0].LotSz)
	return pairInfo, nil
}

func (c *Client) GetFeeFromFilled(symbol, id string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) IceBergOrder(symbol, side, typ, price, size, ice string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) LimitHiddenOrders(symbol string, ol []models.OrderList) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func arrayInGroupsOf(arr []models.OrderInfo, num int64) [][]models.OrderInfo {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]models.OrderInfo{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]models.OrderInfo, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}
