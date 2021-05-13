// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package nbapi

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/hpe/nrm/config"
)

// Server starts the API server to handle REST requests
func Server() {
	listenPort := fmt.Sprintf(":%d", config.Config.APIServerPort)
	glog.Info("starting api server on port", listenPort)
	router := newRouter()
	glog.Fatal(http.ListenAndServe(listenPort, router))

}

func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range nrmroutes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
