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
	"time"

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
