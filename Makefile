# Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim
# GmbH <info@schleibinger.com>. All rights reserved.  Use of this source code is
# governed by a MIT-style license that can be found in the LICENSE file.

# blame: Jan Mercl

install: sio.go sio.c
	go install

example:
	cd example/
	go build

clean:
	go clean

nuke: clean
	rm -f example/example *~
