// Copyright 2018 The Gentee Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"testing"
)

func TestCompile(t *testing.T) {
	for _, item := range forTestCompile {
		vm := NewVM()
		if err := vm.Compile(item.input); err != nil {
			if err.Error() != item.want {
				t.Error(err)
				return
			}
			continue
		}
		if get, err := vm.Run(); err != nil {
			t.Error(err)
			return
		} else if !getWant(t, get, item.want, ``) {
			return
		}
	}
}

var (
	forTestCompile = []inputWant{
		{`run {}`, `<nil>`},
		{`run {return 10}`, ``},
		{`run {return}`, ``},
		{`run int {}`, ``},
		{`run int {return 77}`, ``},
	}
)
