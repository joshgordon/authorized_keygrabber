package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	baseurl    string
	keyfile    string
	sigfile    string
	trustedSig string
}

func main() {
	if len(os.Args) >= 2 {
		fmt.Println("HI. I'm alive")

		var configFile = os.Args[1]
		_, err := os.Stat(configFile)
		if err != nil {
			log.Fatal("Config file is missing", configFile)
		} else {

			var config Config
			if _, err := toml.DecodeFile(configFile, &config); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", config)
		}
	} else {
		log.Fatal("must specify configuration file")
	}
}
