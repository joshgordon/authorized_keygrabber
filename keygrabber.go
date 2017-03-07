package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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
	//Check that we have the right number of command line args.
	if len(os.Args) < 3 {
		log.Fatal("must specify configuration file and destination path")
	}

	//Parse the config file
	config, err := ParseConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	//Pull down the files from the internet.
	filePath, err := GetFileFromURL(config.Baseurl + config.Keyfile)
	if err != nil {
		log.Fatal(err)
	}
	sigFilePath, err := GetFileFromURL(config.Baseurl + config.Sigfile)
	if err != nil {
		log.Fatal(err)
	}

	//Check the signature.
	badSig := CheckSig(filePath, sigFilePath, config.TrustedSig)

	//if it's good, go ahead and copy to the destination.
	if badSig == nil {
		fmt.Println("Good sig!")
		//move the authorized_keys file into place
		fmt.Printf("Moving %s to %s\n", filePath, os.Args[2])
		os.Remove(os.Args[2])
		src, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		dst, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(dst, src)

		//set the permissions on the dest:
		os.Chmod(os.Args[2], 0644)
		//otherwise we just move on.
	} else {
		fmt.Println("Bad sig!")
		fmt.Println("Not moving any files anywhere")
	}
	//clean up.
	os.Remove(filePath)
	os.Remove(sigFilePath)
}
