package main

import (
	"github.com/BurntSushi/toml"
)

type sflow_config struct {
	Address string
	Port    int
}

type app_config struct {
	SFlowConfig sflow_config
}

func ReadConfig(configFile string, AppConfig *app_config) bool {
	if _, err := toml.DecodeFile(configFile, AppConfig); err != nil {
		ErrorLogger.Println("Unable to read config file!")
		ErrorLogger.Println(err)
		return false
	}
	return true
}
