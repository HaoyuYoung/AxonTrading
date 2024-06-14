package okx

import (
	"AxonTrading/base"
	"fmt"

	"net/http"
	"testing"
	"time"
)

var c Client

var (
	baseUrl   string
	apiKey    string
	apiSecret string
	password  string
)

func init() {
	baseUrl = ""
	apiKey = ""
	apiSecret = ""
	password = ""
	c.Client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	c.BaseUrl = baseUrl
	c.AccessKey = apiKey
	c.SecretKey = apiSecret
	c.Password = password
}

/*
	func init() {
		c = &Client{}
		param, _ := json.Marshal(map[string]string{
			"url":       "https://aws.okx.com",
			"apiKey":    "",
			"secretKey": "",
			"password":  "!",
		})
		c.NewFuture(param)
	}
*/
func TestClient_GetBalance(t *testing.T) {
	balance, err := c.GetAccountBalance("USDC")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}
func TestClient_GetBalance1(t *testing.T) {
	balance, err := c.GetAccountBalance("USDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}

func TestClient_MarketOrder(t *testing.T) {
	order, err := c.MarketOrder("BTC-USDT", base.ASK, "10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(order)
}

func TestClient_LimitOrder(t *testing.T) {
	order, err := c.LimitOrder("BTC-USDT", base.ASK, "1", "10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(order)
}

func TestClient_CancelOrders(t *testing.T) {
	err := c.CancelOrders("BTC-USDT")
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_GetOrder(t *testing.T) {
	orderInfo, err := c.GetOrder("BTC-USDT", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderInfo)
}

func TestClient_GetOpenOrders(t *testing.T) {
	orderInfo, err := c.GetOpenOrders("BTC-USDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderInfo)
}
func TestClient_GetMarketPrice(t *testing.T) {
	price, err := c.GetMarketPrice("BTC-USDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(price)
}

func TestClient_Depth(t *testing.T) {
	depth, err := c.Depth("BTC-USDT", "20")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", depth)
}

func TestClient_GetPairInfo(t *testing.T) {
	pairInfo, err := c.GetPairInfo("BTC-USDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", pairInfo)
}

func TestClient_GetTradingFee(t *testing.T) {
	pairInfo, err := c.GetTradingFee("BTC-USDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", pairInfo)
}
func TestClient_GetFutureBalance(t *testing.T) {
	balance, err := c.NewFutureOrder("ETH-USDT", base.ASK, base.SHORT, base.STOP, "1.5", "2500", "3200", base.ISOLATED, false, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}

func TestClient_CancelAll(t *testing.T) {
	err := c.CancelFutureOrders("ETH-USDT")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(info)
}

func TestClient_GetAccountBalance(t *testing.T) {
	balance, err := c.GetAccountBalance("USDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}

func TestClient_GetDepositAddress(t *testing.T) {
	balance, err := c.GetDepositAddress("SUI", "SUI")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}

func TestClient_Withdraw(t *testing.T) {
	balance, err := c.Withdraw("USDT", "TRC20", "", "10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}
