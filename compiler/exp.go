// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

func appendExp(cmpl *compiler, cmd core.ICmd) {
	cmpl.exp = append(cmpl.exp, cmd)
}

func coExpStart(cmpl *compiler) error {
	cmpl.exp = cmpl.exp[:0]
	cmpl.expbuf = cmpl.expbuf[:0]
	return nil
}

func coPush(cmpl *compiler) error {
	var (
		v     interface{}
		vType string
		err   error
	)
	lp := cmpl.getLex()
	token := getToken(lp, cmpl.pos)
	switch lp.Tokens[cmpl.pos].Type {
	case tkFloat:
		if v, err = strconv.ParseFloat(token, 64); err != nil {
			return cmpl.Error(ErrOutOfRange, token)
		}
		vType = `float`
	case tkInt:
		if v, err = strconv.ParseInt(token, 0, 64); err != nil {
			return cmpl.Error(ErrOutOfRange, token)
		}
		vType = `int`
	case tkFalse, tkTrue:
		v = lp.Tokens[cmpl.pos].Type == tkTrue
		vType = `bool`
	case tkChar:
		runes := []rune(token)
		if len(runes) < 3 {
			return cmpl.Error(ErrChar)
		}
		token, err = strconv.Unquote(`"` + strings.Replace(string(runes[1:len(runes)-1]),
			`\'`, `'`, -1) + `"`)
		if err != nil || len([]rune(token)) != 1 {
			return cmpl.Error(ErrChar)
		}
		v = []rune(token)[0]
		vType = `char`
	case tkStr:
		v = lp.Strings[lp.Tokens[cmpl.pos].Index]
		if token[0] == '"' {
			if v, err = unNewLine(v.(string)); err != nil {
				return cmpl.Error(ErrDoubleQuotes)
			}
		}
		vType = `str`
	}
	appendExp(cmpl, &core.CmdValue{Value: v,
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		Result:    cmpl.unit.FindType(vType).(*core.TypeObject)})
	return nil
}

func coExpVar(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos-1)
	init := isInState(cmpl, cmInit, 1)
	if init && cmpl.curOwner().GetResult().Custom != nil &&
		cmpl.getLex().Tokens[cmpl.pos].Type == tkColon {
		if _, ok := cmpl.curOwner().GetResult().Custom.Fields[token]; !ok {
			return cmpl.ErrorPos(cmpl.pos-1, ErrWrongField, token, cmpl.curOwner().GetResult().GetName())
		}
		appendExp(cmpl, &core.CmdValue{Value: token,
			CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)},
			Result:    cmpl.getStrType()})
		return nil
	}
	if cmpl.getLex().Tokens[cmpl.pos].Type == tkCtxEq {
		appendExp(cmpl, &core.CmdValue{Value: token,
			CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)},
			Result:    cmpl.unit.FindType(`str`).(*core.TypeObject)})
		return nil
	}
	if isCapital(token) {
		if token == core.ConstIota && cmpl.curIota == core.NotIota {
			return cmpl.ErrorPos(cmpl.pos-1, ErrIota)
		}
		if constObj := cmpl.unit.FindConst(token); constObj != nil {
			appendExp(cmpl, &core.CmdConst{
				Object:    constObj,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)},
			})
		} else {
			return cmpl.ErrorPos(cmpl.pos-1, ErrUnknownIdent, token)
		}
	} else {
		var fields []string
		if strings.IndexRune(token, '.') >= 0 {
			fields = strings.Split(token, `.`)
			token = fields[0]
			fields = fields[1:]
		}

		block, ind := findVar(cmpl, token)
		if block == nil {
			return cmpl.ErrorPos(cmpl.pos-1, ErrUnknownIdent, token)
		}
		cmdVar := &core.CmdVar{Block: block, Index: ind,
			CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)}}
		typeVar := cmdVar.GetResult()
		for _, field := range fields {
			indField, typeField, err := structIndex(cmpl, typeVar, field)
			if err != nil {
				return err
			}
			index := &core.CmdValue{Value: indField,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)},
				Result:    cmpl.getIntType()}
			cmdVar.Indexes = append(cmdVar.Indexes, core.CmdRet{Cmd: index, Type: typeField})
			typeVar = typeField
		}
		appendExp(cmpl, cmdVar)
	}
	return nil
}

