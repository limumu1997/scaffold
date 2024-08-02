package router

import (
	"net/http"
	"scaffold/app/index/api"
	"scaffold/common/config"
	"scaffold/common/middleware"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func setupRoutes(r *mux.Router) {
	r.HandleFunc("/index", api.IndexHandler).Methods(http.MethodGet)
}

func ListenAndServe() {
	r := mux.NewRouter()

	// 应用日志中间件
	r.Use(middleware.LoggingMiddleware)

	// 应用 CORS 中间件
	r.Use(middleware.CorsMiddleware)

	// 设置路由
	setupRoutes(r)
	srv := &http.Server{
		Handler: r,
		Addr:    config.Config.ListenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	logrus.Infof("listen http://127.0.0.1%s", srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}
