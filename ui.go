package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

type GreeterUI struct {
	config Configuration
	entry  *gtk.Entry
	label  *gtk.Label
}

const CSS_TEMPLATE = `
window {
	-gtk-dpi: %d;
	background-image: url("%s");
	background-size: cover; 
	background-repeat: no-repeat; 
}
label {
	color: %s;
	margin: %dpx;
}
box {
	margin-top: %dpx;
	margin-bottom: %dpx;
	margin-left: %dpx;
	margin-right: %dpx;
}
entry {
	color: %s;
	background-color: %s;
	caret-color: %s;
	border: none;
	box-shadow: none;
}
`

func NewUI(config Configuration, entryCallback func()) (app *GreeterUI, err error) {
	app = &GreeterUI{}
	app.config = config

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
	app.label, err = gtk.LabelNew("")
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
	app.entry.SetAlignment(config.Entry.TextAlignement)
	app.entry.SetWidthChars(config.Entry.WidthChars)
	app.entry.Connect("activate", entryCallback)
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
	css := fmt.Sprintf(CSS_TEMPLATE,
		config.DPI,
		pickWallpaper(BASE_PATH+config.Wallpaper),
		config.Label.Color,
		config.Label.Margin,
		config.Box.MarginTop,
		config.Box.MarginBottom,
		config.Box.MarginLeft,
		config.Box.MarginRight,
		config.Entry.TextColor,
		config.Entry.BackgroundColor,
		config.Entry.CaretColor,
	)

	cssProvider.LoadFromData(css)
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	app.UsernameMode()

	window.ShowAll()
	return
}

func (app *GreeterUI) Start() {
	app.entry.GrabFocus()
	gtk.Main()
}

func (app *GreeterUI) UsernameMode() {
	app.label.SetText(app.config.Label.UsernameText)
	app.entry.SetVisibility(true)
}

func (app *GreeterUI) PasswordMode() {
	app.label.SetText(app.config.Label.PasswordText)
	app.entry.SetVisibility(false)
}

func (app *GreeterUI) DisableEntry() {
	app.entry.SetSensitive(false)
}

func (app *GreeterUI) EnableEntry() {
	app.entry.SetSensitive(true)
	app.entry.GrabFocus()
}

func (app *GreeterUI) PopText() (txt string, err error) {
	txt, err = app.entry.GetText()
	app.entry.SetText("")
	return
}

func pickWallpaper(fpath string) (pickedpath string) {
	filestat, err := os.Stat(fpath)
	if err != nil {
		log.Printf("[load_wallpaper] error opening %s \n(%s)", fpath, err)
		return ""
	}
	if filestat.IsDir() {
		files, _ := os.ReadDir(fpath)
		fpath += files[rand.Intn(len(files))].Name()
	}
	pickedpath = fpath
	return
}
