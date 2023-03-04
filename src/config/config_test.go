package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	configFile := "../../example/patch.yaml"

	_, err := ParseConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

}
