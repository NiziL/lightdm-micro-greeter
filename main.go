package main

/*
#cgo pkg-config: liblightdm-gobject-1 gobject-2.0

#include "lightdm.h"

void greeter_signal_connect(LightDMGreeter* greeter);
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// CONSTANTS
const BASE_PATH = "/etc/lightdm/lightdm-micro-greeter/"
const CONFIG_FILE = BASE_PATH + "config.json"

// GLOBAL VARS
// TODO this is ulgy, can I avoid them ?
var c_username *C.char = nil
var entry *gtk.Entry
var label *gtk.Label

/**********************/
/* JSON Configuration */
/**********************/
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

func loadConfig() (config Configuration, err error) {
	config.Username = ""
	config.Wallpaper = ""
	config.Entry.WidthChars = 10
	config.Entry.Margin = 10
	config.Entry.XLocationRatio = 0.5
	config.Entry.YLocationRatio = 0.5

	file, err := os.Open(CONFIG_FILE)
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

/*****************************/
/* Light DM Server Callbacks */
/*****************************/

//export authentication_complete_cb
func authentication_complete_cb(greeter *C.LightDMGreeter) {
	if C.lightdm_greeter_get_is_authenticated(greeter) == 0 {
		log.Print("[authentication_complete_cb] wrong password")
		if c_username != nil {
			C.lightdm_greeter_authenticate(greeter, c_username, nil)
		} else {
			label.SetText("username")
			entry.SetVisibility(true)
		}
		entry.SetSensitive(true)
		entry.GrabFocus()
	} else {
		log.Print("[authentication_complete_cb] starting session")
		C.lightdm_greeter_start_session_sync(greeter, nil, nil)
	}
}

//export show_prompt_cb
func show_prompt_cb(greeter *C.LightDMGreeter, text *C.char, prompt_type C.LightDMPromptType) {
	log.Print("[show_prompt_cb]")
	label.SetText("password")
	if prompt_type == C.LIGHTDM_PROMPT_TYPE_SECRET {
		entry.SetVisibility(false)
	} else {
		entry.SetVisibility(true)
	}
}

/*************/
/* Utilities */
/*************/

func createEntryCallback(greeter *C.LightDMGreeter) func() {
	return func() {
		input, _ := entry.GetText()
		entry.SetText("")

		c_input := C.CString(input)
		defer C.free(unsafe.Pointer(c_input))

		if C.lightdm_greeter_get_is_authenticated(greeter) != 0 {
			// start_session ?
			log.Print("[entry_callback] authentication ok")
		} else if C.lightdm_greeter_get_in_authentication(greeter) != 0 {
			// give pwd
			log.Print("[entry_callback] giving pwd")
			C.lightdm_greeter_respond(greeter, c_input, nil)
			entry.SetSensitive(false)
		} else {
			// give username
			log.Print("[entry_callback] giving username")
			C.lightdm_greeter_authenticate(greeter, c_input, nil)
		}
	}
}

func loadWallpaper(fpath string, width, height int) (bg *gtk.Image, err error) {
	filestat, err := os.Stat(fpath)
	if err != nil {
		log.Print("[load_wallpaper] error opening " + fpath)
		return
	}
	if filestat.IsDir() {
		files, _ := os.ReadDir(fpath)
		rand.Seed(time.Now().UnixNano())
		fpath += files[rand.Intn(len(files))].Name()
		log.Print("[load_wallpaper] randomly picking " + fpath)
	}
	pixbuf, err := gdk.PixbufNewFromFileAtSize(fpath, width, height)
	if err != nil {
		log.Print("[load_wallpaper] error loading " + fpath)
		return
	}
	bg, err = gtk.ImageNewFromPixbuf(pixbuf)
	return
}

func initGreeter() (greeter *C.LightDMGreeter, err error) {
	greeter = C.lightdm_greeter_new()
	// TODO fix invalid pointer at runtime
	// defer C.free(unsafe.Pointer(greeter))
	// should I really free this ? Does glib free it for me ? should I use g_free ?
	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) == 0 {
		log.Print("[init_greeter] can't connect to LightDM deamon")
		err = fmt.Errorf("can't connect to LightDM deamon")
	} else {
		log.Print("[init_greeter] greeter connected to LightDM deamon")
		C.greeter_signal_connect(greeter)
	}
	return
}

