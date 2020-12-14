package main

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var cliCommands []string
var tos glib.SourceHandle

func setupCliLabel() *gtk.Label {
	o := GetCliOutput(cliCommands)
	label, err := gtk.LabelNew(o)
	label.SetJustify(gtk.JUSTIFY_CENTER)
	Check(err)
	return label
}

func refreshCliLabel(label gtk.Label) {
	//o := GetCliOutput(cliCommands)
	label.SetText(GetCliOutput(cliCommands))
}

func main() {

	// Load Preferences, Icons and Commands from ~/.local/share/nwgcc/preferences.json
	Settings, err := LoadSettings()
	Check(err)

	// Load user-defined CustomRows and Buttons from ~/.config/config.json
	Config, err := LoadConfig()
	Check(err)

	v := reflect.ValueOf(Config)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	fmt.Println(values)

	// Load CLI command toproduce CliLabel content
	cliCommands = LoadCliCommands()

	iconsDir := ""
	if Settings.Preferences.IconSet == "light" {
		iconsDir = filepath.Join(DataDir(), "icons_light")
	} else if Settings.Preferences.IconSet == "dark" {
		iconsDir = filepath.Join(DataDir(), "icons_dark")
	}
	fmt.Println(CreatePixbuf(iconsDir, Settings.Icons.WifiOff))

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	Check(err)

	win.SetTitle("nwgcc: Control Center")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	boxOuterV, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 36)
	win.Add(boxOuterV)

	boxOuterH, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 36)
	boxOuterV.PackStart(boxOuterH, false, false, 10)

	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 36)
	boxOuterH.PackStart(vBox, true, true, 10)

	cliLabel := setupCliLabel()

	vBox.Add(cliLabel)

	win.SetDefaultSize(300, 200)

	glib.TimeoutAdd(uint(Settings.Preferences.RefreshCliSeconds*1000), func() bool {
		refreshCliLabel(*cliLabel)
		return true
	})

	win.ShowAll()

	gtk.Main()
}
