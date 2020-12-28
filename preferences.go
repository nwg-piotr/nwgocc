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
	builder, err := gtk.BuilderNewFromFile("/usr/share/nwgocc/preferences.glade")
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
	btnUser := setButtonImage(builder, "cmd_btn_user", settings.Icons.ClickMe)
	btnUser.Connect("clicked", func() {
		setupCmdDialog(&settings.Preferences.OnClickUser)
	})
	btnWifi := setButtonImage(builder, "cmd_btn_wifi", settings.Icons.ClickMe)
	btnWifi.Connect("clicked", func() {
		setupCmdDialog(&settings.Preferences.OnClickWifi)
	})
	btnBt := setButtonImage(builder, "cmd_btn_bt", settings.Icons.ClickMe)
	btnBt.Connect("clicked", func() {
		setupCmdDialog(&settings.Preferences.OnClickBluetooth)
	})
	btnBattery := setButtonImage(builder, "cmd_btn_battery", settings.Icons.ClickMe)
	btnBattery.Connect("clicked", func() {
		setupCmdDialog(&settings.Preferences.OnClickBattery)
	})

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
	btnUserRows := getButtonFromBuilder(builder, "btn_user_rows")
	btnUserRows.Connect("clicked", func() {
		setupTemplateEditionWindow(&config.CustomRows)
	})

	btnUserButtons := getButtonFromBuilder(builder, "btn_user_buttons")
	btnUserButtons.Connect("clicked", func() {
		setupTemplateEditionWindow(&config.Buttons)
	})

	btnCancel := getButtonFromBuilder(builder, "btn_cancel")
	btnCancel.Connect("clicked", func() {
		prefWindow.Close()
	})

	btnApply := getButtonFromBuilder(builder, "btn_apply")
	btnApply.Connect("clicked", func() {
		saveCliCommands()
		err := SaveSettings()
		Check(err)
		if configChanged {
			err := SaveConfig()
			Check(err)
		}
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

func setButtonImage(builder *gtk.Builder, id string, icon string) *gtk.Button {
	obj, err := builder.GetObject(id)
	if err != nil {
		log.Println(err)
		return nil
	}
	if btn, ok := obj.(*gtk.Button); ok {
		pixbuf := CreatePixbuf(iconsDir, icon, settings.Preferences.IconSizeSmall)
		image, err := gtk.ImageNewFromPixbuf(pixbuf)
		Check(err)
		btn.SetImage(image)
		return btn
	}
	return nil
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

func setupCmdDialog(command *string) {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	Check(err)

	win.SetTransientFor(prefWindow)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
	win.SetTitle("nwgcc: Edit command")
	win.SetProperty("name", "preferences")
	win.Connect("key-release-event", handleEscape)

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox.PackStart(hbox, true, true, 6)
	win.Add(vbox)

	entry, _ := gtk.EntryNew()
	entry.SetWidthChars(25)
	entry.SetText(*command)
	hbox.PackStart(entry, true, true, 3)

	btnCancel, _ := gtk.ButtonNew()
	btnCancel.SetLabel("Cancel")
	btnCancel.Connect("clicked", func() {
		win.Close()
	})
	hbox.PackStart(btnCancel, false, false, 3)

	btnApply, _ := gtk.ButtonNew()
	btnApply.SetLabel("Apply")
	btnApply.Connect("clicked", func() {
		*command, _ = entry.GetText()
		win.Close()
	})
	hbox.PackStart(btnApply, false, false, 3)

	win.ShowAll()
}

func setupTemplateEditionWindow(definitions interface{}) {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	win.SetTransientFor(prefWindow)
	win.SetModal(true)
	win.SetKeepAbove(true)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
	win.SetProperty("name", "preferences")
	win.Connect("key-release-event", handleEscape)

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox.PackStart(hbox, true, true, 20)

	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(10)
	grid.SetRowSpacing(10)
	hbox.PackStart(grid, true, true, 20)

	label, _ := gtk.LabelNew("Label")
	label.SetHAlign(gtk.ALIGN_START)
	grid.Attach(label, 0, 0, 1, 1)

	label, _ = gtk.LabelNew("Command")
	label.SetHAlign(gtk.ALIGN_START)
	grid.Attach(label, 1, 0, 1, 1)

	label, _ = gtk.LabelNew("Icon name or path")
	label.SetHAlign(gtk.ALIGN_START)
	grid.Attach(label, 2, 0, 1, 1)

	lastRow := 0

	switch definitions.(type) {
	case *[]CustomRow:
		win.SetTitle("nwgcc: Edit User Rows")
		for i, d := range *definitions.(*[]CustomRow) {
			entry, _ := gtk.EntryNew()
			entry.SetWidthChars(20)
			entry.SetText(d.Name)
			grid.Attach(entry, 0, i+1, 1, 1)

			entry, _ = gtk.EntryNew()
			entry.SetWidthChars(25)
			entry.SetText(d.Command)
			grid.Attach(entry, 1, i+1, 1, 1)

			iconEntry, _ := gtk.EntryNew()
			iconEntry.SetWidthChars(40)
			iconEntry.SetText(d.Icon)
			iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, CreatePixbuf(iconsDir, d.Icon, settings.Preferences.IconSizeSmall))
			iconEntry.Connect("changed", func() {
				s, _ := iconEntry.GetText()
				iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, CreatePixbuf(iconsDir, s, settings.Preferences.IconSizeSmall))
			})
			grid.Attach(iconEntry, 2, i+1, 1, 1)

			fcButton := setupFileChooserButton()
			fcButton.Connect("file-set", func() {
				iconEntry.SetText(fcButton.GetFilename())
			})
			grid.Attach(fcButton, 3, i+1, 1, 1)

			cb, _ := gtk.CheckButtonNewWithLabel("Delete")
			grid.Attach(cb, 4, i+1, 1, 1)

			lastRow++
		}
	case *[]Button:
		win.SetTitle("nwgcc: Edit User Buttons")
		for i, d := range *definitions.(*[]Button) {
			entry, _ := gtk.EntryNew()
			entry.SetWidthChars(20)
			entry.SetText(d.Name)
			grid.Attach(entry, 0, i+1, 1, 1)

			entry, _ = gtk.EntryNew()
			entry.SetWidthChars(25)
			entry.SetText(d.Command)
			grid.Attach(entry, 1, i+1, 1, 1)

			iconEntry, _ := gtk.EntryNew()
			iconEntry.SetWidthChars(40)
			iconEntry.SetText(d.Icon)
			iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, CreatePixbuf(iconsDir, d.Icon, settings.Preferences.IconSizeSmall))
			iconEntry.Connect("changed", func() {
				s, _ := iconEntry.GetText()
				iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, CreatePixbuf(iconsDir, s, settings.Preferences.IconSizeSmall))
			})
			grid.Attach(iconEntry, 2, i+1, 1, 1)

			fcButton := setupFileChooserButton()
			fcButton.Connect("file-set", func() {
				iconEntry.SetText(fcButton.GetFilename())
			})
			grid.Attach(fcButton, 3, i+1, 1, 1)

			cb, _ := gtk.CheckButtonNewWithLabel("Delete")
			grid.Attach(cb, 4, i+1, 1, 1)

			lastRow++
		}
	default:
		break
	}
	entry, _ := gtk.EntryNew()
	entry.SetWidthChars(20)
	entry.SetPlaceholderText("Enter new label")
	grid.Attach(entry, 0, lastRow+1, 1, 1)

	entry, _ = gtk.EntryNew()
	entry.SetWidthChars(25)
	entry.SetPlaceholderText("Enter new command")
	grid.Attach(entry, 1, lastRow+1, 1, 1)

	iconEntry, _ := gtk.EntryNew()
	iconEntry.SetWidthChars(40)
	iconEntry.SetPlaceholderText("Enter name or choose a file")
	iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, CreatePixbuf(iconsDir, "", settings.Preferences.IconSizeSmall))
	iconEntry.Connect("changed", func() {
		s, _ := iconEntry.GetText()
		iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, CreatePixbuf(iconsDir, s, settings.Preferences.IconSizeSmall))
	})
	grid.Attach(iconEntry, 2, lastRow+1, 1, 1)

	fcButton := setupFileChooserButton()
	fcButton.Connect("file-set", func() {
		iconEntry.SetText(fcButton.GetFilename())
	})
	grid.Attach(fcButton, 3, lastRow+1, 1, 1)

	btn, _ := gtk.ButtonNew()
	btn.SetLabel("Cancel")
	btn.Connect("clicked", func() {
		win.Close()
	})
	grid.Attach(btn, 3, lastRow+2, 1, 1)

	btn, _ = gtk.ButtonNew()
	btn.SetLabel("Apply")

	switch definitions.(type) {
	case *[]CustomRow:
		btn.Connect("clicked", func() {
			var cRows []CustomRow
			var cRow CustomRow
			for row := 1; row < lastRow+1; row++ {
				field, _ := grid.GetChildAt(4, row)
				delete := field.(*gtk.CheckButton).GetActive()
				if !delete {
					field, _ := grid.GetChildAt(0, row)
					text, _ := field.(*gtk.Entry).GetText()
					cRow.Name = text

					field, _ = grid.GetChildAt(1, row)
					text, _ = field.(*gtk.Entry).GetText()
					cRow.Command = text

					field, _ = grid.GetChildAt(2, row)
					text, _ = field.(*gtk.Entry).GetText()
					cRow.Icon = text

					cRows = append(cRows, cRow)
				}
			}
			// We only add a new row if at least name given
			field, _ := grid.GetChildAt(0, lastRow+1)
			text, _ := field.(*gtk.Entry).GetText()
			if text != "" {
				var newRow CustomRow
				newRow.Name = text

				field, _ = grid.GetChildAt(1, lastRow+1)
				text, _ = field.(*gtk.Entry).GetText()
				newRow.Command = text

				field, _ = grid.GetChildAt(2, lastRow+1)
				text, _ = field.(*gtk.Entry).GetText()
				newRow.Icon = text

				cRows = append(cRows, newRow)
			}

			config.CustomRows = cRows
			configChanged = true
			win.Close()
		})
	case *[]Button:
		btn.Connect("clicked", func() {
			var cBtns []Button
			var cBtn Button
			for row := 1; row < lastRow+1; row++ {
				field, _ := grid.GetChildAt(4, row)
				delete := field.(*gtk.CheckButton).GetActive()
				if !delete {
					field, _ := grid.GetChildAt(0, row)
					text, _ := field.(*gtk.Entry).GetText()
					cBtn.Name = text

					field, _ = grid.GetChildAt(1, row)
					text, _ = field.(*gtk.Entry).GetText()
					cBtn.Command = text

					field, _ = grid.GetChildAt(2, row)
					text, _ = field.(*gtk.Entry).GetText()
					cBtn.Icon = text

					cBtns = append(cBtns, cBtn)
				}
			}
			// We only add a new button if at least name given
			field, _ := grid.GetChildAt(0, lastRow+1)
			text, _ := field.(*gtk.Entry).GetText()
			if text != "" {
				var newBtn Button
				newBtn.Name = text

				field, _ = grid.GetChildAt(1, lastRow+1)
				text, _ = field.(*gtk.Entry).GetText()
				newBtn.Command = text

				field, _ = grid.GetChildAt(2, lastRow+1)
				text, _ = field.(*gtk.Entry).GetText()
				newBtn.Icon = text

				cBtns = append(cBtns, newBtn)
			}

			config.Buttons = cBtns
			configChanged = true
			win.Close()
		})
	}
	grid.Attach(btn, 4, lastRow+2, 1, 1)

	// Clear selection on 1st Entry
	field, _ := grid.GetChildAt(0, 1)
	field.(*gtk.Entry).SelectRegion(0, 0)

	field, _ = grid.GetChildAt(0, lastRow+1)
	field.(*gtk.Entry).GrabFocus()

	win.Add(vbox)

	win.ShowAll()
}

func setupFileChooserButton() *gtk.FileChooserButton {
	fcBtn, _ := gtk.FileChooserButtonNew("Choose", gtk.FILE_CHOOSER_ACTION_OPEN)
	filter, _ := gtk.FileFilterNew()
	filter.AddPattern("*.svg")
	filter.AddPattern("*.SVG")
	filter.AddPattern("*.png")
	filter.AddPattern("*.PNG")
	fcBtn.AddFilter(filter)

	return fcBtn
}
