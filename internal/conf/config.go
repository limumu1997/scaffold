package conf

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
	executable, _ := os.Executable()
	res, _ := filepath.EvalSymlinks(filepath.Dir(executable))
	absPath := filepath.Join(res, "config.json")
	dataConfig, err := os.ReadFile(absPath)
	if err != nil {
		slog.Info("prod env can not find config.json file, use embed config")
		dataConfig, _ = c.ReadFile("config.json")
	}
	json.Unmarshal(dataConfig, &Config)
}
