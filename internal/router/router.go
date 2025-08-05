package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"scaffold/internal/config"
	"scaffold/internal/index/api"
	"scaffold/pkg/common/middleware"

	"time"
)

func setupRoutes(r *http.ServeMux) {
	r.HandleFunc("/index", api.IndexHandler)
}

func ListenAndServe() {
	if config.GetConfig().ListenPort == "" {
		return
	}

	r := http.NewServeMux()

	// 设置路由
	setupRoutes(r)

	// 应用中间件
	mux := middleware.LoggingMiddleware(middleware.CorsMiddleware(r))

	srv := &http.Server{
		Handler: mux,
		Addr:    config.GetConfig().ListenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	slog.Info(fmt.Sprintf("Server is listening on http://127.0.0.1%s", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
