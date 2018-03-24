// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

// VirtualMachine contains information of compiled source code
type VirtualMachine struct {
	Packages map[string]*Package
	Runs     []*Run
}

// Unit is a common structure for Package and Run
type Unit struct {
	Objects []*Object
	Names   map[string]*Object
}

// Package contains information about a package library
type Package struct {
	Unit
}

// Run contains information about a compiled script
type Run struct {
	Unit
	RunID int // The index of run function. Undefined (-1) - run has not yet been defined
}

// NewVM returns a new virtual machine
func NewVM() *VirtualMachine {
	vm := VirtualMachine{
		Packages: make(map[string]*Package),
		Runs:     make([]*Run, 0),
	}
	return &vm
}
