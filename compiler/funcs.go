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
	ind := cmpl.appendObj(funcObj)
	return ind
}

func coRun(cmpl *compiler) error {
	if cmpl.runID != core.Undefined {
		return cmpl.Error(ErrRun)
	}
	cmpl.runID = newFunc(cmpl, `run`)
	return nil
}

func coRunBack(cmpl *compiler) error {
	funcObj := cmpl.latestFunc()
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
	funcObj := cmpl.latestFunc()
	funcObj.Block.Result = obj.(*core.TypeObject)
	return coFuncStart(cmpl)
}

func coFuncStart(cmpl *compiler) error {
	funcObj := cmpl.latestFunc()
	funcObj.Block.ParCount = len(funcObj.Block.Vars)
	params := funcObj.GetParams()
	if funcObj.Block.Variadic {
		funcObj.Block.ParCount--
		params = nil
	}
	if obj := getFunc(cmpl, funcObj.Name, params); obj != nil {
		if core.IsVariadic(obj) {
			return cmpl.ErrorFunction(ErrFuncExists, int(funcObj.Block.TokenID), funcObj.Name,
				append(funcObj.GetParams(), nil))
		}
		return cmpl.ErrorFunction(ErrFuncExists, int(funcObj.Block.TokenID), funcObj.Name,
			obj.GetParams())
	}
	if cmpl.runID == cmpl.curFunc {
		return nil
	}
	cmpl.unit.AddFunc(cmpl.curFunc, funcObj, cmpl.unit.Pub != 0)
	if cmpl.unit.Pub == core.PubOne {
		cmpl.unit.Pub = 0
	}
	return nil
}

func getFunc(cmpl *compiler, name string, params []*core.TypeObject) (obj core.IObject) {
	var variadic bool
	obj, variadic = cmpl.unit.FindFunc(name, params)
	if obj == nil || !variadic {
		return
	}
	if params == nil {
		return obj
	}
	objPars := obj.GetParams()
	if len(params) >= len(objPars) {
		equal := true
		for i, typeParam := range objPars {
			if !isEqualTypes(typeParam, params[i]) {
				equal = false
				break
			}
		}
		if equal && obj.GetType() == core.ObjFunc {
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
	return nil
}

func getOperator(cmpl *compiler, name string, left, right core.ICmd) (obj core.IObject) {
	params := []*core.TypeObject{left.GetResult()}
	if right != nil {
		params = append(params, right.GetResult())
	}
	return getFunc(cmpl, name, params)
}

func coFuncBack(cmpl *compiler) error {
	funcObj := cmpl.latestFunc()
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
