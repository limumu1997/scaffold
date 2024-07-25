package router

import (
	"net/http"
	conf "scaffold/config"
	"scaffold/controller"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func ListenAndServe() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/index", controller.IndexHandler).Methods(http.MethodGet)
	srv := &http.Server{
		Handler: r,
		Addr:    conf.Config.ListenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	logrus.Infof("listen http://127.0.0.1%s", srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}
