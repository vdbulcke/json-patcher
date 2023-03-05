package config

import (
	"testing"

	"github.com/vdbulcke/json-patcher/src/logger"
	"go.uber.org/zap"
)

func TestConfig(t *testing.T) {
	configFile := "../../example/patch.yaml"

	cfg, err := ParseConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	l := logger.GetZapLogger(true)
	l.Info("parsed config",zap.Any("config", cfg))

}
