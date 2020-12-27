package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var prefWindow *gtk.Window
var cliTextView *gtk.TextView

func setupPreferencesWindow() {
	builder, err := gtk.BuilderNewFromFile("preferences.glade")
	Check(err)

	obj, err := builder.GetObject("preferences_window")
	Check(err)

	if settings.Preferences.CustomStyling {
		css := filepath.Join(ConfigDir(), *customCSS)
		fmt.Printf("Style: %s\n", css)
		cssProvider, err := gtk.CssProviderNew()
		Check(err)
		err = cssProvider.LoadFromPath(css)
		if err != nil {
			fmt.Println(err)
		}
		screen, _ := gdk.ScreenGetDefault()
		gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
	}

	prefWindow, err := isWindow(obj)
	Check(err)

	// TextView to edit CLI Label command(s)
	cliTextView = setUpCliTextView(builder, "cli_textview")

	// Checkboxes to turn components on/off
	cbCliLabel := setUpCheckButton(builder, "checkbutton_cli_label", settings.Preferences.ShowCliLabel)
	cbCliLabel.Connect("toggled", func() {
		settings.Preferences.ShowCliLabel = cbCliLabel.GetActive()
	})

	cbBrightnessSlider := setUpCheckButton(builder, "checkbutton_brightness_slider", settings.Preferences.ShowBrightnessSlider)
	cbBrightnessSlider.Connect("toggled", func() {
		settings.Preferences.ShowBrightnessSlider = cbBrightnessSlider.GetActive()
	})

	cbVolumeSlider := setUpCheckButton(builder, "checkbutton_volume_slider", settings.Preferences.ShowVolumeSlider)
	cbVolumeSlider.Connect("toggled", func() {
		settings.Preferences.ShowVolumeSlider = cbVolumeSlider.GetActive()
	})

	cbPlayerctl := setUpCheckButton(builder, "checkbutton_playerctl", settings.Preferences.ShowPlayerctl)
	cbPlayerctl.Connect("toggled", func() {
		settings.Preferences.ShowPlayerctl = cbPlayerctl.GetActive()
	})

	cbUserLine := setUpCheckButton(builder, "checkbutton_user_info", settings.Preferences.ShowUserLine)
	cbUserLine.Connect("toggled", func() {
		settings.Preferences.ShowUserLine = cbUserLine.GetActive()
	})

	cbWifiLine := setUpCheckButton(builder, "checkbutton_wifi_status", settings.Preferences.ShowWifiLine)
	cbWifiLine.Connect("toggled", func() {
		settings.Preferences.ShowWifiLine = cbWifiLine.GetActive()
	})

	cbBtLine := setUpCheckButton(builder, "checkbutton_bluetooth_status", settings.Preferences.ShowBtLine)
	cbBtLine.Connect("toggled", func() {
		settings.Preferences.ShowBtLine = cbBtLine.GetActive()
	})

	cbBatteryLine := setUpCheckButton(builder, "checkbutton_battery_level", settings.Preferences.ShowBatteryLine)
	cbBatteryLine.Connect("toggled", func() {
		settings.Preferences.ShowBatteryLine = cbBatteryLine.GetActive()
	})

	cbUserRows := setUpCheckButton(builder, "checkbutton_user_rows", settings.Preferences.ShowUserRows)
	cbUserRows.Connect("toggled", func() {
		settings.Preferences.ShowUserRows = cbUserRows.GetActive()
	})

	cbUserButtons := setUpCheckButton(builder, "checkbutton_user_button", settings.Preferences.ShowUserButtons)
	cbUserButtons.Connect("toggled", func() {
		settings.Preferences.ShowUserButtons = cbUserButtons.GetActive()
	})

	// Buttons to edit user-defined commands assigned to built-in rows
	setButtonImage(builder, "cmd_btn_user", settings.Icons.ClickMe)
	setButtonImage(builder, "cmd_btn_wifi", settings.Icons.ClickMe)
	setButtonImage(builder, "cmd_btn_bt", settings.Icons.ClickMe)
	setButtonImage(builder, "cmd_btn_battery", settings.Icons.ClickMe)

	// Lower checkboxes for various boolean settings
	cbCustomStyling := setUpCheckButton(builder, "checkbutton_custom_css", settings.Preferences.CustomStyling)
	cbCustomStyling.Connect("toggled", func() {
		settings.Preferences.CustomStyling = cbCustomStyling.GetActive()
	})

	cbDontClose := setUpCheckButton(builder, "checkbutton_keep_open", settings.Preferences.DontClose)
	cbDontClose.Connect("toggled", func() {
		settings.Preferences.DontClose = cbDontClose.GetActive()
	})

	cbWindowDecorations := setUpCheckButton(builder, "checkbutton_window_decorations", settings.Preferences.WindowDecorations)
	cbWindowDecorations.Connect("toggled", func() {
		settings.Preferences.WindowDecorations = cbWindowDecorations.GetActive()
	})

	// ComboBox to select active icon set
	cbIconsSet := setUpIconsSetCombo(builder, "combo_box_icons")
	cbIconsSet.Connect("changed", func() {
		settings.Preferences.IconSet = cbIconsSet.GetActiveID()
	})

	sbIconSmall := setUpSpinbutton(builder, "spinbutton_small_icon",
		settings.Preferences.IconSizeSmall, 8, 64)
	sbIconSmall.Connect("value-changed", func() {
		settings.Preferences.IconSizeSmall = int(sbIconSmall.GetValue())
	})

	sbIconLarge := setUpSpinbutton(builder, "spinbutton_large_icon",
		settings.Preferences.IconSizeLarge, 8, 64)
	sbIconLarge.Connect("value-changed", func() {
		settings.Preferences.IconSizeLarge = int(sbIconLarge.GetValue())
	})

	sbRefreshCli := setUpSpinbutton(builder, "spinbutton_refresh_cli",
		settings.Preferences.RefreshCliSeconds, 0, 3600)
	sbRefreshCli.Connect("value-changed", func() {
		settings.Preferences.RefreshCliSeconds = int(sbRefreshCli.GetValue())
	})

	sbRefreshFastMillis := setUpSpinbutton(builder, "spinbutton_refresh_sliders",
		settings.Preferences.RefreshFastMillis, 0, 1000)
	sbRefreshFastMillis.Connect("value-changed", func() {
		settings.Preferences.RefreshFastMillis = int(sbRefreshFastMillis.GetValue())
	})

	sbRefreshSlowSeconds := setUpSpinbutton(builder, "spinbutton_refresh_battery",
		settings.Preferences.RefreshSlowSeconds, 0, 60)
	sbRefreshSlowSeconds.Connect("value-changed", func() {
		settings.Preferences.RefreshSlowSeconds = int(sbRefreshSlowSeconds.GetValue())
	})

	// bottom Buttons
	btnCancel := getButtonFromBuilder(builder, "btn_cancel")
	btnCancel.Connect("clicked", func() {
		prefWindow.Close()
	})

	btnApply := getButtonFromBuilder(builder, "btn_apply")
	btnApply.Connect("clicked", func() {
		saveCliCommands()
		err := SaveSettings()
		Check(err)
		gtk.MainQuit()
	})

	prefWindow.Show()
	prefWindow.Connect("key-release-event", handleEscape)
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func setUpCheckButton(builder *gtk.Builder, id string, active bool) *gtk.CheckButton {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	if cb, ok := obj.(*gtk.CheckButton); ok {
		cb.SetActive(active)
		return cb
	}
	return nil
}

func setUpSpinbutton(builder *gtk.Builder, id string, value int, min, max float64) *gtk.SpinButton {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	if sb, ok := obj.(*gtk.SpinButton); ok {
		sb.SetRange(min, max)
		sb.SetIncrements(1, 1)
		sb.SetValue(float64(value))
		return sb
	}
	return nil
}

func setButtonImage(builder *gtk.Builder, id string, icon string) {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return
	}
	if btn, ok := obj.(*gtk.Button); ok {
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		image, err := gtk.ImageNewFromPixbuf(pixbuf)
		Check(err)
		btn.SetImage(image)
	}
}

func setUpCliTextView(builder *gtk.Builder, id string) *gtk.TextView {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	if tv, ok := obj.(*gtk.TextView); ok {
		buffer, _ := tv.GetBuffer()
		buffer.SetText(strings.Join(cliCommands, "\n"))
		return tv
	}
	return nil
}

func setUpIconsSetCombo(builder *gtk.Builder, id string) *gtk.ComboBoxText {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	if cb, ok := obj.(*gtk.ComboBoxText); ok {
		cb.SetActiveID(settings.Preferences.IconSet)
		return cb
	}
	return nil
}

func getButtonFromBuilder(builder *gtk.Builder, id string) *gtk.Button {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	if btn, ok := obj.(*gtk.Button); ok {
		return btn
	}
	return nil
}

func handleEscape(window *gtk.Window, event *gdk.Event) {
	key := &gdk.EventKey{Event: event}
	if key.KeyVal() == gdk.KEY_Escape {
		window.Close()
	}
}

func saveCliCommands() {
	buffer, err := cliTextView.GetBuffer()
	Check(err)
	start := buffer.GetStartIter()
	end := buffer.GetEndIter()
	s, _ := buffer.GetText(start, end, true)
	saveCliFile(s)
}
