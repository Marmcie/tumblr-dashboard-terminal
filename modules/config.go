package modules

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Consumer_key string
	Secret_key   string
	Debug        bool
}

var debug bool

func GetConfig() Config {

	dir, _ := os.UserHomeDir()

	path := dir + "\\.config\\tumblr-term-dash.json"

	configBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{Debug: false}
	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		print("Error reading config file.\n")
		print("Please refer to readme file for config format.\n")
		log.Fatal(err)
	}

	debug = config.Debug

	return config
}

func IsDebug() bool {
	return debug
}
