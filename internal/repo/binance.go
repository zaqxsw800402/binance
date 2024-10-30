package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/zaqxsw800402/binance/pkg/config"
)

const (
	limit  = 50
	period = "1h"
)

type IBinance interface {
}

type BinanceHttp struct {
	client *futures.Client
}

func NewBinanceHttp(cfg config.Binance) *BinanceHttp {
	futuresClient := binance.NewFuturesClient(cfg.ApiKey, cfg.SecretKey)
	return &BinanceHttp{
		client: futuresClient,
	}
}

func (b *BinanceHttp) GetKLines(ctx context.Context, symbol, timeInterval string) ([]byte, error) {
	res, err := b.client.NewKlinesService().Symbol(symbol).Interval(timeInterval).Limit(limit).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewKlinesService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetPrice24Hr(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewListPriceChangeStatsService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewListPriceChangeStatsService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetPremiumIndex(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewPremiumIndexService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewPremiumIndexService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetAggTrades(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewAggTradesService().Symbol(symbol).Limit(6 * limit).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewAggTradesServiceerr: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetOpenInterestHist(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewOpenInterestStatisticsService().Symbol(symbol).Period(period).Limit(limit).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewOpenInterestStatisticsService: %w", err)
	}

	bytes, err := json.Marshal(res[0])
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetOpenInterest(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewGetOpenInterestService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewGetOpenInterestService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetTopLongShortPosition(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewTopLongShortPositionRatioService().Symbol(symbol).Period(period).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewTopLongShortPositionRatioService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetTopLongShortAccount(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewTopLongShortAccountRatioService().Symbol(symbol).Period(period).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewTopLongShortAccountRatioService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) CheckSymbol(ctx context.Context, symbol string) (bool, error) {
	_, err := b.client.NewPremiumIndexService().Symbol(symbol).Do(ctx)
	if err != nil {
		var apiErr *common.APIError
		if errors.As(err, &apiErr) {
			if apiErr.Code == -1121 {
				return false, nil
			}
		}
		return false, fmt.Errorf("CheckSymbol err: %w", err)
	}

	return true, nil
}

func (b *BinanceHttp) GetTrades(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewHistoricalTradesService().Symbol(symbol).Limit(limit).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewHistoricalTradesService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}

func (b *BinanceHttp) GetOrderBook(ctx context.Context, symbol string) ([]byte, error) {
	res, err := b.client.NewDepthService().Symbol(symbol).Limit(limit).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewDepthService err: %w", err)
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	return bytes, nil
}
