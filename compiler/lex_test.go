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

func getWant(t *testing.T, v interface{}, want string, err error) bool {
	if err != nil {
		t.Error(err)
		return true
	}
	if fmt.Sprint(v) != want {
		t.Errorf(`get != want %s %s`, fmt.Sprint(v), want)
		return true
	}
	return false
}

func TestLex(t *testing.T) {
	for _, item := range forTestLex {
		tokens, err := Lex([]rune(item.input))
		if getWant(t, *tokens, item.want, err) {
			return
		}
	}
}

var (
	forTestLex = []inputWant{
		{"	кириллица55	id_0301 \r\nLongName	", `[{1 1 2} {1 1 14} {2 2 0} {1 2 1}]`},
		{`name`, `[{1 1 1}]`},
		{``, `[]`},
	}
)
