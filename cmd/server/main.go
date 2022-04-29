package main

import (
	"log"
	"os"

	"github.com/alexey-mavrin/graduate-2/cmd/server/internal/config"
	"github.com/alexey-mavrin/graduate-2/internal/server"
)

const (
	defaultConfigFile = "server.cfg"
)

func main() {
	configFile, ok := os.LookupEnv("SERVER_CFG")
	if !ok {
		configFile = defaultConfigFile
	}

	err := config.ParseConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer(config.Cfg.ListenPort, config.Cfg.StoreFile)
	if err != nil {
		log.Fatal(err)
	}
}
