package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
  "io/ioutil"
  "io"
  "net/http"
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

func GetFileFromURL(url string) (string, error) {
  resp, err := http.Get(url)
  defer resp.Body.Close()
  if err != nil {
    return "", err
  }

  tempFile, err := ioutil.TempFile("", "keygrabber")
  defer tempFile.Close()
  if err != nil {
    return "", err
  }

  _, err = io.Copy(tempFile, resp.Body)
  if err != nil {
    return "", err
  }

  return tempFile.Name(), nil

}


func main() {
	if len(os.Args) >= 2 {

    config, err := ParseConfig(os.Args[1])
    if err == nil {
      filePath, err := GetFileFromURL(config.Baseurl + config.Keyfile)
      if err != nil{
        log.Fatal(err)
      }
      fmt.Println(filePath)
      sigFilePath, err := GetFileFromURL(config.Baseurl + config.Sigfile)
      if err != nil{
        log.Fatal(err)
      }
      fmt.Println(sigFilePath)
    } else {
      log.Fatal(err)
    }
	} else {
		log.Fatal("must specify configuration file")
	}
}
