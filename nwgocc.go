package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/allan-simon/go-singleinstance"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/itchyny/volume-go"
)

const version = "0.1.5"
const playing string = "Playing"
const paused string = "Paused"

var (
	cliCommands   []string
	netInterfaces []string
	iconsDir      string
	settings      Settings
	config        Configuration
	win           *gtk.Window
)

var configFile = flag.String("c", "config.json", "user's templates: Config file name")
var cssFile = flag.String("s", "style.css", "custom Styling: css file name")
var debug = flag.Bool("d", false, "Do checks, print results")
var displayVersion = flag.Bool("v", false, "display Version information")
var winPosPointer = flag.Bool("p", false, "place window at the mouse Pointer position (Xorg only)")
var restoreDefaults = flag.Bool("r", false, "Restore defaults (preferences, templates and icons)")

// These values need updates
var (
	wifiIcon  string // to track changes (avoid creating the icon if status unchanged; same below)
	wifiLabel *gtk.Label
	wifiImage *gtk.Image

	interfaceIcon  string
	interfaceLabel *gtk.Label
	interfaceImage *gtk.Image

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
var wayland bool

// Shows output of CLI commands defined in `~/.config/nwgocc/cli_commands` text file
func setupCliLabel() *gtk.Label {
	o := getCliOutput(cliCommands)
	label, err := gtk.LabelNew(o)
	check(err)
	label.SetProperty("name", "cli-label")
	label.SetJustify(gtk.JUSTIFY_CENTER)

	return label
}

func updateCliLabel(label gtk.Label) {
	label.SetText(getCliOutput(cliCommands))
}

// Shows icon + output of `echo $USER`
func setupUserRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	pixbuf := createPixbuf(settings.Icons.User, settings.Preferences.IconSizeSmall)
	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	check(err)
	hBox.PackStart(image, false, false, 2)
	name := fmt.Sprintf("%s@%s", getCommandOutput(settings.Commands.GetUser), getCommandOutput(settings.Commands.GetHost))
	label, _ := gtk.LabelNew(name)
	hBox.PackStart(label, false, false, 2)

	if settings.Preferences.OnClickUser != "" {
		pixbuf := createPixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			launchCommand(settings.Preferences.OnClickUser)
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

// Shows icon appropriate to status + output of `iwgetid -r`
func setupWifiRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	ssid := fmt.Sprintf("%s", getCommandOutput(settings.Commands.GetSsid))
	wifiIcon = settings.Icons.WifiOff
	var wifiText string
	if ssid != "" {
		wifiText = ssid
		wifiIcon = settings.Icons.WifiOn
	} else {
		wifiText = "disconnected"
	}
	pixbuf := createPixbuf(wifiIcon, settings.Preferences.IconSizeSmall)
	wifiImage, _ = gtk.ImageNew()
	wifiImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(wifiImage, false, false, 2)

	wifiLabel, _ = gtk.LabelNew(wifiText)
	wifiLabel.SetText(wifiText)
	hBox.PackStart(wifiLabel, false, false, 2)

	if settings.Preferences.OnClickWifi != "" {
		pixbuf := createPixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			launchCommand(settings.Preferences.OnClickWifi)
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
	ssid := fmt.Sprintf("%s", getCommandOutput(settings.Commands.GetSsid))
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
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
		wifiImage.SetFromPixbuf(pixbuf)
		wifiIcon = icon
	}
	wifiLabel.SetText(status)
}

// Shows icon appropriate to selecten net interface status
func setupInterfaceRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	netInterfaces = listInterfaces()
	isUp := false
	addr := ""
	if settings.Preferences.InterfaceName != "" {
		isUp, addr = interfaceIsUp(settings.Preferences.InterfaceName)
	}

	if isUp {
		interfaceIcon = settings.Icons.NetworkConnected
	} else {
		interfaceIcon = settings.Icons.NetworkDisonnected
	}

	interfaceText := "Not selected"
	if settings.Preferences.InterfaceName != "" {
		if isUp {
			interfaceText = fmt.Sprintf("%s: %s", settings.Preferences.InterfaceName, addr)
		} else {
			interfaceText = settings.Preferences.InterfaceName
		}
	}

	pixbuf := createPixbuf(interfaceIcon, settings.Preferences.IconSizeSmall)
	interfaceImage, _ = gtk.ImageNew()
	interfaceImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(interfaceImage, false, false, 2)

	interfaceLabel, _ = gtk.LabelNew(interfaceText)
	hBox.PackStart(interfaceLabel, false, false, 2)

	if settings.Preferences.OnClickInterface != "" {
		pixbuf := createPixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			launchCommand(settings.Preferences.OnClickInterface)
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

func updateInterfaceRow() {
	var icon string
	isUp, addr := interfaceIsUp(settings.Preferences.InterfaceName)
	if isUp {
		icon = settings.Icons.NetworkConnected
	} else {
		icon = settings.Icons.NetworkDisonnected
	}
	if icon != interfaceIcon {
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
		interfaceImage.SetFromPixbuf(pixbuf)
		interfaceIcon = icon
	}
	var interfaceText string
	if isUp {
		interfaceText = fmt.Sprintf("%s: %s", settings.Preferences.InterfaceName, addr)
	} else {
		interfaceText = settings.Preferences.InterfaceName
	}
	interfaceLabel.SetText(interfaceText)
}

// Shows icon appropriate to status + output of `bluetoothctl show | awk '/Name/{print $2}'`
func setupBluetoothRow() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	btOn := fmt.Sprintf("%s", getCommandOutput(settings.Commands.GetBluetoothStatus)) == "yes"
	var status string
	if btOn {
		btIcon = settings.Icons.BtOn
		status = fmt.Sprintf("%s", getCommandOutput(settings.Commands.GetBluetoothName))
	} else {
		btIcon = settings.Icons.BtOff
		status = "disabled"
	}
	pixbuf := createPixbuf(btIcon, settings.Preferences.IconSizeSmall)
	btImage, _ = gtk.ImageNew()
	btImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(btImage, false, false, 2)

	btLabel, _ = gtk.LabelNew(status)
	btLabel.SetText(status)
	hBox.PackStart(btLabel, false, false, 2)

	if settings.Preferences.OnClickBluetooth != "" {
		pixbuf := createPixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			launchCommand(settings.Preferences.OnClickBluetooth)
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
	btOn := fmt.Sprintf("%s", getCommandOutput(settings.Commands.GetBluetoothStatus)) == "yes"
	icon := ""
	var status string
	if btOn {
		icon = settings.Icons.BtOn
		status = fmt.Sprintf("%s", getCommandOutput(settings.Commands.GetBluetoothName))
	} else {
		icon = settings.Icons.BtOff
		status = "disabled"
	}
	if icon != btIcon {
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
		btImage.SetFromPixbuf(pixbuf)
		btIcon = icon
	}
	btLabel.SetText(status)
}

// Shows icon appropriate to status + output of
// `upower -i $(upower -e | grep BAT) | grep --color=never -E 'state|to\\\\ full|to\\\\ empty|percentage'`
// or `acpi`
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

	pixbuf := createPixbuf(batIcon, settings.Preferences.IconSizeSmall)
	batImage, _ = gtk.ImageNew()
	batImage.SetFromPixbuf(pixbuf)
	hBox.PackStart(batImage, false, false, 2)

	batLabel, _ = gtk.LabelNew(status)
	batLabel.SetText(status)
	hBox.PackStart(batLabel, false, false, 2)

	if settings.Preferences.OnClickBattery != "" {
		pixbuf := createPixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			launchCommand(settings.Preferences.OnClickBattery)
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
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
		batImage, _ = gtk.ImageNew()
		batImage.SetFromPixbuf(pixbuf)
		batIcon = icon
	}

	batLabel.SetText(status)
}

// Creates the brightness slider; getting and setting the value depends on the `light` command
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
	pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
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
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
		briImage.SetFromPixbuf(pixbuf)
	}
	briSlider.SetValue(bri)
}

