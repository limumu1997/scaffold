package cmd

import (
	"scaffold/internal/bootstrap"
	"scaffold/internal/conf"
)

func init() {
	bootstrap.InitLog()
	conf.InitConfig()
}
