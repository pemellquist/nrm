// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package nbapi

import (
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var nrmroutes = routes{
	route{
		"versions",
		"GET",
		"/",
		versions,
	},
	route{
		"get all vpcs",
		"GET",
		currentVersion + "/vpc",
		getAllVPCs,
	},
	route{
		"create vpc",
		"POST",
		currentVersion + "/vpc",
		createVPC,
	},
	route{
		"get vpc",
		"GET",
		currentVersion + "/vpc/{id}",
		getVPC,
	},
	route{
		"delete vpc",
		"DELETE",
		currentVersion + "/vpc/{id}",
		deleteVPC,
	},
}
