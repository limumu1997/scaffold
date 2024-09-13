package app

import (
	"scaffold/internal/config"
	"scaffold/pkg/logger"
)

func InitLog() {
	logger.InitMyLog()
}

func InitConfig() {
	config.InitConfig()
}
