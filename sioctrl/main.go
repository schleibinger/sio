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
	"log"
	"syscall"
	"time"

	"github.com/schleibinger/sio"
)

var bauds = map[int]uint32{
	50:      syscall.B50,
	75:      syscall.B75,
	110:     syscall.B110,
	134:     syscall.B134,
	150:     syscall.B150,
	200:     syscall.B200,
	300:     syscall.B300,
	600:     syscall.B600,
	1200:    syscall.B1200,
	1800:    syscall.B1800,
	2400:    syscall.B2400,
	4800:    syscall.B4800,
	9600:    syscall.B9600,
	19200:   syscall.B19200,
	38400:   syscall.B38400,
	57600:   syscall.B57600,
	115200:  syscall.B115200,
	230400:  syscall.B230400,
	460800:  syscall.B460800,
	500000:  syscall.B500000,
	576000:  syscall.B576000,
	921600:  syscall.B921600,
	1000000: syscall.B1000000,
	1152000: syscall.B1152000,
	1500000: syscall.B1500000,
	2000000: syscall.B2000000,
	2500000: syscall.B2500000,
	3000000: syscall.B3000000,
	3500000: syscall.B3500000,
	4000000: syscall.B4000000,
}

func main() {
	oDTR := flag.Bool("dtr", false, "excercise the DTR control line")
	oDSR := flag.Bool("dsr", false, "excercise the DSR control line")
	oRTS := flag.Bool("rts", false, "excercise the RTS control line")
	oCTS := flag.Bool("cts", false, "excercise the CTS control line")
	oHalf := flag.Duration("t", 100*time.Millisecond, "flipping time")
	baudRate := flag.Int("baud", 57600, "Baud rate")
	flag.Parse()
	if !*oDTR && !*oRTS && !*oDSR && !*oCTS {
		log.Fatal("Please specify at least on of -dtr -rts -dsr -cts")
	}

	dev := "/dev/ttyS0"
	switch flag.NArg() {
	case 0:
		// nop
	case 1:
		dev = flag.Arg(0)
	default:
		log.Fatalf("expected max 1 arg: the serial port device, default is %s", dev)
	}

	port, err := sio.Open(dev, bauds[*baudRate])
	if err != nil {
		log.Fatal(err)
	}

	c := time.Tick(*oHalf)
	var on bool
	for {
		<-c
		if *oDTR {
			if err := port.SetDTR(on); err != nil {
				port.Close()
				log.Fatal(err)
			}
		}

		if *oDSR {
			if err := port.SetDSR(on); err != nil {
				port.Close()
				log.Fatal(err)
			}
		}

		if *oRTS {
			if err := port.SetRTS(on); err != nil {
				port.Close()
				log.Fatal(err)
			}
		}

		if *oCTS {
			if err := port.SetCTS(on); err != nil {
				port.Close()
				log.Fatal(err)
			}
		}

		on = !on
	}
}
