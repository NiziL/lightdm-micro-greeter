package main

/*
#cgo CFLAGS: -I/usr/include/lightdm-gobject-1 -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include
#cgo LDFLAGS: -llightdm-gobject-1
#include "lightdm.h"
*/
import "C"

// We might use this kind of pkg-config later to ensure portability
//#cgo pkg-config: lightdm-gobject-1

import (
	"log"
)

func main() {
	log.Print("lightdm-micro-greeter start up")

	greeter := C.lightdm_greeter_new()
	log.Print("greeter created")

	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) > 0 {
		log.Fatal("Not sync to daemon !")
	}
	log.Print("greeter connected")
}
