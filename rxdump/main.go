// build ignore

/*

Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>.
All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

blame: Jan Mercl

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	".."
)

var bauds = map[int]uint64{
	50:      sio.B50,
	75:      sio.B75,
	110:     sio.B110,
	134:     sio.B134,
	150:     sio.B150,
	200:     sio.B200,
	300:     sio.B300,
	600:     sio.B600,
	1200:    sio.B1200,
	1800:    sio.B1800,
	2400:    sio.B2400,
	4800:    sio.B4800,
	9600:    sio.B9600,
	19200:   sio.B19200,
	38400:   sio.B38400,
	57600:   sio.B57600,
	115200:  sio.B115200,
	230400:  sio.B230400,
	460800:  sio.B460800,
	500000:  sio.B500000,
	576000:  sio.B576000,
	921600:  sio.B921600,
	1000000: sio.B1000000,
	1152000: sio.B1152000,
	1500000: sio.B1500000,
	2000000: sio.B2000000,
	2500000: sio.B2500000,
	3000000: sio.B3000000,
	3500000: sio.B3500000,
	4000000: sio.B4000000,
}

// Removed usage of "log", seems to cause "guru meditation" on Phytec's Phycard
// for reasons not yet investigated.
func logFatalf(format string, args ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		fmt.Printf("%s.%d: ", file, line)
	}
	fmt.Printf(format, args...)
	if format != "" && format[len(format)-1] != '\n' {
		fmt.Println()
	}
	os.Exit(1)
}

func main() {
	baudRate := flag.Int("baud", 57600, "Baud rate")
	dev := "/dev/ttyS0"
	flag.Parse()
	switch flag.NArg() {
	case 0:
		// nop
	case 1:
		dev = flag.Arg(0)
	default:
		logFatalf("expected max 1 arg: the serial port device, default is /dev/ttyS0")
	}

	port, err := sio.Open(dev, bauds[*baudRate])
	if err != nil {
		logFatalf("open: %s", err)
	}

	rxbuf := []byte{0}
	ofs := 0
	for {
		n, err := port.Read(rxbuf)
		if err != nil {
			logFatalf("read: %s", err)
		}

		if n != len(rxbuf) {
			logFatalf("short read: %d %d", n, len(rxbuf))
		}

		if ofs%16 == 0 {
			fmt.Printf("\n%04x: ", ofs)
		}
		fmt.Printf("%02x ", rxbuf[0])
		ofs++
	}
}
