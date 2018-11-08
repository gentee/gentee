// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

const (
	// List of compile states
	cmMain     = iota + 1
	cmRun      // run command
	cmLCurly   // {
	cmBody     // body of the code
	cmExp      // expression
	cmExpIdent // identifier
	cmExpOper  // expecting operator in expression
	cmElseIf   // elif or else
	cmFunc     // func command
	cmParams   // parameters of the function
	cmParam    // getting type name
	cmWantVar
	cmVar // getting var name
	cmWantType
	cmMustVarType    // define variables
	cmVarType        // define variables
	cmConst          // const block
	cmConstDef       // const definitions
	cmConstName      // const identifier
	cmConstListStart // const enum start
	cmConstList      // const enum
	cmInit           // initializing array or map
	cmStruct         // struct definition
	cmStructDef      // struct body
	cmStructFields   // struct fields
	cmStructName     // struct the name of the field

	cmBack // go to back

	// Flags
	cfStopBack = 0x10000 // stop when go to back
	cfStay     = 0x40000 // stay on the current token
)

type compFunc func(*compiler) error

type cmState struct {
	Tokens   interface{} // can be one token or []token
	State    int         // new state
	Func     compFunc
	Callback compFunc
	Flags    int
}

var (
	preCompile = map[int][]cmState{
		cmMain: {
			{tkDefault, ErrDecl, coError, nil, 0},
			{tkLine, 0, nil, nil, 0},
			{tkRun, cmRun, coRun, coRunBack, cfStopBack},
			{tkConst, cmConst, nil, nil, cfStopBack},
			{tkFunc, cmFunc, nil, coFuncBack, cfStopBack},
			{tkStruct, cmStruct, nil, nil, cfStopBack},
		},
		cmRun: {
			{tkDefault, ErrLCurly, coError, nil, 0},
			{tkIdent, 0, coRunName, nil, 0},
			{tkLine, 0, nil, nil, 0},
			{tkLCurly, cmBody, nil, nil, 0},
		},
		cmLCurly: {
			{tkDefault, ErrLCurly, coError, nil, 0},
			{tkLine, 0, nil, nil, 0},
			{tkLCurly, cmBody, nil, nil, 0},
		},
		cmBody: {
			{tkDefault, cmExp, coExpStart, nil, cfStay | cfStopBack},
			{tkLine, 0, nil, nil, 0},
			{tkLCurly, ErrExp, coError, nil, 0},
			{tkRCurly, cmBack, nil, nil, 0},
			{tkType, cmMustVarType, coVarType, nil, cfStopBack},
			{tkIf, cmExp, coIf, coIfBack, cfStopBack},
			{tkWhile, cmExp, coWhile, coWhileBack, cfStopBack},
			{tkFor, cmExp, coFor, coForBack, cfStopBack},
			{tkReturn, cmExp, coReturn, coReturnBack, cfStopBack},
		},
		cmExp: {
			{tkDefault, ErrValue, coError, nil, 0},
			{[]int{tkInt, tkFloat, tkFalse, tkTrue, tkStr, tkChar}, cmExpOper, coPush, nil, cfStopBack},
			{[]int{tkSub, tkMul, tkNot, tkBitNot, tkInc, tkDec}, 0, coUnaryOperator, nil, 0},
			{[]int{tkLPar, tkRPar}, 0, coOperator, nil, 0},
			{[]int{tkLSBracket, tkRSBracket}, 0, coOperator, nil, 0},
			{tkLine, cmBack, coExpEnd, nil, 0},
			{[]int{tkLCurly, tkRCurly, tkColon}, cmBack, coExpEnd, nil, cfStay},
			{tkIdent, cmExpIdent, nil, nil, cfStopBack},
			{tkEnv, cmExpOper, coExpEnv, nil, cfStopBack},
			{tkQuestion, cmExpIdent, nil, nil, cfStopBack},
		},
		cmExpIdent: {
			{tkDefault, cmExpOper, coExpVar, nil, cfStay},
			{tkLPar, cmBack, coCallFunc, nil, cfStay},
			{tkLSBracket, cmBack, coIndex, nil, cfStay},
		},
		cmExpOper: {
			{tkDefault, ErrOper, coError, nil, 0},
			{tkStrExp, cmBack, coOperator, nil, 0},
			{[]int{tkAdd, tkDiv, tkMod, tkMul, tkSub, tkEqual, tkNotEqual, tkGreater, tkGreaterEqual,
				tkLess, tkLessEqual, tkAssign, tkOr, tkAnd, tkBitOr, tkBitAnd, tkBitXor, tkLShift,
				tkRShift, tkAddEq, tkSubEq, tkMulEq, tkDivEq, tkModEq, tkLShiftEq, tkRShiftEq,
				tkBitAndEq, tkBitOrEq, tkBitXorEq, tkRange}, cmBack,
				coOperator, nil, 0},
			{[]int{tkInc, tkDec}, 0, coUnaryPostOperator, nil, 0},
			{[]int{tkRPar, tkRSBracket}, 0, coOperator, nil, 0},
			{[]int{tkLSBracket}, cmBack, coIndex, nil, cfStay},
			{tkComma, cmBack, coComma, nil, 0},
			{[]int{tkLCurly, tkLine, tkRCurly, tkColon}, cmBack, nil, nil, cfStay},
		},
		cmElseIf: {
			{tkDefault, cmBack, coIfEnd, nil, cfStay},
			{tkLine, 0, nil, nil, 0},
			{tkElse, cmLCurly, coElse, nil, 0},
			{tkElif, cmExp, coElif, coElifBack, cfStopBack},
		},
		cmFunc: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkIdent, cmParams, coFuncName, nil, 0},
			{tkLine, 0, nil, nil, 0},
		},
		cmParams: {
			{tkDefault, ErrLCurly, coError, nil, 0},
			{tkIdent, cmLCurly, coRetType, nil, 0},
			{tkLPar, cmParam, nil, nil, cfStopBack},
			{tkLCurly, cmLCurly, coFuncStart, nil, cfStay},
			{tkLine, 0, nil, nil, 0},
		},
		cmParam: {
			{tkDefault, ErrType, coError, nil, 0},
			{tkIdent, cmWantVar, coType, nil, cfStopBack},
			{tkRPar, cmBack, nil, nil, 0},
			{tkLine, 0, nil, nil, 0},
		},
		cmWantVar: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkIdent, cmVar, coVar, nil, 0},
			{tkLine, 0, nil, nil, 0},
		},
		cmVar: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkIdent, 0, coVar, nil, 0},
			{tkComma, cmWantType, nil, nil, 0},
			{tkRPar, cmBack, nil, nil, cfStay},
			{tkLine, 0, nil, nil, 0},
		},
		cmWantType: {
			{tkDefault, ErrType, coError, nil, 0},
			{tkIdent, cmBack, nil, nil, cfStay},
			{tkLine, 0, nil, nil, 0},
		},
		cmMustVarType: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkIdent, cmVarType, coVarExp, nil, 0},
		},
		cmVarType: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkIdent, 0, coVar, nil, 0},
			{tkLine, cmBack, nil, nil, 0},
		},
		cmConst: {
			{tkDefault, cmExp, coConstEnum, coConstEnumBack, cfStay},
			{tkLine, 0, nil, nil, 0},
			{tkLCurly, cmConstDef, nil, nil, 0},
		},
		cmConstDef: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkIdent, cmConstName, coConst, nil, cfStopBack},
			{tkLine, 0, nil, nil, 0},
			{tkRCurly, cmBack, nil, nil, 0},
		},
		cmConstName: {
			{tkDefault, ErrMustAssign, coError, nil, 0},
			{tkAssign, cmExp, coConstExp, coConstExpBack, 0},
		},
		cmConstListStart: {
			{tkDefault, ErrLCurly, coError, nil, 0},
			{tkLine, 0, nil, nil, 0},
			{tkLCurly, cmConstList, nil, nil, 0},
		},
		cmConstList: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkIdent, 0, coConstList, nil, 0},
			{tkLine, 0, nil, nil, 0},
			{tkRCurly, cmBack, nil, nil, 0},
		},
		cmInit: {
			{tkDefault, cmExp, coExpStart, nil, cfStopBack | cfStay},
			{tkLCurly, cmInit, coInitStart, nil, cfStopBack},
			{tkRCurly, cmBack, coInitEnd, nil, 0},
			{tkComma, 0, coInitNext, nil, 0},
			{tkColon, 0, coInitKey, nil, 0},
		},
		cmStruct: {
			{tkDefault, ErrName, coError, nil, 0},
			{tkLine, 0, nil, nil, 0},
			{tkIdent, cmStructDef, coStruct, coStructEnd, 0},
		},
		cmStructDef: {
			{tkDefault, ErrLCurly, coError, nil, 0},
			{tkLCurly, cmStructFields, nil, nil, 0},
			{tkLine, 0, nil, nil, 0},
		},
		cmStructFields: {
			{tkDefault, ErrType, coError, nil, 0},
			{tkLine, 0, coStructLine, nil, 0},
			{tkIdent, cmStructName, coStructType, nil, cfStopBack},
			{tkRCurly, cmBack, coStructLine, nil, 0},
		},
		cmStructName: {
			{tkDefault, ErrLCurly, coError, nil, 0},
			{tkIdent, 0, coStructName, nil, 0},
			{tkRCurly, cmBack, nil, nil, cfStay},
			{tkLine, cmBack, nil, nil, 0},
		},
	}
	compileTable [][tkDefault]*cmState
)

func makeCompileTable() {
	compileTable = make([][tkDefault]*cmState, len(preCompile)+1)

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
