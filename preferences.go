package main

import (
	"errors"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func setupPreferencesWindow() {
	builder, err := gtk.BuilderNewFromFile("preferences.glade")
	Check(err)

	obj, err := builder.GetObject("preferences_window")
	Check(err)

	win, err := isWindow(obj)
	Check(err)

	setCliLabelContent(builder, "cli_textview")

	setCheckButtonState(builder, "checkbutton_cli_label", settings.Preferences.ShowCliLabel)
	setCheckButtonState(builder, "checkbutton_brightness_slider", settings.Preferences.ShowBrightnessSlider)
	setCheckButtonState(builder, "checkbutton_volume_slider", settings.Preferences.ShowVolumeSlider)
	setCheckButtonState(builder, "checkbutton_playerctl", settings.Preferences.ShowPlayerctl)
	setCheckButtonState(builder, "checkbutton_user_info", settings.Preferences.ShowUserLine)
	setCheckButtonState(builder, "checkbutton_wifi_status", settings.Preferences.ShowWifiLine)
	setCheckButtonState(builder, "checkbutton_bluetooth_status", settings.Preferences.ShowBtLine)
	setCheckButtonState(builder, "checkbutton_battery_level", settings.Preferences.ShowBatteryLine)
	setCheckButtonState(builder, "checkbutton_user_rows", settings.Preferences.ShowUserRows)
	setCheckButtonState(builder, "checkbutton_user_button", settings.Preferences.ShowUserButtons)

	setButtonImage(builder, "cmd_btn_user", settings.Icons.ClickMe)
	setButtonImage(builder, "cmd_btn_wifi", settings.Icons.ClickMe)
	setButtonImage(builder, "cmd_btn_bt", settings.Icons.ClickMe)
	setButtonImage(builder, "cmd_btn_battery", settings.Icons.ClickMe)

	setCheckButtonState(builder, "checkbutton_custom_css", settings.Preferences.CustomStyling)
	setCheckButtonState(builder, "checkbutton_keep_open", settings.Preferences.DontClose)
	setCheckButtonState(builder, "checkbutton_window_decorations", settings.Preferences.WindowDecorations)

	setIconsSetComboState(builder, "combo_box_icons")

	setSpinbuttonValue(builder, "spinbutton_small_icon", settings.Preferences.IconSizeSmall, 8, 64)
	setSpinbuttonValue(builder, "spinbutton_large_icon", settings.Preferences.IconSizeLarge, 8, 64)

	setSpinbuttonValue(builder, "spinbutton_refresh_cli", settings.Preferences.RefreshCliSeconds, 0, 3600)
	setSpinbuttonValue(builder, "spinbutton_refresh_sliders", settings.Preferences.RefreshFastMillis, 0, 1000)
	setSpinbuttonValue(builder, "spinbutton_refresh_battery", settings.Preferences.RefreshSlowSeconds, 0, 60)

	win.Show()
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func setCheckButtonState(builder *gtk.Builder, id string, active bool) {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return
	}
	if cb, ok := obj.(*gtk.CheckButton); ok {
		cb.SetActive(active)
	}
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

func setCliLabelContent(builder *gtk.Builder, id string) {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return
	}
	if tv, ok := obj.(*gtk.TextView); ok {
		buffer, _ := tv.GetBuffer()

		buffer.SetText(strings.Join(cliCommands, "\n"))
	}
}

func setIconsSetComboState(builder *gtk.Builder, id string) {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return
	}
	if cb, ok := obj.(*gtk.ComboBoxText); ok {
		cb.SetActiveID(settings.Preferences.IconSet)
	}
}

func setSpinbuttonValue(builder *gtk.Builder, id string, value int, min, max float64) {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return
	}
	if sb, ok := obj.(*gtk.SpinButton); ok {
		sb.SetRange(min, max)
		sb.SetIncrements(1, 1)
		sb.SetValue(float64(value))
	}
}
