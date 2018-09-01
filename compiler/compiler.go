// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

// Compiler contains information of the compilation process
type compiler struct {
	vm       *core.VirtualMachine
	unit     *core.Unit
	owners   []core.ICmd
	exp      []core.ICmd
	expbuf   []ExpBuf
	lexems   []int // stack of lexeme
	runID    int
	pos      int // current position
	newState int // new state from callback function
	callback bool
	states   *[]StateStack
	curType  *core.TypeObject // the current type of parameters or variables
	curConst string
	expConst core.ICmd // expression for constants
	curIota  int64     // current iota
}

// StateStack is used for storing a sequence of states
type StateStack struct {
	Origin *cmState
	Pos    int
	State  int
}

// Priority is a structure for operations in expressions
type Priority struct {
	Priority  int
	RightLeft bool
	Name      string
}

// ExpBuf is a structure for buffer of expression operations
type ExpBuf struct {
	Oper   int
	Pos    int
	LenExp int
}

var (
	priority = map[int]Priority{
		tkAssign:       {5, true, `Assign`},
		tkAddEq:        {5, true, `AssignAdd`},
		tkSubEq:        {5, true, `AssignSub`},
		tkMulEq:        {5, true, `AssignMul`},
		tkDivEq:        {5, true, `AssignDiv`},
		tkModEq:        {5, true, `AssignMod`},
		tkLShiftEq:     {5, true, `AssignLShift`},
		tkRShiftEq:     {5, true, `AssignRShift`},
		tkBitAndEq:     {5, true, `AssignBitAnd`},
		tkBitOrEq:      {5, true, `AssignBitOr`},
		tkBitXorEq:     {5, true, `AssignBitXor`},
		tkAnd:          {7, false, ``},
		tkOr:           {8, false, ``},
		tkBitOr:        {9, false, `BitOr`},
		tkBitXor:       {10, false, `BitXor`},
		tkBitAnd:       {11, false, `BitAnd`},
		tkEqual:        {12, false, `Equal`},
		tkNotEqual:     {12, false, `Equal`},
		tkLess:         {12, false, `Less`},
		tkLessEqual:    {12, false, `Greater`},
		tkGreater:      {12, false, `Greater`},
		tkGreaterEqual: {12, false, `Less`},
		tkLShift:       {14, false, `LShift`},
		tkRShift:       {14, false, `RShift`},
		tkAdd:          {15, false, `Add`},
		tkSub:          {15, false, `Sub`},
		tkDiv:          {20, false, `Div`},
		tkMod:          {20, false, `Mod`},
		tkMul:          {20, false, `Mul`},
		tkInc | tkUnary | tkPost: {29, false, ``},
		tkDec | tkUnary | tkPost: {29, false, ``},
		tkBitNot | tkUnary:       {30, true, `BitNot`},
		tkSub | tkUnary:          {30, true, `Sign`},
		tkNot | tkUnary:          {30, true, `Not`},
		tkMul | tkUnary:          {30, true, `Len`},
		tkInc | tkUnary:          {30, true, ``},
		tkDec | tkUnary:          {30, true, ``},
		tkStrExp:                 {35, false, `ExpStr`},
		tkLPar:                   {50, true, ``},
		tkRPar:                   {50, true, ``},
	}
)

func init() {
	makeParseTable()
	makeCompileTable()
}

func (cmpl *compiler) curOwner() *core.CmdBlock {
	return cmpl.owners[len(cmpl.owners)-1].(*core.CmdBlock)
}

