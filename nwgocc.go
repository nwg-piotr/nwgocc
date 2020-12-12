package main

import (
	"fmt"
	"reflect"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	fmt.Printf("Config dir: %s\n", ConfigDir())
	fmt.Printf("Data dir: %s\n", DataDir())

	// Load Preferences, Icons and Commands from ~/.local/share/nwgcc/preferences.json
	Settings, err := LoadSettings()
	Check(err)

	v := reflect.ValueOf(Settings)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	fmt.Println(values)

	// Load user-defined CustomRows and Buttons from ~/.config/config.json
	Config, err := LoadConfig()
	Check(err)

	v = reflect.ValueOf(Config)

	values = make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	fmt.Println(values)

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

	l, _ := gtk.LabelNew("Hello, gotk3!")

	vBox.Add(l)

	win.SetDefaultSize(300, 200)

	win.ShowAll()

	gtk.Main()
}
