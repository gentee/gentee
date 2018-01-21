// Copyright 2018 The Gentee Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"fmt"
	"testing"
)

func TestExp(t *testing.T) {
	for _, item := range forTestExp {
		vm := NewVM()
		if err := vm.Compile(fmt.Sprintf(`run {
			return %s
		}`, item.input)); err != nil {
			t.Error(err)
			return
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
	forTestExp = []inputWant{
		{`101`, ``},
	}
)
