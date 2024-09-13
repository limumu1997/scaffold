package router

import (
	"net/http"
	"scaffold/internal/config"
	"scaffold/internal/index/api"
	"scaffold/pkg/common/middleware"

	"time"

	"github.com/sirupsen/logrus"
)

func setupRoutes(r *http.ServeMux) {
	r.HandleFunc("/index", api.IndexHandler)
}

func ListenAndServe() {
	r := http.NewServeMux()

	// 设置路由
	setupRoutes(r)

	// 应用中间件
	mux := middleware.LoggingMiddleware(middleware.CorsMiddleware(r))

	srv := &http.Server{
		Handler: mux,
		Addr:    config.Config.ListenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	logrus.Infof("listen http://127.0.0.1%s", srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}
