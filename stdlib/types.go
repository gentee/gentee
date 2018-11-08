// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"reflect"

	"github.com/gentee/gentee/core"
)

type initType struct {
	name     string
	original reflect.Type
	index    string // support index of
}

// InitTypes appends stdlib types to the virtual machine
func InitTypes(vm *core.VirtualMachine) {
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
		{`keyval`, reflect.TypeOf(core.KeyValue{}), ``},
		{`struct`, typeStruct, ``},
		// arr* is for embedded array funcs. It means array of any type
		{`arr*`, typeArr, ``},
		{`arr`, typeArr, `str`},
		{`arr`, typeArr, `int`},
		{`arr`, typeArr, `bool`},
		// map* is for embedded map funcs. It means map of any type
		{`map*`, typeMap, ``},
		{`map`, typeMap, `str`},
		{`map`, typeMap, `int`},
		{`map`, typeMap, `bool`},
	} {
		var indexOf core.IObject //*core.TypeObject
		if len(item.index) > 0 {
			indexOf = vm.StdLib().Names[item.index]
		}
		vm.StdLib().NewType(item.name, item.original, indexOf)
	}
	defType := func(name string, original reflect.Type) {
		typeObject := core.TypeObject{
			Object: core.Object{
				Name: name,
			},
			Original: original,
		}
		vm.StdLib().NewObject(&typeObject)
		typeObject.IndexOf = vm.StdLib().Names[`str`].(*core.TypeObject)
	}
	defType(`arr`, typeArr)
	defType(`map`, typeMap)
}
