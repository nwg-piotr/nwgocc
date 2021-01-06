package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func configDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return (fmt.Sprintf("%s/nwgocc", os.Getenv("XDG_CONFIG_HOME")))
	}
	return (fmt.Sprintf("%s/.config/nwgocc", os.Getenv("HOME")))
}

func dataDir() string {
	return (fmt.Sprintf("%s/.local/share/nwgocc", os.Getenv("HOME")))
}

func setupDirs() {
	cDir := configDir()
	dDir := dataDir()
	iconsLightDir := fmt.Sprintf("%s/icons_light", dDir)
	iconsDarkDir := fmt.Sprintf("%s/icons_dark", dDir)

	// Create config dir if not found (contains CLI commands, templates, CSS)
	createDir(cDir)
	// copy files if not found
	copyFile("/usr/share/nwgocc/cli_commands", fmt.Sprintf("%s/cli_commands", cDir))
	copyFile("/usr/share/nwgocc/%s", fmt.Sprintf("%s/%s", cDir, *configFile))
	copyFile("/usr/share/nwgocc/style.css", fmt.Sprintf("%s/style.css", cDir))

	// Create data dir if not found (contains icons_light/, icons_dark/, preferences.json)
	createDir(dDir)
	copyFile("/usr/share/nwgocc/preferences.json", fmt.Sprintf("%s/preferences.json", dDir))

	createDir(iconsLightDir)

	// Copy missing icons
	files, err := ioutil.ReadDir("/usr/share/nwgocc/icons_light")
	check(err)
	for _, file := range files {
		copyFile(fmt.Sprintf("/usr/share/nwgocc/icons_light/%s", file.Name()), fmt.Sprintf("%s/%s", iconsLightDir, file.Name()))
	}

	createDir(iconsDarkDir)

	files, err = ioutil.ReadDir("/usr/share/nwgocc/icons_dark")
	check(err)
	for _, file := range files {
		copyFile(fmt.Sprintf("/usr/share/nwgocc/icons_dark/%s", file.Name()), fmt.Sprintf("%s/%s", iconsDarkDir, file.Name()))
	}
}

// Create dir if it doesn't exist
func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		check(err)
		if err == nil {
			fmt.Println("Creating dir:", dir)
		}
	}
}

func copyFile(src, dst string) error {
	if !*restoreDefaults {
		if _, err := os.Stat(dst); !os.IsNotExist(err) {
			return err
		}
	}
	fmt.Println("Copying file:", dst)

	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
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
	SetBrightness      string `json:"set_brightness"`
	Systemctl          string `json:"systemctl"`
	Playerctl          string `json:"playerctl"`
}

// Settings store user preferecnces, icon definitions and external commands
type Settings struct {
	Preferences Preferences `json:"preferences"`
	Icons       Icons       `json:"icons"`
	Commands    Commands    `json:"commands"`
}

// Parses the config.json file and returns Configuration instance
func loadConfig() (Configuration, error) {
	path := fmt.Sprintf("%s/%s", configDir(), *configFile)
	bytes, err := ioutil.ReadFile(path)
	check(err)
	//if err != nil {
	//	return Configuration{}, err
	//}

	var c Configuration
	err = json.Unmarshal(bytes, &c)
	check(err)
	//if err != nil {
	//	return Configuration{}, err
	//}

	return c, nil
}

// Saves current Configuration to a json file
func saveConfig() error {
	path := fmt.Sprintf("%s/%s", configDir(), *configFile)
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

// Parses the preferences.json file and returns Settings instance
func loadSettings() (Settings, error) {
	path := fmt.Sprintf("%s/preferences.json", dataDir())
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

// Saves current settings to a json file
func saveSettings() error {
	path := fmt.Sprintf("%s/preferences.json", dataDir())
	bytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

// Parses the cli_commands txt file and returns shell commands as []string slice
func loadCliCommands() []string {
	path := fmt.Sprintf("%s/cli_commands", configDir())
	bytes, err := ioutil.ReadFile(path)
	check(err)
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

func saveCliFile(s string) {
	path := fmt.Sprintf("%s/cli_commands", configDir())
	b := []byte(s)
	err := ioutil.WriteFile(path, b, 0644)
	check(err)
}

// Returns output of each command as a string, separated with new lines, ready for use in cliLabel
func getCliOutput(commands []string) string {
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
			o = o[0:38] + "…"
		}
		o = strings.TrimSpace(o)
		output = append(output, o)
	}

	return string(strings.Join(output, "\n"))
}

// Returns output of a CLI command with optional arguments
func getCommandOutput(command string) string {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// Checks external commands availability
func checkCommands(commands Commands) {
	fmt.Println("Checking commands availability:")
	v := reflect.ValueOf(commands)
	values := make([]interface{}, v.NumField())

	m := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		cmd, _ := values[i].(string)
		cmd = strings.Split(cmd, " ")[0]
		available := getCommandOutput(fmt.Sprintf("command -v %s ", cmd)) != ""
		if !keyFound(m, cmd) {
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
	return getCommandOutput(fmt.Sprintf("command -v %s ", cmd)) != ""
}

// This needs more work: currently won't work on systemd-less systems.
func btServiceEnabled() bool {
	if isCommand(settings.Commands.Systemctl) {
		return getCommandOutput("systemctl is-enabled bluetooth.service") == "enabled" &&
			getCommandOutput("systemctl is-active bluetooth.service") == "active"
	}
	return false
}