func coExpEnd(cmpl *compiler) error {
	for len(cmpl.expbuf) > 0 {
		if err := popBuf(cmpl); err != nil {
			return err
		}
	}
	init := isInState(cmpl, cmInit, 0)
	if len(cmpl.exp) > 1 && !init && !isCase(cmpl) {
		return cmpl.Error(ErrCompiler, `coExpEnd`)
	}
	for len(cmpl.exp) > 0 {
		cmpl.curOwner().Children = append(cmpl.curOwner().Children, cmpl.exp[0])
		cmpl.exp = cmpl.exp[1:]
	}
	return nil
}

func popBuf(cmpl *compiler) error {
	var obj core.IObject
	expBuf := cmpl.expbuf[len(cmpl.expbuf)-1]
	prior := priority[expBuf.Oper]
	switch expBuf.Oper {
	case tkAnd, tkOr:
		if len(cmpl.exp) < 2 {
			return cmpl.Error(ErrValue)
		}
		right := cmpl.exp[len(cmpl.exp)-1]
		left := cmpl.exp[len(cmpl.exp)-2]
		if !isBoolResult(left) || !isBoolResult(right) {
			return cmpl.ErrorPos(expBuf.Pos, ErrBoolOper)
		}
		id := core.StackAnd
		if expBuf.Oper == tkOr {
			id = core.StackOr
		}
		icmd := &core.CmdBlock{ID: uint32(id),
			Result: right.GetResult(), CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Children: []core.ICmd{left, right}}
		cmpl.exp[len(cmpl.exp)-2] = icmd
		cmpl.exp = cmpl.exp[:len(cmpl.exp)-1]
	case tkStrExp:
		if len(cmpl.exp) < 2 {
			return cmpl.Error(ErrValue)
		}
		right := cmpl.exp[len(cmpl.exp)-1]
		left := cmpl.exp[len(cmpl.exp)-2]
		obj = getOperator(cmpl, prior.Name, left, right)
		if obj == nil {
			return cmpl.ErrorFunction(ErrFunction, expBuf.Pos, prior.Name, []*core.TypeObject{
				left.GetResult(), right.GetResult()})
		}
		icmd := &core.CmdBinary{CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Object: obj, Result: obj.Result(), Left: left, Right: right}
		cmpl.exp[len(cmpl.exp)-2] = icmd
		cmpl.exp = cmpl.exp[:len(cmpl.exp)-1]
	case tkAssign, tkAddEq, tkSubEq, tkMulEq, tkDivEq, tkModEq, tkLShiftEq, tkRShiftEq, tkBitAndEq,
		tkBitOrEq, tkBitXorEq:
		if len(cmpl.exp) < 2 {
			return cmpl.Error(ErrValue)
		}
		right := cmpl.exp[len(cmpl.exp)-1]
		left := cmpl.exp[len(cmpl.exp)-2]

		if expBuf.Oper == tkAssign && left.GetType() == core.CtUnary {
			if left.GetObject() == cmpl.vm.StdLib().FindObj(core.DefGetEnv) {
				setEnv := cmpl.vm.StdLib().FindObj(core.DefSetEnv)
				icmd := &core.CmdBinary{CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
					Object: setEnv,
					Result: setEnv.Result(), Left: left.(*core.CmdUnary).Operand, Right: right}
				cmpl.exp[len(cmpl.exp)-2] = icmd
				cmpl.exp = cmpl.exp[:len(cmpl.exp)-1]
				break
			}
		}
		if left.GetType() != core.CtVar {
			return cmpl.ErrorPos(expBuf.Pos, ErrLValue)
		}

		obj = getOperator(cmpl, prior.Name, left, right)
		if obj == nil {
			if expBuf.Oper == tkAssign {
				if left.GetResult().Custom != nil {
					if left.GetResult() != right.GetResult() {
						return cmpl.ErrorPos(expBuf.Pos, ErrStructAssign, right.GetResult().GetName(),
							left.GetResult().GetName())
					}
					obj = cmpl.vm.StdLib().FindObj(core.DefAssignStructStruct)
				} else if left.GetResult().Func != nil {
					if !isEqualTypes(left.GetResult(), right.GetResult()) {
						return cmpl.ErrorPos(expBuf.Pos, ErrStructAssign, right.GetResult().GetName(),
							left.GetResult().GetName())
					}
					obj = cmpl.vm.StdLib().FindObj(core.DefAssignFnFn)
				} else if left.GetResult() == right.GetResult() {
					if left.GetResult().Original == reflect.TypeOf(core.Array{}) {
						obj = cmpl.unit.FindObj(core.DefAssignArr)
					} else if left.GetResult().Original == reflect.TypeOf(core.Map{}) {
						obj = cmpl.unit.FindObj(core.DefAssignMap)
					}
				}
			} else if expBuf.Oper == tkBitAndEq {
				if left.GetResult().Custom != nil {
					if left.GetResult() != right.GetResult() {
						return cmpl.ErrorPos(expBuf.Pos, ErrStructAssign, right.GetResult().GetName(),
							left.GetResult().GetName())
					}
					obj = cmpl.vm.StdLib().FindObj(core.DefAssignBitAndStructStruct)
				} else if left.GetResult() == right.GetResult() {
					if left.GetResult().Original == reflect.TypeOf(core.Array{}) {
						obj = cmpl.unit.FindObj(core.DefAssignBitAndArrArr)
					} else if left.GetResult().Original == reflect.TypeOf(core.Map{}) {
						obj = cmpl.unit.FindObj(core.DefAssignBitAndMapMap)
					}
				}
			}
			if expBuf.Oper == tkAddEq && left.GetResult().Original == reflect.TypeOf(core.Array{}) {
				if left.GetResult().IndexOf == right.GetResult() {
					if right.GetResult().Original == reflect.TypeOf(core.Array{}) {
						obj = cmpl.unit.FindObj(core.DefAssignAddArr)
					} else if right.GetResult().Original == reflect.TypeOf(core.Map{}) {
						obj = cmpl.unit.FindObj(core.DefAssignAddMap)
					}
				}
			}
			if obj == nil {
				return cmpl.ErrorFunction(ErrFunction, expBuf.Pos, prior.Name, []*core.TypeObject{
					left.GetResult(), right.GetResult()})
			}
		}
		icmd := &core.CmdBlock{ID: core.StackAssign, Object: obj,
			Result: left.GetResult(), CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Children: []core.ICmd{left, right}}
		cmpl.exp[len(cmpl.exp)-2] = icmd
		cmpl.exp = cmpl.exp[:len(cmpl.exp)-1]
	case tkAdd, tkSub, tkMul, tkMod, tkDiv, tkEqual, tkNotEqual, tkLess, tkLessEqual, tkGreater,
		tkGreaterEqual, tkBitOr, tkBitXor, tkBitAnd, tkLShift, tkRShift, tkRange, tkCtxEq:
		if len(cmpl.exp) < 2 {
			return cmpl.Error(ErrValue)
		}
		right := cmpl.exp[len(cmpl.exp)-1]
		left := cmpl.exp[len(cmpl.exp)-2]
		obj = getOperator(cmpl, prior.Name, left, right)
		if obj == nil {
			return cmpl.ErrorFunction(ErrFunction, expBuf.Pos, prior.Name, []*core.TypeObject{
				left.GetResult(), right.GetResult()})
		}
		icmd := &core.CmdBinary{CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Object: obj, Result: obj.Result(), Left: left, Right: right}
		if expBuf.Oper == tkNotEqual || expBuf.Oper == tkLessEqual || expBuf.Oper == tkGreaterEqual {
			objNot := getFunc(cmpl, `Not`, []*core.TypeObject{obj.Result()})
			if objNot == nil {
				return cmpl.ErrorFunction(ErrFunction, expBuf.Pos, `Not`,
					[]*core.TypeObject{obj.Result()})
			}
			cmdNot := &core.CmdUnary{CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
				Object: objNot, Result: objNot.Result(), Operand: icmd}
			cmpl.exp[len(cmpl.exp)-2] = cmdNot
		} else {
			cmpl.exp[len(cmpl.exp)-2] = icmd
		}
		cmpl.exp = cmpl.exp[:len(cmpl.exp)-1]
	case tkInc | tkUnary, tkDec | tkUnary, tkInc | tkUnary | tkPost, tkDec | tkUnary | tkPost:
		if len(cmpl.exp) == 0 {
			return cmpl.Error(ErrValue)
		}
		top := cmpl.exp[len(cmpl.exp)-1]
		if !isIntResult(top) {
			return cmpl.ErrorPos(expBuf.Pos, ErrIntOper)
		}
		if top.GetType() != core.CtVar {
			return cmpl.ErrorPos(expBuf.Pos, ErrLValue)
		}
		val := 1
		if (expBuf.Oper & 0xff) == tkDec {
			val = -1
		}
		if (expBuf.Oper & tkPost) > 0 {
			val *= 2
		}
		icmd := &core.CmdBlock{ID: core.StackIncDec, ParCount: val,
			Result: top.GetResult(), CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Children: []core.ICmd{top}}
		cmpl.exp[len(cmpl.exp)-1] = icmd
	case tkSub | tkUnary, tkMul | tkUnary, tkNot | tkUnary, tkBitXor | tkUnary, tkCtx | tkUnary,
		tkDoubleCtx | tkUnary, tkBitOr | tkUnary:
		if len(cmpl.exp) == 0 {
			return cmpl.Error(ErrValue)
		}
		top := cmpl.exp[len(cmpl.exp)-1]
		if expBuf.Oper == (tkMul | tkUnary) {
			switch top.GetResult().Original {
			case reflect.TypeOf(core.Array{}):
				obj = cmpl.unit.FindObj(core.DefLenArr)
			case reflect.TypeOf(core.Map{}):
				obj = cmpl.unit.FindObj(core.DefLenMap)
			default:
				obj = getOperator(cmpl, prior.Name, top, nil)
			}
		} else {
			obj = getOperator(cmpl, prior.Name, top, nil)
		}
		if obj == nil {
			return cmpl.ErrorFunction(ErrFunction, expBuf.Pos, prior.Name,
				[]*core.TypeObject{top.GetResult()})
		}
		icmd := &core.CmdUnary{CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Object: obj, Result: obj.Result(), Operand: cmpl.exp[len(cmpl.exp)-1]}
		cmpl.exp[len(cmpl.exp)-1] = icmd
	case tkLPar:
		return cmpl.Error(ErrLPar)
	case tkLSBracket:
		return cmpl.Error(ErrLSBracket)
	default:
		return cmpl.Error(ErrCompiler, fmt.Sprintf(`popBuf unknown token %d`, expBuf.Oper))
	}
	cmpl.expbuf = cmpl.expbuf[:len(cmpl.expbuf)-1]
	return nil
}

