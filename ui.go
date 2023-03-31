package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"syscall"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	_ "embed"
)

type GreeterUI struct {
	config Configuration
	entry  *gtk.Entry
	label  *gtk.Label
}

//go:embed template.css
var CSS_TEMPLATE string

func NewUI(config Configuration, entryCallback func()) (app *GreeterUI, err error) {
	app = &GreeterUI{}
	app.config = config

	gtk.Init(nil)

	// Create fullscreen window
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return
	}
	// get screen information to resize as full screen
	// Fullscreen() and Maximize() are not working here
	display, err := window.GetDisplay()
	if err != nil {
		return
	}
	monitor, err := display.GetPrimaryMonitor()
	if err != nil {
		return
	}
	window.Resize(monitor.GetGeometry().GetWidth(), monitor.GetGeometry().GetHeight())

	// Quit gracefully on destroy
	window.Connect("destroy", func() {
		// looks like dead code, never really called
		log.Printf("greeter window destroyed")
		gtk.MainQuit()
	})

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

	// init callback for power/suspend/hibernate/reboot
	window.Connect("key_press_event", func(win *gtk.Window, e *gdk.Event) {
		event := gdk.EventKeyNewFromEvent(e)
		if event.State() == gdk.CONTROL_MASK+uint(gdk.SHIFT_MASK) {
			switch event.KeyVal() {
			case gdk.KEY_R:
				log.Printf("Triggering restart")
				syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
			case gdk.KEY_P:
				log.Printf("Triggering poweroff")
				syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
			case gdk.KEY_H:
				log.Printf("Triggering hibernation")
				syscall.Reboot(syscall.LINUX_REBOOT_CMD_SW_SUSPEND)
			}
		}
	})

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
