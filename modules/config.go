package modules

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Consumer_key string
	Secret_key   string
	Timezone     string
	Debug        bool
}

var config Config
var initialized bool

func GetConfig() Config {

	if !initialized {
		loadConfig()
	}
	return config
}
func loadConfig() {
	dir, _ := os.UserHomeDir()

	path := dir + "\\.config\\tumblr-dt.json"

	configBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	c := Config{Debug: false}
	err = json.Unmarshal(configBytes, &c)

	if err != nil {
		print("Error reading config file.\n")
		print("Please refer to readme file for config format.\n")
		log.Fatal(err)
	}

	config = c

	if len(config.Timezone) == 0 {
		config.Timezone = "Local"
	}

	initialized = true
}

func IsDebug() bool {
	return config.Debug
}