func coFnOperator(cmpl *compiler) error {
	lp := cmpl.getLex()
	if len(lp.Tokens) == cmpl.pos+1 {
		return cmpl.ErrorPos(len(lp.Tokens), ErrEnd)
	}
	if lp.Tokens[cmpl.pos+1].Type != tkIdent {
		return cmpl.ErrorPos(cmpl.pos+1, ErrName)
	}
	token := getToken(lp, cmpl.pos+1)
	faddr := strings.Split(token, `.`)
	if len(faddr) != 2 {
		return cmpl.Error(ErrAddrFunc)
	}
	ftype := cmpl.unit.FindType(faddr[1])
	if ftype == nil || ftype.(*core.TypeObject).Func == nil {
		return cmpl.Error(ErrNoFuncType, faddr[1])
	}
	obj := getFunc(cmpl, faddr[0], ftype.(*core.TypeObject).Func.Params)
	if obj == nil {
		return cmpl.ErrorFunction(ErrFunction, cmpl.pos+1, faddr[0], ftype.(*core.TypeObject).Func.Params)
	}
	if !isEqualTypes(obj.Result(), ftype.(*core.TypeObject).Func.Result) {
		return cmpl.Error(ErrFnReturn, faddr[0], faddr[1])
	}
	appendExp(cmpl, &core.CmdValue{Value: &core.Fn{Func: obj},
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos + 1)},
		Result:    ftype.(*core.TypeObject)})
	cmpl.newPos = cmpl.pos + 1
	return nil
}

