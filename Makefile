get:
	go get github.com/gotk3/gotk3/gdk
	go get github.com/gotk3/gotk3/glib
	go get github.com/gotk3/gotk3/gtk
	go get github.com/itchyny/volume-go
	go get github.com/allan-simon/go-singleinstance

build:
	go build -o bin/nwgocc *.go

install:
	mkdir -p /usr/share/nwgocc
	cp configs/* /usr/share/nwgocc
	cp preferences.glade /usr/share/nwgocc
	cp -R icons_light /usr/share/nwgocc
	cp -R icons_dark /usr/share/nwgocc
	cp bin/nwgocc /usr/bin

uninstall:
	rm -r /usr/share/nwgocc
	rm /usr/bin/nwgocc

run:
	go run *.go
