package main

/*
#cgo pkg-config: liblightdm-gobject-1 gobject-2.0

#include "lightdm.h"

void greeter_signal_connect(LightDMGreeter* greeter);
*/
import "C"

import (
	"encoding/json"
	"log"
	"os"
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

func load_config() (config Configuration) {
	config.Username = ""
	config.Wallpaper = ""
	config.Entry.WidthChars = 10
	config.Entry.Margin = 10
	config.Entry.XLocationRatio = 0.5
	config.Entry.YLocationRatio = 0.5

	file, err := os.Open(CONFIG_FILE)
	defer file.Close()
	if err != nil {
		log.Fatal("[load_config] Does /etc/lightdm/lightdm-micro-greeter/config.json exist ?")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("[load_config] Does not look like a valid json syntax.")
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

// TODO investigate lightdm server, do I really need this ?
//export show_message_cb
func show_message_cb(greeter *C.LightDMGreeter, text *C.char, msg_type C.LightDMMessageType) {
	log.Print("[show_message_cb]")
	txt := C.GoString(text)
	label.SetText(txt)
}

/*****************/
/* GTK Callbacks */
/*****************/

func create_entry_cb(greeter *C.LightDMGreeter) func() {
	return func() {
		input, _ := entry.GetText()
		entry.SetText("")

		c_input := C.CString(input)
		defer C.free(unsafe.Pointer(c_input))

		if C.lightdm_greeter_get_is_authenticated(greeter) != 0 {
			// start_session ?
			log.Print("authentication ok")
		} else if C.lightdm_greeter_get_in_authentication(greeter) != 0 {
			// give pwd
			C.lightdm_greeter_respond(greeter, c_input, nil)
			entry.SetSensitive(false)
		} else {
			// give username
			C.lightdm_greeter_authenticate(greeter, c_input, nil)
		}
	}
}

func main() {
	log.Print("lightdm-micro-greeter start up")

	// Reading configuration file
	config := load_config()

	// Connect to LightDM deamon
	greeter := C.lightdm_greeter_new()
	// TODO fix invalid pointer at runtime
	// defer C.free(unsafe.Pointer(greeter))
	// should I really free this ? Does glib free it for me ? should I use g_free ?
	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) == 0 {
		log.Fatal("greeter can't connect to daemon")
	} else {
		log.Print("greeter connected to lightdm deamon")
	}
	C.greeter_signal_connect(greeter)

	// GTK
	log.Print("gtk init")
	gtk.Init(nil)

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	// seems useless, never called when session starts
	win.Connect("destroy", func() {
		log.Print("destroy signal called: quitting gtk")
		gtk.MainQuit()
	})

	// get screen information to resize as full screen
	// .Fullscreen() is not working here
	display, _ := win.GetDisplay()
	monitor, _ := display.GetPrimaryMonitor()
	rect := monitor.GetGeometry()
	win.Resize(rect.GetWidth(), rect.GetHeight())

	// simple fixed layout
	layout, _ := gtk.FixedNew()
	win.Add(layout)

	// create a box for label and entry
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, config.Entry.Margin)
	box.SetHAlign(gtk.ALIGN_CENTER)
	box.SetVAlign(gtk.ALIGN_CENTER)

	label, _ = gtk.LabelNew("username")
	label.SetHAlign(gtk.ALIGN_CENTER)
	label.SetVAlign(gtk.ALIGN_CENTER)
	box.Add(label)

	entry, _ = gtk.EntryNew()
	entry.SetHAlign(gtk.ALIGN_CENTER)
	entry.SetVAlign(gtk.ALIGN_CENTER)
	entry.SetWidthChars(config.Entry.WidthChars)
	entry.Connect("activate", create_entry_cb(greeter))
	box.Add(entry)

	// set background image, auto scaling while preserving aspect ratio
	if config.Wallpaper != "" {
		pixbuf, _ := gdk.PixbufNewFromFileAtSize(BASE_PATH+config.Wallpaper, rect.GetWidth(), rect.GetHeight())
		bg, _ := gtk.ImageNewFromPixbuf(pixbuf)
		layout.Put(bg, 0, 0)
	}

	// set box
	// TODO find a cleaner way to acheive this, might induce flickering
	// I have to put the box and render it before having access to its size
	layout.Add(box)
	win.ShowAll()
	// now that size is known, compute center
	center_x := int(float32(rect.GetWidth()) * config.Entry.XLocationRatio)
	center_y := int(float32(rect.GetHeight()) * config.Entry.YLocationRatio)
	offset_x := box.GetAllocatedWidth() / 2
	offset_y := box.GetAllocatedWidth() / 2
	// center box
	layout.Remove(box)
	layout.Put(box, center_x-offset_x, center_y-offset_y)
	entry.GrabFocus()

	// Starts autologin if provided
	if config.Username != "" {
		c_username = C.CString(config.Username)
		defer C.free(unsafe.Pointer(c_username))
		C.lightdm_greeter_authenticate(greeter, c_username, nil)
	}

	log.Print("gtk start")
	gtk.Main()
}
