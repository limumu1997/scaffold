package app

import (
	"scaffold/internal/config"
	"scaffold/internal/router"
	"scaffold/pkg/common/util"

	"github.com/sirupsen/logrus"
)

func init() {
	InitLog()
	InitConfig()
}

func Start() {
	logrus.Infof("start %s version %s", config.Config.Service.Name, Version)
	if util.IsRunInDocker() {
		run()
	} else {
		startDaemon()
	}
}

func run() {
	router.ListenAndServe()
}
