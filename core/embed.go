// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
	"runtime"
	"strings"
)

// NewEmbed adds a new EmbedObject to Unit
func (unit *Unit) NewEmbed(Func interface{}) {
	var (
		outType *TypeObject
		inTypes []*TypeObject
	)
	name := runtime.FuncForPC(reflect.ValueOf(Func).Pointer()).Name()
	name = name[strings.LastIndexByte(name, '.')+1:]
	if isLow := strings.Index(name, `º`); isLow >= 0 {
		name = name[:isLow] // Cut off ºType in the case like AddºStr
	}

	t := reflect.TypeOf(Func)
	if t.NumOut() >= 1 {
		outType = unit.TypeByGoType(t.Out(0))
	}
	if inCount := t.NumIn(); inCount > 0 {
		inTypes = make([]*TypeObject, inCount)
		for i := 0; i < inCount; i++ {
			inTypes[i] = unit.TypeByGoType(t.In(i))
		}
		if strings.HasPrefix(name, `Assign`) {
			inTypes = inTypes[1:]
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
