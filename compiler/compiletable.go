// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

const (
	// List of compile states
	cmMain           = iota
	cmConst          // const block
	cmConstDef       // const definitions
	cmConstName      // const identifier
	cmConstListStart // const enum start
	cmConstList      // const enum
	cmRun            // run command
	cmFunc           // func command
	cmParams         // parameters of the function
	cmParam          // getting type name
	cmWantVar
	cmVar // getting var name
	cmWantType
	cmLCurly      // {
	cmBody        // body of the code
	cmExp         // expression
	cmExpIdent    // identifier
	cmExpOper     // expecting operator in expression
	cmElseIf      // elif or else
	cmMustVarType // define variables
	cmVarType     // define variables
	cmQuestion    // ?(condition, exp1, exp2)

	// Flags
	cfSkip     = 0x10000  // stay on the current state
	cfBack     = 0x20000  // go to the previous state
	cfStay     = 0x40000  // stay on the current token
	cfError    = 0x80000  // return error
	cfCallBack = 0x100000 // call func when we come back
)

type compFunc func(*compiler) error

type cmState struct {
	Tokens interface{} // can be one token or []token
	Action int
	Func   compFunc
	BackTo int
}

var (
	preCompile = map[int][]cmState{
		cmMain: {
			{tkDefault, cfError | ErrDecl, nil, 0},
			{tkLine, cfSkip, nil, 0},
			{tkConst, cmConst, nil, 1},
			{tkRun, cmRun | cfCallBack, coRun, 1},
			{tkFunc, cmFunc | cfCallBack, coFunc, 1},
		},
		cmConst: {
			{tkDefault, cmExp | cfCallBack | cfStay, coConstEnum, 0},
			{tkLine, cfSkip, nil, 0},
			{tkLCurly, cmConstDef, nil, 0},
		},
		cmConstDef: {
			{tkDefault, cfError | ErrName, nil, 0},
			{tkIdent, cmConstName, coConst, 1},
			{tkLine, cfSkip, nil, 0},
			{tkRCurly, cfBack, nil, 0},
		},
		cmConstName: {
			{tkDefault, cfError | ErrMustAssign, nil, 0},
			{tkAssign, cmExp | cfCallBack, coConstExp, 0},
		},
		cmConstListStart: {
			{tkDefault, cfError | ErrLCurly, nil, 0},
			{tkLine, cfSkip, nil, 0},
			{tkLCurly, cmConstList, nil, 0},
		},
		cmConstList: {
			{tkDefault, cfError | ErrName, nil, 0},
			{tkIdent, cfSkip, coConstList, 0},
			{tkLine, cfSkip, nil, 0},
			{tkRCurly, cfBack, nil, 0},
		},
		cmRun: {
			{tkDefault, cfError | ErrLCurly, nil, 0},
			{tkIdent, cfSkip, coRunName, 0},
			{tkLine, cfSkip, nil, 0},
			{tkLCurly, cmBody, nil, 0},
		},
		cmFunc: {
			{tkDefault, cfError | ErrName, nil, 0},
			{tkIdent, cmParams, coFuncName, 0},
			{tkLine, cfSkip, nil, 0},
		},
		cmParams: {
			{tkDefault, cfError | ErrLCurly, nil, 0},
			{tkIdent, cmLCurly, coRetType, 0},
			{tkLPar, cmParam, nil, 1},
			{tkLCurly, cfStay | cmLCurly, coFuncStart, 0},
			{tkLine, cfSkip, nil, 0},
		},
		cmParam: {
			{tkDefault, cfError | ErrType, nil, 0},
			{tkIdent, cmWantVar, coType, 1},
			{tkRPar, cfBack, nil, 0},
			{tkLine, cfSkip, nil, 0},
		},
		cmWantVar: {
			{tkDefault, cfError | ErrName, nil, 0},
			{tkIdent, cmVar, coVar, 0},
			{tkLine, cfSkip, nil, 0},
		},
		cmVar: {
			{tkDefault, cfError | ErrName, nil, 0},
			{tkIdent, cfSkip, coVar, 0},
			{tkComma, cmWantType, nil, 0},
			{tkRPar, cfStay | cfBack, nil, 0},
			{tkLine, cfSkip, nil, 0},
		},
		cmWantType: {
			{tkDefault, cfError | ErrType, nil, 0},
			{tkIdent, cfStay | cfBack, nil, 0},
			{tkLine, cfSkip, nil, 0},
		},
		cmLCurly: {
			{tkDefault, cfError | ErrLCurly, nil, 0},
			{tkLine, cfSkip, nil, 0},
			{tkLCurly, cmBody, nil, 0},
		},
		cmBody: {
			{tkDefault, cmExp | cfStay, coExpStart, 1},
			{tkLine, cfSkip, nil, 0},
			{tkLCurly, cfError | ErrExp, nil, 0},
			{tkRCurly, cfBack, nil, 0},
			{tkType, cmMustVarType, coVarType, 1},
			{tkIf, cmExp | cfCallBack, coIf, 1},
			{tkWhile, cmExp | cfCallBack, coWhile, 1},
			{tkFor, cmExp | cfCallBack, coFor, 1},
			{tkReturn, cmExp | cfCallBack, coReturn, 1},
		},
		cmExp: {
			{tkDefault, cfError | ErrValue, nil, 0},
			{[]int{tkInt, tkFalse, tkTrue, tkStr, tkChar}, cmExpOper, coPush, 1},
			{[]int{tkSub, tkMul, tkNot, tkBitNot, tkInc, tkDec}, cfSkip, coUnaryOperator, 0},
			{[]int{tkLPar, tkRPar}, cfSkip, coOperator, 0},
			{[]int{tkLSBracket, tkRSBracket}, cfSkip, coOperator, 0},
			{tkLine, cfBack, coExpEnd, 0},
			{[]int{tkLCurly, tkRCurly}, cfStay | cfBack, coExpEnd, 0},
			{tkIdent, cmExpIdent, nil, 1},
			{tkEnv, cmExpOper, coExpEnv, 1},
			{tkQuestion, cmExpIdent, nil, 1},
		},
		cmExpIdent: {
			{tkDefault, cmExpOper | cfStay, coExpVar, 0},
			{tkLPar, cfBack | cfStay, coCallFunc, 0},
			{tkLSBracket, cfBack | cfStay, coIndex, 0},
		},
		cmExpOper: {
			{tkDefault, cfError | ErrOper, nil, 0},
			{tkStrExp, cfBack, coOperator, 0},
			{[]int{tkAdd, tkDiv, tkMod, tkMul, tkSub, tkEqual, tkNotEqual, tkGreater, tkGreaterEqual,
				tkLess, tkLessEqual, tkAssign, tkOr, tkAnd, tkBitOr, tkBitAnd, tkBitXor, tkLShift,
				tkRShift, tkAddEq, tkSubEq, tkMulEq, tkDivEq, tkModEq, tkLShiftEq, tkRShiftEq,
				tkBitAndEq, tkBitOrEq, tkBitXorEq, tkRange}, cfBack,
				coOperator, 0},
			{[]int{tkInc, tkDec}, cfSkip, coUnaryPostOperator, 0},
			{[]int{tkRPar, tkRSBracket}, cfSkip, coOperator, 0},
			{tkComma, cfBack, coComma, 0},
			{[]int{tkLCurly, tkLine, tkRCurly}, cfStay | cfBack, nil, 0},
		},
		cmElseIf: {
			{tkDefault, cfBack | cfStay, coIfEnd, 0},
			{tkLine, cfSkip, nil, 0},
			{tkElse, cmLCurly, coElse, 0},
			{tkElif, cmExp | cfCallBack, coElif, 1},
		},
		cmMustVarType: {
			{tkDefault, cfError | ErrName, nil, 0},
			{tkIdent, cmVarType, coVarExp, 0},
		},
		cmVarType: {
			{tkDefault, cfError | ErrName, nil, 0},
			{tkIdent, cfSkip, coVar, 0},
			{tkLine, cfBack, nil, 0},
		},
		cmQuestion: {
			{tkDefault, cfError | ErrCompiler, nil, 0},
			{tkLPar, cfBack | cfStay, coCallFunc, 0},
		},
	}
	compileTable [][tkDefault]*cmState
)

func makeCompileTable() {
	compileTable = make([][tkDefault]*cmState, len(preCompile))

	for state, items := range preCompile {
		for i, item := range items {
			ptr := &preCompile[state][i]
			switch v := item.Tokens.(type) {
			case int:
				if v == tkDefault {
					for i := 0; i < tkDefault; i++ {
						compileTable[state][i] = ptr
					}
				} else {
					compileTable[state][v] = ptr
				}
			case []int:
				for _, id := range v {
					compileTable[state][id] = ptr
				}
			default:
				panic(`corrupted preCompile table`)
			}
		}

	}
}