// Compile compiles the source code
func Compile(vm *core.VirtualMachine, input, path string) error {
	var (
		state int
	)

	lp, errID := LexParsing([]rune(input))
	lp.Path = path
	cmpl := &compiler{
		vm: vm,
		unit: &core.Unit{
			Objects: make([]core.IObject, 0),
			Names:   make(map[string]core.IObject),
			Lexeme:  []*core.Lex{lp},
		},
		lexems:  []int{0}, // added lp in Lexeme
		runID:   core.Undefined,
		owners:  make([]core.ICmd, 0, 128),
		exp:     make([]core.ICmd, 0, 128),
		expbuf:  make([]ExpBuf, 0, 128),
		curIota: core.NotIota,
	}
	if len(lp.Tokens) == 0 {
		return cmpl.Error(ErrEmptyCode)
	}
	if errID != 0 {
		cmpl.pos = len(lp.Tokens) - 1
		return cmpl.Error(errID)
	}

	stackState := make([]StateStack, 0, 32)
main:
	for i := 0; i < len(lp.Tokens); i++ {
		cmpl.pos = i
		token := lp.Tokens[i]
		if state == cmBody && token.Type == tkIdent {
			obj, _ := getType(cmpl)
			if obj != nil {
				token.Type = tkType
			}
		}
		next := compileTable[state][token.Type]
		cmpl.states = &stackState
		cmpl.newState = 0
		//fmt.Printf("NEXT i=%d state=%d token=%d v=%v flag=%x nextstate=%v\r\n", i, state, token.Type,
		//	getToken(cmpl.getLex(), i), next.Action&0xff0000, next.Action&0xffff)
		flag := next.Action & 0xff0000
		if flag&cfError != 0 {
			return cmpl.Error(next.Action & 0xffff)
		}
		if (state == cmExp || state == cmExpOper) && token.Type == tkLine {
			if state == cmExp && lp.Tokens[i-1].Type >= tkAdd && lp.Tokens[i-1].Type <= tkComma {
				continue
			}
			for _, expBuf := range cmpl.expbuf {
				if expBuf.Oper == tkLPar {
					continue main
				}
			}
		}
		if next.Func != nil {
			if err := next.Func(cmpl); err != nil {
				return err
			}
			if cmpl.newState != 0 {
				state = cmpl.newState & 0xffff
				stackState = append(stackState, StateStack{Origin: next, Pos: i, State: state})
				if cmpl.newState&cfStay != 0 {
					i--
				}
				continue
			}
		}
		if flag&cfSkip != 0 {
			continue
		}
		if flag&cfStay != 0 {
			i--
		}
		if flag&cfBack != 0 {
			if len(stackState) == 0 {
				return cmpl.Error(ErrCompiler, `Compile`)
			}
			for len(stackState) > 0 {
				prev := stackState[len(stackState)-1]
				state = prev.State
				cmpl.newState = 0
				if prev.Origin.Action&cfCallBack != 0 {
					cmpl.callback = true
					cmpl.pos = prev.Pos
					if err := prev.Origin.Func(cmpl); err != nil {
						return err
					}
					if cmpl.newState > 0 {
						if cmpl.newState&cfStay != 0 {
							i--
						}
						state = cmpl.newState & 0xffff
					}
					cmpl.callback = false
				}
				if cmpl.newState == 0 {
					stackState = stackState[:len(stackState)-1]
				}
				if prev.Origin.BackTo == 1 || cmpl.newState != 0 {
					break
				}
			}
			continue
		}
		stackState = append(stackState, StateStack{Origin: next, Pos: i, State: state})
		state = next.Action & 0xffff
	}
	if len(stackState) > 0 {
		cmpl.pos = stackState[len(stackState)-1].Pos + 1
		return cmpl.Error(ErrEnd)
	}

	if cmpl.runID != core.Undefined {
		cmpl.unit.Type = core.UnitRun
		cmpl.unit.RunID = cmpl.runID
		if len(cmpl.unit.Name) == 0 {
			cmpl.unit.Name = path
		}
		if unitIndex, ok := vm.Names[cmpl.unit.Name]; ok {
			if vm.Units[unitIndex].Lexeme[0].Path != path {
				return cmpl.Error(ErrLink, cmpl.unit.Name)
			}
			vm.Units[unitIndex] = cmpl.unit
			vm.Compiled = unitIndex
		} else {
			vm.Units = append(vm.Units, cmpl.unit)
			vm.Compiled = len(vm.Units) - 1
			vm.Names[cmpl.unit.Name] = vm.Compiled
		}
	} else {
		cmpl.unit.Type = core.UnitPackage
		// TODO: append to vm.Packages
	}

	return nil
}

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

