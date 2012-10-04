// +build linux

/*

1) Original: Copyright (c) 2005-2008 Dustin Sallings <dustin@spy.net>. 

2) Mods: Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>. Blame: Jan Mercl.

All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

*/

// Package sio supports communication using a serial port. Currently works only
// on Linux. Cgo is not used.
package sio

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

// Addr represents a network end point address.
type Addr interface {
	Network() string // name of the network
	String() string  // string form of address
}

type addr struct {
	net string
	str string
}

// Implementation of Addr
func (a *addr) Network() string {
	return a.net
}

// Implementation of Addr
func (a *addr) String() string {
	return a.str
}

type Port struct {
	f *os.File
	a *addr
}

// Open returns a Port implementing net.Conn or an error if any. The Port
// behavior is like of the merged returns of net.DialTCP and
// net.ListenTCP.Accept, i.e. the net.Conn represents a bidirectional byte
// stream. The only supported mode ATM is 8N1. The serial line is put into raw
// mode (e.g. no HW nor XON/XOFF flow control).
//
// Ex.: sio.Open("/dev/ttyS0", syscall.B115200)
func Open(dev string, rate uint32) (p *Port, err error) {
	var f *os.File
	f, err = os.OpenFile(dev, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY, 0666)
	if err != nil {
		return nil, err
	}

	fd := f.Fd()
	t := syscall.Termios{
		Iflag:  syscall.IGNPAR,
		Cflag:  syscall.CS8 | syscall.CREAD | syscall.CLOCAL | rate,
		Cc:     [32]uint8{syscall.VMIN: 1},
		Ispeed: rate,
		Ospeed: rate,
	}
	if _, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	); err != 0 {
		return nil, err
	}

	if err = syscall.SetNonblock(int(fd), false); err != nil {
		return
	}

	return &Port{f, &addr{dev, dev}}, nil
}

// Implementation of net.Conn
func (p *Port) Read(b []byte) (n int, err error) {
	return p.f.Read(b)
}

// Implementation of net.Conn
func (p *Port) Write(b []byte) (n int, err error) {
	return p.f.Write(b)
}

// Implementation of net.Conn
func (p *Port) Close() error {
	return p.f.Close()
}

// Implementation of net.Conn
func (p *Port) LocalAddr() Addr {
	return p.a
}

// Implementation of net.Conn
func (p *Port) RemoteAddr() Addr {
	return &addr{} // Ignored
}

// Implementation of net.Conn
func (p *Port) SetDeadline(t time.Time) error {
	return nil // Ignored
}

// Implementation of net.Conn
func (p *Port) SetReadDeadline(t time.Time) error {
	return nil // Ignored
}

// Implementation of net.Conn
func (p *Port) SetWriteDeadline(t time.Time) error {
	return nil // Ignored
}
