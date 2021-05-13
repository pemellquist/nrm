// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package apiutils

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
)

// RESTError represents an error response from backend
type RESTError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// JSONWrite respond in JSON format
func JSONWrite(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		glog.Errorln(err.Error())
		return
	}
}

// JSONWriteError respond with error in JSON format
func JSONWriteError(w http.ResponseWriter, status int) {
	glog.Errorf("failure, response: %d", status)
	JSONWrite(w, RESTError{Code: status, Text: http.StatusText(status)}, status)
}

// JSONWriteErrorText respond with error in JSON format with text message
func JSONWriteErrorText(w http.ResponseWriter, status int, msg string) {
	glog.Errorf("failure, response: %d", status)
	JSONWrite(w, RESTError{Code: status, Text: msg}, status)
}
