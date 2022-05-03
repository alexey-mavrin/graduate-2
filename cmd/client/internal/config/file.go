package config

import (
	"encoding/json"
	"os"
)

// Config contains client config parameters set in the config file
type Config struct {
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	ServerAddr string `json:"server_address"`
	CacheFile  string `json:"cache_file"`
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
