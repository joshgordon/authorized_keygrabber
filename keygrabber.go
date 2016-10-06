package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
  "http"
)

type Config struct {
	Baseurl    string
	Keyfile    string
	Sigfile    string
	TrustedSig string
}


//Parse the config file and return a config struct.
func ParseConfig(configFile string) (Config, error) {
  _, err := os.Stat(configFile)
  if err != nil {
    return Config{}, err
  } else {
    var config Config
    if _, err := toml.DecodeFile(configFile, &config); err != nil {
      return Config{}, err
    } else {
      return config, err
    }
  }
}

func GetFileFromURL(url string) (err) {
  resp, err := http.Get(url)
  defer resp.close()

  tempFile, err := ioutil.TempFile("", "keygrabber")
  defer tempFile.close()

  
}


func main() {
	if len(os.Args) >= 2 {
		fmt.Println("HI. I'm alive")

    config, err := ParseConfig(os.Args[1])
    if err == nil {
      fmt.Printf("%+v\n", config)
    } else {
      log.Fatal(err)
    }
	} else {
		log.Fatal("must specify configuration file")
	}
}
