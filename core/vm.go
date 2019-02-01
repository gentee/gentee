// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
)

// VirtualMachine contains information of compiled source code
type VirtualMachine struct {
	Units     []*Unit
	UnitNames map[string]int
	Objects   []IObject
	Linked    map[string]int // compiled files
}

const (
	// DefName is the key name for stdlib
	DefName = `stdlib`

	// PubOne means the only next object is public
	PubOne = 1
	// PubAll means all objects are public
	PubAll = 2
)

// Unit is a common structure for source code
type Unit struct {
	VM        *VirtualMachine
	Index     uint32            // Index of the Unit
	NameSpace map[string]uint32 // name space of the unit
	Included  map[uint32]bool   // false - included or true - imported units
	Lexeme    []*Lex            // The array of source code
	RunID     int               // The index of run function. Undefined (-1) - run has not yet been defined
	Name      string            // The name of the unit
	Pub       int               // Public mode
}

// NewVM returns a new virtual machine
func NewVM() *VirtualMachine {
	vm := VirtualMachine{
		UnitNames: make(map[string]int),
		Units:     make([]*Unit, 0, 32),
		Objects:   make([]IObject, 0, 500),
		Linked:    make(map[string]int),
	}
	return &vm
}

// InitUnit initialize a unit structure
func (vm *VirtualMachine) InitUnit() *Unit {
	return &Unit{
		VM:        vm,
		RunID:     Undefined,
		NameSpace: make(map[string]uint32),
		Included:  make(map[uint32]bool),
	}
}

// TypeByGoType returns the type by the go type name
func (unit *Unit) TypeByGoType(goType reflect.Type) *TypeObject {
	var name string
	switch goType.String() {
	case `int64`:
		name = `int`
	case `float64`:
		name = `float`
	case `bool`:
		name = `bool`
	case `string`:
		name = `str`
	case `int32`:
		name = `char`
	case `core.KeyValue`:
		name = `keyval`
	case `core.Range`:
		name = `range`
	case `*core.Buffer`:
		name = `buf`
	case `*core.Array`:
		name = `arr`
	case `*core.Map`:
		name = `map`
	default:
		return nil
	}
	if obj := unit.FindType(name); obj != nil {
		return obj.(*TypeObject)
	}
	return nil
}

// StdLib returns the pointer to Standard Library Unit
func (vm *VirtualMachine) StdLib() *Unit {
	return vm.Unit(DefName)
}

// Unit returns the pointer to Unit by its name
func (vm *VirtualMachine) Unit(name string) *Unit {
	return vm.Units[vm.UnitNames[name]]
}

// Run executes run block
func (vm *VirtualMachine) Run(unitID int) (interface{}, error) {
	rt := newRunTime(vm)
	if unitID < 0 || unitID >= len(vm.Units) {
		return nil, runtimeError(rt, nil, ErrRunIndex)
	}
	unit := vm.Units[unitID]
	if unit.RunID == Undefined {
		return nil, runtimeError(rt, nil, ErrNotRun)
	}
	funcRun := vm.Objects[unit.RunID].(*FuncObject)
	if err := rt.runCmd(&funcRun.Block); err != nil {
		return nil, err
	}
	var result interface{}
	if funcRun.Block.Result != nil {
		if len(rt.Stack) == 0 {
			return nil, runtimeError(rt, nil, ErrRuntime)
		}
		result = rt.Stack[len(rt.Stack)-1]
	}
	return result, nil
}
