/*

1) Original: Copyright (c) 2005-2008 Dustin Sallings <dustin@spy.net>. 

2) Mods: Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>. Blame: Jan Mercl.

All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

*/

// +build ignore

package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/schleibinger/sio"
)

func main() {
	baudRate := flag.Int("baud", 57600, "Baud rate")
	mode := flag.String("mode", "8N1", "8N1 | 7E1 | 7O1")
	if flag.Parse(); flag.NArg() != 1 {
		log.Fatal("expected 1 arg: the serial port device")
	}

	port, err := sio.Open(flag.Arg(0), *baudRate, sio.ParseMode(*mode))
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, port)
}
