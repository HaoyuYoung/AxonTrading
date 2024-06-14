package binance

import (
	"AxonTrading/base"
	"AxonTrading/models"
	"fmt"
	"sync"
	"testing"
)

var c Client

var (
	apiKey    string
	apiSecret string
)

func init() {
	apiKey = ""
	apiSecret = ""

	err := c.NewWithParams(apiKey, apiSecret)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestClient_Cancel(t *testing.T) {

	err1, err := c.GetFutureTradingFee("BTCUSDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err1, err)
}

func TestClient_GetAccountBalance(t *testing.T) {

	balance, err := c.GetAccountBalance("USDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}

func TestClient_Withdraw(t *testing.T) {

	balance, err := c.Withdraw("USDT", "TRX", "", "10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balance)
}

func TestClient_MarketOrder(t *testing.T) {
	id, err := c.MarketOrder("BUSDUSDT", base.ASK, "111")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}

func TestClient_LimitOrder(t *testing.T) {
	id, err := c.LimitOrder("BTCUSDT", base.ASK, "5", "10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}

func TestClient_TakerOrder(t *testing.T) {
	id, err := c.TakerOrder("BTCUSDT", base.BID, "274", "0.3")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}

func TestClient_CancelOrder(t *testing.T) {
	order, err := c.CancelOrder("BTCUSDT", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(order)
}

func TestClient_CancelOrders(t *testing.T) {
	err := c.CancelOrders("BTCUSDT")
	if err != nil {
		fmt.Println(err)
	}

}

func TestClient_GetOrder(t *testing.T) {
	info, err := c.GetOrder("BTCUSDT", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

func TestClient_GetOpenOrders(t *testing.T) {
	orders, err := c.GetOpenOrders("BTCUSDT")
	if err != nil {
		return
	}
	fmt.Println(orders)
}

func TestClient_GetOpenOrdersWithSide(t *testing.T) {
	infos, err := c.GetOpenOrdersWithSide("BNBUSDT", base.ASK)
	if err != nil {
		return
	}
	fmt.Println(infos)
}

func TestClient_GetOpenSplitOrders(t *testing.T) {
	info1, info2, err := c.GetOpenSplitOrders("BNBUSDT")
	if err != nil {
		return
	}
	fmt.Println(info1, info2)
}

func TestClient_GetMarketPrice(t *testing.T) {
	price, err := c.GetMarketPrice("BNBUSDT")
	if err != nil {
		return
	}
	fmt.Println(price)
}

func TestClient_Depth(t *testing.T) {
	depth, err := c.Depth("BNBUSDT", "10")
	if err != nil {
		return
	}
	fmt.Println(depth.Bids)
}

func TestClient_LimitOrders(t *testing.T) {
	var ol []models.OrderList
	ol = append(ol, models.OrderList{
		Side:  base.ASK,
		Price: "300",
		Size:  "0.1",
	}, models.OrderList{
		Side:  base.BID,
		Price: "100",
		Size:  "0.1",
	})
	orders, err := c.LimitOrders("BNBUSDT", ol)
	if err != nil {
		return
	}
	fmt.Println(orders)

}

func TestClient_MakerOrder(t *testing.T) {
	order, err := c.MakerOrder("BNBUSDT", base.ASK, "270", "0.1")
	if err != nil {
		return
	}
	fmt.Println(order)
}

func TestClient_GetPairInfo(t *testing.T) {
	info, err := c.GetPairInfo("PONDUSDT")
	if err != nil {
		return
	}
	fmt.Println(info)
}

func TestClient_GetTradingFee(t *testing.T) {
	fee, err := c.GetTradingFee("BNBUSDT")
	if err != nil {
		return
	}
	fmt.Println(fee)
}

func TestEDU(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			order, err := c.MarketOrder("EDUUSDT", base.BID, "200")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(order)
		}
	}()
	wg.Add(1)
	go func() {
		for {
			order, err := c.MarketOrder("EDUUSDT", base.BID, "50")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(order)
		}
	}()
	wg.Wait()

}
