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

// Compile compiles the source code
func (workspace *Workspace) Compile(input, path string) (int, error) {
	return compiler.Compile(workspace.VM, input, path)
}

// CompileFile compiles the source file
func (workspace *Workspace) CompileFile(filename string) (int, error) {
	return compiler.CompileFile(workspace.VM, filename)
}

// Unit returns the unit structure by its index
func (workspace *Workspace) Unit(unitID int) *core.Unit {
	return workspace.VM.Units[unitID]
}

// Run executes the unit with specified name
func (workspace *Workspace) Run(unitID int) (interface{}, error) {
	return workspace.VM.Run(unitID)
}

// Version returns th ecurrent version of the Gentee compiler
func (workspace *Workspace) Version() string {
	return core.Version
}
