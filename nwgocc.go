package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	cliCommands []string
	iconsDir    string
	settings    Settings
)

// These values need updates
var (
	wifiIcon  string // to track changes
	wifiLabel *gtk.Label
	wifiImage *gtk.Image

	btIcon  string
	btLabel *gtk.Label
	btImage *gtk.Image
)

func setupCliLabel() *gtk.Label {
	o := GetCliOutput(cliCommands)
	label, err := gtk.LabelNew(o)
	label.SetJustify(gtk.JUSTIFY_CENTER)
	Check(err)
	return label
}

func updateCliLabel(label gtk.Label) {
	label.SetText(GetCliOutput(cliCommands))
}

func setupUserRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	pixbuf := CreatePixbuf(iconsDir, settings.Icons.User, settings.Preferences.IconSizeSmall)
	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	Check(err)
	hBox.PackStart(image, false, false, 2)
	name := fmt.Sprintf("%s@%s", GetCommandOutput(settings.Commands.GetUser), GetCommandOutput(settings.Commands.GetHost))
	label, _ := gtk.LabelNew(name)
	hBox.PackStart(label, false, false, 2)

	if settings.Preferences.OnClickUser != "" {
		pixbuf := CreatePixbuf(iconsDir, settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			LaunchCommand(settings.Preferences.OnClickUser)
		})
	}

	eventBox.Connect("enter-notify-event", func() {
		if settings.Preferences.CustomStyling {
			hBox.SetProperty("name", "row-selected")
		} else {
			styleContext.SetState(gtk.STATE_FLAG_SELECTED)
		}
	})

	eventBox.Connect("leave-notify-event", func() {
		if settings.Preferences.CustomStyling {
			hBox.SetProperty("name", "row-normal")
		} else {
			styleContext.SetState(gtk.STATE_FLAG_NORMAL)
		}
	})

	eventBox.Add(hBox)

	return eventBox
}

func setupWifiRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	ssid := fmt.Sprintf("%s", GetCommandOutput(settings.Commands.GetSsid))
	wifiIcon = settings.Icons.WifiOff
	var wifiText string
	if ssid != "" {
		wifiText = ssid
		wifiIcon = settings.Icons.WifiOn
	} else {
		wifiText = "disconnected"
	}
	pixbuf := CreatePixbuf(iconsDir, wifiIcon, settings.Preferences.IconSizeSmall)
	wifiImage, _ = gtk.ImageNew()
	wifiImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(wifiImage, false, false, 2)

	wifiLabel, _ = gtk.LabelNew(wifiText)
	wifiLabel.SetText(wifiText)
	hBox.PackStart(wifiLabel, false, false, 2)

	if settings.Preferences.OnClickWifi != "" {
		pixbuf := CreatePixbuf(iconsDir, settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			LaunchCommand(settings.Preferences.OnClickWifi)
		})
	}

	eventBox.Connect("enter-notify-event", func() {
		if settings.Preferences.CustomStyling {
			hBox.SetProperty("name", "row-selected")
		} else {
			styleContext.SetState(gtk.STATE_FLAG_SELECTED)
		}
	})
	eventBox.Connect("leave-notify-event", func() {
		if settings.Preferences.CustomStyling {
			hBox.SetProperty("name", "row-normal")
		} else {
			styleContext.SetState(gtk.STATE_FLAG_NORMAL)
		}
	})
	eventBox.Add(hBox)

	return eventBox
}

func updateWifiRow() {
	ssid := fmt.Sprintf("%s", GetCommandOutput(settings.Commands.GetSsid))
	icon := ""
	var status string
	if ssid != "" {
		status = ssid
		icon = settings.Icons.WifiOn
	} else {
		status = "disconnected"
		icon = settings.Icons.WifiOff
	}
	if icon != wifiIcon {
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		wifiImage.SetFromPixbuf(pixbuf)
		wifiIcon = icon
	}
	wifiLabel.SetText(status)
}

func setupBluetoothRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	btOn := fmt.Sprintf("%s", GetCommandOutput(settings.Commands.GetBluetoothStatus)) == "yes"
	var status string
	if btOn {
		btIcon = settings.Icons.BtOn
		status = fmt.Sprintf("%s", GetCommandOutput(settings.Commands.GetBluetoothName))
	} else {
		btIcon = settings.Icons.BtOff
		status = "disabled"
	}
	pixbuf := CreatePixbuf(iconsDir, btIcon, settings.Preferences.IconSizeSmall)
	btImage, _ = gtk.ImageNew()
	btImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(btImage, false, false, 2)

	btLabel, _ = gtk.LabelNew(status)
	btLabel.SetText(status)
	hBox.PackStart(btLabel, false, false, 2)

	if settings.Preferences.OnClickBluetooth != "" {
		pixbuf := CreatePixbuf(iconsDir, settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			LaunchCommand(settings.Preferences.OnClickBluetooth)
		})
	}

	eventBox.Connect("enter-notify-event", func() {
		if settings.Preferences.CustomStyling {
			hBox.SetProperty("name", "row-selected")
		} else {
			styleContext.SetState(gtk.STATE_FLAG_SELECTED)
		}
	})
	eventBox.Connect("leave-notify-event", func() {
		if settings.Preferences.CustomStyling {
			hBox.SetProperty("name", "row-normal")
		} else {
			styleContext.SetState(gtk.STATE_FLAG_NORMAL)
		}
	})
	eventBox.Add(hBox)

	return eventBox
}

func updateBluetoothRow() {
	btOn := fmt.Sprintf("%s", GetCommandOutput(settings.Commands.GetBluetoothStatus)) == "yes"
	icon := ""
	var status string
	if btOn {
		icon = settings.Icons.BtOn
		status = fmt.Sprintf("%s", GetCommandOutput(settings.Commands.GetBluetoothName))
	} else {
		icon = settings.Icons.BtOff
		status = "disabled"
	}
	if icon != btIcon {
		fmt.Printf("Status: %s, Icon: %s\n", status, icon)
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		btImage.SetFromPixbuf(pixbuf)
		btIcon = icon
	}
	btLabel.SetText(status)
}

func handleKeyboard(window *gtk.Window, event *gdk.Event) {
	key := &gdk.EventKey{Event: event}
	if key.KeyVal() == gdk.KEY_Escape {
		gtk.MainQuit()
	}
}

func main() {
	timeStart := time.Now()

	// Load Preferences, Icons and Commands from ~/.local/share/nwgcc/preferences.json
	settings, _ = LoadSettings()

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

	// Empty means: gtk icons in use
	iconsDir = ""
	if settings.Preferences.IconSet == "light" {
		iconsDir = filepath.Join(DataDir(), "icons_light")
	} else if settings.Preferences.IconSet == "dark" {
		iconsDir = filepath.Join(DataDir(), "icons_dark")
	}

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	Check(err)

	win.SetTitle("nwgcc: Control Center")
	win.SetDecorated(settings.Preferences.WindowDecorations)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.Connect("key-release-event", handleKeyboard)

	boxOuterV, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 36)
	win.Add(boxOuterV)

	boxOuterH, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 36)
	boxOuterV.PackStart(boxOuterH, false, false, 10)

	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	boxOuterH.PackStart(vBox, true, true, 10)

	var cliLabel *gtk.Label
	if settings.Preferences.ShowCliLabel {
		cliLabel = setupCliLabel()
		vBox.PackStart(cliLabel, true, true, 4)
		sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		vBox.PackStart(sep, true, true, 4)
	}

	if settings.Preferences.ShowUserLine {
		userRow := setupUserRow()
		vBox.PackStart(userRow, false, false, 4)
	}

	var wifiRow *gtk.EventBox
	if settings.Preferences.ShowWifiLine {
		wifiRow = setupWifiRow()
		vBox.PackStart(wifiRow, false, false, 4)
	}

	var btRow *gtk.EventBox
	if settings.Preferences.ShowBtLine {
		btRow = setupBluetoothRow()
		vBox.PackStart(btRow, false, false, 4)
	}

	win.SetDefaultSize(300, 200)

	glib.TimeoutAdd(uint(settings.Preferences.RefreshCliSeconds*1000), func() bool {
		if cliLabel != nil {
			updateCliLabel(*cliLabel)
		}
		if wifiRow != nil {
			updateWifiRow()
		}
		if btRow != nil {
			updateBluetoothRow()
		}
		return true
	})

	win.ShowAll()

	fmt.Printf("Time: %v ms\n", time.Now().Sub(timeStart).Milliseconds())
	gtk.Main()
}
