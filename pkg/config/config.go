package config

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Configuration sources.
type Config struct {
	Places map[string]string
}

// LoadConfig loads the goplaces configuration file from a given file location.
// If the location starts with a http, it will be downloaded over HTTP.
func LoadConfig(location string) (*Config, error) {
	var contents []byte

	if strings.HasPrefix(location, "http") {
		log.Println("Reading config from network")

		response, err := http.Get(location)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		contents, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Reading config from local file")

		var err error
		contents, err = os.ReadFile(location)
		if err != nil {
			return nil, err
		}
	}

	config := &Config{}
	err := yaml.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
