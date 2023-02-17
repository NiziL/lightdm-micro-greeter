package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

type AppUI struct {
	entry *gtk.Entry
	label *gtk.Label
}

const CSS_DATA = `
window {
	-gtk-dpi: 384;
	background-image: url("/etc/lightdm/lightdm-micro-greeter/wallpapers/earth.png");
	background-size: 100%; 
	background-repeat: no-repeat; 
}
label {
	color: #ffffff;
	margin: 10;
}
box {
	margin-top: 0;
	margin-bottom: 0;
}
`

func (app *AppUI) Init(config Configuration) (err error) {
	gtk.Init(nil)

	// Create fullscreen window
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return
	}
	window.Connect("destroy", func() {
		// looks like dead code, never really called
		gtk.MainQuit()
	})
	// get screen information to resize as full screen
	// window.Fullscreen() is not working here
	display, err := window.GetDisplay()
	if err != nil {
		return
	}
	monitor, err := display.GetPrimaryMonitor()
	if err != nil {
		return
	}
	window.Resize(monitor.GetGeometry().GetWidth(), monitor.GetGeometry().GetHeight())

	// Create UI layout
	// init box
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}
	box.SetHAlign(gtk.ALIGN_CENTER)
	box.SetVAlign(gtk.ALIGN_CENTER)
	window.Add(box)

	// init label
	app.label, err = gtk.LabelNew("username:")
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

	// Setup CSS provider
	screen := window.GetScreen()
	if err != nil {
		return
	}
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		return
	}
	cssProvider.LoadFromData(CSS_DATA) // TODO forge CSS_DATA with config
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	window.ShowAll()
	return
}

func (app *AppUI) Start(entryCallback func()) {
	app.entry.Connect("activate", entryCallback)
	app.entry.GrabFocus()
	gtk.Main()
}

func (app *AppUI) UsernameMode() {
	app.label.SetText("username:")
	app.entry.SetVisibility(true)
}

func (app *AppUI) PasswordMode() {
	app.label.SetText("password:")
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

func pickWallpaper(fpath string, width, height int) (pickedpath string, err error) {
	filestat, err := os.Stat(fpath)
	if err != nil {
		err = fmt.Errorf("[load_wallpaper] error opening %s \n(%s)", fpath, err)
		return
	}
	if filestat.IsDir() {
		files, _ := os.ReadDir(fpath)
		fpath += files[rand.Intn(len(files))].Name()
	}
	pickedpath = fpath
	return
}
