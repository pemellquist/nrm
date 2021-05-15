// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package nbapi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hpe/nrm/apiutils"
	"github.com/hpe/nrm/config"
	"github.com/hpe/nrm/sbapi"
	"github.com/hpe/nrm/vpc"
)

const vpcMaxBodySize = 1024

var sb sbapi.SBAPI

// init SBAPI based on config setting
func init() {
	glog.Info("initializing vpc handlers")

	if !config.LoadConfig() {
		glog.Error("unable to load config, exiting")
		os.Exit(2)
	}

	glog.Info("nbapi mode: " + config.Config.SBAPI)
	if config.Config.SBAPI == "map" {
		glog.Info("using map mode sbapi")
		sb = sbapi.MapImpl{}
	} else {
		if config.Config.SBAPI == "psm" {
			glog.Info("using psm mode sbapi")
			glog.Info("psm base url: " + config.Config.PSMBaseURL)
			glog.Info("psm global tenantID: " + config.Config.PSMTenantID)
			sb = sbapi.PSMImpl{}
		} else {
			glog.Error("nbapi not supported: " + config.Config.SBAPI)
			os.Exit(3)
		}
	}

}

// get all VPCs
func getAllVPCs(w http.ResponseWriter, r *http.Request) {
	glog.Info("get vpcs")
	vpcs, err := sb.GetVPCs()

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vpcs); err != nil {
		glog.Error(err.Error())
		return
	}
}

// get a specific vpc
func getVPC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	glog.Info("get vpc id = ", id)

	exists, err := sb.ExistsVPC(id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err.Error())
		return
	}

	if !exists {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		glog.Info("not found id = ", id)
		return
	}

	vpc, err := sb.GetVPC(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vpc); err != nil {
		glog.Error(err.Error())
	}
}

// delete an existing VPC
func deleteVPC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	glog.Info("delete vpc id = ", id)
	exists, err := sb.ExistsVPC(id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err.Error())
		return
	}

	if !exists {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		glog.Info("not found id = ", id)
		return
	}

	err = sb.DeleteVPC(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// create a new VPC
func createVPC(w http.ResponseWriter, r *http.Request) {
	glog.Info("create vpc")

	// only allow POSTs to this URI
	if r.Method != http.MethodPost {
		glog.Warningf("method: %s not allowed to %s", r.Method, r.URL)
		apiutils.JSONWriteError(w, http.StatusMethodNotAllowed)
		return
	}

	// read request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, vpcMaxBodySize))
	if err != nil {
		glog.Errorf("unable to read request body: %v", err)
		apiutils.JSONWriteError(w, http.StatusInternalServerError)
		return
	}

	// unmarshall request into vpc struct
	var vpc vpc.Vpc
	if err := json.Unmarshal(body, &vpc); err != nil {
		glog.Errorf("unable to unmarshall request body: %v", err)
		apiutils.JSONWriteError(w, http.StatusUnprocessableEntity)
		return
	}

	// name required
	if vpc.Name == "" {
		glog.Errorf("VPC name must be specified")
		apiutils.JSONWriteErrorText(w, http.StatusUnprocessableEntity, "VPC name required")
		return
	}

	// tenant id required
	if vpc.TenantID == "" {
		glog.Errorf("VPC TenantID must be specified")
		apiutils.JSONWriteErrorText(w, http.StatusUnprocessableEntity, "VPC tenantId required")
		return
	}

	// look for global tenant ID override
	if config.Config.PSMTenantID != "" {
		glog.Info("over riding tenantId with global value: " + config.Config.PSMTenantID)
		vpc.TenantID = config.Config.PSMTenantID
	}

	// fill out creation fields
	vpc.ID = uuid.New().String()
	t := time.Now()
	vpc.CreatedDate = t.Format(time.RFC3339)
	vpc.LastUpdateDate = vpc.CreatedDate

	// log out orginal vpc request before calling handler
	buff, _ := json.Marshal(&vpc)
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, buff, "", "\t")
	if error != nil {
		glog.Error("JSON parse error: ", error)
		apiutils.JSONWriteError(w, http.StatusUnprocessableEntity)
		return
	}
	glog.Info(string(prettyJSON.Bytes()))

	// call interface handler to carry out creation
	err = sb.CreateVPC(vpc)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error(err.Error())
		return
	}

	// return back new VPC which was created
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vpc); err != nil {
		glog.Error(err.Error())
	}
}
