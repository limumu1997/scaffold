package config

import (
	"embed"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

//go:embed config.json
var embeddedConfig embed.FS

var (
	cfg      *Config
	initOnce sync.Once
)

// 导出的配置结构
type Config struct {
	Service    ServiceConfig `json:"service"`
	DataPath   string        `json:"data_path"`
	ListenPort string        `json:"listen_port"`
}

type ServiceConfig struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// InitConfig 初始化配置，只会执行一次
func InitConfig() error {
	var err error
	initOnce.Do(func() {
		cfg, err = loadConfig()
	})
	return err
}

// loadConfig 加载配置文件
func loadConfig() (*Config, error) {
	const configFilename = "config.json"

	// 查找运行目录下是否有配置文件
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	baseDir, err := filepath.EvalSymlinks(filepath.Dir(execPath))
	if err != nil {
		return nil, err
	}
	fullPath := filepath.Join(baseDir, configFilename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		slog.Warn("Cannot find config.json on disk, fallback to embedded config")
		data, err = embeddedConfig.ReadFile(configFilename)
		if err != nil {
			return nil, errors.New("cannot read config.json from embedded FS or disk")
		}
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// GetConfig 获取全局配置单例
func GetConfig() *Config {
	if cfg == nil {
		panic("config not initialized. call InitConfig() first")
	}
	return cfg
}
