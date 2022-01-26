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

//export authentication_complete_cb
func authentication_complete_cb(greeter *C.LightDMGreeter) {
	log.Print("authentication_complete_cb called!")

	if C.lightdm_greeter_get_is_authenticated(greeter) == 0 {
		log.Print("wrong password dumbass")
	} else {
		C.lightdm_greeter_start_session_sync(greeter, nil, nil)
	}
}

//export show_prompt_cb
func show_prompt_cb(greeter *C.LightDMGreeter, text *C.char, prompt_type C.LightDMPromptType) {
	log.Print("show_prompt_cb called!")

	// TODO let user give password
	//c_password := C.CString(password)
	//defer C.free(unsafe.Pointer(c_password))
	//C.lightdm_greeter_respond(greeter, c_password, nil)
}

func main() {
	log.Print("lightdm-micro-greeter start up")

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

    grid, _ := gtk.GridNew()
    grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	l, _ := gtk.LabelNew("username")
	grid.Add(l)

	entry, _ := gtk.EntryNew()
	entry.Connect("activate", func() {
		username, _ := entry.GetText()
		c_username := C.CString(username)
		defer C.free(unsafe.Pointer(c_username))
		log.Print("starting authentication for ", username)
		C.lightdm_greeter_authenticate(greeter, c_username, nil)
	})
    grid.AttachNextTo(entry, l, gtk.POS_RIGHT, 1, 1)

	l_pwd, _ := gtk.LabelNew("password")
	grid.Add(l_pwd)

	pwd_entry, _ := gtk.EntryNew()
    pwd_entry.SetVisibility(false)
	pwd_entry.Connect("activate", func() {
		pwd, _ := pwd_entry.GetText()
		c_pwd := C.CString(pwd)
		defer C.free(unsafe.Pointer(c_pwd))
        log.Print("sending password")
		C.lightdm_greeter_respond(greeter, c_pwd, nil)
	})
    grid.AttachNextTo(pwd_entry, l_pwd, gtk.POS_RIGHT, 1, 1)
    
    win.Add(grid)
	win.ShowAll()
	gtk.Main()
}
