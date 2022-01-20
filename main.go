package main

// We might use this kind of pkg-config later to ensure portability
//#cgo pkg-config: lightdm-gobject-1
// But for now, let's stick on the old and ugly way

/*
#cgo CFLAGS: -I/usr/include/lightdm-gobject-1 -I/usr/include/glib-2.0 -I/usr/lib/glib-2.0/include
#cgo LDFLAGS: -llightdm-gobject-1 -lgobject-2.0

#include "lightdm.h"

extern void test_cb(void);

void connect_wrapper(GObject* instance, const gchar* detailed_signal, void (*c_handler)(void)) {
    g_signal_connect(instance, detailed_signal, G_CALLBACK(c_handler), NULL);
}

void authenticate_signal_connect(LightDMGreeter* greeter) {
    g_signal_connect(greeter, "authentication-complete", G_CALLBACK(test_cb), NULL);
}
*/
import "C"

import (
    "unsafe"
	"log"
)

//export test_cb
func test_cb() {
	log.Print("callback called!")
}
/*
seems to be the easiest way to create g_callback from Go
but in this case, the preambule should only contains declaration, no definition
either I can add __attribute__((weak)) to C definition so the linker ignore duplicated definitions
or I can put the C function in a .c file and compile it first
*/

func signal_connect(object *C.LightDMGreeter, signal string, callback func()) {
    /*
	C.g_object_connect(greeter, "authentication-complete", nil, nil)
	// unexpected type: ...
	// => cannot use variadic C func with cgo
	// => have to write a C wrapper

	C.g_signal_connect(greeter, "authentication-complete", C.G_CALLBACK(test_cb), nil)
	// "could not determine kind of name for C.g_signal_connect"
    // => cgo cannot use/load C macro 
    */

    signal_name := C.CString(signal)
    defer C.free(unsafe.Pointer(signal_name))

    C.connect_wrapper((*C.GObject)(unsafe.Pointer(object)), signal_name, callback)
}

func main() {
	log.Print("lightdm-micro-greeter start up")

	greeter := C.lightdm_greeter_new()
    defer C.free(unsafe.Pointer(greeter))
	log.Print("greeter created")

	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) > 0 {
		log.Fatal("Not sync to daemon !")
	}
	log.Print("greeter connected to lightdm deamon")


    C.authenticate_signal_connect(greeter)
    //signal_connect(greeter, "authentication-complete", test_cb)
}
