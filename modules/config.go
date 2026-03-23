package modules

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Consumer_key          string
	Secret_key            string
	Use_blog_avatar_color bool
	Timezone              string
	Debug                 bool
	Testing               bool
	Blacklist             []string
	Redirect_port         string
	Colors                struct {
		Bg           string
		Focus        string
		Focus_border string
		White        string
		Grey         string
		H1           string
		H2           string
		Image        string
		Quote        string
		Filtered     string
		Blacklisted  string
	}
}

func makeConfig() Config {
	con := Config{}
	con.Colors.Bg = "#060616"
	con.Colors.Focus = "#135366"
	con.Colors.Focus_border = "#30c0f0"
	con.Colors.White = "#ffffff"
	con.Colors.Grey = "#aaaaaa"
	con.Colors.H1 = "#40f0f0"
	con.Colors.H2 = "#a0f000"
	con.Colors.Image = "#40a0f0"
	con.Colors.Quote = "#f0f000"
	con.Colors.Filtered = "#ff0000"
	con.Colors.Blacklisted = "#ff00ff"

	con.Debug = false
	con.Testing = false

	return con
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

	c := makeConfig()
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
