package binance

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
)

const btc = "BTC"

type Trader struct {
	client    *binance.Client
	btcAmount string
}

func NewTrader(APIKey, secretKey string, BTCAmount string) *Trader {
	return &Trader{
		client:    binance.NewClient(APIKey, secretKey),
		btcAmount: BTCAmount,
	}
}

func (t *Trader) Trade(coin string) {
	var b strings.Builder
	b.Grow(10)
	fmt.Fprintf(&b, coin)
	fmt.Fprintf(&b, btc)
	symbol := b.String()
	buyOrder, err := t.client.NewCreateOrderService().Symbol(symbol).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).
		QuoteOrderQty(t.btcAmount).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(1 * time.Second)
	sellOrder, err := t.client.NewCreateOrderService().Symbol(symbol).
		Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).
		Quantity(buyOrder.ExecutedQuantity).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Spent BTC: %s\n", buyOrder.CummulativeQuoteQuantity)
	fmt.Printf("Earnt BTC: %s\n", sellOrder.CummulativeQuoteQuantity)
	fmt.Printf("Qty      : %s\n", buyOrder.ExecutedQuantity)
	spentBTC, err := strconv.ParseFloat(buyOrder.CummulativeQuoteQuantity, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	qty, err := strconv.ParseFloat(buyOrder.ExecutedQuantity, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	price := spentBTC / qty
	fmt.Printf("Price    : %f\n", price)
}
