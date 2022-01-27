package main

// TODO use pkg-config instead of flags to ensure portability
//#cgo pkg-config: lightdm-gobject-1 (?)

/*
#cgo CFLAGS: -I/usr/include/lightdm-gobject-1 -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include
#cgo LDFLAGS: -llightdm-gobject-1 -lgobject-2.0 -lglib-2.0

#include "lightdm.h"

extern void authentication_complete_cb(LightDMGreeter *greeter);
extern void show_prompt_cb(LightDMGreeter *greeter, char *text, LightDMPromptType type);

extern void greeter_signal_connect(LightDMGreeter* greeter);
*/
import "C"

import (
	"log"
	"unsafe"
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
	password := "password"

	c_password := C.CString(password)
	C.lightdm_greeter_respond(greeter, c_password, nil)
}

func main() {
	log.Print("lightdm-micro-greeter start up")

	loop := C.g_main_loop_new(nil, 0)
	// should i defer free loop ?

	greeter := C.lightdm_greeter_new()
	//defer C.free(unsafe.Pointer(greeter)) // TODO fix invalid pointer at runtime, should I really free this ? Does glib free it for me ?
	log.Print("greeter created")

	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) == 0 {
		log.Fatal("Not sync to daemon !")
	}
	log.Print("greeter connected to lightdm deamon")

	C.greeter_signal_connect(greeter)
	log.Print("greeter callbacks binded")

	// TODO let user give username
	username := "username"

	c_username := C.CString(username)
	defer C.free(unsafe.Pointer(c_username))

	C.lightdm_greeter_authenticate(greeter, c_username, nil)
	log.Print("authentication started")

	C.g_main_loop_run(loop)
}
