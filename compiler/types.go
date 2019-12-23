// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"reflect"
	"strings"

	"github.com/gentee/gentee/core"
)

type initType struct {
	name     string
	original reflect.Type
	index    string // support index of
}

// InitTypes appends stdlib types to the virtual machine
func InitTypes(ws *core.Workspace) {
	typeArr := reflect.TypeOf(core.Array{})
	typeMap := reflect.TypeOf(core.Map{})
	typeStruct := reflect.TypeOf(core.Struct{})
	for _, item := range []initType{
		{`int`, reflect.TypeOf(int64(0)), ``},
		{`float`, reflect.TypeOf(float64(0.0)), ``},
		{`bool`, reflect.TypeOf(true), ``},
		{`char`, reflect.TypeOf('a'), ``},
		{`str`, reflect.TypeOf(``), `char`},
		{`range`, reflect.TypeOf(core.Range{}), `int`},
		{`buf`, reflect.TypeOf(core.Buffer{}), `int`},
		{`set`, reflect.TypeOf(core.Set{}), `bool`},
		{`keyval`, reflect.TypeOf(core.KeyValue{}), ``},
		{`struct`, typeStruct, ``},
		{`fn`, reflect.TypeOf(core.Fn{}), ``},
		{`thread`, reflect.TypeOf(int64(0)), ``},
		{`error`, reflect.TypeOf(core.RuntimeError{}), ``},
		{`obj`, reflect.TypeOf(core.Obj{}), ``},
		// arr* is for embedded array funcs. It means array of any type
		{`arr*`, typeArr, ``},
		{`arr.str`, typeArr, `str`},
		{`arr.int`, typeArr, `int`},
		{`arr.bool`, typeArr, `bool`},
		// map* is for embedded map funcs. It means map of any type
		{`map*`, typeMap, ``},
		{`map.str`, typeMap, `str`},
		{`map.int`, typeMap, `int`},
		{`map.bool`, typeMap, `bool`},
	} {
		var indexOf core.IObject
		if len(item.index) > 0 {
			indexOf = ws.StdLib().FindType(item.index)
		}
		tobj := ws.StdLib().NewType(item.name, item.original, indexOf)
		if item.name == `obj` {
			tobj.(*core.TypeObject).IndexOf = tobj.(*core.TypeObject)
		}
	}
	// Define aliases
	ws.StdLib().NameSpace[`@arr`] = ws.StdLib().NameSpace[`@arr.str`]
	ws.StdLib().NameSpace[`@map`] = ws.StdLib().NameSpace[`@map.str`]
}

// NewStructType adds a new struct type to Unit
func NewStructType(ws *core.Workspace, name string, fields []string) *core.TypeObject {
	names := make(map[string]int64)
	types := make([]*core.TypeObject, len(fields))
	for i, item := range fields {
		itype := strings.SplitN(item, `:`, 2)
		names[itype[0]] = int64(i)
		types[i] = ws.StdLib().FindType(itype[1]).(*core.TypeObject)
	}
	pType := ws.StdLib().NewType(name, reflect.TypeOf(core.Struct{}), nil).(*core.TypeObject)
	pType.Custom = &core.StructType{
		Fields: names,
		Types:  types,
	}
	return pType
}
