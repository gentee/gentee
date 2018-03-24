// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

const (
	// DefName is the key name for stdlib
	DefName = ``
)

// InitStdlib appends stdlib types and fucntions to the virtual machine
func InitStdlib(vm *core.VirtualMachine) {
	vm.Packages[DefName] = &core.Package{}
	InitTypes(vm)
}
