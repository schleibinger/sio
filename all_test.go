/*

Copyright (c) 2012 Schleibinger Geräte Teubert u. Greim GmbH
<info@schleibinger.com>.
All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

blame: Jan Mercl

*/

package sio

import (
	"regexp"
	"testing"
)

func TestParseMode(t *testing.T) {
	data := []struct {
		in string
		e  Conf
	}{
		{"S8N1", NOSUPP},
		{"S8E1", NOSUPP},
		{"S8O1", NOSUPP},
		{"8N1", S8N1},
		{"7E1", S7E1},
		{"701", NOSUPP},
		{"7O1", S7O1},
		{"7N1", NOSUPP},
	}
	for _, test := range data {
		in := test.in
		g, e := ParseMode(in), test.e
		if g != e {
			t.Errorf("%q: got %d, exp %d", in, g, e)
		}
	}
}

func TestOpen(t *testing.T) {
	data := []struct { // all must fail
		rate int
		mode Conf
		re   string
	}{
		{-1, S8N1, "[Bb]aud|[Rr]ate"},
		{9600, NOSUPP, "[Cc]onf|[Mmode]"},
		{-1, NOSUPP, "[Bb]aud|[Rr]ate|[Cc]onf|[Mmode]"},
	}
	e := Error("")
	for _, test := range data {
		rate, mode, re := test.rate, test.mode, test.re
		_, g := Open("", rate, mode)
		if g == nil {
			t.Errorf("Open(\"\", %d, %d) exp. err == nil, got (%T)(%#v)", rate, mode, g, g)
		}

		if _, ok := g.(Error); !ok {
			t.Errorf("Open(\"\", %d, %d) exp. err of type %T, got (%T)(%#v)", rate, mode, e, g, g)
		}

		sg := g.Error()
		m, err := regexp.MatchString(re, sg)
		if err != nil {
			t.Fatal(err)
		}

		if !m {
			t.Errorf("Open(\"\", %d, %d): %q doesn't match %q", rate, mode, sg, re)
		}
	}
}
