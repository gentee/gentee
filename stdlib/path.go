// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"path/filepath"

	"github.com/gentee/gentee/core"
)

// InitPath appends stdlib filepath functions to the virtual machine
func InitPath(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		JoinPath, // JoinPath()
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// JoinPath joins any number of path elements into a single path.
func JoinPath(pars ...interface{}) string {
	names := make([]string, len(pars))
	for i, name := range pars {
		names[i] = fmt.Sprint(name)
	}
	return filepath.Join(names...)
}
