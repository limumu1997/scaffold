package app

import (
	"scaffold/cmd"
	"scaffold/router"
)

func init() {
	cmd.InitLog()
	cmd.InitConfig()
}

func start() {
	router.ListenAndServe()
}
