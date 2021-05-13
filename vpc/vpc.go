// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package vpc

// Vpc is the defintion of a Virtual Priavte Cloud
type Vpc struct {
	TenantName     string `json:"tenantName"`     // name of owning tenant
	TenantID       string `json:"tenantId"`       // tenant id of owning tenant
	Name           string `json:"name"`           // name of VPC, specified by creator
	ID             string `json:"id"`             // id of VPC, created by service at create time
	CreatedDate    string `json:"createdDate"`    // date VPC was created
	LastUpdateDate string `json:"lastUpdateDate"` // date VPC last updated
	DefaultCIDR    string `json:"defaultCIDR"`    // default CIDR for this VPC when created
}
