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
		WidthChars      int
		TextAlignement  float32
		TextColor       string
		BackgroundColor string
		CaretColor      string
	}
	Box struct {
		MarginLeft   int
		MarginTop    int
		MarginBottom int
		MarginRight  int
	}
}

func loadConfig(fpath string) (config Configuration, err error) {
	// default configuration
	config.Username = ""
	config.Wallpaper = ""
	config.DPI = 96
	config.Label.Color = "#000000"
	config.Label.Margin = 0
	config.Label.UsernameText = "username:"
	config.Label.PasswordText = "password:"
	config.Entry.WidthChars = 10
	config.Entry.TextColor = "#000000"
	config.Entry.BackgroundColor = "#ffffff"
	config.Entry.CaretColor = "#000000"
	config.Entry.TextAlignement = 0.5
	config.Box.MarginTop = 0
	config.Box.MarginBottom = 0
	config.Box.MarginLeft = 0
	config.Box.MarginRight = 0

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
