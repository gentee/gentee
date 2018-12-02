// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
)

const (
	// ConstDepth is the name of the max depth of calling functions
	ConstDepth = `DEPTH`
	// ConstCycle is the name of the max count of cycle
	ConstCycle = `CYCLE`
	// ConstIota is the name of iota for constants
	ConstIota = `IOTA`
	// ConstVersion is the version of Gentee compiler
	ConstVersion = `VERSION`

	// NotIota means that constant doesn't use IOTA
	NotIota = -1

	// Version is the current version of the compiler
	Version = `1.0.0-beta.1`
)

// NewConst adds a new ConstObject to Unit
func (unit *Unit) NewConst(name string, value interface{}, redefined bool) {
	result := unit.TypeByGoType(reflect.TypeOf(value))
	unit.NewObject(&ConstObject{
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
	})
}
