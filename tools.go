package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// Check ...
func Check(e error) {
	if e != nil {
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
		pixbuf, err = gdk.PixbufNewFromFileAtSize(filepath.Join(DataDir(),
			"icons_light/icon-missing.svg"), size, size)
		Check(err)
	}
	return pixbuf
}
