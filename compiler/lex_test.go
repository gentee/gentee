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
		{"	Aufzählung кириллица55	id_0301 \r\nLongName	",
			`[{1 1 10} {1 12 11} {1 24 7} {2 32 1} {1 34 8}]`},
		{`name ®`, `[{1 0 4} {3 5 0}] 1:6: unknown character`},
		{``, `[]`},
	}
)
