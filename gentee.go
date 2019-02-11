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
}

// New creates a new Gentee workspace
func New() *Gentee {
	g := Gentee{
		core.NewVM(),
	}
	stdlib.InitStdlib(g.VirtualMachine)
	return &g
}

// Compile compiles the Gentee source code. It returns id of the compiled unit or error.
func (g *Gentee) Compile(input, path string) (int, error) {
	return compiler.Compile(g.VirtualMachine, input, path)
}

// CompileFile compiles the source file. It returns id of the compiled unit or error.
func (g *Gentee) CompileFile(filename string) (int, error) {
	return compiler.CompileFile(g.VirtualMachine, filename)
}

// Unit returns the unit structure by its index.
func (g *Gentee) Unit(unitID int) *core.Unit {
	return g.Units[unitID]
}

// Run executes the unit with specified name.
func (g *Gentee) Run(unitID int) (interface{}, error) {
	return g.VirtualMachine.Run(unitID)
}

// Version returns th ecurrent version of the Gentee compiler
func (g *Gentee) Version() string {
	return core.Version
}