func coUnaryOperator(cmpl *compiler) error {
	return appendExpBuf(cmpl, int(cmpl.getLex().Tokens[cmpl.pos].Type)|tkUnary)
}

func coCtxOperator(cmpl *compiler) error {
	lp := cmpl.getLex()
	if len(lp.Tokens) == cmpl.pos+1 || lp.Tokens[cmpl.pos+1].Type != tkIdent {
		return cmpl.ErrorPos(cmpl.pos+1, ErrName)
	}
	token := getToken(lp, cmpl.pos+1)
	lp.Tokens[cmpl.pos+1].Type = tkStr
	lp.Strings = append(lp.Strings, `#`+token+`#`)
	lp.Tokens[cmpl.pos+1].Index = int32(len(lp.Strings) - 1)
	return appendExpBuf(cmpl, int(lp.Tokens[cmpl.pos].Type)|tkUnary)
}

func appendExpBuf(cmpl *compiler, operation int) error {
	expBuf := ExpBuf{
		Oper:   operation,
		Pos:    cmpl.pos,
		LenExp: len(cmpl.exp),
	}
	if len(cmpl.expbuf) == 0 || operation == tkCallFunc || operation == tkIndex {
		cmpl.expbuf = append(cmpl.expbuf, expBuf)
		return nil
	}
	for len(cmpl.expbuf) > 0 {
		oper := cmpl.expbuf[len(cmpl.expbuf)-1].Oper
		if oper == tkLSBracket {
			if operation == tkRSBracket {
				cmpl.expbuf = cmpl.expbuf[:len(cmpl.expbuf)-1]
				if len(cmpl.expbuf) == 0 || len(cmpl.exp) < 2 {
					return cmpl.Error(ErrNoIndex)
				}
				prevToken := cmpl.expbuf[len(cmpl.expbuf)-1]
				if prevToken.Oper != tkIndex {
					return cmpl.ErrorPos(cmpl.pos-1, ErrVarIndex)
				}
				if cmpl.exp[len(cmpl.exp)-2].GetType() != core.CtVar {
					return cmpl.ErrorPos(prevToken.Pos-1, ErrVarIndex)
				}
				if len(cmpl.exp)-prevToken.LenExp == 0 {
					return cmpl.Error(ErrNoIndex)
				}
				return setIndex(cmpl)
			}
			break
		}
		if oper == tkLPar {
			if operation == tkRPar {
				cmpl.expbuf = cmpl.expbuf[:len(cmpl.expbuf)-1]
				if len(cmpl.expbuf) > 0 {
					prevToken := cmpl.expbuf[len(cmpl.expbuf)-1]
					if prevToken.Oper == tkCallFunc {
						nameFunc := getToken(cmpl.getLex(), prevToken.Pos-1)
						numParams := len(cmpl.exp) - prevToken.LenExp
						params := make([]*core.TypeObject, 0)
						for i := 0; i < numParams; i++ {
							params = append(params, cmpl.exp[prevToken.LenExp+i].GetResult())
						}
						if nameFunc == `$` {
							if len(cmpl.expbuf) == 1 && cmpl.curOwner().ID != core.StackReturn {
								nameFunc = `Command`
							} else {
								nameFunc = `CommandOutput`
							}
						}
						if nameFunc == `?` {
							if numParams != 3 || !isBoolResult(cmpl.exp[prevToken.LenExp]) {
								return cmpl.ErrorPos(prevToken.Pos-1, ErrQuestionPars)
							}
							if params[1] != params[2] {
								return cmpl.ErrorPos(prevToken.Pos-1, ErrQuestion)
							}
							icmd := &core.CmdBlock{ID: uint32(core.StackQuestion),
								Result:    params[1],
								CmdCommon: core.CmdCommon{TokenID: uint32(prevToken.Pos - 1)}}
							for i := prevToken.LenExp; i < len(cmpl.exp); i++ {
								icmd.Children = append(icmd.Children, cmpl.exp[i])
							}
							cmpl.exp = cmpl.exp[:len(cmpl.exp)-numParams]
							cmpl.exp = append(cmpl.exp, icmd)
						} else {
							var (
								result *core.TypeObject
								pobj   core.IObject
								fnVar  core.ICmd
							)
							obj := getFunc(cmpl, nameFunc, params)
							if obj == nil {
								var isMatch bool
								block, ind := findVar(cmpl, nameFunc)
								if block != nil {
									fnVar = &core.CmdVar{Block: block, Index: ind,
										CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)}}
									if typeVar := fnVar.GetResult(); typeVar.Func != nil {
										if len(typeVar.Func.Params) == len(params) {
											isMatch = true
											for i, par := range typeVar.Func.Params {
												if !isEqualTypes(par, params[i]) {
													isMatch = false
													break
												}
											}
										}
										if isMatch {
											result = typeVar.Func.Result
										} else {
											return cmpl.ErrorFunction(ErrFnCall, prevToken.Pos-1,
												typeVar.GetName(), params)
										}
									}
								}
								if !isMatch {
									return cmpl.ErrorFunction(ErrFunction, prevToken.Pos-1, nameFunc, params)
								}
							} else {
								result = obj.Result()
								pobj = obj
							}
							icmd := &core.CmdAnyFunc{CmdCommon: core.CmdCommon{TokenID: uint32(prevToken.Pos - 1)},
								Object: pobj, Result: result, FnVar: fnVar}
							for i := prevToken.LenExp; i < len(cmpl.exp); i++ {
								icmd.Children = append(icmd.Children, cmpl.exp[i])
							}
							cmpl.exp = cmpl.exp[:len(cmpl.exp)-numParams]
							cmpl.exp = append(cmpl.exp, icmd)
						}
						cmpl.expbuf = cmpl.expbuf[:len(cmpl.expbuf)-1]
					}
				}
				return nil
			}
			break
		}
		if operation == tkRPar || operation == tkRSBracket {
			if err := popBuf(cmpl); err != nil {
				return err
			}
			continue
		}
		prev := priority[oper]
		if priority[operation].Priority > prev.Priority {
			cmpl.expbuf = append(cmpl.expbuf, expBuf)
			return nil
		}
		if priority[operation].Priority == prev.Priority && prev.RightLeft {
			cmpl.expbuf = append(cmpl.expbuf, expBuf)
			return nil
		}
		if priority[operation].Priority < prev.Priority ||
			(priority[operation].Priority == prev.Priority && !prev.RightLeft) {
			if err := popBuf(cmpl); err != nil {
				return err
			}
		}
	}
	if len(cmpl.expbuf) == 0 {
		if operation == tkRPar {
			return cmpl.Error(ErrRPar)
		}
		if operation == tkRSBracket {
			return cmpl.Error(ErrRSBracket)
		}
	}
	cmpl.expbuf = append(cmpl.expbuf, expBuf)
	return nil
}

