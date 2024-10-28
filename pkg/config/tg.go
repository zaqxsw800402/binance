package config

type Telegram struct {
	Token  string `mapstructure:"TOKEN"`
	ChatID int64  `mapstructure:"CHAT_ID"`
}
