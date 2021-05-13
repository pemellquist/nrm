package sbapi

import (
	"github.com/hpe/nrm/vpc"
)

// SBAPI Interface
type SBAPI interface {
	GetVPCs() ([]vpc.Vpc, error)
	GetVPC(id string) (vpc.Vpc, error)
	CreateVPC(vpc vpc.Vpc) error
	ExistsVPC(id string) (bool, error)
	DeleteVPC(id string) error
}
