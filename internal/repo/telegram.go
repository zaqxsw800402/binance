package repo

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/zaqxsw800402/binance/pkg/config"
)

type Telegram struct {
	token  string
	chatID int64
}

func NewTelegram(cfg config.Telegram) *Telegram {
	return &Telegram{
		token:  cfg.Token,
		chatID: cfg.ChatID,
	}
}

func (t Telegram) Send(msg string) error {
	data := url.Values{}
	data.Set("chat_id", fmt.Sprintf("-%d", t.chatID))
	data.Set("text", msg)
	telegramUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	_, err := http.PostForm(telegramUrl, data)
	if err != nil {
		return fmt.Errorf("http PostForm error: %w", err)
	}

	return nil
}
