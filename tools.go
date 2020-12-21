package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Check ...
func Check(e error) {
	if e != nil {
		log.Println(e)
		panic(e)
	}
}

// CreatePixbuf ...
func CreatePixbuf(iconsDir, icon string, size int) *gdk.Pixbuf {
	// full path given
	iconPath := ""
	if strings.HasPrefix(icon, "/") {
		iconPath = icon
		pixbuf, err := gdk.PixbufNewFromFileAtSize(iconPath, size, size)
		if err != nil {
			fmt.Println(err)
			pixbuf, err = gdk.PixbufNewFromFileAtSize(filepath.Join(DataDir(),
				"icons_light/icon-missing.svg"), size, size)
			Check(err)
		}
		return pixbuf
	}

	// gtk icons in use - just name given
	if iconsDir == "" {
		iconTheme, _ := gtk.IconThemeGetDefault()
		pixbuf, _ := iconTheme.LoadIcon(icon, size, gtk.ICON_LOOKUP_FORCE_SIZE)

		return pixbuf
	}

	// just name given, and we don't use gtk icons
	iconPath = filepath.Join(iconsDir, fmt.Sprintf("%s.svg", icon))
	pixbuf, err := gdk.PixbufNewFromFileAtSize(iconPath, size, size)
	if err != nil {
		iconTheme, _ := gtk.IconThemeGetDefault()
		pixbuf, err := iconTheme.LoadIcon(icon, size, gtk.ICON_LOOKUP_FORCE_SIZE)

		if err != nil {
			pixbuf, err = gdk.PixbufNewFromFileAtSize(filepath.Join(DataDir(),
				"icons_light/icon-missing.svg"), size, size)
			Check(err)
		}
		return pixbuf

	}
	return pixbuf
}

// LaunchCommand starts external command and quits
func LaunchCommand(command string) {
	elements := strings.Split(command, " ")
	cmd := exec.Command(elements[0], elements[1:]...)
	go cmd.Run()
	if !settings.Preferences.DontClose {
		glib.TimeoutAdd(uint(100), func() bool {
			gtk.MainQuit()
			return false
		})

	}
}

// KeyFound checks if key in map[string]string
func KeyFound(m map[string]string, key string) bool {
	for k := range m {
		if k == key {
			return true
		}
	}
	return false
}

func getBattery(command string) (string, int) {
	msg := ""
	perc := 0
	if strings.Fields(command)[0] == "upower" {
		bat := strings.Split(GetCommandOutput(command), "\n")
		var state, time, percentage string
		for _, line := range bat {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "time to empty") {
				strings.Replace(line, "time to empty", "time_to_empty", 0)
			}
			parts := strings.Fields(line)
			for i, l := range parts {
				if strings.Contains(l, "state:") {
					state = parts[i+1]
				}
				if strings.Contains(l, "time_to_empty") {
					time = parts[i+1]
				}
				if strings.Contains(l, "percentage") {
					pl := len(parts[i+1])
					percentage = parts[i+1][:pl-1]
					p, err := strconv.Atoi(percentage)
					if err == nil {
						perc = p
					}
				}
			}
		}
		msg = fmt.Sprintf("%d%% %s %s", perc, state, time)

	} else if strings.Fields(command)[0] == "acpi" {
		bat := strings.Fields(GetCommandOutput(command))
		msg = strings.Join(bat[2:], " ")
		pl := len(bat[3])
		percentage := bat[3][:pl-2]
		p, err := strconv.Atoi(percentage)
		if err == nil {
			perc = p
		}
	}

	return msg, perc
}

func getBrightness() float64 {
	brightness := 0.0
	output := GetCommandOutput(settings.Commands.GetBrightness)
	bri, e := strconv.ParseFloat(output, 64)
	if e == nil {
		brightness = math.Round(bri)
	}

	return brightness
}

func setBrightness(value int) {
	cmd := exec.Command("light", "-S", fmt.Sprint(value))
	cmd.Run()
}
