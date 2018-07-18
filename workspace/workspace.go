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
func (workspace *Workspace) Compile(input string) error {
	return compiler.Compile(workspace.VM, input, ``)
}

// CompileFile compiles the source file
func (workspace *Workspace) CompileFile(filename string) error {
	absname, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	input, err := ioutil.ReadFile(absname)
	if err != nil {
		return err
	}
	return compiler.Compile(workspace.VM, string(input), absname)
}

func (workspace *Workspace) Run(name string) (interface{}, error) {
	return workspace.VM.Run(name)
}
