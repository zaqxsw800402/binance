package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zaqxsw800402/binance/internal/repo"
	"github.com/zaqxsw800402/binance/pkg/config"
)

// 定義你的測試套件
type TestBinanceService struct {
	suite.Suite
	service *BinanceService
}

// SetupSuite 在測試套件開始前執行
func (s *TestBinanceService) SetupSuite() {
	// 初始化邏輯，例如資料庫連線
	config.SetFolderPath("../../configs")
	cfg := config.InitYaml("binance")

	binanceHttp := repo.NewBinanceHttp(cfg.Binance)
	gpt := repo.NewChatGpt(cfg.OpenAI)
	tg := repo.NewTelegram(cfg.Telegram)
	s.service = NewBinanceService(binanceHttp, gpt, tg)
}

func (s *TestBinanceService) TestGetData() {
	ctx := context.Background()
	symbol := "ETHUSDT"
	data, err := s.service.GetData(ctx, symbol)
	s.Require().NoError(err)
	s.T().Logf("data: %s", data)

	// save to file
	//err = os.WriteFile("data.json", data, 0644)
	//s.Require().NoError(err)
}

// TestSuite 啟動測試套件
func TestTestBinanceService(t *testing.T) {
	suite.Run(t, new(TestBinanceService))
}
