// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

const (
	// List of states
	stMain = iota
	stIdent
	stInt
	stHexOct
	stHex
	stOct
	stDiv
	stCommentLine
	stComment
	stCommentEnd
	stEqual
	stNot
	stGreater
	stLess
	stOr
	stAnd
	stAdd
	stSub
	stMul
	stRShift
	stLShift
	stXor
	stMod

	stError // it must be the last state

	// Flags for lexical parser
	fStart = 0x10000 // the beginning of the tken
	fToken = 0x20000 // tken has been parsed
	fNext  = 0x40000 // stay on the state and get the next character
	fStay  = 0x80000 // stay on the current character

	alphabet = 128
)

/* Alphabet for preTable
as is: _ 0 + - * / ( ) { } = ; , | & < > ! ? ^ % ~

7 0-7
9 0-9
a a-fA-F
n \n
r \r
s space
t \t
x xX
z a-zA-Z and unicode letter
*/

type preState struct {
	Key    string
	Action int
}

var (
	keywords = map[string]int{
		`elif`:   tkElif,
		`else`:   tkElse,
		`false`:  tkFalse,
		`func`:   tkFunc,
		`if`:     tkIf,
		`while`:  tkWhile,
		`return`: tkReturn,
		`run`:    tkRun,
		`true`:   tkTrue,
		`const`:  tkConst,
	}
	preTable = map[int][]preState{
		stMain: {
			{`z`, fStart | stIdent},
			{`+`, fStart | stAdd},
			{`-`, fStart | stSub},
			{`*`, fStart | stMul},
			{`/`, fStart | stDiv},
			{`(`, fToken | tkLPar},
			{`)`, fToken | tkRPar},
			{`{`, fToken | tkLCurly},
			{`}`, fToken | tkRCurly},
			{`=`, fStart | stEqual},
			{`&`, fStart | stAnd},
			{`|`, fStart | stOr},
			{`^`, fStart | stXor},
			{`%`, fStart | stMod},
			{`>`, fStart | stGreater},
			{`<`, fStart | stLess},
			{`~`, fToken | tkBitNot},
			{`!`, fStart | stNot},
			{`?`, fToken | tkQuestion},
			{`,`, fToken | tkComma},
			{`srt`, fNext},
			{`9`, fStart | stInt},
			{`0`, fStart | stHexOct},
			{`n;`, fToken | tkLine},
		},
		stIdent: {
			{``, fToken | tkIdent},
			{`z9_`, fNext},
		},
		stInt: {
			{``, fToken | tkInt},
			{`9`, fNext},
		},
		stHexOct: {
			{``, fToken | tkInt},
			{`x`, stHex},
			{`7`, stOct},
		},
		stHex: {
			{``, fToken | tkIntHex},
			{`z_`, stError},
			{`9a`, fNext},
		},
		stOct: {
			{``, fToken | tkIntOct},
			{`z9_`, stError},
			{`7`, fNext},
		},
		stDiv: {
			{``, fToken | tkDiv},
			{`=`, fNext | fToken | tkDivEq},
			{`/`, stCommentLine},
			{`*`, stComment},
		},
		stCommentLine: {
			{``, fNext},
			{`n`, fStay | stMain},
		},
		stComment: {
			{``, fNext},
			{`*`, stCommentEnd},
		},
		stCommentEnd: {
			{``, fStay | stComment},
			{`/`, stMain},
		},
		stEqual: {
			{``, fToken | tkAssign},
			{`=`, fNext | fToken | tkEqual},
		},
		stNot: {
			{``, fToken | tkNot},
			{`=`, fNext | fToken | tkNotEqual},
		},
		stGreater: {
			{``, fToken | tkGreater},
			{`>`, stRShift},
			{`=`, fNext | fToken | tkGreaterEqual},
		},
		stLess: {
			{``, fToken | tkLess},
			{`<`, stLShift},
			{`=`, fNext | fToken | tkLessEqual},
		},
		stOr: {
			{``, fToken | tkBitOr},
			{`=`, fNext | fToken | tkBitOrEq},
			{`|`, fNext | fToken | tkOr},
		},
		stAnd: {
			{``, fToken | tkBitAnd},
			{`=`, fNext | fToken | tkBitAndEq},
			{`&`, fNext | fToken | tkAnd},
		},
		stAdd: {
			{``, fToken | tkAdd},
			{`=`, fNext | fToken | tkAddEq},
			{`+`, fNext | fToken | tkInc},
		},
		stSub: {
			{``, fToken | tkSub},
			{`=`, fNext | fToken | tkSubEq},
			{`-`, fNext | fToken | tkDec},
		},
		stMul: {
			{``, fToken | tkMul},
			{`=`, fNext | fToken | tkMulEq},
		},
		stRShift: {
			{``, fToken | tkRShift},
			{`=`, fNext | fToken | tkRShiftEq},
		},
		stLShift: {
			{``, fToken | tkLShift},
			{`=`, fNext | fToken | tkLShiftEq},
		},
		stXor: {
			{``, fToken | tkBitXor},
			{`=`, fNext | fToken | tkBitXorEq},
		},
		stMod: {
			{``, fToken | tkMod},
			{`=`, fNext | fToken | tkModEq},
		},
	}

	parseTable [][alphabet]int
)

func makeParseTable() {

	fromto := func(state, jump int, from, to rune) {
		for i := from; i <= to; i++ {
			parseTable[state][i] = jump
		}
	}
	parseTable = make([][alphabet]int, stError)
	for state, items := range preTable {
		var (
			def int
		)
		if items[0].Key == `` {
			def = items[0].Action
		} else {
			def = stError
		}
		for i := 0; i < alphabet; i++ {
			parseTable[state][i] = def
		}
		for _, item := range items {
			jump := item.Action
			for _, ch := range item.Key {
				switch ch {
				case '7':
					fromto(state, jump, '0', '7')
				case '9':
					fromto(state, jump, '0', '9')
				case 'a':
					fromto(state, jump, 'a', 'f')
					fromto(state, jump, 'A', 'F')
				case 'n':
					parseTable[state][0xa] = jump
				case 'r':
					parseTable[state][0xd] = jump
				case 's':
					parseTable[state][' '] = jump
				case 't':
					parseTable[state][0x9] = jump
				case 'x':
					parseTable[state]['x'] = jump
					parseTable[state]['X'] = jump
				case 'z':
					fromto(state, jump, 'a', 'z')
					fromto(state, jump, 'A', 'Z')
					parseTable[state][alphabet-1] = jump
				default:
					parseTable[state][ch] = jump
				}
			}
		}
	}
}
