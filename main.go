// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hpe/nrm/config"
	"github.com/hpe/nrm/nbapi"

	"github.com/golang/glog"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: nrm -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
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

func main() {
	glog.Info("starting network resource manager ... ")

	if !config.LoadConfig() {
		glog.Error("unable to load config, exiting")
		return
	}

	go nbapi.Server()

	mainloop()

}
