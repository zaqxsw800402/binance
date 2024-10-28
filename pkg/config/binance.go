package config

type Binance struct {
	ApiKey    string `mapstructure:"API_KEY"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}
