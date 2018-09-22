// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitRange appends stdlib int functions to the virtual machine
func InitRange(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		NewRange, // binary ..
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// NewRange adds two rune values
func NewRange(left, right int64) core.Range {
	return core.Range{From: left, To: right}
}
