/*

1) Original: Copyright (c) 2005-2008 Dustin Sallings <dustin@spy.net>. 

2) Mods: Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>. Blame: Jan Mercl.

All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

*/

// Package sio lets you access old serial junk. It's a go-gettable forkmod of
// dustin's rs232 package found at https://github.com/dustin/rs232.go.
package sio

/*

#include <stdlib.h>
#include <fcntl.h>
#include <termios.h>

void initRates();

*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
)

var (
	rates = map[int]_Ctype_speed_t{}
	modes = map[string]Conf{
		"8N1": S8N1,
		"7E1": S7E1,
		"7O1": S7O1,
	}
	configs = map[Conf]struct{ reset, set C.tcflag_t }{
		S8N1: {
			C.PARENB | C.CSTOPB | C.CSIZE,
			C.CS8,
		},
		S7E1: {
			C.PARODD | C.CSTOPB | C.CSIZE,
			C.CS7 | C.PARENB,
		},
		S7O1: {
			C.CSTOPB | C.CSIZE,
			C.CS7 | C.PARENB | C.PARODD,
		},
	}
)

func init() {
	C.initRates()
}

//export addRate
func addRate(num int, val _Ctype_speed_t) {
	rates[num] = val
}

// ParseMode converts 's' to a Conf value. NOSUPP is returned for invalid mode strings.
func ParseMode(s string) Conf {
	return modes[s]
}

// Conf type represents the basic serial configuration to provide to Open.
type Conf int

// Conf typed configurations for Open
const (
	NOSUPP Conf = iota
	S8N1
	S7E1
	S7O1
)

// Error type of errors returned by Open.
type Error string

const (
	ErrGetAttr           Error = "sio.Open: C.tcgetattr failed"
	ErrSetInputSpeed     Error = "sio.Open: C.cfsetispeed failed"
	ErrSetOutputSpeed    Error = "sio.Open: C.cfsetospeed failed"
	ErrSetAttr           Error = "sio.Open: C.tcsetattr failed"
	ErrDisablingBlocking Error = "sio.Open: Failed to set non-blocking mode"
)

// Error implements error.
func (e Error) Error() string {
	return string(e)
}

// Port type represents a serial port.
type Port struct {
	f *os.File
}

// Open returns a serial port implementing io.ReadWriteCloser or an error if
// any.
//
// Example:  sio.Open("/dev/ttyS0", 115200, sio.S8N1)
func Open(dev string, rate int, conf Conf) (p *Port, err error) {
	config, ok := configs[conf]
	if !ok {
		return nil, Error(fmt.Sprintf("sio.Open: Configuration Conf(%d) not supported", conf))
	}

	bauds, ok := rates[rate]
	if !ok {
		return nil, Error(fmt.Sprintf("sio.Open: Baud rate %d not supported", rate))
	}

	f, err := os.OpenFile(dev, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY, 0666)
	if err != nil {
		return nil, Error("sio.Open: " + err.Error())
	}

	options := &C.struct_termios{}
	p, fd := &Port{f}, f.Fd()

	if C.tcgetattr(C.int(fd), options) < 0 {
		return nil, ErrGetAttr
	}

	if C.cfsetispeed(options, bauds) < 0 {
		return nil, ErrSetInputSpeed
	}

	if C.cfsetospeed(options, bauds) < 0 {
		return nil, ErrSetOutputSpeed
	}

	options.c_cflag &^= config.reset | C.CRTSCTS       // off
	options.c_cflag |= config.set | C.CLOCAL | C.CREAD // on
	options.c_cc[C.VMIN] = 1                           // Don't EOF set a zero read, just block

	if C.tcsetattr(C.int(fd), C.TCSANOW, options) < 0 {
		return nil, ErrSetAttr
	}

	if syscall.SetNonblock(int(fd), false) != nil {
		return nil, ErrDisablingBlocking
	}

	return
}

// Close implements io.Close
func (p *Port) Close() error {
	return p.f.Close()
}

// File returns port's underlying os.File.
func (p *Port) File() *os.File {
	return p.f
}

// Read implements io.Read
func (p *Port) Read(b []byte) (n int, err error) {
	return p.f.Read(b)
}

// Write implements io.Write
func (p *Port) Write(b []byte) (n int, err error) {
	return p.f.Write(b)
}
