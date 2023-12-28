package app

import (
	"net/http"
	"scaffold/internal/conf"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func ListenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/index", IndexHandler).Methods(http.MethodPost)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    conf.Config.ListenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	logrus.Fatal(srv.ListenAndServe())
}
