package app

import "scaffold/internal/router"

func init() {
	InitLog()
	InitConfig()
}

func start() {
	router.ListenAndServe()
}
