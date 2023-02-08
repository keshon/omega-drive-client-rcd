package state

import (
	"app/src/conf"
	"encoding/base64"
)

/*
	State package contains structs and vars that temporary store information during app run
*/

// Settings values
type SettingsValuesStruct struct {
	General struct {
		StoreLocationPortable bool   `json:"store_location_portable"`
		AppdataPath           string `json:"app_data_path"`
	}
	Cache struct {
		DefaultPath  string `json:"default_path"`
		OverridePath string `json:"override_path"`
		Disabled     bool   `json:"disabled"`
	}
	Remote struct {
		ReconnectRate  string `json:"reconnect_rate"`
		RandServerPort bool   `json:"rand_server_port"`
	}
}

var (
	SettingsValues SettingsValuesStruct // settings

	// Encoded `username:password` pair for basic auths
	RcAuthEncoded = base64.StdEncoding.EncodeToString([]byte(conf.RcUsername + ":" + conf.RcPassword)) // rcd
)
