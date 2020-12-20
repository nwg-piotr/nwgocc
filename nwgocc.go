package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/itchyny/volume-go"
)

var (
	cliCommands []string
	iconsDir    string
	settings    Settings
	config      Configuration
)

// These values need updates
var (
	wifiIcon  string // to track changes
	wifiLabel *gtk.Label
	wifiImage *gtk.Image

	btIcon  string
	btLabel *gtk.Label
	btImage *gtk.Image

	batIcon  string
	batLabel *gtk.Label
	batImage *gtk.Image

	briIcon   string
	briSlider *gtk.Scale
	briImage  *gtk.Image

	volIcon   string
	volSlider *gtk.Scale
	volImage  *gtk.Image
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
	}

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
	}

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
	}

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
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		btImage.SetFromPixbuf(pixbuf)
		btIcon = icon
	}
	btLabel.SetText(status)
}

func setupBatteryRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	status := ""
	val := 0
	if isCommand(settings.Commands.GetBattery) {
		status, val = getBattery(settings.Commands.GetBattery)
	} else if isCommand(settings.Commands.GetBatteryAlt) {
		status, val = getBattery(settings.Commands.GetBatteryAlt)
	}

	switch {
	case val > 95:
		batIcon = settings.Icons.BatteryFull
	case val > 50:
		batIcon = settings.Icons.BatteryGood
	case val > 20:
		batIcon = settings.Icons.BatteryLow
	default:
		batIcon = settings.Icons.BatteryEmpty
	}

	pixbuf := CreatePixbuf(iconsDir, batIcon, settings.Preferences.IconSizeSmall)
	batImage, _ = gtk.ImageNew()
	batImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(batImage, false, false, 2)

	batLabel, _ = gtk.LabelNew(status)
	batLabel.SetText(status)
	hBox.PackStart(batLabel, false, false, 2)

	if settings.Preferences.OnClickBattery != "" {
		pixbuf := CreatePixbuf(iconsDir, settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			LaunchCommand(settings.Preferences.OnClickBattery)
		})
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
	}
	eventBox.Add(hBox)

	return eventBox
}

func updateBatteryRow() {
	status := ""
	val := 0
	if isCommand(settings.Commands.GetBattery) {
		status, val = getBattery(settings.Commands.GetBattery)
	} else if isCommand(settings.Commands.GetBatteryAlt) {
		status, val = getBattery(settings.Commands.GetBatteryAlt)
	}
	icon := ""
	switch {
	case val > 95:
		icon = settings.Icons.BatteryFull
	case val > 50:
		icon = settings.Icons.BatteryGood
	case val > 20:
		icon = settings.Icons.BatteryLow
	default:
		icon = settings.Icons.BatteryEmpty
	}

	if icon != batIcon {
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		batImage, _ = gtk.ImageNew()
		batImage.SetFromPixbuf(pixbuf)
		batIcon = icon
	}

	batLabel.SetText(status)
}

func setupBrightnessBox() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	bri := getBrightness()
	icon := ""
	switch {
	case bri > 70:
		icon = settings.Icons.BrightnessHigh
	case bri > 30:
		icon = settings.Icons.BrightnessMedium
	default:
		icon = settings.Icons.BrightnessLow
	}
	pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
	briImage, _ = gtk.ImageNew()
	briImage.SetFromPixbuf(pixbuf)
	box.PackStart(briImage, false, false, 2)

	briSlider, _ = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	briSlider.SetValue(bri)
	briSlider.Connect("value-changed", func() {
		b := briSlider.GetValue()
		setBrightness(int(b))
	})

	box.PackStart(briSlider, true, true, 2)

	return box
}

func updateBrightnessRow() {
	bri := getBrightness()
	icon := ""
	switch {
	case bri > 70:
		icon = settings.Icons.BrightnessHigh
	case bri > 30:
		icon = settings.Icons.BrightnessMedium
	default:
		icon = settings.Icons.BrightnessLow
	}
	if icon != briIcon {
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		briImage.SetFromPixbuf(pixbuf)
	}
	briSlider.SetValue(bri)
}

func setupVolumeBox() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	vol, _ := volume.GetVolume()
	muted, err := volume.GetMuted()
	Check(err)
	icon := ""
	if !muted {
		switch {
		case vol > 70:
			icon = settings.Icons.VolumeHigh
		case vol > 30:
			icon = settings.Icons.VolumeMedium
		default:
			icon = settings.Icons.VolumeLow
		}
	} else {
		icon = settings.Icons.VolumeMuted
	}

	pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
	volImage, _ = gtk.ImageNew()
	volImage.SetFromPixbuf(pixbuf)
	box.PackStart(volImage, false, false, 2)

	volSlider, _ = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	volSlider.SetValue(float64(vol))
	volSlider.Connect("value-changed", func() {
		b := volSlider.GetValue()
		err := volume.SetVolume(int(b))
		Check(err)
	})

	box.PackStart(volSlider, true, true, 2)

	return box
}

func updateVolumeRow() {
	vol, _ := volume.GetVolume()
	muted, err := volume.GetMuted()
	Check(err)
	icon := ""
	if !muted {
		switch {
		case vol > 70:
			icon = settings.Icons.VolumeHigh
		case vol > 30:
			icon = settings.Icons.VolumeMedium
		default:
			icon = settings.Icons.VolumeLow
		}
	} else {
		icon = settings.Icons.VolumeMuted
	}

	if icon != volIcon {
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		volImage.SetFromPixbuf(pixbuf)
	}

	volSlider.SetValue(float64(vol))
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

	//CheckCommands(settings.Commands)

	// Load user-defined CustomRows and Buttons from ~/.config/config.json
	config, _ = LoadConfig()

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
		cliCommands = LoadCliCommands()
		if len(cliCommands) > 0 {
			cliLabel = setupCliLabel()
			vBox.PackStart(cliLabel, true, true, 4)
			sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
			vBox.PackStart(sep, true, true, 4)
		}
	}

	var briRow *gtk.Box
	if settings.Preferences.ShowBrightnessSlider {
		briRow = setupBrightnessBox()
		vBox.PackStart(briRow, false, false, 4)
	}

	var volRow *gtk.Box
	if settings.Preferences.ShowVolumeSlider {
		volRow = setupVolumeBox()
		vBox.PackStart(volRow, false, false, 4)
	}

	if settings.Preferences.ShowBrightnessSlider || settings.Preferences.ShowVolumeSlider {
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
	if settings.Preferences.ShowBtLine && btServiceEnabled() {
		btRow = setupBluetoothRow()
		vBox.PackStart(btRow, false, false, 4)
	}

	var batRow *gtk.EventBox
	if settings.Preferences.ShowBatteryLine {
		batRow = setupBatteryRow()
		vBox.PackStart(batRow, false, false, 4)
	}

	win.SetDefaultSize(300, 200)

	glib.TimeoutAdd(uint(settings.Preferences.RefreshCliSeconds*1000), func() bool {
		if cliLabel != nil {
			updateCliLabel(*cliLabel)
		}
		return true
	})
	glib.TimeoutAdd(uint(settings.Preferences.RefreshSlowSeconds*1000), func() bool {
		if batRow != nil {
			updateBatteryRow()
		}
		return true
	})
	glib.TimeoutAdd(uint(settings.Preferences.RefreshFastMillis), func() bool {
		if briRow != nil {
			updateBrightnessRow()
		}

		if volRow != nil {
			updateVolumeRow()
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

	fmt.Printf("Ready in %v ms\n", time.Now().Sub(timeStart).Milliseconds())
	gtk.Main()
}
