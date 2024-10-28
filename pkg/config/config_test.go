package config

import (
	"testing"
)

func TestInitYaml(t *testing.T) {
	SetFolderPath("../../configs")
	cfg := InitYaml("binance")
	t.Logf("%+v", cfg)
}
