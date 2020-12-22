# nwgocc

For now this code is an extension to the [nwgcc](https://github.com/nwg-piotr/nwgcc) project written in python.
It uses the [gotk3](https://github.com/gotk3/gotk3) library to provide the `nwgocc` binary, as an alternative,
faster command to use instead of `nwgcc`. **It will not work standalone:** `nwgcc` must be installed first
and started at least once, to init necessary directories and files.

At the moment just the main program window has been implemented. Clicking the Preferences button will terminate
`nwgocc` and run the python code with the settings window open. As long as I have enough time and ~~beer~~ coffee,
the preferences-related part may be written in the future.

**Note** This is my very first more serious Go code. I'm 100% sure it may (should) be improved in many ways. 

## Installation

Install [nwgcc](https://github.com/nwg-piotr/nwgcc) first. Take a close look at optional dependencies. In case you
decide not to use the `nwgcc` command any longer, you may omit the `python-pyalsa` package.

1. Clone the repository, install the `go` package (make dependency).
2. Install Go packages. Be patient, it may take some time.

```
go get github.com/gotk3/gotk3/gtk
go get github.com/gotk3/gotk3/gdk
go get github.com/gotk3/gotk3/glib
go get github.com/itchyny/volume-go/cmd/volume
```

3. cd to your clone directory, build the binary:

```
go build -o nwgocc
```

Or you may give a try to the x86_64 binary in the `bin` folder. Download it and rename to `nwgocc`.
