package config

import (
	"encoding/json"
	"os"
)

// Config contains client config parameters set in the config file
type Config struct {
	StoreFile  string `json:"store_file"`
	ServerKey  string `json:"server_key"`
	ServerCRT  string `json:"server_crt"`
	ListenPort int    `json:"listen_port"`
}

// Cfg holds global parameters from config file
var Cfg Config

// ParseConfigFile parses the named config file
func ParseConfigFile(file string) error {

	cFileData, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(cFileData, &Cfg)
	if err != nil {
		return err
	}

	return nil
}
