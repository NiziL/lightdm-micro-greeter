package main

import (
	"encoding/json"
	"log"
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
	config.Username = ""
	config.Wallpaper = ""
	config.Entry.WidthChars = 10
	config.Entry.Margin = 10
	config.Entry.XLocationRatio = 0.5
	config.Entry.YLocationRatio = 0.5

	file, err := os.Open(fpath)
	if err != nil {
		log.Print("[load_config] error opening " + CONFIG_FILE)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Print("[load_config] " + CONFIG_FILE + "is not a valid JSON")
	}

	return
}
