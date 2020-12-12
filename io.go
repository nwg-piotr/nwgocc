package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ConfigDir returns the .config dir path
func ConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return (fmt.Sprintf("%s/nwgcc", os.Getenv("XDG_CONFIG_HOME")))
	}
	return (fmt.Sprintf("%s/.config/nwgocc", os.Getenv("HOME")))
}

// DataDir returns data directory path
func DataDir() string {
	return (fmt.Sprintf("%s/.local/share/nwgcc", os.Getenv("HOME")))
}

// CustomRow contains fields of a single user-defined row
type CustomRow struct {
	Name    string `json:"name"`
	Command string `json:"cmd"`
	Icon    string `json:"icon"`
}

// Button contains fields of a single user-defined button
type Button struct {
	Name    string `json:"name"`
	Command string `json:"cmd"`
	Icon    string `json:"icon"`
}

// Configuration stores all the user-defined content: custom rows and buttons
type Configuration struct {
	CustomRows []CustomRow `json:"custom_rows"`
	Buttons    []Button    `json:"buttons"`
}

// Preferences store program settings
type Preferences struct {
	IconSet              string `json:"icon_set"`
	CustomStyling        bool   `json:"custom_styling"`
	DontClose            bool   `json:"dont_close"`
	WindowDecorations    bool   `json:"window_decorations"`
	ShowCliLabel         bool   `json:"show_cli_label"`
	ShowBrightnessSlider bool   `json:"show_brightness_slider"`
	ShowVolumeSlider     bool   `json:"show_volume_slider"`
	ShowPlayerctl        bool   `json:"show_playerctl"`
	ShowUserLine         bool   `json:"show_user_line"`
	ShowWifiLine         bool   `json:"show_wifi_line"`
	ShowBtLine           bool   `json:"show_bt_line"`
	ShowBatteryLine      bool   `json:"show_battery_line"`
	ShowUserRows         bool   `json:"show_user_rows"`
	ShowUserButtons      bool   `json:"show_user_buttons"`
	IconSizeSmall        int    `json:"icon_size_small"`
	IconSizeLarge        int    `json:"icon_size_large"`
	RefreshFastMillis    int    `json:"refresh_fast_millis"`
	RefreshSlowSeconds   int    `json:"refresh_slow_seconds"`
	RefreshCliSeconds    int    `json:"refresh_cli_seconds"`
	OnClickUser          string `json:"on-click-user"`
	OnClickWifi          string `json:"on-click-wifi"`
	OnClickBluetooth     string `json:"on-click-bluetooth"`
	OnClickBattery       string `json:"on-click-battery"`
}

// Icons store icon definitions
type Icons struct {
	BatteryEmpty       string `json:"battery-empty"`
	BatteryLow         string `json:"battery-low"`
	BatteryGood        string `json:"battery-good"`
	BatteryFull        string `json:"battery-full"`
	User               string `json:"user"`
	WifiOn             string `json:"wifi-on"`
	WifiOff            string `json:"wifi-off"`
	BrightnessLow      string `json:"brightness-low"`
	BrightnessMedium   string `json:"brightness-medium"`
	BrightnessHigh     string `json:"brightness-high"`
	BtOn               string `json:"bt-on"`
	BtOff              string `json:"bt-off"`
	VolumeLow          string `json:"volume-low"`
	VolumeMedium       string `json:"volume-medium"`
	VolumeHigh         string `json:"volume-high"`
	VolumeMuted        string `json:"volume-muted"`
	MediaPlaybackPause string `json:"media-playback-pause"`
	MediaPlaybackStart string `json:"media-playback-start"`
	MediaPlaybackStop  string `json:"media-playback-stop"`
	MediaSeekBackward  string `json:"media-seek-backward"`
	MediaSeekForward   string `json:"media-seek-forward"`
	MediaSkipBackward  string `json:"media-skip-backward"`
	MediaSkipForward   string `json:"media-skip-forward"`
	ClickMe            string `json:"click-me"`
}

// Commands store external commands
type Commands struct {
	GetBattery         string `json:"get_battery"`
	GetBatteryAlt      string `json:"get_battery_alt"`
	GetBluetoothName   string `json:"get_bluetooth_name"`
	GetBluetoothStatus string `json:"get_bluetooth_status"`
	GetBrightness      string `json:"get_brightness"`
	GetHost            string `json:"get_host"`
	GetSsid            string `json:"get_ssid"`
	GetUser            string `json:"get_user"`
	GetVolumeAlt       string `json:"get_volume_alt"`
	SetBrightness      string `json:"set_brightness"`
	SetVolumeAlt       string `json:"set_volume_alt"`
	Systemctl          string `json:"systemctl"`
	Playerctl          string `json:"playerctl"`
}

// Settings store user preferecnces, icon definitions and external commands
type Settings struct {
	Preferences Preferences `json:"preferences"`
	Icons       Icons       `json:"icons"`
	Commands    Commands    `json:"commands"`
}

// LoadConfig parses the config.json file and returns Configuration instance
func LoadConfig() (Configuration, error) {
	path := fmt.Sprintf("%s/config.json", ConfigDir())
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Configuration{}, err
	}

	var c Configuration
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}

	return c, nil
}

// SaveConfig saves current Configuration to a json file
func SaveConfig(c Configuration, path string) error {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

// LoadSettings parses the preferences.json file and returns Settings instance
func LoadSettings() (Settings, error) {
	path := fmt.Sprintf("%s/preferences.json", DataDir())
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}

	var s Settings
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return Settings{}, err
	}

	return s, nil
}
