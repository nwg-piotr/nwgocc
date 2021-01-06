# nwg Control Center (Go version)

nwg Control Center is a highly customisable, GTK-based GUI, intended for use with window managers. It may serve as an
extension to bars / panels, providing built-in and user-defined controls. Default theme may be overridden with custom
css style sheets.

Main window, Nordic-bluish-accent GTK theme, icons Custom light:

![Main window](https://scrot.cloud/images/2020/12/31/main_window-3.png)

[More screenshots](https://scrot.cloud/album/nwgocc.tDg)

This code was primarily written in python and published as [nwgcc](https://github.com/nwg-piotr/nwgcc). As I wanted it
to work a bit faster, and also to learn something new, I translated it into Go, using the
[gotk3](https://github.com/gotk3/gotk3) library. This is my very first more serious code in golang. I'm 100% sure it
may be improved in many ways.

## Dependencies

- go
- gtk3

The above seems to be enough on Arch Linux. On other distributions you may also need `libgtk-3`, `libglib2` and
`libgdk-pixbuf2`. Check the reference on https://github.com/gotk3/gotk3.

### Components

For built-in components to work, you need external commands / dependencies as below. If you don't need one, you may
skip installing related packages (e.g. on a desktop machine, you probably don't need the brightness slider).

- `light`: for Brightness slider
- `alsa`, `alsa-utils`: for Volume slider
- `playerctl`: for mpris media player controller buttons
- `wireless_tools`: for Wi-fi status
- `bluez`, `bluez-utils`: for Bluetooth status

Sample user defined commands use `blueman` and `NetworkManager`.

## Installation

0. Install dependencies and your selection of optional dependencies.

1. `git clone https://github.com/nwg-piotr/nwgocc.git`
2. `cd nwgocc`
3. `make get` (This may take some time; you may skip this step next time.)
4. `make build`
5. `sudo make install`


## To uninstall

`sudo make uninstall`

## Usage

```text
Usage of nwgocc:
  -c string
    	user's templates: Config file name (default "config.json")
  -d	Do checks, print results
  -p	place window at the mouse Pointer position (Xorg only)
  -r	Restore defaults (preferences, templates and icons)
  -s string
    	custom Styling: css file name (default "style.css")
  -v	display Version information
 ```

 Click the Preferences button to adjust the window to your needs. For your own custom styling, either modify the
 `~/.config/nwgocc/style.css` file, or place your own `whatever.css` in the same folder, and use the `-s` flag.

 You may also make a copy of `~/.config/nwgocc/config.json` under another name, for further use with the `-c` flag.

## Credits

- GUI uses the [gotk3](https://github.com/gotk3/gotk3) package, Copyright (c) 2013-2014 Conformal Systems LLC,
Copyright (c) 2015-2018 gotk3 contributors, released under the terms of the
[ISC License](https://github.com/gotk3/gotk3/blob/master/LICENSE).

- Sound control relies on the [volume-go](https://github.com/itchyny/volume-go) package, Copyright (c) 2017-2020 itchyny,
released under the terms of the [MIT License](https://github.com/itchyny/volume-go/blob/master/LICENSE).

- Most of custom icons come from my favorite [Papirus icon theme](https://github.com/PapirusDevelopmentTeam/papirus-icon-theme),
released under the terms of the
[GNU General Public License, version 3](https://github.com/PapirusDevelopmentTeam/papirus-icon-theme/blob/master/LICENSE).
