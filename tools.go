package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
)

// Check ...
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// CreatePixbuf ...
func CreatePixbuf(iconsDir, icon string) *gdk.Pixbuf {
	iconPath := ""
	// just name given
	if !strings.HasPrefix(icon, "/") {
		iconPath = filepath.Join(iconsDir, fmt.Sprintf("%s.svg", icon))
	// full path given
	} else {
		iconPath = icon
	}
	pixbuf, err := gdk.PixbufNewFromFile(iconPath)
	if err != nil {
		pixbuf, err = gdk.PixbufNewFromFile(filepath.Join(DataDir(), "icons_light/icon-missing.svg"))
		Check(err)
	}

	return pixbuf
}
