// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
	"strings"
)

const (
	// DefAssignAddArr appends the array to array
	DefAssignAddArr = `AssignAddºArr`
	// DefAssignAddArrArr appends the array to array
	DefAssignAddArrArr = `AssignAddºArrArr`
	// DefAssignAddMap appends the map to array
	DefAssignAddMap = `AssignAddºArrMap`
	// DefAssignArr assigns one array to another
	DefAssignArr = `AssignºArrArr`
	// DefAssignMap assigns one map to another
	DefAssignMap = `AssignºMapMap`
	// DefLenArr returns the length of the array
	DefLenArr = `LenºArr`
	// DefLenMap returns the length of the map
	DefLenMap = `LenºMap`
	// DefAssignIntInt equals int = int
	DefAssignIntInt = `#Assign#int#int`
	// DefAssignStructStruct equals struct = struct
	DefAssignStructStruct = `AssignºStructStruct`
	// DefAssignBitAndStructStruct equals struct &= struct
	DefAssignBitAndStructStruct = `AssignBitAndºStructStruct`
	// DefAssignFnFn equals fn = fn
	DefAssignFnFn = `AssignºFnFn`
	// DefAssignFileFile equals file = file
	DefAssignFileFile = `AssignºFileFile`
	// DefAssignBitAndArrArr equals arr &= arr
	DefAssignBitAndArrArr = `AssignBitAndºArrArr`
	// DefAssignBitAndMapMap equals map &= map
	DefAssignBitAndMapMap = `AssignBitAndºMapMap`
	// DefNewKeyValue returns a pair of key value
	DefNewKeyValue = `NewKeyValue`
	// DefGetEnv returns an environment variable
	DefGetEnv = `GetEnv`
)

var (
	defFuncs = map[string]bool{
		DefAssignAddArr:             true,
		DefAssignAddArrArr:          true,
		DefAssignAddMap:             true,
		DefAssignArr:                true,
		DefAssignMap:                true,
		DefLenArr:                   true,
		DefLenMap:                   true,
		DefAssignIntInt:             true,
		DefAssignStructStruct:       true,
		DefAssignFileFile:           true,
		DefAssignFnFn:               true,
		DefAssignBitAndStructStruct: true,
		DefAssignBitAndArrArr:       true,
		DefAssignBitAndMapMap:       true,
		DefNewKeyValue:              true,
		DefGetEnv:                   true,
	}
)

// NameToType searches the type by its name. It accepts names like name.name.name.
// It creates a new type if it absents.
func (unit *Unit) NameToType(name string) IObject {
	obj := unit.FindType(name)
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

// ImportEmbed imports Embed funcs to Unit
func (unit *Unit) ImportEmbed(embed Embed) {
	var (
		code     []Bcode
		retType  *TypeObject
		parTypes []*TypeObject
		fnc      interface{}
	)
	if embed.Func == nil {
		code = []Bcode{Bcode(embed.Code)}
	} else {
		fnc = int32(embed.Code)
	}
	if len(embed.Ret) > 0 {
		retType = unit.NameToType(embed.Ret).(*TypeObject)
	}
	if len(embed.Pars) > 0 {
		pars := strings.Split(embed.Pars, `,`)
		parTypes = make([]*TypeObject, len(pars))
		for i, item := range pars {
			parTypes[i] = unit.NameToType(strings.TrimSpace(item)).(*TypeObject)
		}
	}
	obj := unit.NewObject(&EmbedObject{
		Object: Object{
			Name: embed.Name,
			Unit: unit,
			BCode: Bytecode{
				Code: code,
			},
		},
		Func:     fnc,
		Return:   retType,
		Params:   parTypes,
		Variadic: embed.Variadic,
		Runtime:  embed.Runtime,
		CanError: embed.CanError,
	})
	ind := len(unit.VM.Objects) - 1
	if defFuncs[embed.Name] {
		unit.NameSpace[embed.Name] = uint32(ind) | NSPub
		if embed.Name == DefGetEnv {
			unit.AddFunc(ind, obj, true)
		}
		return
	}
	if strings.HasSuffix(embed.Name, `Auto`) {
		unit.NameSpace[`?`+embed.Name[:len(embed.Name)-4]] = uint32(ind) | NSPub
		return
	}
	unit.AddFunc(ind, obj, true)
}
