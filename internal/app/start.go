package app

import (
	"fmt"
	"log/slog"
	"scaffold/internal/config"
	"scaffold/internal/router"
	"scaffold/pkg/common/util"
)

func init() {
	InitLog()
	InitConfig()
}

func Start() {
	slog.Info(fmt.Sprintf("Start %s version %s", config.Config.Service.Name, Version))
	if util.IsRunInDocker() {
		run()
	} else {
		startDaemon()
	}
}

func run() {
	router.ListenAndServe()
}