// Creates the volume slider; depends on the `volume-go` package.
func setupVolumeRow() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	vol, _ := volume.GetVolume()
	muted, err := volume.GetMuted()
	check(err)
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

	pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
	volImage, _ = gtk.ImageNew()
	volImage.SetFromPixbuf(pixbuf)
	box.PackStart(volImage, false, false, 2)

	volSlider, _ = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 100, 1)
	volSlider.SetValue(float64(vol))
	volSlider.Connect("value-changed", func() {
		b := volSlider.GetValue()
		err := volume.SetVolume(int(b))
		check(err)
	})

	box.PackStart(volSlider, true, true, 2)

	if settings.Preferences.ShowPlayerctl && isCommand(settings.Commands.Playerctl) {
		playerctlStatus := getCommandOutput("playerctl status /dev/null 2>&1")
		if playerctlStatus == playing || playerctlStatus == paused {
			icon := settings.Icons.MediaSkipBackward
			pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
			image, _ := gtk.ImageNew()
			image.SetFromPixbuf(pixbuf)
			eb, _ := gtk.EventBoxNew()
			eb.Connect("button-press-event", func() {
				cmd := exec.Command("playerctl", "previous")
				cmd.Run()
			})
			eb.Add(image)
			box.PackStart(eb, false, false, 0)

			if playerctlStatus == playing {
				icon = settings.Icons.MediaPlaybackPause
			} else {
				icon = settings.Icons.MediaPlaybackStart
			}
			pixbuf = createPixbuf(icon, settings.Preferences.IconSizeSmall)
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
			pixbuf = createPixbuf(icon, settings.Preferences.IconSizeSmall)
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
	check(err)
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
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
		volImage.SetFromPixbuf(pixbuf)
	}

	volSlider.SetValue(float64(vol))

	if settings.Preferences.ShowPlayerctl {
		if getCommandOutput("playerctl status /dev/null 2>&1") == playing {
			icon = settings.Icons.MediaPlaybackPause
		} else {
			icon = settings.Icons.MediaPlaybackStart
		}
		if icon != playIcon {
			pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
			if playImage != nil {
				playImage.SetFromPixbuf(pixbuf)
			}
			playIcon = icon
		}
	}
}

