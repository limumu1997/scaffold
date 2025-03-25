package app

import (
	"scaffold/internal/config"
	"scaffold/pkg/logger"
)

func InitLog() {
	// 初始化日志
	logger.InitMyLog()
}

func InitConfig() {
	config.InitConfig()
}