func setIndex(cmpl *compiler) error {
	cmdVar := cmpl.exp[len(cmpl.exp)-2].(*core.CmdVar)
	typeObject := cmdVar.GetResult()
	varIndex := cmpl.getIntType()
	if typeObject.Original == reflect.TypeOf(core.Map{}) {
		varIndex = cmpl.getStrType()
	}
	if typeObject.IndexOf == nil {
		return cmpl.ErrorPos(cmpl.expbuf[len(cmpl.expbuf)-1].Pos-1, ErrSupportIndex,
			typeObject.GetName())
	}
	index := cmpl.exp[len(cmpl.exp)-1]
	if index.GetResult() != varIndex {
		return cmpl.ErrorPos(cmpl.pos, ErrTypeIndex, varIndex.GetName())
	}
	(*cmdVar).Indexes = append((*cmdVar).Indexes, core.CmdRet{Cmd: index, Type: typeObject.IndexOf})
	cmpl.exp = cmpl.exp[:len(cmpl.exp)-1]
	cmpl.expbuf = cmpl.expbuf[:len(cmpl.expbuf)-1]
	return nil
}

func coOperator(cmpl *compiler) error {
	return appendExpBuf(cmpl, int(cmpl.getLex().Tokens[cmpl.pos].Type))
}

