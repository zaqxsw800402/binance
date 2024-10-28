package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zaqxsw800402/binance/pkg/config"
)

func TestBinanceHttp_GetOrderBook(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	orderBook, err := binance.GetOrderBook(ctx, symbol)
	require.NoError(t, err)
	t.Logf("orderBook: %s", orderBook)
}

func TestBinanceHttp_GetKLines(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	timeInterval := "5m"
	klines, err := binance.GetKLines(ctx, symbol, timeInterval)
	require.NoError(t, err)
	t.Logf("klines: %s", klines)
}

func TestBinanceHttp_GetTrades(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")
	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	trades, err := binance.GetTrades(ctx, symbol)
	require.NoError(t, err)
	t.Logf("trades: %s", trades)
}

func TestBinanceHttp_GetTicker(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	price, err := binance.GetPrice24Hr(ctx, symbol)
	require.NoError(t, err)
	t.Logf("price: %s", price)
}

func TestBinanceHttp_GetPremiumIndex(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT123"
	price, err := binance.GetPremiumIndex(ctx, symbol)
	require.NoError(t, err)
	t.Logf("price: %s", price)
}

func TestBinanceHttp_GetAggTrades(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	price, err := binance.GetAggTrades(ctx, symbol)
	require.NoError(t, err)
	t.Logf("price: %s", price)
}

func TestBinanceHttp_GetOpenInterestHist(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	price, err := binance.GetOpenInterestHist(ctx, symbol)
	require.NoError(t, err)
	t.Logf("price: %s", price)
}

func TestBinanceHttp_GetTopLongShortPosition(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	price, err := binance.GetTopLongShortPosition(ctx, symbol)
	require.NoError(t, err)
	t.Logf("price: %s", price)
}

func TestBinanceHttp_CheckSymbol(t *testing.T) {
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binance := NewBinanceHttp(cfg.Binance)

	ctx := context.Background()
	symbol := "BTCUSDT"
	exist, err := binance.CheckSymbol(ctx, symbol)
	require.NoError(t, err)
	t.Logf("exist: %v", exist)
}
