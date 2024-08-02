package app

import (
	"scaffold/app/router"
	"scaffold/cmd"
)

func init() {
	cmd.InitLog()
	cmd.InitConfig()
}

func start() {
	router.ListenAndServe()
}
