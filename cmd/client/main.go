package main

import (
	"log"
	"os"

	"github.com/alexey-mavrin/graduate-2/cmd/client/internal/action"
	"github.com/alexey-mavrin/graduate-2/cmd/client/internal/config"
)

const defaultConfigFile = "gosecret.cfg"

func main() {
	var err error
	configFile, ok := os.LookupEnv("GOSECRET_CFG")
	if !ok {
		configFile = defaultConfigFile
	}

	err = config.ParseConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	err = action.ChooseAct()
	if err != nil {
		log.Fatal(err)
	}
}
