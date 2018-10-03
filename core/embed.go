// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
	"runtime"
	"strings"
)

// NewEmbedTypes adds a new EmbedObject to Unit with types
func (unit *Unit) NewEmbedTypes(Func interface{}, inTypes []*TypeObject, outType *TypeObject) {
	name := runtime.FuncForPC(reflect.ValueOf(Func).Pointer()).Name()
	name = name[strings.LastIndexByte(name, '.')+1:]
	if isLow := strings.Index(name, `º`); isLow >= 0 {
		name = name[:isLow] // Cut off ºType in the case like AddºStr
	}

	t := reflect.TypeOf(Func)
	if t.NumOut() >= 1 && outType == nil {
		outType = unit.TypeByGoType(t.Out(0))
	}
	if inCount := t.NumIn(); inCount > 0 && inTypes == nil {
		inTypes = make([]*TypeObject, inCount)
		for i := 0; i < inCount; i++ {
			inTypes[i] = unit.TypeByGoType(t.In(i))
		}
		if strings.HasPrefix(name, `Assign`) {
			inTypes[0] = outType
		}
	}
	unit.NewObject(&EmbedObject{
		Object: Object{
			Name: name,
			Unit: unit,
		},
		Func:   Func,
		Return: outType,
		Params: inTypes,
	})
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
	unit.NewEmbedTypes(Func, inTypes, unit.NameToType(out).(*TypeObject))
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
