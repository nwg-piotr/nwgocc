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
	check(err)

	obj, err := builder.GetObject("preferences_window")
	check(err)

	if settings.Preferences.CustomStyling {
		css := filepath.Join(configDir(), *customCSS)
		fmt.Printf("Style: %s\n", css)
		cssProvider, err := gtk.CssProviderNew()
		check(err)
		err = cssProvider.LoadFromPath(css)
		if err != nil {
			fmt.Println(err)
		}
		screen, _ := gdk.ScreenGetDefault()
		gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
	}

	prefWindow, err := isWindow(obj)
	check(err)

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

	btnIcons := getButtonFromBuilder(builder, "btn_icons")
	btnIcons.Connect("clicked", func() {
		setupIconsEditionWindow()
	})

	btnCancel := getButtonFromBuilder(builder, "btn_cancel")
	btnCancel.Connect("clicked", func() {
		prefWindow.Close()
	})

	btnApply := getButtonFromBuilder(builder, "btn_apply")
	btnApply.Connect("clicked", func() {
		saveCliCommands()
		err := saveSettings()
		check(err)
		if configChanged {
			err := saveConfig()
			check(err)
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
		pixbuf := createPixbuf(icon, settings.Preferences.IconSizeSmall)
		image, err := gtk.ImageNewFromPixbuf(pixbuf)
		check(err)
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
	check(err)
	start := buffer.GetStartIter()
	end := buffer.GetEndIter()
	s, _ := buffer.GetText(start, end, true)
	saveCliFile(s)
}

func setupCmdDialog(command *string) {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	check(err)

	win.SetTransientFor(prefWindow)
	win.SetModal(true)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
	win.SetTitle("nwgocc: Edit command")
	win.SetProperty("name", "preferences")
	win.Connect("key-release-event", handleEscape)

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox.PackStart(hbox, true, true, 6)
	win.Add(vbox)

	entry, _ := gtk.EntryNew()
	entry.SetProperty("name", "edit-field")
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
		win.SetTitle("nwgocc: Edit User Rows")
		for i, d := range *definitions.(*[]CustomRow) {
			entry, _ := gtk.EntryNew()
			entry.SetProperty("name", "edit-field")
			entry.SetWidthChars(20)
			entry.SetText(d.Name)
			grid.Attach(entry, 0, i+1, 1, 1)

			entry, _ = gtk.EntryNew()
			entry.SetProperty("name", "edit-field")
			entry.SetWidthChars(25)
			entry.SetText(d.Command)
			grid.Attach(entry, 1, i+1, 1, 1)

			iconEntry, _ := gtk.EntryNew()
			iconEntry.SetProperty("name", "edit-field")
			iconEntry.SetWidthChars(40)
			iconEntry.SetText(d.Icon)
			iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf(d.Icon, settings.Preferences.IconSizeSmall))
			iconEntry.Connect("changed", func() {
				s, _ := iconEntry.GetText()
				iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf(s, settings.Preferences.IconSizeSmall))
			})
			grid.Attach(iconEntry, 2, i+1, 1, 1)

			fcButton := setupFCButton(iconEntry)
			grid.Attach(fcButton, 3, i+1, 1, 1)

			cb, _ := gtk.CheckButtonNewWithLabel("Delete")
			grid.Attach(cb, 4, i+1, 1, 1)

			lastRow++
		}
	case *[]Button:
		win.SetTitle("nwgocc: Edit User Buttons")
		for i, d := range *definitions.(*[]Button) {
			entry, _ := gtk.EntryNew()
			entry.SetProperty("name", "edit-field")
			entry.SetWidthChars(20)
			entry.SetText(d.Name)
			grid.Attach(entry, 0, i+1, 1, 1)

			entry, _ = gtk.EntryNew()
			entry.SetProperty("name", "edit-field")
			entry.SetWidthChars(25)
			entry.SetText(d.Command)
			grid.Attach(entry, 1, i+1, 1, 1)

			iconEntry, _ := gtk.EntryNew()
			iconEntry.SetProperty("name", "edit-field")
			iconEntry.SetWidthChars(40)
			iconEntry.SetText(d.Icon)
			iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf(d.Icon, settings.Preferences.IconSizeSmall))
			iconEntry.Connect("changed", func() {
				s, _ := iconEntry.GetText()
				iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf(s, settings.Preferences.IconSizeSmall))
			})
			grid.Attach(iconEntry, 2, i+1, 1, 1)

			fcButton := setupFCButton(iconEntry)
			grid.Attach(fcButton, 3, i+1, 1, 1)

			cb, _ := gtk.CheckButtonNewWithLabel("Delete")
			grid.Attach(cb, 4, i+1, 1, 1)

			lastRow++
		}
	default:
		break
	}
	entry, _ := gtk.EntryNew()
	entry.SetProperty("name", "edit-field")
	entry.SetWidthChars(20)
	entry.SetPlaceholderText("Enter new label")
	grid.Attach(entry, 0, lastRow+1, 1, 1)

	entry, _ = gtk.EntryNew()
	entry.SetProperty("name", "edit-field")
	entry.SetWidthChars(25)
	entry.SetPlaceholderText("Enter new command")
	grid.Attach(entry, 1, lastRow+1, 1, 1)

	iconEntry, _ := gtk.EntryNew()
	iconEntry.SetProperty("name", "edit-field")
	iconEntry.SetWidthChars(40)
	iconEntry.SetPlaceholderText("Enter name or choose a file")
	iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf("", settings.Preferences.IconSizeSmall))
	iconEntry.Connect("changed", func() {
		s, _ := iconEntry.GetText()
		iconEntry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf(s, settings.Preferences.IconSizeSmall))
	})
	grid.Attach(iconEntry, 2, lastRow+1, 1, 1)

	fcButton := setupFCButton(iconEntry)
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

