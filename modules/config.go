package modules

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Consumer_key string
	Secret_key   string
}

func GetConfig() Config {

	dir, _ := os.UserHomeDir()

	path := dir + "\\.config\\tumblr-terminal.json"

	configBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		print("Error reading config file.\n")
		print("Please refer to readme file for config format.\n")
		log.Fatal(err)
	}

	return config
}
