package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Username  string
	Wallpaper string
	DPI       int
	Label     struct {
		Color        string
		Margin       int
		UsernameText string
		PasswordText string
	}
	Entry struct {
		WidthChars int
	}
	Box struct {
		OffsetTop    int
		OffsetBottom int
	}
}

func loadConfig(fpath string) (config Configuration, err error) {
	// default configuration
	config.Username = ""
	config.Wallpaper = ""
	config.DPI = 96
	config.Label.Color = "#ffffff"
	config.Label.Margin = 10
	config.Label.UsernameText = "username:"
	config.Label.PasswordText = "password:"
	config.Entry.WidthChars = 10
	config.Box.OffsetTop = 0
	config.Box.OffsetBottom = 0

	// loading conf file
	file, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	return
}
