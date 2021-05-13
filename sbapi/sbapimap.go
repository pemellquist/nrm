// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package sbapi

import (
	"github.com/golang/glog"
	"github.com/hpe/nrm/vpc"
)

// MapImpl implements Persist interface
type MapImpl struct {
}

var m = make(map[string]vpc.Vpc)

// GetVPCs return list of created VPCs
func (mi MapImpl) GetVPCs() ([]vpc.Vpc, error) {
	vpcs := make([]vpc.Vpc, len(m), len(m))
	idx := 0
	for _, value := range m {
		glog.Info("add VPC")
		vpcs[idx] = value
		idx++
	}
	return vpcs, nil
}

// CreateVPC persists VPC
func (mi MapImpl) CreateVPC(vpc vpc.Vpc) error {
	glog.Info("create VPC")
	m[vpc.ID] = vpc
	return nil
}

// ExistsVPC detects if VPC already exists by ID
func (mi MapImpl) ExistsVPC(id string) (bool, error) {
	if _, ok := m[id]; ok {
		return true, nil
	}
	return false, nil
}

// GetVPC returns VPC based on ID
func (mi MapImpl) GetVPC(id string) (vpc.Vpc, error) {
	return m[id], nil
}

// DeleteVPC removes existing VPC based on ID
func (mi MapImpl) DeleteVPC(id string) error {
	delete(m, id)
	return nil
}
