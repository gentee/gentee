// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

const (
	// Offsets of command's ID  Cmd.ID == (cmf[Flag] << 24) + Index
	cmfStack   = iota // Stack command
	cmfStdlib         // Stdlib command
	cmfPackage        // Package command
	cmfFunc           // Code command

	Undefined = -1 // Undefined index
)

// Cmd is a byte-code command
type Cmd struct {
	ID      int
	Value   interface{}
	TokenID int // the index of the token
}

// Object is used for getting vars, funcs etc. by the name
type Object struct {
	Type int
	Ptr  interface{} // The pointer to the object
}

// Code contains information about a compiled block
type Code struct {
	Owner    *Code
	ByteCode []Cmd
	Children []*Code
}

type Func struct {
	Code
	Name   string // the name of the function
	Return int    // the type of the result
	LexID  int    // the identifier of source code in VirtualMachine
}

type Type struct {
	Name string // the name of the type
}

// VirtualMachine contains information of compiled source code
type VirtualMachine struct {
	Funcs    []*Code // The array of functions. Its index is used as ID in Cmd
	RunID    int     // The index of run function. Undefined (-1) - run has not yet been defined
	Lexeme   []*Lex  // The array of source code
	Root     Code
	Compiler Compiler // The structure which is being used during the compilation
}

// NewVM returns a new virtual machine
func NewVM() *VirtualMachine {
	vm := VirtualMachine{
		Funcs:  make([]*Code, 0, 32),
		Lexeme: make([]*Lex, 0, 4),
		RunID:  Undefined,
	}
	return &vm
}
