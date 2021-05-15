// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package config

import (
	"io/ioutil"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

// NRMConfig is a representation of a config data for this service which is loaded
type NRMConfig struct {
	APIServerPort int    `yaml:"api_server_port"`
	SBAPI         string `yaml:"sbapi"`
	PSMBaseURL    string `yaml:"psm_base_url"`
	PSMTenantID   string `yaml:"psm_tenant_id"`
	PSMUserName   string `yaml:"psm_username"`
	PSMUserPwd    string `yaml:"psm_pwd"`
}

var configfile = "config.yaml"

// Config is the currently loaded configuration info
var Config NRMConfig

// LoadConfig loads config data from the config file
func LoadConfig() bool {
	glog.Info("loading NRM config from: ", configfile)
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		glog.Error("could not open: ", configfile)
		return false
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		glog.Error("could not parse config: ", configfile)
		return false
	}

	marshalled, err := yaml.Marshal(&Config)
	if err != nil {
		glog.Warning("could not marshall config for logging")
	} else {
		glog.Info("\n----- NRM config ------\n", string(marshalled), "------------------")
	}

	return true
}
