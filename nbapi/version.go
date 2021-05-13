// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package nbapi

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
)

const currentVersion = "/v1"

type version struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func versions(w http.ResponseWriter, r *http.Request) {
	glog.Info("get versions")
	version := version{"Network Resource Manager API", currentVersion}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(version); err != nil {
		glog.Error(err.Error())
		panic(err)
	}
}
