// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"github.com/gentee/gentee/compiler"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/stdlib"
)

// Gentee is a common structure for compiling and executing Gentee source code
type Gentee struct {
	*core.VirtualMachine
	cmdLine []string
}

// New creates a new Gentee workspace
func New() *Gentee {
	g := Gentee{
		VirtualMachine: core.NewVM(),
	}
	stdlib.InitStdlib(g.VirtualMachine)
	return &g
}

// Compile compiles the Gentee source code.
// The function returns id of the compiled unit and error code.
func (g *Gentee) Compile(input, path string) (int, error) {
	return compiler.Compile(g.VirtualMachine, input, path)
}

// CompileFile compiles the specified Gentee source file.
// The function returns id of the compiled unit and error code.
func (g *Gentee) CompileFile(filename string) (int, error) {
	return compiler.CompileFile(g.VirtualMachine, filename)
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

// Run executes the unit by its identifier.
func (g *Gentee) Run(unitID int) (interface{}, error) {
	return g.VirtualMachine.Run(unitID, g.cmdLine)
}

// Version returns the current version of the Gentee compiler.
func (g *Gentee) Version() string {
	return core.Version
}
