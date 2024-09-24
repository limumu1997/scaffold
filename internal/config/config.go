package config

import (
	"embed"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
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
	cnf := "config.json"
	executable, _ := os.Executable()
	res, _ := filepath.EvalSymlinks(filepath.Dir(executable))
	absPath := filepath.Join(res, cnf)
	dataConfig, err := os.ReadFile(absPath)
	if err != nil {
		slog.Info("prod env can not find config.json file, use embed config")
		dataConfig, _ = c.ReadFile(cnf)
	}
	json.Unmarshal(dataConfig, &Config)
}
