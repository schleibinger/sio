// +build linux

/*

1) Original: Copyright (c) 2005-2008 Dustin Sallings <dustin@spy.net>. 

2) Mods: Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>. Blame: Jan Mercl.

All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

*/

// Package sio lets you access old serial junk. It's a go-gettable forkmod of
// dustin's rs232 package found at https://github.com/dustin/rs232.go.

// Package sio supports communication using a serial port. Currently works only
// on Linux. Cgo is not used.
package sio

import (
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type addr struct {
	net string
	str string
}

// Implementation of net.Addr
func (a *addr) Network() string {
	return a.net
}

// Implementation of net.Addr
func (a *addr) String() string {
	return a.str
}

// Tunnel represents a bidirectional TCP Tunnel of a serial line.
type Tunnel struct {
	RxConn *net.TCPConn
	TxConn *net.TCPConn
}

// TCP returns one end of a bidirectional TCP tunnel of a serial line or an
// error if any. It implements net.Conn in the same manner as the one returned
// by Open and is intended as a drop in replacement for the same.
func TCP(rxConn, txConn *net.TCPConn) (*Tunnel, error) {
	return &Tunnel{rxConn, txConn}, nil
}

// Implementation of net.Conn
func (t *Tunnel) Read(b []byte) (n int, err error) {
	return t.RxConn.Read(b)
}

// Implementation of net.Conn
func (t *Tunnel) Write(b []byte) (n int, err error) {
	return t.TxConn.Write(b)
}

// Implementation of net.Conn
func (t *Tunnel) Close() (err error) {
	err = t.TxConn.Close()
	if e2 := t.RxConn.Close(); err == nil && e2 != nil {
		err = e2
	}
	return
}

// Implementation of net.Conn. Note that the returned address is the one from
// the TxConn passed to TCP().
func (t *Tunnel) LocalAddr() net.Addr {
	return t.TxConn.LocalAddr()
}

// Implementation of net.Conn. Note that the returned address is the one from
// TxConn passed to TCP().
func (t *Tunnel) RemoteAddr() net.Addr {
	return t.TxConn.RemoteAddr()
}

// Implementation of net.Conn
func (tu *Tunnel) SetDeadline(t time.Time) (err error) {
	err = tu.SetReadDeadline(t)
	if e2 := tu.SetWriteDeadline(t); err == nil && e2 != nil {
		err = e2
	}
	return
}

// Implementation of net.Conn
func (tu *Tunnel) SetReadDeadline(t time.Time) error {
	return tu.RxConn.SetDeadline(t)
}

// Implementation of net.Conn
func (tu *Tunnel) SetWriteDeadline(t time.Time) error {
	return tu.TxConn.SetDeadline(t)
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
	//if _, ok := bauds[rate]; !ok {
	//return nil, fmt.Errorf("Unknown baud rate 0x%x", rate)
	//}

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
func (p *Port) LocalAddr() net.Addr {
	return p.a
}

// Implementation of net.Conn
func (p *Port) RemoteAddr() net.Addr {
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
