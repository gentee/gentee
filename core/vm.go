// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
)

const (
	// DefName is the key name for stdlib
	DefName = `stdlib`
)

// VirtualMachine contains information of compiled source code
type VirtualMachine struct {
	Units    []*Unit
	Names    map[string]int
	Compiled int // the index of the latest compiled unit
}

// UnitType is used for types of runs or packages
type UnitType int

const (
	// UnitPackage is a package
	UnitPackage UnitType = iota + 1
	// UnitRun is an executing module
	UnitRun
)

// Unit is a common structure for Library and Run packages
type Unit struct {
	Type    UnitType
	Objects []IObject
	Names   map[string]IObject
	Lexeme  []*Lex // The array of source code
	RunID   int    // The index of run function. Undefined (-1) - run has not yet been defined
	Name    string // The name of the unit
}

// NewVM returns a new virtual machine
func NewVM() *VirtualMachine {
	vm := VirtualMachine{
		Names: make(map[string]int),
		Units: make([]*Unit, 0, 32),
	}
	return &vm
}

// InitUnit initialize a unit structure
func InitUnit(unitType UnitType) *Unit {
	return &Unit{
		Type:    unitType,
		Objects: make([]IObject, 0),
		Names:   make(map[string]IObject),
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
	case `*core.Array`:
		name = `arr`
	case `*core.Map`:
		name = `map`
	default:
		return nil
	}
	if obj, ok := unit.Names[name]; ok && obj.GetType() == ObjType {
		return obj.(*TypeObject)
	}
	return nil
}

// StdLib returns the pointer to Standard Library Unit
func (vm *VirtualMachine) StdLib() *Unit {
	return vm.Units[vm.Names[DefName]]
}

// Unit returns the pointer to Unit by its name
func (vm *VirtualMachine) Unit(name string) *Unit {
	return vm.Units[vm.Names[name]]
}

// Run executes run block
func (vm *VirtualMachine) Run(name string) (interface{}, error) {
	rt := newRunTime(vm)
	unit := vm.Unit(name)
	if unit == nil || unit.Type == UnitPackage {
		return nil, runtimeError(rt, nil, ErrRunIndex)
	}
	funcRun := unit.Objects[unit.RunID].(*FuncObject)
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
