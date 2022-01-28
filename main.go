package main

/*
#cgo CFLAGS: -I/usr/include/lightdm-gobject-1 -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include
#cgo LDFLAGS: -llightdm-gobject-1 -lgobject-2.0 -lglib-2.0

#include "lightdm.h"

void greeter_signal_connect(LightDMGreeter* greeter);
*/
import "C"

import (
	"log"
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

// TODO I'd like to avoid these global vars
var entry *gtk.Entry
var label *gtk.Label

//export authentication_complete_cb
func authentication_complete_cb(greeter *C.LightDMGreeter) {
	if C.lightdm_greeter_get_is_authenticated(greeter) == 0 {
		log.Print("[authentication_complete_cb] wrong password")
		label.SetText("username:")
		entry.SetVisibility(true)
		entry.SetSensitive(true)
		entry.GrabFocus()
	} else {
		log.Print("[authentication_complete_cb] starting session")
		C.lightdm_greeter_start_session_sync(greeter, nil, nil)
	}
}

//export show_prompt_cb
func show_prompt_cb(greeter *C.LightDMGreeter, text *C.char, prompt_type C.LightDMPromptType) {
	// text is the lightdm deamon answer, have seen nothing else than "password:"
	log.Print("[show_prompt_cb]")
	txt := C.GoString(text)
	label.SetText(txt)
	if prompt_type == C.LIGHTDM_PROMPT_TYPE_SECRET {
		entry.SetVisibility(false)
	} else {
		entry.SetVisibility(true)
	}
}

//export show_message_cb
func show_message_cb(greeter *C.LightDMGreeter, text *C.char, msg_type C.LightDMMessageType) {
	log.Print("[show_message_cb]")
	txt := C.GoString(text)
	label.SetText(txt)
}

// may I use this builder to avoid global var (at least entry and label) ?
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

	var err error

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

	gtk.Init(nil)
	log.Print("gtk init")

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("unable to create window", err)
	}
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	display, err := win.GetDisplay()
	monitor, err := display.GetPrimaryMonitor()
	rect := monitor.GetGeometry()
	win.Resize(rect.GetWidth(), rect.GetHeight())

	grid, err := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	label, err = gtk.LabelNew("username:")
	grid.Add(label)

	entry, err = gtk.EntryNew()
	grid.AttachNextTo(entry, label, gtk.POS_RIGHT, 1, 1)

	entry.Connect("activate", create_entry_cb(greeter))

	win.Add(grid)
	win.ShowAll()

	gtk.Main()
}
