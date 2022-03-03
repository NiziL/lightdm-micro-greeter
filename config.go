package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Username  string
	Wallpaper string
	Entry     struct {
		WidthChars     int
		Margin         int
		XLocationRatio float32
		YLocationRatio float32
	}
}

func loadConfig(fpath string) (config Configuration, err error) {
	// default configuration
	config.Username = ""
	config.Wallpaper = ""
	config.Entry.WidthChars = 10
	config.Entry.Margin = 10
	config.Entry.XLocationRatio = 0.5
	config.Entry.YLocationRatio = 0.5

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
