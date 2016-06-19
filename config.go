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

func ReadConfig(configFile string) (app_config, error) {
	var AppConfig app_config
	if _, err := toml.DecodeFile(configFile, &AppConfig); err != nil {
		ErrorLogger.Println("Unable to read config file!")
		ErrorLogger.Println(err)
		return AppConfig, err
	}
	return AppConfig, nil
}
