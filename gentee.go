// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"github.com/gentee/gentee/compiler"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/vm"
)

// Exec is a structure with a bytecode that is ready to run
type Exec struct {
	*core.Exec
}

// Unit is a structure describing source code unit
type Unit struct {
	*core.Unit
}

// Settings is a structure with parameters for running bytecode
type Settings struct {
	vm.Settings
}

// Gentee is a common structure for compiling and executing Gentee source code
type Gentee struct {
	*core.Workspace
}

// New creates a new Gentee workspace
func New() *Gentee {
	g := Gentee{
		Workspace: core.NewVM(vm.EmbedFuncs),
	}
	compiler.InitStdlib(g.Workspace)
	return &g
}

// Compile compiles the Gentee source code.
// The function returns bytecode, id of the compiled unit and error code.
func (g *Gentee) Compile(input, path string) (*Exec, int, error) {
	unitID, err := compiler.Compile(g.Workspace, input, path)
	if err != nil {
		return nil, 0, err
	}
	exec, err := compiler.Link(g.Workspace, unitID)
	return &Exec{Exec: exec}, unitID, err
}

// CompileAndRun compiles the specified Gentee source file and run it.
func (g *Gentee) CompileAndRun(filename string) (interface{}, error) {
	exec, _, err := g.CompileFile(filename)
	if err != nil {
		return nil, err
	}
	return exec.Run(Settings{})
}

// CompileFile compiles the specified Gentee source file.
// The function returns bytecode, id of the compiled unit and error code.
func (g *Gentee) CompileFile(filename string) (*Exec, int, error) {
	unitID, err := compiler.CompileFile(g.Workspace, filename)
	if err != nil {
		return nil, 0, err
	}
	exec, err := compiler.Link(g.Workspace, unitID)
	return &Exec{Exec: exec}, unitID, err
}

// Unit returns the unit structure by its index.
func (g *Gentee) Unit(unitID int) Unit {
	return Unit{Unit: g.Units[unitID]}
}

// Run executes the bytecode.
func (exec *Exec) Run(settings Settings) (interface{}, error) {
	return vm.Run(exec.Exec, settings.Settings)
}

// Version returns the current version of the Gentee compiler.
func (g *Gentee) Version() string {
	return core.Version
}
