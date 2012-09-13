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
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/schleibinger/sio"
)

const serve = "The quick brown fox jumps over the lazy dog"

func main() {
	baudRate := flag.Int("baud", 57600, "Baud rate")
	mode := flag.String("mode", "8N1", "8N1 | 7E1 | 7O1")
	server := flag.Bool("server", false, "Send dummy text to the serial device continuously")
	if flag.Parse(); flag.NArg() != 1 {
		log.Fatal("expected 1 arg: the serial port device")
	}

	port, err := sio.Open(flag.Arg(0), *baudRate, sio.ParseMode(*mode))
	if err != nil {
		log.Fatal(err)
	}

	state, err := terminal.MakeRaw(int(port.File().Fd()))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := terminal.Restore(int(port.File().Fd()), state); err != nil {
			log.Println(err)
		}
	}()

	switch *server {
	case true:
		var s string
		var sum int
		pass := 0
		for {
			pass++
			s = fmt.Sprintf("%d(%d): %s @ %s.\n", pass, sum, serve, time.Now())
			n, err := port.Write([]byte(s))
			if err != nil {
				log.Fatal(err)
			}

			if n != len(s) {
				log.Fatal(n, len(serve))
			}

			fmt.Print(s)
			sum += n
		}
	case false:
		io.Copy(os.Stdout, port)
	}
}
