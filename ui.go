package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/***************/
/* GLOBAL VARS */
/***************/

// gtk widget needed by exported function
var entry *gtk.Entry = nil
var label *gtk.Label = nil

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

func initUI(config Configuration, entryCallback func()) {
	log.Print("[init_ui] gtk init")
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
	box, _ := initEntryBox(config.Entry.Margin, config.Entry.WidthChars, entryCallback)

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

	log.Print("[init_ui] gtk start main loop")
	gtk.Main()
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
