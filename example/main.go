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
