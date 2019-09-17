// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"github.com/gentee/gentee/compiler"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/stdlib"
	"github.com/gentee/gentee/vm"
)

// Gentee is a common structure for compiling and executing Gentee source code
type Gentee struct {
	*core.Workspace
	cmdLine []string
}

// New creates a new Gentee workspace
func New() *Gentee {
	g := Gentee{
		Workspace: core.NewVM(),
	}
	stdlib.InitStdlib(g.Workspace)
	return &g
}

// Compile compiles the Gentee source code.
// The function returns id of the compiled unit and error code.
func (g *Gentee) Compile(input, path string) (*core.Exec, int, error) {
	unitID, err := compiler.Compile(g.Workspace, input, path)
	if err != nil {
		return nil, 0, err
	}
	exec, err := compiler.Link(g.Workspace, unitID)
	return exec, unitID, err
}

// CompileFile compiles the specified Gentee source file.
// The function returns id of the compiled unit and error code.
func (g *Gentee) CompileFile(filename string) (*core.Exec, int, error) {
	unitID, err := compiler.CompileFile(g.Workspace, filename)
	if err != nil {
		return nil, 0, err
	}
	exec, err := compiler.Link(g.Workspace, unitID)
	return exec, unitID, err
}

// Unit returns the unit structure by its index.
func (g *Gentee) Unit(unitID int) *core.Unit {
	return g.Units[unitID]
}

// CmdLine sets command-line parameters.
func (g *Gentee) CmdLine(args ...string) {
	g.cmdLine = make([]string, 0, len(args))
	if len(args) > 0 {
		g.cmdLine = append(g.cmdLine, args...)
	}
}

// Run executes the bytecode.
func (g *Gentee) Run(exec *core.Exec) (interface{}, error) {
	return vm.Run(exec, vm.Settings{CmdLine: g.cmdLine})
}

/*func (g *Gentee) Run(unitID int) (interface{}, error) {
	return g.Workspace.Run(unitID, g.cmdLine)
}*/

// Version returns the current version of the Gentee compiler.
func (g *Gentee) Version() string {
	return core.Version
}