func (cmpl *compiler) getLex() *core.Lex {
	return cmpl.unit.Lexeme[len(cmpl.unit.Lexeme)-1]
}

func coRun(cmpl *compiler) error {
	if cmpl.callback {
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
	if cmpl.runID != core.Undefined {
		return cmpl.Error(ErrRun)
	}
	cmpl.runID = newFunc(cmpl, `run`)
	return nil
}

func coRunName(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if len(cmpl.unit.Name) != 0 {
		cmpl.newState = cmLCurly
		return coRetType(cmpl)
	} else if _, err := getType(cmpl); err != nil {
		cmpl.unit.Name = token
	} else {
		cmpl.newState = cmLCurly
		return coRetType(cmpl)
	}
	return nil
}

func getFunc(cmpl *compiler, name string, params []*core.TypeObject) (obj core.IObject) {
	checkUnit := func(unit *core.Unit) core.IObject {
		obj = unit.Names[name]
		for obj != nil {
			objPars := obj.GetParams()
			if len(params) == len(objPars) {
				equal := true
				for i, typeParam := range objPars {
					if typeParam != params[i] {
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
	if obj = checkUnit(cmpl.vm.StdLib()); obj == nil {
		obj = checkUnit(cmpl.unit)
	}
	return obj
}

func coFunc(cmpl *compiler) error {
	if cmpl.callback {
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
	return nil
}

func coFuncStart(cmpl *compiler) error {
	funcObj := cmpl.unit.Objects[len(cmpl.unit.Objects)-1].(*core.FuncObject)
	funcObj.Block.ParCount = len(funcObj.Block.Vars)
	if obj := getFunc(cmpl, funcObj.Name, funcObj.GetParams()); obj != nil &&
		obj != cmpl.unit.Objects[len(cmpl.unit.Objects)-1] {
		return cmpl.ErrorFunction(ErrFuncExists, int(funcObj.Block.TokenID), funcObj.Name,
			funcObj.GetParams())
	}
	return nil
}

func coFuncName(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if isCapital(token) {
		return cmpl.Error(ErrCapitalLetters)
	}
	newFunc(cmpl, token)
	return nil
}

func getType(cmpl *compiler) (obj core.IObject, err error) {
	findType := func(src core.IObject) {
		for obj = src; obj != nil && obj.GetType() != core.ObjType; {
			obj = obj.GetNext()
		}
	}
	token := getToken(cmpl.getLex(), cmpl.pos)
	findType(cmpl.unit.Names[token])
	if obj == nil {
		findType(cmpl.vm.StdLib().Names[token])
	}
	if obj == nil {
		return nil, cmpl.Error(ErrType)
	}
	return
}

func coType(cmpl *compiler) error {
	obj, err := getType(cmpl)
	if err != nil {
		return err
	}
	cmpl.curType = obj.(*core.TypeObject)
	return nil
}

func coVarType(cmpl *compiler) error {
	obj, err := getType(cmpl)
	if err != nil {
		return err
	}
	cmpl.curType = obj.(*core.TypeObject)
	return nil
}

func coVar(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if isCapital(token) {
		return cmpl.Error(ErrCapitalLetters)
	}
	if cmpl.vm.StdLib().Names[token] != nil ||
		cmpl.unit.Names[token] != nil {
		return cmpl.Error(ErrUsedName, token)
	}
	block := cmpl.curOwner()
	for block != nil {
		if _, ok := block.VarNames[token]; ok {
			return cmpl.Error(ErrUsedName, token)
		}
		block = block.Parent
	}

	block = cmpl.owners[0].(*core.CmdBlock)
	if block.VarNames == nil {
		block.VarNames = make(map[string]int)
	}
	block.VarNames[token] = len(block.Vars)
	block.Vars = append(block.Vars, cmpl.curType)
	return nil
}

func coConst(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if !isCapital(token) {
		return cmpl.Error(ErrConstName)
	}
	cmpl.curConst = token
	return nil
}

func coConstEnum(cmpl *compiler) error {
	if cmpl.callback {
		if cmpl.expConst == nil {
			owner := cmpl.curOwner()
			cmpl.expConst = owner.Children[0]
			cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
			cmpl.newState = cmConstListStart
		} else { // const finishes
			cmpl.curIota = core.NotIota
		}
		return nil
	}
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackBlock, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmpl.owners = append(cmpl.owners, &cmd)
	cmpl.curIota = 0
	cmpl.expConst = nil
	return nil
}

func coConstList(cmpl *compiler) error {
	if err := coConst(cmpl); err != nil {
		return err
	}
	constObj := &core.ConstObject{
		Object: core.Object{
			Name:  cmpl.curConst,
			LexID: len(cmpl.unit.Lexeme) - 1,
			Unit:  cmpl.unit,
		},
		Redefined: false,
		Exp:       cmpl.expConst,
		Return:    cmpl.expConst.GetResult(),
		Iota:      cmpl.curIota,
	}
	cmpl.curIota++
	cmpl.unit.Objects = append(cmpl.unit.Objects, constObj)
	if curName := cmpl.unit.Names[cmpl.curConst]; curName == nil {
		cmpl.unit.Names[cmpl.curConst] = constObj
	} else {
		curName.SetNext(constObj)
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

func appendCmd(cmpl *compiler, cmd core.ICmd) {
	owner := cmpl.curOwner()
	if cmd.GetType() == core.CtStack {
		cmd.(*core.CmdBlock).Parent = owner
	}
	owner.Children = append(owner.Children, cmd)
}

func appendExp(cmpl *compiler, cmd core.ICmd) {
	cmpl.exp = append(cmpl.exp, cmd)
}

func appendExpBuf(cmpl *compiler, operation int) error {
	expBuf := ExpBuf{
		Oper:   operation,
		Pos:    cmpl.pos,
		LenExp: len(cmpl.exp),
	}
	if len(cmpl.expbuf) == 0 || operation == tkCallFunc {
		cmpl.expbuf = append(cmpl.expbuf, expBuf)
		return nil
	}
	for len(cmpl.expbuf) > 0 {
		oper := cmpl.expbuf[len(cmpl.expbuf)-1].Oper
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
							obj := getFunc(cmpl, nameFunc, params)
							if obj == nil {
								return cmpl.ErrorFunction(ErrFunction, prevToken.Pos-1, nameFunc, params)
							}
							icmd := &core.CmdAnyFunc{CmdCommon: core.CmdCommon{TokenID: uint32(prevToken.Pos - 1)},
								Object: obj, Result: obj.Result()}
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
		if operation == tkRPar {
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
	if operation == tkRPar && len(cmpl.expbuf) == 0 {
		return cmpl.Error(ErrRPar)
	}

	cmpl.expbuf = append(cmpl.expbuf, expBuf)
	return nil
}

func popBuf(cmpl *compiler) error {
	expBuf := cmpl.expbuf[len(cmpl.expbuf)-1]
	prior := priority[expBuf.Oper]
	obj := cmpl.vm.StdLib().Names[prior.Name]
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
		for obj != nil {
			params := obj.GetParams()
			if len(params) == 2 && left.GetResult() == params[0] &&
				right.GetResult() == params[1] {
				break
			}
			obj = obj.GetNext()
		}
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
			if left.GetObject() == cmpl.vm.StdLib().Names[`GetEnv`] {
				setEnv := cmpl.vm.StdLib().Names[`SetEnv`]
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
		for obj != nil {
			params := obj.GetParams()
			if len(params) == 2 && left.GetResult() == params[0] &&
				right.GetResult() == params[1] {
				break
			}
			obj = obj.GetNext()
		}
		if obj == nil {
			return cmpl.ErrorFunction(ErrFunction, expBuf.Pos, prior.Name, []*core.TypeObject{
				left.GetResult(), right.GetResult()})
		}
		//			right = &core.CmdBinary{CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
		//				Object: obj, Result: obj.Result(), Left: left, Right: right}
		//		}
		icmd := &core.CmdBlock{ID: core.StackAssign, Object: obj,
			Result: right.GetResult(), CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Children: []core.ICmd{left, right}}
		cmpl.exp[len(cmpl.exp)-2] = icmd
		cmpl.exp = cmpl.exp[:len(cmpl.exp)-1]
	case tkAdd, tkSub, tkMul, tkMod, tkDiv, tkEqual, tkNotEqual, tkLess, tkLessEqual, tkGreater,
		tkGreaterEqual, tkBitOr, tkBitXor, tkBitAnd, tkLShift, tkRShift:
		if len(cmpl.exp) < 2 {
			return cmpl.Error(ErrValue)
		}
		right := cmpl.exp[len(cmpl.exp)-1]
		left := cmpl.exp[len(cmpl.exp)-2]
		for obj != nil {
			params := obj.GetParams()
			if len(params) == 2 && left.GetResult() == params[0] &&
				right.GetResult() == params[1] {
				break
			}
			obj = obj.GetNext()
		}
		if obj == nil {
			return cmpl.ErrorFunction(ErrFunction, expBuf.Pos, prior.Name, []*core.TypeObject{
				left.GetResult(), right.GetResult()})
		}
		//		fmt.Println(`Bin`, prior.Name, obj, left.GetResult(), right.GetResult(), obj.GetParams())
		icmd := &core.CmdBinary{CmdCommon: core.CmdCommon{TokenID: uint32(expBuf.Pos)},
			Object: obj, Result: obj.Result(), Left: left, Right: right}
		if expBuf.Oper == tkNotEqual || expBuf.Oper == tkLessEqual || expBuf.Oper == tkGreaterEqual {
			objNot := cmpl.vm.StdLib().Names[`Not`]
			for objNot != nil {
				params := objNot.GetParams()
				if len(params) == 1 && obj.Result() == params[0] {
					break
				}
				objNot = objNot.GetNext()
			}
			if objNot == nil {
				return cmpl.Error(ErrCompiler, `popBuf not`)
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
	case tkSub | tkUnary, tkMul | tkUnary, tkNot | tkUnary, tkBitNot | tkUnary:
		if len(cmpl.exp) == 0 {
			return cmpl.Error(ErrValue)
		}
		top := cmpl.exp[len(cmpl.exp)-1]
		for obj != nil {
			params := obj.GetParams()
			if len(params) == 1 && top.GetResult() == params[0] {
				break
			}
			obj = obj.GetNext()
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
	default:
		return cmpl.Error(ErrCompiler, fmt.Sprintf(`popBuf unknown token %d`, expBuf.Oper))
	}
	cmpl.expbuf = cmpl.expbuf[:len(cmpl.expbuf)-1]
	return nil
}

func coExpStart(cmpl *compiler) error {
	cmpl.exp = cmpl.exp[:0]
	cmpl.expbuf = cmpl.expbuf[:0]
	return nil
}

func coExpEnd(cmpl *compiler) error {
	for len(cmpl.expbuf) > 0 {
		if err := popBuf(cmpl); err != nil {
			return err
		}
	}
	if len(cmpl.exp) > 1 {
		return cmpl.Error(ErrCompiler, `coExpEnd`)
	}
	if len(cmpl.exp) > 0 {
		cmpl.curOwner().Children = append(cmpl.curOwner().Children, cmpl.exp[len(cmpl.exp)-1])
	}
	return nil
}

func coWhile(cmpl *compiler) error {
	if cmpl.callback {
		cmd := cmpl.curOwner()
		if cmd.ID == core.StackWhile {
			if len(cmd.Children) == 1 {
				if !isBoolResult(cmd.Children[0]) {
					cmpl.pos = cmd.Children[0].GetToken()
					return cmpl.Error(ErrBoolExp)
				}
				cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
					CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
				cmd.Children = append(cmd.Children, &cmdIf)
				cmpl.owners = append(cmpl.owners, &cmdIf)
				cmpl.newState = cmLCurly // | cfStay
			}
		} else {
			cmpl.owners = cmpl.owners[:len(cmpl.owners)-2]
		}
		return nil
	}
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackWhile, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coIf(cmpl *compiler) error {
	if cmpl.callback {
		cmd := cmpl.curOwner()
		if cmd.ID == core.StackIf {
			if len(cmd.Children) == 1 {
				if !isBoolResult(cmd.Children[0]) {
					cmpl.pos = cmd.Children[0].GetToken()
					return cmpl.Error(ErrBoolExp)
				}
				cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
					CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
				cmd.Children = append(cmd.Children, &cmdIf)
				cmpl.owners = append(cmpl.owners, &cmdIf)
				cmpl.newState = cmLCurly // | cfStay
			} else {
				cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
			}
		} else {
			cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
			cmpl.newState = cmElseIf
		}
		return nil
	}
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackIf, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coElse(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID != core.StackIf {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		cmd = cmpl.curOwner()
	}
	cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmd.Children = append(cmd.Children, &cmdIf)
	cmpl.owners = append(cmpl.owners, &cmdIf)
	return nil
}

func coElif(cmpl *compiler) error {
	if cmpl.callback {
		cmd := cmpl.curOwner()
		if cmd.ID == core.StackIf {
			if !isBoolResult(cmd.Children[len(cmd.Children)-1]) {
				cmpl.pos = cmd.Children[len(cmd.Children)-1].GetToken()
				return cmpl.Error(ErrBoolExp)
			}
			cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
			cmd.Children = append(cmd.Children, &cmdIf)
			cmpl.owners = append(cmpl.owners, &cmdIf)
			cmpl.newState = cmLCurly
		} else {
			cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		}
		return nil
	}
	coExpStart(cmpl)
	return nil
}

func coIfEnd(cmpl *compiler) error {
	if cmpl.curOwner().ID != core.StackIf {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	}
	return nil
}

func coVarExp(cmpl *compiler) error {
	if err := coVar(cmpl); err != nil {
		return err
	}
	if cmpl.getLex().Tokens[cmpl.pos+1].Type == tkAssign {
		coExpStart(cmpl)
		cmpl.newState = cmExp | cfStay
	}
	return nil
}

func coReturn(cmpl *compiler) error {
	if cmpl.callback {
		owner := cmpl.curOwner()
		funcObj := cmpl.unit.Objects[len(cmpl.unit.Objects)-1].(*core.FuncObject)
		switch len(owner.Children) {
		case 0:
			if funcObj.Block.Result != nil {
				return cmpl.Error(ErrCompiler, `coReturn 0`)
			}
		case 1:
			if funcObj.Block.Result == nil {
				return cmpl.Error(ErrReturn)
			}
			if funcObj.Block.Result != owner.Children[0].GetResult() {
				return cmpl.Error(ErrReturnType)
			}
		default:
			return cmpl.Error(ErrCompiler, `coReturn 1`)
		}
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		return nil
	}
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackReturn, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coConstExp(cmpl *compiler) error {
	if cmpl.callback {
		owner := cmpl.curOwner()
		constObj := &core.ConstObject{
			Object: core.Object{
				Name:  cmpl.curConst,
				LexID: len(cmpl.unit.Lexeme) - 1,
				Unit:  cmpl.unit,
			},
			Redefined: false,
			Exp:       owner.Children[0],
			Return:    owner.Children[0].GetResult(),
			Iota:      core.NotIota,
		}

		cmpl.unit.Objects = append(cmpl.unit.Objects, constObj)
		if curName := cmpl.unit.Names[cmpl.curConst]; curName == nil {
			cmpl.unit.Names[cmpl.curConst] = constObj
		} else {
			curName.SetNext(constObj)
		}
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		return nil
	}
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackBlock, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmpl.owners = append(cmpl.owners, &cmd)
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
	case tkInt, tkIntHex, tkIntOct:
		if v, err = strconv.ParseInt(token, 0, 64); err != nil {
			return cmpl.Error(ErrOutOfRange, token)
		}
		vType = `int`
	case tkFalse, tkTrue:
		v = lp.Tokens[cmpl.pos].Type == tkTrue
		vType = `bool`
	case tkStr:
		v = lp.Strings[lp.Tokens[cmpl.pos].Index]
		if token[0] == '"' {
			v, err = strconv.Unquote(`"` + strings.Replace(strings.Replace(v.(string),
				"\n", `\n`, -1), "\r", `\r`, -1) + `"`)
			if err != nil {
				return cmpl.Error(ErrDoubleQuotes)
			}
		}
		vType = `str`
	}
	appendExp(cmpl, &core.CmdValue{Value: v,
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		Result:    cmpl.vm.StdLib().Names[vType].(*core.TypeObject)})
	return nil
}

func coUnaryOperator(cmpl *compiler) error {
	return appendExpBuf(cmpl, int(cmpl.getLex().Tokens[cmpl.pos].Type)|tkUnary)
}

func coUnaryPostOperator(cmpl *compiler) error {
	return appendExpBuf(cmpl, int(cmpl.getLex().Tokens[cmpl.pos].Type)|tkUnary|tkPost)
}

func coOperator(cmpl *compiler) error {
	return appendExpBuf(cmpl, int(cmpl.getLex().Tokens[cmpl.pos].Type))
}

func coCallFunc(cmpl *compiler) error {
	return appendExpBuf(cmpl, tkCallFunc)
}

func coExpVar(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos-1)
	if isCapital(token) {
		var (
			constObj core.IObject
			ok       bool
		)
		if token == core.ConstIota && cmpl.curIota == core.NotIota {
			return cmpl.ErrorPos(cmpl.pos-1, ErrIota)
		}
		if constObj, ok = cmpl.unit.Names[token]; !ok {
			constObj, _ = cmpl.vm.StdLib().Names[token]
		}
		if constObj != nil {
			appendExp(cmpl, &core.CmdConst{
				Object:    constObj,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)},
			})
		} else {
			return cmpl.ErrorPos(cmpl.pos-1, ErrUnknownIdent, token)
		}
	} else {
		block := cmpl.curOwner()
		for block != nil {
			if ind, ok := block.VarNames[token]; ok {
				appendExp(cmpl, &core.CmdVar{Block: block, Index: ind,
					CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos - 1)}})
				break
			}
			block = block.Parent
		}
		if block == nil {
			return cmpl.ErrorPos(cmpl.pos-1, ErrUnknownIdent, token)
		}
	}
	return nil
}

func coExpEnv(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	getEnv := cmpl.vm.StdLib().Names[`GetEnv`]
	icmd := &core.CmdValue{Value: token[1:],
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		Result:    getEnv.Result()}
	appendExp(cmpl, &core.CmdUnary{CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		Object: getEnv, Result: getEnv.Result(), Operand: icmd})
	return nil
}

func coComma(cmpl *compiler) error {
	for len(cmpl.expbuf) > 0 && cmpl.expbuf[len(cmpl.expbuf)-1].Oper != tkLPar {
		if err := popBuf(cmpl); err != nil {
			return err
		}
	}
	if len(cmpl.expbuf) < 2 || (cmpl.expbuf[len(cmpl.expbuf)-1].Oper != tkLPar &&
		cmpl.expbuf[len(cmpl.expbuf)-2].Oper != tkCallFunc) {
		return cmpl.Error(ErrCompiler, `Comma`)
	}
	return nil
}
