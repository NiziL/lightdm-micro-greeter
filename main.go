package main

// We might use this kind of pkg-config later to ensure portability
//#cgo pkg-config: lightdm-gobject-1
// But for now, let's stick on the old and ugly way

/*
#cgo CFLAGS: -I/usr/include/lightdm-gobject-1 -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include
#cgo LDFLAGS: -llightdm-gobject-1 -lgobject-2.0

#include "lightdm.h"

extern void authentication_complete_cb(void);
extern void show_prompt_cb(void);

__attribute__((weak)) // TODO that's ugly and I should fix it
void greeter_signal_connect(LightDMGreeter* greeter) {
    g_signal_connect(greeter, "authentication-complete", G_CALLBACK(authentication_complete_cb), NULL);
    g_signal_connect(greeter, "show-prompt", G_CALLBACK(show_prompt_cb), NULL);
}
*/
import "C"

import (
    "unsafe"
	"log"
)

//export authentication_complete_cb
func authentication_complete_cb() {
	log.Print("authentication cb called!")
}

//export show_prompt_cb
func show_prompt_cb() {
    log.Print("prompt cb called!")
}

func main() {
	log.Print("lightdm-micro-greeter start up")

	greeter := C.lightdm_greeter_new()
    //defer C.free(unsafe.Pointer(greeter)) // TODO invalid pointer ?
	log.Print("greeter created")

	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) > 0 {
		log.Fatal("Not sync to daemon !")
	}
	log.Print("greeter connected to lightdm deamon")

    C.greeter_signal_connect(greeter)

    username := C.CString("user")
    defer C.free(unsafe.Pointer(username))
    C.lightdm_greeter_authenticate(greeter, username, nil)

    for true {
    }
}
