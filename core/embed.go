// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
	"runtime"
	"strings"
)

const (
	// DefAssignIntInt equals int = int
	DefAssignIntInt = `AssignºIntInt`
	// DefAssignStructStruct equals struct = struct
	DefAssignStructStruct = `AssignºStructStruct`
	// DefNewKeyValue returns a pair of key value
	DefNewKeyValue = `NewKeyValue`
)

var (
	defFuncs = map[string]bool{
		DefAssignIntInt:       true,
		DefAssignStructStruct: true,
		DefNewKeyValue:        true,
	}
)

// NewEmbedTypes adds a new EmbedObject to Unit with types
func (unit *Unit) NewEmbedTypes(Func interface{}, inTypes []*TypeObject, outType *TypeObject) {
	var variadic, isRuntime bool

	name := runtime.FuncForPC(reflect.ValueOf(Func).Pointer()).Name()
	name = name[strings.LastIndexByte(name, '.')+1:]
	originalName := name
	if isLow := strings.Index(name, `º`); isLow >= 0 {
		name = name[:isLow] // Cut off ºType in the case like AddºStr
	}

	t := reflect.TypeOf(Func)
	if t.NumOut() >= 1 && outType == nil {
		outType = unit.TypeByGoType(t.Out(0))
	}
	inCount := t.NumIn()
	if variadic = t.IsVariadic(); variadic {
		inCount--
	}
	if inCount > 0 {
		isRuntime = t.In(0) == reflect.TypeOf(&RunTime{})
		if inTypes == nil {
			inTypes = make([]*TypeObject, inCount)
			for i := 0; i < inCount; i++ {
				inTypes[i] = unit.TypeByGoType(t.In(i))
			}
			if strings.HasPrefix(name, `Assign`) {
				inTypes[0] = outType
			}
		}
	}
	obj := unit.NewObject(&EmbedObject{
		Object: Object{
			Name: name,
			Unit: unit,
		},
		Func:     Func,
		Return:   outType,
		Params:   inTypes,
		Variadic: variadic,
		Runtime:  isRuntime,
	})
	if defFuncs[originalName] {
		unit.Names[originalName] = obj
	}
}

// NewEmbed adds a new EmbedObject to Unit
func (unit *Unit) NewEmbed(Func interface{}) {
	unit.NewEmbedTypes(Func, nil, nil)
}

// NewEmbedExt adds a new EmbedObject to Unit with string types
func (unit *Unit) NewEmbedExt(Func interface{}, in string, out string) {
	ins := strings.Split(in, `,`)
	inTypes := make([]*TypeObject, len(ins))
	for i, item := range ins {
		inTypes[i] = unit.NameToType(item).(*TypeObject)
	}
	var retType *TypeObject
	if len(out) > 0 {
		retType = unit.NameToType(out).(*TypeObject)
	}
	unit.NewEmbedTypes(Func, inTypes, retType)
}

// NameToType searches the type by its name. It accepts names like name.name.name.
// It creates a new type if it absents.
func (unit *Unit) NameToType(name string) IObject {
	obj := unit.Names[name]
	for obj != nil && obj.GetType() != ObjType {
		obj = obj.GetNext()
	}
	if obj == nil {
		ins := strings.SplitN(name, `.`, 2)
		if len(ins) == 2 {
			if ins[0] == `arr` {
				indexOf := unit.NameToType(ins[1])
				if indexOf != nil {
					obj = unit.NewType(name, reflect.TypeOf(Array{}), indexOf.(*TypeObject))
				}
			} else if ins[0] == `map` {
				indexOf := unit.NameToType(ins[1])
				if indexOf != nil {
					obj = unit.NewType(name, reflect.TypeOf(Map{}), indexOf.(*TypeObject))
				}
			}
		}
	}
	return obj
}
