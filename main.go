package main

// TODO use pkg-config instead of flags to ensure portability
//#cgo pkg-config: lightdm-gobject-1 (?)

/*
#cgo CFLAGS: -I/usr/include/lightdm-gobject-1 -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include
#cgo LDFLAGS: -llightdm-gobject-1 -lgobject-2.0 -lglib-2.0

#include "lightdm.h"

extern void authentication_complete_cb(LightDMGreeter *greeter);
extern void show_prompt_cb(LightDMGreeter *greeter, char *text, LightDMPromptType type);

// TODO that's ugly and I should fix it
// preamble should not contain any definition when using //export later
__attribute__((weak))
void greeter_signal_connect(LightDMGreeter* greeter) {
	g_signal_connect(greeter, "authentication-complete", G_CALLBACK(authentication_complete_cb), NULL);
	g_signal_connect(greeter, "show-prompt", G_CALLBACK(show_prompt_cb), NULL);
}
*/
import "C"

import (
	"log"
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

// TODO I'd like to avoid this global var
var waiting_pwd, waiting_resp bool
var entry *gtk.Entry
var label *gtk.Label

//export authentication_complete_cb
func authentication_complete_cb(greeter *C.LightDMGreeter) {
	log.Print("authentication_complete_cb called!")

	if C.lightdm_greeter_get_is_authenticated(greeter) == 0 {
		log.Print("wrong password")
		waiting_pwd = false
	} else {
		C.lightdm_greeter_start_session_sync(greeter, nil, nil)
	}
	waiting_resp = false
}

//export show_prompt_cb
func show_prompt_cb(greeter *C.LightDMGreeter, text *C.char, prompt_type C.LightDMPromptType) {
	// text is the lightdm deamon answer, have seen nothing else than "password:"
	log.Print("show_prompt_cb called!")
	waiting_pwd = true
}

// may I use this builder to avoid global var (at least entry and label) ?
func create_entry_cb(greeter *C.LightDMGreeter) func() {
	return func() {
		input, _ := entry.GetText()
		entry.SetText("")

		c_input := C.CString(input)
		defer C.free(unsafe.Pointer(c_input))

		if waiting_pwd {
			log.Print("pwd entered")
			C.lightdm_greeter_respond(greeter, c_input, nil)
			waiting_resp = true
			label.SetText("username: ")
			entry.SetVisibility(true)
		} else if !waiting_resp {
			log.Print("starting authentication for ", input)
			C.lightdm_greeter_authenticate(greeter, c_input, nil)
			label.SetText("password: ")
			entry.SetVisibility(false)
		}
	}
}

func main() {
	log.Print("lightdm-micro-greeter start up")

	var err error

	waiting_pwd = false
	waiting_resp = false

	greeter := C.lightdm_greeter_new()
	// TODO fix invalid pointer at runtime
	// defer C.free(unsafe.Pointer(greeter))
	// should I really free this ? Does glib free it for me ? should I use g_free ?
	log.Print("greeter created")

	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) == 0 {
		log.Fatal("Not sync to daemon !")
	}
	log.Print("greeter connected to lightdm deamon")

	C.greeter_signal_connect(greeter)
	log.Print("greeter callbacks binded")

	gtk.Init(nil)
	log.Print("gtk init")

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("unable to create window:", err)
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