func initWindow() (window *gtk.Window, width, height int, err error) {
	window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return
	}
	// seems useless, never called when session starts
	window.Connect("destroy", func() {
		log.Print("destroy signal called: quitting gtk")
		gtk.MainQuit()
	})

	// get screen information to resize as full screen
	// .Fullscreen() is not working here
	display, err := window.GetDisplay()
	if err != nil {
		return
	}
	monitor, err := display.GetPrimaryMonitor()
	if err != nil {
		return
	}
	rect := monitor.GetGeometry()
	width = rect.GetWidth()
	height = rect.GetHeight()

	window.Resize(width, height)

	return
}

func initEntryBox(margin, widthChars int, callback func()) (box *gtk.Box, err error) {
	box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, margin)
	if err != nil {
		return
	}
	box.SetHAlign(gtk.ALIGN_CENTER)
	box.SetVAlign(gtk.ALIGN_CENTER)

	label, err = gtk.LabelNew("username")
	if err != nil {
		return
	}
	label.SetHAlign(gtk.ALIGN_CENTER)
	label.SetVAlign(gtk.ALIGN_CENTER)
	box.Add(label)

	entry, err = gtk.EntryNew()
	if err != nil {
		return
	}
	entry.SetHAlign(gtk.ALIGN_CENTER)
	entry.SetVAlign(gtk.ALIGN_CENTER)
	entry.SetWidthChars(widthChars)
	entry.Connect("activate", callback)
	box.Add(entry)

	return
}

func main() {
	log.Print("[main] start up")

	// Reading configuration file
	config, _ := loadConfig()

	// Start greeter
	greeter, err := initGreeter()
	if err != nil {
		log.Fatalf("[start_greeter] fatal error: %s", err)
	}

	// Init GUI
	log.Print("[main] gtk init")
	gtk.Init(nil)

	// fullscreen window
	win, width, height, err := initWindow()
	if err != nil {
		log.Fatalf("[init_window] fatal error: %s", err)
	}

	// simple fixed layout
	layout, _ := gtk.FixedNew()
	win.Add(layout)

	// set background image, auto scaling while preserving aspect ratio
	bg, err := loadWallpaper(BASE_PATH+config.Wallpaper, width, height)
	if err != nil {
		log.Print("[load_wallpaper] default white screen")
	} else {
		log.Print("[load_wallpaper] wallpaper loaded")
		layout.Put(bg, 0, 0)
	}

	// init entry box
	box, _ := initEntryBox(
		config.Entry.Margin, config.Entry.WidthChars,
		createEntryCallback(greeter))
	/* TODO find a cleaner way to acheive this, might induce flickering
	   for now, I have to put the box and render it before having access to its size */
	// set box centered
	layout.Add(box)
	win.ShowAll()
	// now that size is known, compute center
	center_x := int(float32(width) * config.Entry.XLocationRatio)
	center_y := int(float32(height) * config.Entry.YLocationRatio)
	offset_x := box.GetAllocatedWidth() / 2
	offset_y := box.GetAllocatedWidth() / 2
	// center box
	layout.Remove(box)
	layout.Put(box, center_x-offset_x, center_y-offset_y)
	// set cursor in entry
	entry.GrabFocus()

	// Starts autologin if provided
	if config.Username != "" {
		c_username = C.CString(config.Username)
		defer C.free(unsafe.Pointer(c_username))
		C.lightdm_greeter_authenticate(greeter, c_username, nil)
	}

	log.Print("[main] gtk start main loop")
	gtk.Main()
}
