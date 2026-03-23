package modules

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Consumer_key string
	Secret_key   string
	Oauth_token  string
	Oauth_secret string
}

func GetConfig() Config {

	dir, _ := os.UserHomeDir()

	path := dir + "\\.config\\tumblr2.json"

	configBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func (c *Config) ToString() string {
	str := ""
	str += "consumer_key : " + c.Consumer_key + "\n"
	str += "secret_key  : " + c.Secret_key + "\n"
	str += "oauth_token : " + c.Oauth_token + "\n"
	str += "oauth_secret : " + c.Oauth_secret
	return str
}
