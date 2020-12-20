package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

// ConfigDir returns the .config dir path
func ConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return (fmt.Sprintf("%s/nwgcc", os.Getenv("XDG_CONFIG_HOME")))
	}
	return (fmt.Sprintf("%s/.config/nwgcc", os.Getenv("HOME")))
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
	GetBluetoothName   string `json:"get_bt_name"`
	GetBluetoothStatus string `json:"get_bt_status"`
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

// LoadCliCommands parses the cli_commands txt file and returns shell commands as []string slice
func LoadCliCommands() []string {
	path := fmt.Sprintf("%s/cli_commands", ConfigDir())
	bytes, err := ioutil.ReadFile(path)
	Check(err)
	lines := strings.Split(string(bytes), "\n")
	// trim whitespaces, remove commented out and empty lines
	var output []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") && line != "" {
			output = append(output, line)
		}
	}
	return output
}

// GetCliOutput returns output of each command as a string, separated with new lines, ready for use in cliLabel
func GetCliOutput(commands []string) string {
	var output []string
	for _, command := range commands {
		out, err := exec.Command("sh", "-c", command).Output()
		var o string
		if err == nil {
			o = string(out)
		} else {
			o = fmt.Sprintf("%s", err)
		}
		if len(o) > 38 {
			o = o[0:38] + "..."
		}
		o = strings.TrimSpace(o)
		output = append(output, o)
	}

	return string(strings.Join(output, "\n"))
}

// GetCommandOutput returns output of a CLI command with optional arguments
func GetCommandOutput(command string) string {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// CheckCommands checks external commands availability
func CheckCommands(commands Commands) {
	fmt.Println("Checking commands availability:")
	v := reflect.ValueOf(commands)
	values := make([]interface{}, v.NumField())

	m := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		cmd, _ := values[i].(string)
		cmd = strings.Split(cmd, " ")[0]
		available := GetCommandOutput(fmt.Sprintf("command -v %s ", cmd)) != ""
		if !KeyFound(m, cmd) {
			if available {
				m[cmd] = "available"
			} else {
				m[cmd] = "not found"
			}
		}
	}
	for key, value := range m {
		fmt.Printf("  '%s' %s\n", key, value)
	}
}

func isCommand(command string) bool {
	cmd := strings.Fields(command)[0]
	return GetCommandOutput(fmt.Sprintf("command -v %s ", cmd)) != ""
}

func btServiceEnabled() bool {
	if isCommand(settings.Commands.Systemctl) {
		return GetCommandOutput("systemctl is-enabled bluetooth.service") == "enabled"
	}
	return false
}
