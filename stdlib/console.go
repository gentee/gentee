// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// InitConsole appends stdlib console functions to the virtual machine
func InitConsole(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		Print,   // Print()
		Println, // Println()
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// Print writes to standard output.
func Print(pars ...interface{}) (int64, error) {
	n, err := fmt.Print(pars...)
	return int64(n), err
}

// Println writes to standard output.
func Println(pars ...interface{}) (int64, error) {
	n, err := fmt.Println(pars...)
	return int64(n), err
}
