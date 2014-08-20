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
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
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

const serve = "The quick brown fox jumps over the lazy dog"

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
	server := flag.Bool("server", false, "Send dummy text to the serial device continuously")
	echo := flag.Bool("echo", false, "ping client/server, combine w/ -server")
	verbose := flag.Bool("v", false, "verbose (applies to -echo)")
	slow := flag.Bool("slow", false, "slow echo server to 1 byte/sec")
	flag.Parse()
	dev := "/dev/ttyS0"
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

	switch *echo {
	case false:
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
					logFatalf("write: %s", err)
				}

				if n != len(s) {
					logFatalf("short write: %d %d", n, len(serve))
				}

				fmt.Print(s)
				sum += n
			}
		case false:
			io.Copy(os.Stdout, port)
		}
	case true:
		switch *server {
		case true:
			pass := 0
			for {
				pass++
				rxbuf := []byte{0}
				n, err := port.Read(rxbuf)
				if err != nil {
					logFatalf("read: %s", err)
				}

				if n != len(rxbuf) {
					logFatalf("short read: %d %d", n, len(rxbuf))
				}

				if *verbose {
					fmt.Printf("rx(%d): %02xH\n", pass, rxbuf[0])
				}

				if *slow {
					<-time.After(time.Second)
				}

				txbuf := []byte{0}
				copy(txbuf, rxbuf)
				if n, err = port.Write(txbuf); err != nil {
					logFatalf("write: %s", err)
				}

				if n != len(txbuf) {
					logFatalf("short write: %d %d", n, len(txbuf))
				}

				if *verbose {
					fmt.Printf("tx(%d): %02xH\n", pass, txbuf[0])
				}

			}
		case false:
			pass := 0
			for txbuf, rxbuf := []byte{' '}, []byte{0}; ; txbuf[0]++ {
				pass++

				n, err := port.Write(txbuf)
				if err != nil {
					logFatalf("write: %s", err)
				}

				if n != len(txbuf) {
					logFatalf("short write: %d %d", n, len(txbuf))
				}

				if *verbose {
					fmt.Printf("tx(%d): %02xH\n", pass, txbuf)
					os.Stdout.Sync()
				}

				if n, err = port.Read(rxbuf); err != nil {
					logFatalf("read: %s", err)
				}

				if n != len(rxbuf) {
					logFatalf("short read: %d %d", n, len(rxbuf))
				}

				if !bytes.Equal(txbuf, rxbuf) {
					logFatalf(
						"echoed data mismatch:\nTX:\n%s\nRX:\n%s",
						hex.Dump(txbuf),
						hex.Dump(rxbuf),
					)
				}

				if *verbose {
					fmt.Printf("rx(%d): %02xH\n", pass, rxbuf[0])
				}

			}
		}
	}
}
