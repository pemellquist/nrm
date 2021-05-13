// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package sbapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/hpe/nrm/config"
	"github.com/hpe/nrm/psm"
	"github.com/hpe/nrm/vpc"
)

// PSMImpl implements SBAPI into Pensando PSM interface
type PSMImpl struct {
}

// GetVPCs return list of created VPCs
func (psmi PSMImpl) GetVPCs() ([]vpc.Vpc, error) {
	vpcs := make([]vpc.Vpc, len(m), len(m))

	return vpcs, nil
}

// CreateVPC creates VPC
func (psmi PSMImpl) CreateVPC(vpc vpc.Vpc) error {

	glog.Info("create vpc using psm")

	var psmVpc psm.PSMvpc
	psmVpc.APIVersion = "v1"
	psmVpc.Kind = "VirtualRouter"

	psmVpc.Meta.Name = vpc.Name
	psmVpc.Meta.Tenant = vpc.TenantID

	psmVpc.Spec.Type = "tenant"
	psmVpc.Spec.VXLANID = "foo"              // TODO: may need to move this to NRM VPC
	psmVpc.Spec.RouterMac = "don't know yet" // TODO ?

	psmVpc.Spec.RouteIE.AddrFamily = "l2vpn-evpn" // TODO ?
	psmVpc.Spec.RouteIE.RdAuto = true

	var rtExport psm.RTExport
	rtExport.AdminValue = "??"
	rtExport.AssignedValue = "??"
	rtExport.Type = "type0"
	psmVpc.Spec.RouteIE.RTExports = make([]psm.RTExport, 1, 1)
	psmVpc.Spec.RouteIE.RTExports[0] = rtExport

	var rtImport psm.RTImport
	rtImport.AdminValue = "??"
	rtImport.AssignedValue = "??"
	rtImport.Type = "type0"
	psmVpc.Spec.RouteIE.RTImports = make([]psm.RTImport, 1, 1)
	psmVpc.Spec.RouteIE.RTImports[0] = rtImport

	// log out request vpc
	buff, _ := json.Marshal(&psmVpc)
	var prettyJSONReq bytes.Buffer
	error := json.Indent(&prettyJSONReq, buff, "", "\t")
	if error != nil {
		glog.Error("JSON parse error: ", error)
		return error
	}
	glog.Info("request vpc to psm :")
	glog.Info(string(prettyJSONReq.Bytes()))

	// make request to PSM
	request, err := json.Marshal(psmVpc)
	if err != nil {
		glog.Error(err)
		return err
	}

	// post request to PSM
	url := "https://" + config.Config.PSMBaseURL + "/psm/configs/network/v1/tenant/" +
		vpc.TenantID +
		"/virtualrouters"
	glog.Info("psm vpc create url: " + url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // TODO: ignore certs for now
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(request))
	if err != nil {
		glog.Error(err)
		return err
	}

	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// unmarshall request into psm vpc struct
	var psmVpcResponse psm.PSMvpc
	if err := json.Unmarshal(body, &psmVpcResponse); err != nil {
		glog.Errorf("unable to unmarshall body: %v", err)
		return err
	}

	// log out psm response
	respbuff, _ := json.Marshal(&psmVpc)
	var prettyJSONResp bytes.Buffer
	error = json.Indent(&prettyJSONResp, respbuff, "", "\t")
	if error != nil {
		glog.Error("JSON parse error: ", error)
		return err
	}
	glog.Info("response vpc from psm :")
	glog.Info(string(prettyJSONResp.Bytes()))

	return nil
}

// ExistsVPC detects if VPC already exists by ID
func (psmi PSMImpl) ExistsVPC(id string) (bool, error) {

	return false, nil
}

// GetVPC returns VPC based on ID
func (psmi PSMImpl) GetVPC(id string) (vpc.Vpc, error) {
	var v vpc.Vpc
	return v, nil
}

// DeleteVPC removes existing VPC based on ID
func (psmi PSMImpl) DeleteVPC(id string) error {

	return nil
}