// User-defined rows; name, command and icon defined in `~/.config/nwgocc/config.json`
func setupCustomRow(icon, name, cmd string) *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	styleContext, _ := eventBox.GetStyleContext()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if settings.Preferences.CustomStyling {
		hBox.SetProperty("name", "row-normal")
	}

	if icon != "" {
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
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
		pixbuf := createPixbuf(settings.Icons.ClickMe, settings.Preferences.IconSizeSmall)
		image, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(image, false, false, 2)

		eventBox.Connect("button-press-event", func() {
			launchCommand(cmd)
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

// Built-in Preferences button
func setupPreferencesButton() *gtk.Button {
	button, _ := gtk.ButtonNew()
	if settings.Preferences.CustomStyling {
		button.SetProperty("name", "custom-button")
	}
	ico := "nwgocc"
	if iconsDir != "" {
		ico = fmt.Sprintf("%s/nwgocc-symbolic.svg", iconsDir)
	}
	pixbuf := createPixbuf(ico, settings.Preferences.IconSizeLarge)
	image, _ := gtk.ImageNewFromPixbuf(pixbuf)
	button.SetImage(image)
	button.SetAlwaysShowImage(true)
	button.SetTooltipText("Preferences")
	button.Connect("clicked", func() {
		setupPreferencesWindow()
	})

	return button
}

// User-defined buttons; name, command and icon defined in `~/.config/nwgocc/config.json`
func setupCustomButton(icon, name, cmd string) *gtk.Button {
	button, _ := gtk.ButtonNew()
	if settings.Preferences.CustomStyling {
		button.SetProperty("name", "custom-button")
	}
	pixbuf := createPixbuf(icon, settings.Preferences.IconSizeLarge)
	image, _ := gtk.ImageNewFromPixbuf(pixbuf)
	button.SetImage(image)
	button.SetAlwaysShowImage(true)
	if name != "" {
		button.SetTooltipText(name)
	}
	button.Connect("clicked", func() {
		launchCommand(cmd)
	})

	return button
}

// Exit on Esc key
func handleKeyboard(window *gtk.Window, event *gdk.Event) {
	key := &gdk.EventKey{Event: event}
	if key.KeyVal() == gdk.KEY_Escape {
		gtk.MainQuit()
	}
}

func main() {
	timeStart := time.Now()

	// Gentle SIGTERM handler thanks to reiki4040 https://gist.github.com/reiki4040/be3705f307d3cd136e85
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM {
				fmt.Println("SIGTERM received, bye bye!")
				gtk.MainQuit()
			}
		}
	}()

	// We don't want multiple instances. For better user experience (when nwgocc attached to a button or a key binding),
	// let's kill the running instance and exit.
	lockFilePath := fmt.Sprintf("%s/nwgocc.lock", tempDir())
	lockFile, err := singleinstance.CreateLockFile(lockFilePath)
	if err != nil {
		pid, err := readTextFile(lockFilePath)
		if err == nil {
			i, err := strconv.Atoi(pid)
			if err == nil {
				fmt.Println("Running instance found, sending SIGTERM and exiting...")
				syscall.Kill(i, syscall.SIGTERM)
			}
		}
		os.Exit(0)
	}
	defer lockFile.Close()

	flag.Parse()

	if *displayVersion {
		fmt.Printf("nwgocc version %s\n", version)
		os.Exit(0)
	}

	wayland = isWayland()
	fmt.Printf("Wayland: %t\n", wayland)

	setupDirs()

	// Load Preferences, Icons and Commands from ~/.local/share/nwgocc/preferences.json
	settings, _ = loadSettings()
	checkMissingSettings()

	// On `-d` check and print commands availability
	if *debug {
		checkCommands(settings.Commands)
	}

	// Load user-defined CustomRows and Buttons from ~/.config/config.json
	config, _ = loadConfig()
	fmt.Printf("Templates: '%s'\n", *configFile)

	// Empty means: gtk icons in use
	iconsDir = ""
	if settings.Preferences.IconSet == "light" {
		iconsDir = filepath.Join(dataDir(), "icons_light")
		fmt.Println("Icons: Custom light")
	} else if settings.Preferences.IconSet == "dark" {
		iconsDir = filepath.Join(dataDir(), "icons_dark")
		fmt.Println("Icons: Custom dark")
	} else {
		fmt.Println("Icons: GTK")
	}

	gtk.Init(nil)

	if settings.Preferences.CustomStyling {
		css := filepath.Join(configDir(), *cssFile)
		fmt.Printf("Style: '%s'\n", css)
		cssProvider, err := gtk.CssProviderNew()
		check(err)
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
	check(err)

	win.SetTitle("nwgocc: Control Center")
	if !wayland {
		err = win.SetIconFromFile("/usr/share/pixmaps/nwgocc.svg")
		if err != nil {
			win.SetIconName("nwgocc")
		}
	}
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

	cliCommands = loadCliCommands()
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

	var interfaceRow *gtk.EventBox
	if settings.Preferences.ShowInterfaceLine {
		interfaceRow = setupInterfaceRow()
		vBox.PackStart(interfaceRow, false, false, 4)
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

		if interfaceRow != nil {
			updateInterfaceRow()
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
