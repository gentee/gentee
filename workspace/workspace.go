// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package workspace

import (
	"github.com/gentee/gentee/compiler"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/stdlib"
)

// Workspace is a common structure for compiling and executing
type Workspace struct {
	VM *core.VirtualMachine
}

// New creates a new Workspace structure
func New() *Workspace {
	workspace := Workspace{
		VM: core.NewVM(),
	}
	stdlib.InitStdlib(workspace.VM)
	return &workspace
}

func (workspace *Workspace) Compile(input, name string) error {
	return compiler.Compile(workspace.VM, input, name)
}
