package app

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/zaqxsw800402/binance/internal/controller"
	"github.com/zaqxsw800402/binance/internal/repo"
	"github.com/zaqxsw800402/binance/internal/service"
	"github.com/zaqxsw800402/binance/pkg/config"
	"github.com/zaqxsw800402/binance/pkg/logger"
)

type Application struct {
	fileName string
}

func NewApplication(fileName string) *Application {
	return &Application{
		fileName: fileName,
	}
}

func (a *Application) Run() {
	logger.Init()
	err := a.init(a.fileName)
	if err != nil {
		err = fmt.Errorf("init %w", err)
		panic(err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (a *Application) init(fileName string) error {
	cfg := config.InitYaml(fileName)

	gpt := repo.NewChatGpt(cfg.OpenAI)
	binanceHttp := repo.NewBinanceHttp(cfg.Binance)
	telegram := repo.NewTelegram(cfg.Telegram)

	binanceService := service.NewBinanceService(binanceHttp, gpt, telegram)
	telegramController, err := controller.NewTelegramController(cfg.Telegram, binanceService)
	if err != nil {
		return fmt.Errorf("NewTelegramController err: %w", err)
	}

	telegramController.Start()

	return nil
}
