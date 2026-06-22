package modules

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
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
	Post_theme            string
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

	Keymaps struct {
		Navigation struct {
			Up         string
			Down       string
			Left       string
			Right      string
			JumpNext   string
			JumpPrev   string
			JumpTop    string
			JumpBottom string
		}

		Switcher struct {
			Open       string
			Close      string
			Up         string
			Down       string
			Suggestion string
		}

		Links struct {
			Open  string
			Close string
		}
		IncreaseSize string
		DecreaseSize string

		ToggleFeed string

		ControlHelp string
		OpenLink    string
		LoadMore    string
		LoadBlog    string
		Confirm     string
		Quit        string
		LogOut      string
		Log         string
	}
	Initialized bool
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
	con.Post_theme = "default"
	con.Initialized = false

	con.Keymaps.Navigation = struct {
		Up         string
		Down       string
		Left       string
		Right      string
		JumpNext   string
		JumpPrev   string
		JumpTop    string
		JumpBottom string
	}{
		Up:         "k",
		Down:       "j",
		Left:       "h",
		Right:      "l",
		JumpNext:   "d",
		JumpPrev:   "u",
		JumpTop:    "g",
		JumpBottom: "G",
	}
	con.Keymaps.Switcher = struct {
		Open       string
		Close      string
		Up         string
		Down       string
		Suggestion string
	}{
		Open:       "]",
		Close:      "esc",
		Up:         "up",
		Down:       "down",
		Suggestion: "right",
	}
	con.Keymaps.Links = struct {
		Open  string
		Close string
	}{
		Open:  "L",
		Close: "esc",
	}

	con.Keymaps.OpenLink = "o"
	con.Keymaps.LoadMore = "r"
	con.Keymaps.LoadBlog = "b"
	con.Keymaps.Confirm = "enter"
	con.Keymaps.Quit = "q"
	con.Keymaps.LogOut = "delete"
	con.Keymaps.IncreaseSize = "right"
	con.Keymaps.DecreaseSize = "left"
	con.Keymaps.ToggleFeed = "t"
	con.Keymaps.ControlHelp = "?"
	con.Keymaps.Log = "ctrl+l"
	return con
}

var config Config
var initialized bool

func GetConfig() Config {

	if !initialized {
		result := loadConfig()
		if !result {
			return makeConfig()
		}
	}
	return config
}

func resolveConfigPath() string {

	flagPath := flag.String("config", "", "config path")
	flag.Parse()
	if len(*flagPath) > 0 {
		_, err := os.ReadFile(*flagPath)
		if err == nil {
			return *flagPath
		}
	}

	configPath, _ := os.UserConfigDir()
	filePath := fmt.Sprintf("%s\\tumblr-dt\\tumblr-dt.toml", configPath)

	_, err := os.ReadFile(filePath)
	if err == nil {
		return filePath
	}

	dir, _ := os.UserHomeDir()
	filePath = fmt.Sprintf("%s\\.config\\tumblr-dt.toml", dir)

	_, err = os.ReadFile(filePath)

	if err == nil {
		return filePath
	}
	return ""
}

func loadConfig() bool {
	path := resolveConfigPath()
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	c := makeConfig()
	_, err = toml.Decode(string(configBytes), &c)

	if err != nil {
		print("Error reading config file.\n")
		print("Please refer to readme file for config format.\n")
		log.Fatal(err)
	}

	config = c

	if len(config.Timezone) == 0 {
		config.Timezone = "Local"
	}

	config.Initialized = true
	initialized = true
	return true
}

func IsDebug() bool {
	return config.Debug
}
