package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/itchyny/volume-go"
)

var version = "0.1.0"

var (
	cliCommands []string
	iconsDir    string
	settings    Settings
	config      Configuration
)

var customCSS = flag.String("css", "style.css", "custom css file name")
var debug = flag.Bool("d", false, "do checks, print results")
var displayVersion = flag.Bool("v", false, "display version information")
var winPosPointer = flag.Bool("p", false, "place window at the mouse pointer position (Xorg only)")
var restoreDefaults = flag.Bool("r", false, "restore defaults (preferences, templates and icons)")

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
	playIcon  string
	volSlider *gtk.Scale
	volImage  *gtk.Image
	playImage *gtk.Image
)

var configChanged = false

func setupCliLabel() *gtk.Label {
	o := GetCliOutput(cliCommands)
	label, err := gtk.LabelNew(o)
	label.SetProperty("name", "cli-label")
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

	pixbuf := CreatePixbuf(settings.Icons.User, settings.Preferences.IconSizeSmall)
	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	Check(err)
	hBox.PackStart(image, false, false, 2)
	name := fmt.Sprintf("%s@%s", GetCommandOutput(settings.Commands.GetUser), GetCommandOutput(settings.Commands.GetHost))
	label, _ := gtk.LabelNew(name)
	hBox.PackStart(label, false, false, 2)

	if settings.Preferences.OnClickUser != "" {
		pixbuf := CreatePixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
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
	pixbuf := CreatePixbuf(wifiIcon, settings.Preferences.IconSizeSmall)
	wifiImage, _ = gtk.ImageNew()
	wifiImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(wifiImage, false, false, 2)

	wifiLabel, _ = gtk.LabelNew(wifiText)
	wifiLabel.SetText(wifiText)
	hBox.PackStart(wifiLabel, false, false, 2)

	if settings.Preferences.OnClickWifi != "" {
		pixbuf := CreatePixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
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
		pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
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
	pixbuf := CreatePixbuf(btIcon, settings.Preferences.IconSizeSmall)
	btImage, _ = gtk.ImageNew()
	btImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(btImage, false, false, 2)

	btLabel, _ = gtk.LabelNew(status)
	btLabel.SetText(status)
	hBox.PackStart(btLabel, false, false, 2)

	if settings.Preferences.OnClickBluetooth != "" {
		pixbuf := CreatePixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
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
		pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
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

	pixbuf := CreatePixbuf(batIcon, settings.Preferences.IconSizeSmall)
	batImage, _ = gtk.ImageNew()
	batImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(batImage, false, false, 2)

	batLabel, _ = gtk.LabelNew(status)
	batLabel.SetText(status)
	hBox.PackStart(batLabel, false, false, 2)

	if settings.Preferences.OnClickBattery != "" {
		pixbuf := CreatePixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
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
		pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
		batImage, _ = gtk.ImageNew()
		batImage.SetFromPixbuf(pixbuf)
		batIcon = icon
	}

	batLabel.SetText(status)
}

func setupBrightnessRow() *gtk.Box {
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
	pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
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
		pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
		briImage.SetFromPixbuf(pixbuf)
	}
	briSlider.SetValue(bri)
}

func setupVolumeRow() *gtk.Box {
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

	pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
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

	if settings.Preferences.ShowPlayerctl && isCommand(settings.Commands.Playerctl) {
		playerctlStatus := GetCommandOutput("playerctl status /dev/null 2>&1")
		fmt.Println(playerctlStatus)
		if playerctlStatus == "Playing" || playerctlStatus == "Paused" {
			icon := settings.Icons.MediaSkipBackward
			pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
			image, _ := gtk.ImageNew()
			image.SetFromPixbuf(pixbuf)
			eb, _ := gtk.EventBoxNew()
			eb.Connect("button-press-event", func() {
				cmd := exec.Command("playerctl", "previous")
				cmd.Run()
			})
			eb.Add(image)
			box.PackStart(eb, false, false, 0)

			if playerctlStatus == "Playing" {
				icon = settings.Icons.MediaPlaybackPause
			} else {
				icon = settings.Icons.MediaPlaybackStart
			}
			pixbuf = CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
			playImage, _ = gtk.ImageNew()
			playImage.SetFromPixbuf(pixbuf)
			eb, _ = gtk.EventBoxNew()
			eb.Connect("button-press-event", func() {
				cmd := exec.Command("playerctl", "play-pause")
				cmd.Run()
			})
			eb.Add(playImage)
			box.PackStart(eb, false, false, 0)

			icon = settings.Icons.MediaSkipForward
			pixbuf = CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
			image, _ = gtk.ImageNew()
			image.SetFromPixbuf(pixbuf)
			eb, _ = gtk.EventBoxNew()
			eb.Connect("button-press-event", func() {
				cmd := exec.Command("playerctl", "next")
				cmd.Run()
			})
			eb.Add(image)
			box.PackStart(eb, false, false, 0)
		}
	}

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
		pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
		volImage.SetFromPixbuf(pixbuf)
	}

	volSlider.SetValue(float64(vol))

	if settings.Preferences.ShowPlayerctl {
		if GetCommandOutput("playerctl status /dev/null 2>&1") == "Playing" {
			icon = settings.Icons.MediaPlaybackPause
		} else {
			icon = settings.Icons.MediaPlaybackStart
		}
		if icon != playIcon {
			pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
			if playImage != nil {
				playImage.SetFromPixbuf(pixbuf)
			}
			playIcon = icon
		}
	}
}

func setupCustomRow(icon, name, cmd string) *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	if icon != "" {
		pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNew()
		image.SetFromPixbuf(pixbuf)
		hBox.PackStart(image, false, false, 2)
	}

	if name != "" {
		label, _ := gtk.LabelNew(name)
		label.SetText(name)
		hBox.PackStart(label, false, false, 2)
	}

	if cmd != "" {
		pixbuf := CreatePixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			LaunchCommand(cmd)
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

func setupPreferencesButton() *gtk.Button {
	button, _ := gtk.ButtonNew()
	if settings.Preferences.CustomStyling {
		button.SetProperty("name", "custom-button")
	}
	pixbuf := CreatePixbuf("emblem-system-symbolic", settings.Preferences.IconSizeLarge)
	image, _ := gtk.ImageNewFromPixbuf(pixbuf)
	button.SetImage(image)
	button.SetAlwaysShowImage(true)
	button.SetTooltipText("Preferences")
	button.Connect("clicked", func() {
		setupPreferencesWindow()
	})

	return button
}

func setupCustomButton(icon, name, cmd string) *gtk.Button {
	button, _ := gtk.ButtonNew()
	if settings.Preferences.CustomStyling {
		button.SetProperty("name", "custom-button")
	}
	pixbuf := CreatePixbuf(icon, settings.Preferences.IconSizeLarge)
	image, _ := gtk.ImageNewFromPixbuf(pixbuf)
	button.SetImage(image)
	button.SetAlwaysShowImage(true)
	if name != "" {
		button.SetTooltipText(name)
	}
	button.Connect("clicked", func() {
		LaunchCommand(cmd)
	})

	return button
}

func handleKeyboard(window *gtk.Window, event *gdk.Event) {
	key := &gdk.EventKey{Event: event}
	if key.KeyVal() == gdk.KEY_Escape {
		gtk.MainQuit()
	}
}

func main() {
	timeStart := time.Now()

	flag.Parse()

	if *displayVersion {
		fmt.Printf("nwgocc version %s\n", version)
		os.Exit(0)
	}

	setupDirs()

	// Load Preferences, Icons and Commands from ~/.local/share/nwgocc/preferences.json
	settings, _ = LoadSettings()

	if *debug {
		CheckCommands(settings.Commands)
	}

	// Load user-defined CustomRows and Buttons from ~/.config/config.json
	config, _ = LoadConfig()

	// Empty means: gtk icons in use
	iconsDir = ""
	if settings.Preferences.IconSet == "light" {
		iconsDir = filepath.Join(DataDir(), "icons_light")
		fmt.Println("Icons: Custom light")
	} else if settings.Preferences.IconSet == "dark" {
		iconsDir = filepath.Join(DataDir(), "icons_dark")
		fmt.Println("Icons: Custom dark")
	} else {
		fmt.Println("Icons: GTK")
	}

	gtk.Init(nil)

	if settings.Preferences.CustomStyling {
		css := filepath.Join(ConfigDir(), *customCSS)
		fmt.Printf("Style: '%s'\n", css)
		cssProvider, err := gtk.CssProviderNew()
		Check(err)
		err = cssProvider.LoadFromPath(css)
		if err != nil {
			fmt.Println(err)
		}
		screen, _ := gdk.ScreenGetDefault()
		gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
	} else {
		fmt.Println("Style: GTK")
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	Check(err)

	win.SetTitle("nwgocc: Control Center")
	win.SetProperty("name", "window")
	win.SetDecorated(settings.Preferences.WindowDecorations)
	if *winPosPointer {
		win.SetPosition(gtk.WIN_POS_MOUSE)
	}
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

	cliCommands = LoadCliCommands()
	var cliLabel *gtk.Label
	if settings.Preferences.ShowCliLabel {
		if len(cliCommands) > 0 {
			cliLabel = setupCliLabel()
			vBox.PackStart(cliLabel, true, true, 4)
			sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
			vBox.PackStart(sep, true, true, 6)
		}
	}

	var briRow *gtk.Box
	if settings.Preferences.ShowBrightnessSlider {
		briRow = setupBrightnessRow()
		vBox.PackStart(briRow, false, false, 4)
	}

	var volRow *gtk.Box
	if settings.Preferences.ShowVolumeSlider {
		volRow = setupVolumeRow()
		vBox.PackStart(volRow, false, false, 4)
	}

	if settings.Preferences.ShowBrightnessSlider || settings.Preferences.ShowVolumeSlider {
		sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		vBox.PackStart(sep, true, true, 6)
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

	if settings.Preferences.ShowUserRows {
		sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		vBox.PackStart(sep, true, true, 6)

		for _, item := range config.CustomRows {
			customRow := setupCustomRow(item.Icon, item.Name, item.Command)
			vBox.PackStart(customRow, false, false, 4)
		}
	}

	sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	vBox.PackStart(sep, true, true, 6)

	buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	preferencesButton := setupPreferencesButton()
	buttonBox.PackStart(preferencesButton, true, false, 4)

	if settings.Preferences.ShowUserButtons {
		for _, item := range config.Buttons {
			customBtn := setupCustomButton(item.Icon, item.Name, item.Command)
			buttonBox.PackStart(customBtn, true, false, 4)
		}
	}

	vBox.PackStart(buttonBox, false, false, 8)

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
