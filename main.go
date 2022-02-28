package main

/*
#cgo pkg-config: liblightdm-gobject-1 gobject-2.0

#include "lightdm.h"

void greeter_signal_connect(LightDMGreeter* greeter);
*/
import "C"

import (
	"fmt"
	"log"
	"unsafe"
)

/*************/
/* CONSTANTS */
/*************/

const BASE_PATH = "/etc/lightdm/lightdm-micro-greeter/"
const CONFIG_FILE = BASE_PATH + "config.json"

/***************/
/* GLOBAL VARS */
/***************/

// flag, nil in multi user mode
var c_username *C.char = nil

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

/************************/
/* GTK Callback Factory */
/************************/

func createEntryCallback(greeter *C.LightDMGreeter) func() {
	return func() {
		input, _ := entry.GetText()
		entry.SetText("")

		c_input := C.CString(input)
		defer C.free(unsafe.Pointer(c_input))

		if C.lightdm_greeter_get_is_authenticated(greeter) != 0 {
			// session starting
			// looks like dead code, not printed in log
			label.SetText("session starts...")
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

/*************************/
/* lightdm-micro-greeter */
/*************************/

func initGreeter(config Configuration) (greeter *C.LightDMGreeter, err error) {
	greeter = C.lightdm_greeter_new()
	if C.lightdm_greeter_connect_to_daemon_sync(greeter, nil) == 0 {
		log.Print("[init_greeter] can't connect to LightDM deamon")
		err = fmt.Errorf("can't connect to LightDM deamon")
	} else {
		log.Print("[init_greeter] greeter connected to LightDM deamon")
		C.greeter_signal_connect(greeter)
	}

	// Autologin
	if config.Username != "" {
		c_username = C.CString(config.Username)
		defer C.free(unsafe.Pointer(c_username))
		C.lightdm_greeter_authenticate(greeter, c_username, nil)
	}

	return
}

func main() {
	config, err := loadConfig(CONFIG_FILE)
	if err != nil {
		log.Printf("[load_config] fallback on default configuration (%s)", err)
	}

	greeter, err := initGreeter(config)
	if err != nil {
		log.Fatalf("[start_greeter] fatal error: %s", err)
	}

	err = initUI(config, createEntryCallback(greeter))
	if err != nil {
		log.Fatalf("[init_ui] fatal error: %s", err)
	}
}
