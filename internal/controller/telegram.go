package controller

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zaqxsw800402/binance/internal/service"
	"github.com/zaqxsw800402/binance/pkg/config"
)

type TelegramController struct {
	chatID  int64
	bot     *tgbotapi.BotAPI
	binance *service.BinanceService
}

func NewTelegramController(cfg config.Telegram, service *service.BinanceService) (*TelegramController, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("NewTelegramController err: %w", err)
	}

	bot.Debug = true

	return &TelegramController{
		bot:     bot,
		chatID:  cfg.ChatID,
		binance: service,
	}, nil
}

func (t *TelegramController) Start() {
	go t.ReceiveData()
}

func (t *TelegramController) Stop() {
	// 停止接收消息
}

func (t *TelegramController) ReceiveData() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bot := t.bot
	targetChatID := t.chatID // 獲取機器人的用戶名
	botUsername := bot.Self.UserName

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Chat.IsGroup() && update.Message.Chat.ID == -1*targetChatID {
			ctx := context.Background()
			// 取得用戶的消息
			userMessage := update.Message.Text

			// 檢查消息中是否包含 @機器人名稱
			if !strings.Contains(userMessage, "@"+botUsername) {
				continue
			}

			splitN := strings.SplitN(userMessage, " ", 2)
			if len(splitN) < 2 {
				continue
			}

			symbol := strings.TrimSpace(strings.ToUpper(splitN[1]))

			// 處理並回應包含 @提到的消息
			// replyMessage := fmt.Sprintf("Hello, %s! You mentioned me.", update.Message.From.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, t.PrepareMsg(ctx, symbol))

			// 發送回應訊息
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}
	}
}

func (t *TelegramController) PrepareMsg(ctx context.Context, symbol string) string {
	// 檢查交易對是否存在
	exist, err := t.binance.CheckExist(ctx, symbol)
	if err != nil {
		return fmt.Sprintf("CheckExist err: %v", err)
	}

	if !exist {
		return fmt.Sprintf("symbol %s not exist", symbol)
	}

	// 取得交易對的數據
	data, err := t.binance.GetData(ctx, symbol)
	if err != nil {
		return fmt.Sprintf("GetData err: %v", err)
	}

	return data
}
