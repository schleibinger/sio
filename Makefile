# Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim
# GmbH <info@schleibinger.com>. All rights reserved.  Use of this source code is
# governed by a MIT-style license that can be found in the LICENSE file.

# blame: Jan Mercl

test:
	go test -i
	go test

install: sio.go sio.c
	go install

example: example/main.go sio.c sio.go
	cd example/
	go build

clean:
	go clean

nuke: clean
	rm -f example/example *~
