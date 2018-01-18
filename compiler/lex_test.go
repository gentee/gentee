// Copyright 2018 The Gentee Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"fmt"
	"testing"
)

type inputWant struct {
	input string
	want  string
}

func getWant(t *testing.T, v interface{}, want, ext string) bool {
	get := fmt.Sprint(v) + ext
	if get != want {
		t.Errorf(`get != want; %s != %s`, get, want)
		return true
	}
	return false
}

func TestLex(t *testing.T) {
	for _, item := range forTestLex {
		var ext string
		lp, err := LexParsing([]rune(item.input))
		if err != nil {
			line, column := lp.LineColumn(len(lp.Tokens) - 1)
			ext = fmt.Sprintf(` %d:%d: %s`, line, column, err)
		}
		if getWant(t, lp.Tokens, item.want, ext) {
			return
		}
	}
}

var (
	forTestLex = []inputWant{
		{`main ok`,
			`[{128 0 4} {1 5 2}]`},
		{`+ - * / () {}=`,
			`[{64 0 0} {65 2 0} {66 4 0} {67 6 0} {68 8 0} {69 9 0} {70 11 0} {71 12 0} {72 13 0}]`},
		{`0 0xaB78f 16780 0756 0779`,
			`[{3 0 1} {4 2 7} {3 10 5} {5 16 4} {6 24 0}] 1:25: wrong sequence of characters`},
		{"	Aufzählung кириллица55	id_0301 \r\nLongName	",
			`[{1 1 10} {1 12 11} {1 24 7} {2 33 0} {1 34 8}]`},
		{`name; b ®`, `[{1 0 4} {2 4 0} {1 6 1} {6 8 0}] 1:9: unknown character`},
		{``, `[]`},
	}
)
