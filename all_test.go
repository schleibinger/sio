/*

Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>.
All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

blame: Jan Mercl

*/

package sio

import (
	"syscall"
	"testing"
)

func Test(t *testing.T) {
	p, err := Open("/dev/ttyS0", syscall.B57600)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", p)

	n, err := p.Write([]byte{' '})
	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Fatal(n)
	}
}
