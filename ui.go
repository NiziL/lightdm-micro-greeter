package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type AppUI struct {
	entry *gtk.Entry
	label *gtk.Label
}

func (app *AppUI) Init(config Configuration) (err error) {
	gtk.Init(nil)

	// fullscreen window
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return
	}
	// looks like dead code, never really called
	window.Connect("destroy", func() {
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
	width := rect.GetWidth()
	height := rect.GetHeight()

	window.Resize(width, height)
	if err != nil {
		return
	}

	// simple fixed layout
	layout, err := gtk.FixedNew()
	if err != nil {
		return
	}
	window.Add(layout)

	// init box for label and entry
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, config.Entry.Margin)
	if err != nil {
		return
	}
	box.SetHAlign(gtk.ALIGN_CENTER)
	box.SetVAlign(gtk.ALIGN_CENTER)

	// init label
	app.label, err = gtk.LabelNew("username")
	if err != nil {
		return
	}
	app.label.SetHAlign(gtk.ALIGN_CENTER)
	app.label.SetVAlign(gtk.ALIGN_CENTER)
	box.Add(app.label)

	// init entry
	app.entry, err = gtk.EntryNew()
	if err != nil {
		return
	}
	app.entry.SetHAlign(gtk.ALIGN_CENTER)
	app.entry.SetVAlign(gtk.ALIGN_CENTER)
	app.entry.SetWidthChars(config.Entry.WidthChars)
	box.Add(app.entry)

	// set box centered
	// TODO find a cleaner way to acheive this, might induce flickering
	// for now, I have to put the box and render it before having access to its size
	layout.Add(box)
	window.ShowAll()
	// now that size is known, compute center
	center_x := int(float32(width) * config.Entry.XLocationRatio)
	center_y := int(float32(height) * config.Entry.YLocationRatio)
	offset_x := box.GetAllocatedWidth() / 2
	offset_y := box.GetAllocatedWidth() / 2
	// center box
	layout.Remove(box)

	// set background image, auto scaling while preserving aspect ratio
	bg, err := loadWallpaper(BASE_PATH+config.Wallpaper, width, height)
	if err != nil {
		err = fmt.Errorf("[load_wallpaper] error loading wallpaper \n(%s)", err)
	} else {
		layout.Put(bg, 0, 0)
	}
	layout.Put(box, center_x-offset_x, center_y-offset_y)
	window.ShowAll()

	return
}

func (app *AppUI) Start(entryCallback func()) {
	app.entry.Connect("activate", entryCallback)
	app.entry.GrabFocus()
	gtk.Main()
}

func (app *AppUI) UsernameMode() {
	app.label.SetText("username")
	app.entry.SetVisibility(true)
}

func (app *AppUI) PasswordMode() {
	app.label.SetText("password")
	app.entry.SetVisibility(false)
}

func (app *AppUI) DisableEntry() {
	app.entry.SetSensitive(false)
}

func (app *AppUI) EnableEntry() {
	app.entry.SetSensitive(true)
	app.entry.GrabFocus()
}

func (app *AppUI) PopText() (txt string, err error) {
	txt, err = app.entry.GetText()
	app.entry.SetText("")
	return
}

func loadWallpaper(fpath string, width, height int) (bg *gtk.Image, err error) {
	filestat, err := os.Stat(fpath)
	if err != nil {
		err = fmt.Errorf("[load_wallpaper] error opening %s \n(%s)", fpath, err)
		return
	}
	if filestat.IsDir() {
		files, _ := os.ReadDir(fpath)
		rand.Seed(time.Now().UnixNano())
		fpath += files[rand.Intn(len(files))].Name()
	}
	pixbuf, err := gdk.PixbufNewFromFileAtSize(fpath, width, height)
	if err != nil {
		err = fmt.Errorf("[load_wallpaper] error loading %s \n(%s)", fpath, err)
		return
	}
	bg, err = gtk.ImageNewFromPixbuf(pixbuf)
	return
}
