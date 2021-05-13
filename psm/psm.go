// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package psm

// PMeta is meta data
type PMeta struct {
	Name   string `json:"name"`
	Tenant string `json:"tenant"`
}

// RTExport route export
type RTExport struct {
	Type          string `json:"type"`
	AdminValue    string `json:"admin-value"`
	AssignedValue string `json:"assigned-value"`
}

// RTImport route import
type RTImport struct {
	Type          string `json:"type"`
	AdminValue    string `json:"admin-value"`
	AssignedValue string `json:"assigned-value"`
}

// RouteIE route import exports
type RouteIE struct {
	AddrFamily string     `json:"address-family"`
	RdAuto     bool       `json:"rd-auto"`
	RTExports  []RTExport `json:"rt-export"` // array
	RTImports  []RTImport `json:"rt-import"` // array
}

// PSpec spec block
type PSpec struct {
	Type      string  `json:"type"`
	RouterMac string  `json:"router-mac-address"`
	VXLANID   string  `json:"vxlan-vni"`
	RouteIE   RouteIE `json:"route-import-export"`
}

// PSMvpc main PSM VPC
type PSMvpc struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"api-version"`
	Meta       PMeta  `json:"meta"`
	Spec       PSpec  `json:"spec"`
}
