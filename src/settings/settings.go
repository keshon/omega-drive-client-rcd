package settings

import (
	"app/src/conf"
	"app/src/state"
	"app/src/utils"
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadSettingsValues() {
	// Find config file location and override the default settings
	appdataConf := getAppdataDir() + "conf.json"
	state.SettingsValues.General.AppdataPath = createAppdataDir()
	state.SettingsValues.Cache.DefaultPath = createAppdataDir() + "cache"

	localConfExist, err := utils.PathExist("conf.json")
	if err != nil {
		panic("err reading conf")
	}

	if localConfExist {
		confData := utils.ReadFile("conf.json")
		json.Unmarshal(confData, &state.SettingsValues)
	} else {
		appdataConfExists, err := utils.PathExist(appdataConf)
		if err != nil {
			panic("err reading appdata conf")
		}

		if appdataConfExists {
			confData := utils.ReadFile(appdataConf)
			json.Unmarshal(confData, &state.SettingsValues)
		} else {
			// err
		}
	}
}

func getAppdataDir() string {
	appdataPath, err := os.UserConfigDir()
	if err != nil {
		panic("appdata path is invalid")
	}

	path := filepath.Join(appdataPath+"/", conf.AppName)

	return path + "/"
}

func createAppdataDir() string {
	appdataPath, err := os.UserConfigDir()
	if err != nil {
		panic("appdata path is invalid")
	}

	path := filepath.Join(appdataPath+"/", conf.AppName)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic("err creating new appdata path")
	}

	return path + "/"
}
