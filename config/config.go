package config

import (
	"embed"
	"encoding/json"
	"os"
)

var (
	//go:embed config.json
	c      embed.FS
	Config config
)

type serviceConfig struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type config struct {
	Service    serviceConfig `json:"service"`
	ListenPort string        `json:"listen_port"`
}

func InitConfig() {
	conf := "config.json"
	dataConfig, err := os.ReadFile(conf)
	if err != nil {
		dataConfig, _ = c.ReadFile("config.json")
	}
	json.Unmarshal(dataConfig, &Config)
}
