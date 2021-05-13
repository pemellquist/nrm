// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/hpe/nrm/apiutils"
	"github.com/hpe/nrm/psm"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var psmroutes = routes{
	route{
		"create psm vpc",
		"POST",
		"/psm/configs/network/v1/tenant/{tid}/virtualrouters",
		createVPC,
	},
}

func createVPC(w http.ResponseWriter, r *http.Request) {
	glog.Info("create PSM vpc ..")

	vars := mux.Vars(r)
	tid := vars["tid"]
	glog.Info("psm tenant id = ", tid)

	if r.Method != http.MethodPost {
		glog.Warningf("method: %s not allowed to %s", r.Method, r.URL)
		apiutils.JSONWriteError(w, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		glog.Errorf("unable to read request body: %v", err)
		apiutils.JSONWriteError(w, http.StatusInternalServerError)
		return
	}

	var psmvpc psm.PSMvpc
	if err := json.Unmarshal(body, &psmvpc); err != nil {
		glog.Errorf("unable to unmarshall request body: %v", err)
		apiutils.JSONWriteError(w, http.StatusUnprocessableEntity)
		return
	}

	// print it and turn it around
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		glog.Error("JSON parse error: ", error)
		apiutils.JSONWriteError(w, http.StatusUnprocessableEntity)
		return
	}
	glog.Info(string(prettyJSON.Bytes()))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(psmvpc); err != nil {
		glog.Error(err.Error())
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: mockpsm -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
	flag.Parse()
}

func mainloop() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func server() {
	listenPort := fmt.Sprintf(":%d", 8089)
	glog.Info("starting api server on port: ", listenPort)
	router := newRouter()
	glog.Fatal(http.ListenAndServeTLS(listenPort, "https-server.crt", "https-server.key", router))
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range psmroutes {
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

func main() {
	glog.Info("starting mock PSM ... ")
	go server()
	mainloop()
}
