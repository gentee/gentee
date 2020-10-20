// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
)

const (
	ConstIotaID = iota
	ConstDepthID
	ConstCycleID
	ConstScriptID
)

const (
	// ConstDepth is the name of the max depth of calling functions
	ConstDepth = `DEPTH`
	// ConstCycle is the name of the max count of cycle
	ConstCycle = `CYCLE`
	// ConstIota is the name of iota for constants
	ConstIota = `IOTA`
	// ConstScript is the script path
	ConstScript = `SCRIPT`
	// ConstVersion is the version of Gentee compiler
	ConstVersion   = `VERSION`
	ConstRecursive = `RECURSIVE`
	ConstOnlyFiles = `ONLYFILES`
	ConstRegExp    = `REGEXP`

	// NotIota means that constant doesn't use IOTA
	NotIota = -1

	// Version is the current version of the compiler
	Version = `1.15.1+2`
)

// NewConst adds a new ConstObject to Unit
func (unit *Unit) NewConst(name string, value interface{}, redefined bool) int32 {
	result := unit.TypeByGoType(reflect.TypeOf(value))
	obj := &ConstObject{
		Object: Object{
			Name: name,
			Unit: unit,
		},
		Redefined: redefined,
		Exp: &CmdValue{
			Value:  value,
			Result: result,
		},
		Return: result,
		Iota:   NotIota,
	}
	unit.NewObject(obj)
	ind := uint32(len(unit.VM.Objects) - 1)
	obj.ObjID = int32(ind)
	if obj.Pub {
		ind |= NSPub
	}
	unit.NameSpace[npConst+name] = ind
	return obj.ObjID
}
