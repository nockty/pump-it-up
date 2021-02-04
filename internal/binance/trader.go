package binance

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
)

const btc = "BTC"

type Trader struct {
	client        *binance.Client
	btcAmount     string
	symbolBuilder bytes.Buffer
}

func NewTrader(APIKey, secretKey string, BTCAmount string) *Trader {
	var symbolBuilder strings.Builder
	symbolBuilder.Grow(100)
	return &Trader{
		client:    binance.NewClient(APIKey, secretKey),
		btcAmount: BTCAmount,
	}
}

func (t *Trader) Trade(coin string) {
	symbolBuilder := t.symbolBuilder
	fmt.Fprintf(&symbolBuilder, coin)
	fmt.Fprintf(&symbolBuilder, btc)
	symbol := symbolBuilder.String()
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
}
