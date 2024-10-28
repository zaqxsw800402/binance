package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/zaqxsw800402/binance/internal/model"
	"github.com/zaqxsw800402/binance/internal/repo"
	"golang.org/x/sync/errgroup"
)

type BinanceService struct {
	mu        sync.Mutex
	symbolMap map[string]bool

	binanceHttp *repo.BinanceHttp
	openai      *repo.ChatGpt
	tg          *repo.Telegram
}

func NewBinanceService(
	binanceHttp *repo.BinanceHttp,
	openai *repo.ChatGpt,
	tg *repo.Telegram,
) *BinanceService {
	return &BinanceService{
		symbolMap:   map[string]bool{},
		binanceHttp: binanceHttp,
		openai:      openai,
		tg:          tg,
	}
}

func (b *BinanceService) checkSymbol(symbol string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.symbolMap[symbol]
}

func (b *BinanceService) addSymbol(symbol string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.symbolMap[symbol] = true
}

func (b *BinanceService) CheckExist(ctx context.Context, symbol string) (bool, error) {
	if b.checkSymbol(symbol) {
		return true, nil
	}

	exist, err := b.binanceHttp.CheckSymbol(ctx, symbol)
	if err != nil {
		return false, fmt.Errorf("binanceHttp.CheckSymbol err: %w", err)
	}

	if !exist {
		return false, fmt.Errorf("symbol %s not exist", symbol)
	}

	b.addSymbol(symbol)

	return true, nil
}

func (b *BinanceService) GetData(ctx context.Context, symbol string) (string, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	data := model.NewAllPrompt()

	eg.Go(func() error {
		klines, err := b.binanceHttp.GetKLines(egCtx, symbol, "15m")
		if err != nil {
			return fmt.Errorf("binanceHttp.GetKLines err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("15分鐘的k line", klines))
		return nil
	})

	eg.Go(func() error {
		klines, err := b.binanceHttp.GetKLines(egCtx, symbol, "4h")
		if err != nil {
			return fmt.Errorf("binanceHttp.GetKLines err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("4小時的k line", klines))
		return nil
	})

	eg.Go(func() error {
		klines, err := b.binanceHttp.GetKLines(egCtx, symbol, "1w")
		if err != nil {
			return fmt.Errorf("binanceHttp.GetKLines err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("1星期的k line", klines))
		return nil
	})

	// eg.Go(func() error {
	//	trades, err := b.binanceHttp.GetTrades(egCtx, symbol)
	//	if err != nil {
	//		return fmt.Errorf("binanceHttp.GetTrades err: %w", err)
	//	}
	//
	//	data.Trades = trades
	//data.AddPrompt(model.NewPrompt("近期交易單", trades))
	//return nil
	//})

	eg.Go(func() error {
		price24Hr, err := b.binanceHttp.GetPrice24Hr(egCtx, symbol)
		if err != nil {
			return fmt.Errorf("binanceHttp.GetPrice24Hr err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("24小時的價格變動", price24Hr))
		return nil
	})

	eg.Go(func() error {
		premiumIndex, err := b.binanceHttp.GetPremiumIndex(egCtx, symbol)
		if err != nil {
			return fmt.Errorf("binanceHttp.GetPremiumIndex err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("資金費率", premiumIndex))
		return nil
	})

	eg.Go(func() error {
		aggTrades, err := b.binanceHttp.GetAggTrades(egCtx, symbol)
		if err != nil {
			return fmt.Errorf("binanceHttp.GetAggTrades err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("最近的交易數據", aggTrades))
		return nil
	})

	eg.Go(func() error {
		openInterestHist, err := b.binanceHttp.GetOpenInterestHist(egCtx, symbol)
		if err != nil {
			return fmt.Errorf("binanceHttp.GetOpenInterestHist err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("未平倉合約量統計數據", openInterestHist))
		return nil
	})

	eg.Go(func() error {
		interest, err := b.binanceHttp.GetOpenInterest(egCtx, symbol)
		if err != nil {
			return fmt.Errorf("binanceHttp.GetOpenInterest err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("即時未平倉合約量", interest))
		return nil
	})

	eg.Go(func() error {
		position, err := b.binanceHttp.GetTopLongShortPosition(egCtx, symbol)
		if err != nil {
			return fmt.Errorf("binanceHttp.GetTopLongShortPosition err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("多空雙方的持倉比例", position))
		return nil
	})

	eg.Go(func() error {
		account, err := b.binanceHttp.GetTopLongShortAccount(egCtx, symbol)
		if err != nil {
			return fmt.Errorf("binanceHttp.GetTopLongShortAccount err: %w", err)
		}

		data.AddPrompt(model.NewPrompt("多空雙方賬戶的數量比例", account))
		return nil
	})

	if err := eg.Wait(); err != nil {
		return "", err
	}

	content, err := b.openai.Chat(ctx, data)
	if err != nil {
		return "", fmt.Errorf("openai.Chat err: %w", err)
	}

	return content, nil
}