func setupFCButton(entry *gtk.Entry) *gtk.Button {
	btn, _ := gtk.ButtonNew()
	imgOpen, _ := gtk.ImageNewFromPixbuf(createPixbuf("document-open-symbolic", settings.Preferences.IconSizeSmall))
	btn.SetImage(imgOpen)
	btn.Connect("clicked", func() {
		dlg, _ := gtk.FileChooserDialogNewWith2Buttons(
			"Choose an image", nil, gtk.FILE_CHOOSER_ACTION_OPEN,
			"Open", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL,
		)
		dlg.SetDefaultResponse(gtk.RESPONSE_OK)
		filter, _ := gtk.FileFilterNew()
		filter.SetName("images")
		filter.AddMimeType("image/png")
		filter.AddMimeType("image/svg")
		filter.AddPattern("*.png")
		filter.AddPattern("*.svg")
		dlg.SetFilter(filter)
		response := dlg.Run()
		if response == gtk.RESPONSE_OK {
			filename := dlg.GetFilename()
			entry.SetText(filename)
		}
		dlg.Destroy()
	})
	return btn
}

func setupIconsEditionWindow() {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	win.SetTransientFor(prefWindow)
	win.SetModal(true)
	win.SetKeepAbove(true)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
	win.SetProperty("name", "preferences")
	win.SetDefaultSize(300, 720)
	win.Connect("key-release-event", handleEscape)

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vbox.PackStart(hbox, true, true, 10)

	scrolledWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	scrolledWindow.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(10)
	grid.SetRowSpacing(3)
	scrolledWindow.Add(grid)
	hbox.PackStart(scrolledWindow, true, true, 20)

	label, _ := gtk.LabelNew("Icon")
	label.SetHAlign(gtk.ALIGN_START)
	grid.Attach(label, 0, 0, 1, 1)

	label, _ = gtk.LabelNew("Name or path")
	label.SetHAlign(gtk.ALIGN_START)
	grid.Attach(label, 1, 0, 1, 1)

	win.SetTitle("nwgocc: Edit Icons Dictionary")

	lbl, entry, fcBtn := iconEditionFields("battery-empty", settings.Icons.BatteryEmpty)
	grid.Attach(lbl, 0, 1, 1, 1)
	grid.Attach(entry, 1, 1, 1, 1)
	grid.Attach(fcBtn, 2, 1, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("battery-low", settings.Icons.BatteryLow)
	grid.Attach(lbl, 0, 2, 1, 1)
	grid.Attach(entry, 1, 2, 1, 1)
	grid.Attach(fcBtn, 2, 2, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("battery-good", settings.Icons.BatteryGood)
	grid.Attach(lbl, 0, 3, 1, 1)
	grid.Attach(entry, 1, 3, 1, 1)
	grid.Attach(fcBtn, 2, 3, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("battery-full", settings.Icons.BatteryFull)
	grid.Attach(lbl, 0, 4, 1, 1)
	grid.Attach(entry, 1, 4, 1, 1)
	grid.Attach(fcBtn, 2, 4, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("user", settings.Icons.User)
	grid.Attach(lbl, 0, 5, 1, 1)
	grid.Attach(entry, 1, 5, 1, 1)
	grid.Attach(fcBtn, 2, 5, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("wifi-on", settings.Icons.WifiOn)
	grid.Attach(lbl, 0, 6, 1, 1)
	grid.Attach(entry, 1, 6, 1, 1)
	grid.Attach(fcBtn, 2, 6, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("wifi-off", settings.Icons.WifiOff)
	grid.Attach(lbl, 0, 7, 1, 1)
	grid.Attach(entry, 1, 7, 1, 1)
	grid.Attach(fcBtn, 2, 7, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("brightness-low", settings.Icons.BrightnessLow)
	grid.Attach(lbl, 0, 8, 1, 1)
	grid.Attach(entry, 1, 8, 1, 1)
	grid.Attach(fcBtn, 2, 8, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("brightness-medium", settings.Icons.BrightnessMedium)
	grid.Attach(lbl, 0, 9, 1, 1)
	grid.Attach(entry, 1, 9, 1, 1)
	grid.Attach(fcBtn, 2, 9, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("brightness-high", settings.Icons.BrightnessHigh)
	grid.Attach(lbl, 0, 10, 1, 1)
	grid.Attach(entry, 1, 10, 1, 1)
	grid.Attach(fcBtn, 2, 10, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("bt-on", settings.Icons.BtOn)
	grid.Attach(lbl, 0, 11, 1, 1)
	grid.Attach(entry, 1, 11, 1, 1)
	grid.Attach(fcBtn, 2, 11, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("bt-off", settings.Icons.BtOff)
	grid.Attach(lbl, 0, 12, 1, 1)
	grid.Attach(entry, 1, 12, 1, 1)
	grid.Attach(fcBtn, 2, 12, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("volume-low", settings.Icons.VolumeLow)
	grid.Attach(lbl, 0, 13, 1, 1)
	grid.Attach(entry, 1, 13, 1, 1)
	grid.Attach(fcBtn, 2, 13, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("volume-medium", settings.Icons.VolumeMedium)
	grid.Attach(lbl, 0, 14, 1, 1)
	grid.Attach(entry, 1, 14, 1, 1)
	grid.Attach(fcBtn, 2, 14, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("volume-high", settings.Icons.VolumeHigh)
	grid.Attach(lbl, 0, 15, 1, 1)
	grid.Attach(entry, 1, 15, 1, 1)
	grid.Attach(fcBtn, 2, 15, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("volume-muted", settings.Icons.VolumeMuted)
	grid.Attach(lbl, 0, 16, 1, 1)
	grid.Attach(entry, 1, 16, 1, 1)
	grid.Attach(fcBtn, 2, 16, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("media-playback-pause", settings.Icons.MediaPlaybackPause)
	grid.Attach(lbl, 0, 17, 1, 1)
	grid.Attach(entry, 1, 17, 1, 1)
	grid.Attach(fcBtn, 2, 17, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("media-playback-start", settings.Icons.MediaPlaybackStart)
	grid.Attach(lbl, 0, 18, 1, 1)
	grid.Attach(entry, 1, 18, 1, 1)
	grid.Attach(fcBtn, 2, 18, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("media-playback-stop", settings.Icons.MediaPlaybackStop)
	grid.Attach(lbl, 0, 19, 1, 1)
	grid.Attach(entry, 1, 19, 1, 1)
	grid.Attach(fcBtn, 2, 19, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("media-skip-backward", settings.Icons.MediaSkipBackward)
	grid.Attach(lbl, 0, 20, 1, 1)
	grid.Attach(entry, 1, 20, 1, 1)
	grid.Attach(fcBtn, 2, 20, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("media-skip-forward", settings.Icons.MediaSkipForward)
	grid.Attach(lbl, 0, 21, 1, 1)
	grid.Attach(entry, 1, 21, 1, 1)
	grid.Attach(fcBtn, 2, 21, 1, 1)

	lbl, entry, fcBtn = iconEditionFields("click-me", settings.Icons.ClickMe)
	grid.Attach(lbl, 0, 22, 1, 1)
	grid.Attach(entry, 1, 22, 1, 1)
	grid.Attach(fcBtn, 2, 22, 1, 1)

	hbox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	btn, _ := gtk.ButtonNew()
	btn.SetLabel("Apply")
	btn.Connect("clicked", func() {

		settings.Icons.BatteryEmpty = getTextFromGrid(grid, 1, 1)
		settings.Icons.BatteryLow = getTextFromGrid(grid, 1, 2)
		settings.Icons.BatteryGood = getTextFromGrid(grid, 1, 3)
		settings.Icons.BatteryFull = getTextFromGrid(grid, 1, 4)
		settings.Icons.User = getTextFromGrid(grid, 1, 5)
		settings.Icons.WifiOn = getTextFromGrid(grid, 1, 6)
		settings.Icons.WifiOff = getTextFromGrid(grid, 1, 7)
		settings.Icons.BrightnessLow = getTextFromGrid(grid, 1, 8)
		settings.Icons.BrightnessMedium = getTextFromGrid(grid, 1, 9)
		settings.Icons.BrightnessHigh = getTextFromGrid(grid, 1, 10)
		settings.Icons.BtOn = getTextFromGrid(grid, 1, 11)
		settings.Icons.BtOff = getTextFromGrid(grid, 1, 12)
		settings.Icons.VolumeLow = getTextFromGrid(grid, 1, 13)
		settings.Icons.VolumeMedium = getTextFromGrid(grid, 1, 14)
		settings.Icons.VolumeHigh = getTextFromGrid(grid, 1, 15)
		settings.Icons.VolumeMuted = getTextFromGrid(grid, 1, 16)
		settings.Icons.MediaPlaybackPause = getTextFromGrid(grid, 1, 17)
		settings.Icons.MediaPlaybackStart = getTextFromGrid(grid, 1, 18)
		settings.Icons.MediaPlaybackStop = getTextFromGrid(grid, 1, 19)
		settings.Icons.MediaSkipBackward = getTextFromGrid(grid, 1, 20)
		settings.Icons.MediaSkipForward = getTextFromGrid(grid, 1, 21)
		settings.Icons.ClickMe = getTextFromGrid(grid, 1, 22)

		win.Close()
	})

	hbox.PackEnd(btn, false, false, 20)
	vbox.PackStart(hbox, false, true, 10)

	// Clear selection on 1st Entry
	field, err := grid.GetChildAt(1, 1)
	if err == nil {
		field.(*gtk.Entry).SelectRegion(0, 0)
	}

	btn, _ = gtk.ButtonNew()
	btn.SetLabel("Cancel")
	btn.Connect("clicked", func() {
		defer win.Close()
	})
	hbox.PackEnd(btn, false, false, 0)

	btn.GrabFocus()

	win.Add(vbox)

	win.ShowAll()
}

func getTextFromGrid(grid *gtk.Grid, col, row int) string {
	text := ""
	field, err := grid.GetChildAt(col, row)
	if err == nil {
		if entry, ok := field.(*gtk.Entry); ok {
			text, _ = entry.GetText()
		}
	}
	return text
}

func iconEditionFields(name, value string) (*gtk.Label, *gtk.Entry, *gtk.Button) {
	label, err := gtk.LabelNew(name)
	check(err)
	label.SetHAlign(gtk.ALIGN_START)

	entry, err := gtk.EntryNew()
	check(err)
	entry.SetProperty("name", "edit-field")
	entry.SetWidthChars(40)
	entry.SetText(value)
	entry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf(value, settings.Preferences.IconSizeSmall))
	entry.Connect("changed", func() {
		s, err := entry.GetText()
		check(err)
		entry.SetIconFromPixbuf(gtk.ENTRY_ICON_PRIMARY, createPixbuf(s, settings.Preferences.IconSizeSmall))
	})

	button := setupFCButton(entry)

	return label, entry, button
}