func coRPar(cmpl *compiler) error {
	if cmpl.getLex().Tokens[cmpl.pos-1].Type == tkLPar {
		// This is for calling functions without parameters - 2 * sum() + 1
		cmpl.dynamic = &cmState{tkRPar, cmExpOper, nil, nil, cfStay | cfStopBack}
		return nil
	}
	return appendExpBuf(cmpl, int(cmpl.getLex().Tokens[cmpl.pos].Type))
}

func coCallFunc(cmpl *compiler) error {
	return appendExpBuf(cmpl, tkCallFunc)
}

func coComma(cmpl *compiler) error {
	if isCase(cmpl) {
		for len(cmpl.expbuf) > 0 {
			if err := popBuf(cmpl); err != nil {
				return err
			}
		}
		return nil
	}
	for len(cmpl.expbuf) > 0 && cmpl.expbuf[len(cmpl.expbuf)-1].Oper != tkLPar {
		if err := popBuf(cmpl); err != nil {
			return err
		}
	}
	if isInState(cmpl, cmInit, 1) || isInState(cmpl, cmInit, 2) {
		return nil
	}
	if len(cmpl.expbuf) < 2 || (cmpl.expbuf[len(cmpl.expbuf)-1].Oper != tkLPar &&
		cmpl.expbuf[len(cmpl.expbuf)-2].Oper != tkCallFunc) {
		return cmpl.Error(ErrOper)
	}
	return nil
}

func coUnaryPostOperator(cmpl *compiler) error {
	return appendExpBuf(cmpl, int(cmpl.getLex().Tokens[cmpl.pos].Type)|tkUnary|tkPost)
}
