// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"reflect"
	"strings"

	"github.com/gentee/gentee/core"
)

func newFunc(cmpl *compiler, name string) int {
	funcObj := &core.FuncObject{
		Object: core.Object{
			Name:  name,
			LexID: len(cmpl.unit.Lexeme) - 1,
			Unit:  cmpl.unit,
		},
		Block: core.CmdBlock{
			CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
			ID:        core.StackBlock,
		},
	}
	funcObj.Block.Object = funcObj
	cmpl.owners = []core.ICmd{&funcObj.Block}
	cmpl.unit.Objects = append(cmpl.unit.Objects, funcObj)
	if curName := cmpl.unit.Names[name]; curName == nil {
		cmpl.unit.Names[name] = funcObj
	} else {
		curName.SetNext(funcObj)
	}
	return len(cmpl.unit.Objects) - 1
}

func coRun(cmpl *compiler) error {
	if cmpl.runID != core.Undefined {
		return cmpl.Error(ErrRun)
	}
	cmpl.runID = newFunc(cmpl, `run`)
	return nil
}

func coRunBack(cmpl *compiler) error {
	funcObj := cmpl.unit.Objects[len(cmpl.unit.Objects)-1].(*core.FuncObject)
	if funcObj.Block.Result == nil {
		return nil
	}
	if len(funcObj.Block.Children) == 0 {
		return cmpl.Error(ErrMustReturn)
	}
	last := funcObj.Block.Children[len(funcObj.Block.Children)-1]
	if last.GetType() != core.CtStack ||
		last.(*core.CmdBlock).ID != core.StackReturn {
		return cmpl.Error(ErrMustReturn)
	}
	return nil
}

func coRunName(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if len(cmpl.unit.Name) != 0 {
		cmpl.dynamic = &cmState{tkLCurly, cmLCurly, nil, nil, 0}
		return coRetType(cmpl)
	} else if _, err := getType(cmpl); err != nil {
		if strings.IndexRune(token, '.') >= 0 {
			return cmpl.Error(ErrIdent)
		}
		cmpl.unit.Name = token
	} else {
		cmpl.dynamic = &cmState{tkLCurly, cmLCurly, nil, nil, 0}
		return coRetType(cmpl)
	}
	return nil
}

func coRetType(cmpl *compiler) error {
	obj, err := getType(cmpl)
	if err != nil {
		return err
	}
	funcObj := cmpl.unit.Objects[len(cmpl.unit.Objects)-1].(*core.FuncObject)
	funcObj.Block.Result = obj.(*core.TypeObject)
	return coFuncStart(cmpl)
}

func coFuncStart(cmpl *compiler) error {
	funcObj := cmpl.unit.Objects[len(cmpl.unit.Objects)-1].(*core.FuncObject)
	funcObj.Block.ParCount = len(funcObj.Block.Vars)
	params := funcObj.GetParams()
	if funcObj.Block.Variadic {
		funcObj.Block.ParCount--
		params = nil
	}
	if obj := getFunc(cmpl, funcObj.Name, params, true); obj != nil &&
		obj != cmpl.unit.Objects[len(cmpl.unit.Objects)-1] {
		if isVariadic(obj) {
			return cmpl.ErrorFunction(ErrFuncExists, int(funcObj.Block.TokenID), funcObj.Name,
				append(funcObj.GetParams(), nil))
		}
		return cmpl.ErrorFunction(ErrFuncExists, int(funcObj.Block.TokenID), funcObj.Name,
			obj.GetParams())
	}
	return nil
}

func getFunc(cmpl *compiler, name string, params []*core.TypeObject, isFunc bool) (obj core.IObject) {
	checkUnit := func(unit *core.Unit) core.IObject {
		obj = unit.Names[name]
		for obj != nil {
			if params == nil && (obj.GetType() == core.ObjFunc || obj.GetType() == core.ObjEmbedded) {
				return obj
			}
			objPars := obj.GetParams()
			if isVariadic(obj) && len(params) >= len(objPars) {
				equal := true
				for i, typeParam := range objPars {
					if !isEqualTypes(typeParam, params[i]) {
						equal = false
						break
					}
				}
				if obj.GetType() == core.ObjFunc {
					block := obj.(*core.FuncObject).Block
					for i := len(objPars); i < len(params); i++ {
						ptype := params[i]
						if !isEqualTypes(block.Vars[len(objPars)].IndexOf, ptype) {
							if ptype.Original == reflect.TypeOf(core.Array{}) {
								if isEqualTypes(block.Vars[len(objPars)].IndexOf, ptype.IndexOf) {
									continue
								}
							}
							equal = false
							break
						}
					}
				}
				if equal {
					return obj
				}
			}
			if len(params) == len(objPars) {
				equal := true
				for i, typeParam := range objPars {
					if !isEqualTypes(typeParam, params[i]) {
						equal = false
						break
					}
				}
				if equal {
					return obj
				}
			}
			obj = obj.GetNext()
		}
		return nil
	}
	if obj = checkUnit(cmpl.vm.StdLib()); obj == nil && isFunc {
		obj = checkUnit(cmpl.unit)
	}
	if obj != nil && params != nil {
		if name == `AssignAdd` && strings.HasPrefix(params[0].GetName(), `arr.arr`) &&
			!isEqualTypes(params[0].IndexOf, params[1]) {
			return nil
		}
		if name == `Assign` && (strings.HasPrefix(params[0].GetName(), `arr`) ||
			strings.HasPrefix(params[0].GetName(), `map`)) &&
			!isEqualTypes(params[0], params[1]) {
			return nil
		}
	}
	return obj
}

func getOperator(cmpl *compiler, name string, left, right core.ICmd) (obj core.IObject) {
	params := []*core.TypeObject{left.GetResult()}
	if right != nil {
		params = append(params, right.GetResult())
	}
	return getFunc(cmpl, name, params, false)
}

func coFuncBack(cmpl *compiler) error {
	funcObj := cmpl.unit.Objects[len(cmpl.unit.Objects)-1].(*core.FuncObject)
	if funcObj.Block.Result != nil {
		if len(funcObj.Block.Children) == 0 {
			return cmpl.Error(ErrMustReturn)
		}
		last := funcObj.Block.Children[len(funcObj.Block.Children)-1]
		if last.GetType() != core.CtStack ||
			last.(*core.CmdBlock).ID != core.StackReturn {
			return cmpl.Error(ErrMustReturn)
		}
	}
	return nil
}

func coFuncName(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if isCapital(token) {
		return cmpl.Error(ErrCapitalLetters)
	}
	if strings.IndexRune(token, '.') >= 0 {
		return cmpl.Error(ErrIdent)
	}
	newFunc(cmpl, token)
	return nil
}
