// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package workspace

import (
	"io/ioutil"
	"path/filepath"

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
func (workspace *Workspace) Compile(input, path string) (*core.Unit, error) {
	if err := compiler.Compile(workspace.VM, input, path); err != nil {
		return nil, err
	}
	return workspace.VM.Units[workspace.VM.Compiled], nil
}

// CompileFile compiles the source file
func (workspace *Workspace) CompileFile(filename string) (*core.Unit, error) {
	absname, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	input, err := ioutil.ReadFile(absname)
	if err != nil {
		return nil, err
	}
	return workspace.Compile(string(input), absname)
}

// Run executes the unit with specified name
func (workspace *Workspace) Run(name string) (interface{}, error) {
	return workspace.VM.Run(name)
}

// Version returns th ecurrent version of the Gentee compiler
func (workspace *Workspace) Version() string {
	return core.Version
}
