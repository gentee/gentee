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
	Units map[string]*Unit
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
}

// NewVM returns a new virtual machine
func NewVM() *VirtualMachine {
	vm := VirtualMachine{
		Units: make(map[string]*Unit),
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
	case `bool`:
		name = `bool`
	case `string`:
		name = `str`
	default:
		return nil
	}
	if obj, ok := unit.Names[name]; ok && obj.GetType() == ObjType {
		return obj.(*TypeObject)
	}
	return nil
}

// Run executes run block
func (vm *VirtualMachine) Run(name string) (interface{}, error) {
	rt := newRunTime(vm)
	unit := vm.Units[name]
	if unit == nil || unit.Type == UnitPackage {
		return nil, runtimeError(rt, ErrRunIndex)
	}
	funcRun := unit.Objects[unit.RunID].(*FuncObject)
	if err := rt.runCmd(&funcRun.Block); err != nil {
		return nil, err
	}
	var result interface{}
	if funcRun.Block.Result != nil {
		if len(rt.Stack) == 0 {
			return nil, runtimeError(rt, ErrRuntime)
		}
		result = rt.Stack[len(rt.Stack)-1]
	}
	return result, nil
}
