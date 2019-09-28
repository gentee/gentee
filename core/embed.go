// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
	"runtime"
	"strings"
)

type Link struct {
	Func interface{}
	Code Bcode
}

const (
	// DefAssignAddArr appends the array to array
	DefAssignAddArr = `AssignAddºArrArr`
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
	// DefAssignBitAndArrArr equals arr &= arr
	DefAssignBitAndArrArr = `AssignBitAndºArrArr`
	// DefAssignBitAndMapMap equals map &= map
	DefAssignBitAndMapMap = `AssignBitAndºMapMap`
	// DefNewKeyValue returns a pair of key value
	DefNewKeyValue = `NewKeyValue`
	// DefSetEnv sets an environment variable
	DefSetEnv = `#SetEnv#str#int`
	// DefGetEnv returns an environment variable
	DefGetEnv = `GetEnv`
)

var (
	defFuncs = map[string]bool{
		DefAssignAddArr:             true,
		DefAssignAddMap:             true,
		DefAssignArr:                true,
		DefAssignMap:                true,
		DefLenArr:                   true,
		DefLenMap:                   true,
		DefAssignIntInt:             true,
		DefAssignStructStruct:       true,
		DefAssignFnFn:               true,
		DefAssignBitAndStructStruct: true,
		DefAssignBitAndArrArr:       true,
		DefAssignBitAndMapMap:       true,
		DefNewKeyValue:              true,
		DefSetEnv:                   true,
		DefGetEnv:                   true,
	}
)

// NewEmbedTypes adds a new EmbedObject to Unit with types
func (unit *Unit) NewEmbedTypes(Func interface{}, inTypes []*TypeObject, outType *TypeObject) {
	var (
		variadic, isRuntime, isError bool
		code                         []Bcode
	)

	if v, ok := Func.(Link); ok {
		Func = v.Func
		code = []Bcode{v.Code}
	}
	name := runtime.FuncForPC(reflect.ValueOf(Func).Pointer()).Name()
	name = name[strings.LastIndexByte(name, '.')+1:]
	originalName := name
	if isLow := strings.Index(name, `º`); isLow >= 0 {
		name = name[:isLow] // Cut off ºType in the case like AddºStr
	}

	t := reflect.TypeOf(Func)
	if t.NumOut() >= 1 {
		if outType == nil {
			outType = unit.TypeByGoType(t.Out(0))
		}
		typeOut := t.Out(t.NumOut() - 1)
		isError = typeOut.Kind() == reflect.Interface &&
			typeOut.Implements(reflect.TypeOf((*error)(nil)).Elem())
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
			BCode: Bytecode{
				Code: code,
			},
		},
		Func:     Func,
		Return:   outType,
		Params:   inTypes,
		Variadic: variadic,
		Runtime:  isRuntime,
		CanError: isError,
	})
	ind := len(unit.VM.Objects) - 1
	if defFuncs[originalName] {
		unit.NameSpace[originalName] = uint32(ind) | NSPub
		return
	}
	if strings.HasSuffix(originalName, `Auto`) {
		unit.NameSpace[`?`+name] = uint32(ind) | NSPub
		return
	}
	unit.AddFunc(ind, obj, true)
}

// NewEmbed adds a new EmbedObject to Unit
func (unit *Unit) NewEmbed(Func interface{}) {
	unit.NewEmbedTypes(Func, nil, nil)
}

// NewEmbedExt adds a new EmbedObject to Unit with string types
func (unit *Unit) NewEmbedExt(Func interface{}, in string, out string) {
	var ins []string
	if len(in) > 0 {
		ins = strings.Split(in, `,`)
	}
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
		return
	}
	if strings.HasSuffix(embed.Name, `Auto`) {
		unit.NameSpace[`?`+embed.Name[:len(embed.Name)-4]] = uint32(ind) | NSPub
		return
	}
	unit.AddFunc(ind, obj, true)
}
