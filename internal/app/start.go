package app

import (
	"fmt"
	"log/slog"
	"scaffold/internal/config"
	"scaffold/internal/router"
	"scaffold/pkg/common"
	"scaffold/pkg/common/util"
)

func init() {
	InitLog()
	InitConfig()
}

func Start() {
	slog.Info(fmt.Sprintf("Start %s version %s", config.GetConfig().Service.Name, common.Version))
	if util.IsRunInDocker() {
		run()
	} else {
		startDaemon()
	}
}

func run() {
	router.ListenAndServe()
}
