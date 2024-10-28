package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Binance  Binance  `mapstructure:"binance"`
	OpenAI   OpenAI   `mapstructure:"openai"`
	Telegram Telegram `mapstructure:"telegram"`
}

var folderPath = "configs"

func SetFolderPath(path string) {
	folderPath = path
}

func InitYaml(filename string) *Config {
	slog.Info("[Conf@Init]Conf Init start")
	// 设置配置文件的搜索路径和名称
	viper.AddConfigPath(folderPath) // 指向项目根目录下的config文件夹
	viper.SetConfigName(filename)   // 配置文件名称(不含后缀)

	// 设置配置文件类型和环境变量
	viper.SetConfigType("yaml") // 根据实际配置文件修改

	var conf Config
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("读取配置文件失败", "error", err, "path", folderPath, "filename", filename)
		os.Exit(1)
	}

	// 将配置文件的内容映射到Server变量
	if err := viper.Unmarshal(&conf); err != nil {
		slog.Error("映射配置文件失败", "error", err)
		os.Exit(1)
	}

	return &conf
}
