// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"github.com/gentee/gentee/core"
)

// Compiler contains information of the compilation process
type compiler struct {
	vm   *core.VirtualMachine
	unit *core.Unit
}

func init() {
	makeParseTable()
	makeCompileTable()
}

// Compile compiles the source code
func Compile(vm *core.VirtualMachine, input, name string) error {
	compiler := Compiler{
		vm: vm,
		unit: &core.Unit{
			Objects: make([]*Object, 0),
			Names:   make(map[string]*Object),
		},
	}
	return nil
}
