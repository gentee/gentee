// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"strings"
	"unicode"

	"github.com/gentee/gentee/core"
)

type lexItem struct {
	Pattern interface{} // rune or string or []string or []rune
	Action  int
	Func    lexFunc
}

type lexFunc func(*lexEngine, int, int)

const (
	lexMain = iota + 1
	lexIdent
	lexOctHex
	lexInt
	lexOct
	lexHex
	lexFloat
	lexFloatExp
	lexFloatSignExp
	lexChar
	lexStrQuote
	lexStrDouble
	lexCommentLine
	lexComment
	lexCmd
	lexCmdLine
	lexMustIdent

	lexBack
	lexBackNext
	lexError

	alphabet = 128
)

const (
	fBack   = 0x010000 << iota // go to back further when callback
	fNewBuf                    // Start a new buffer
	fSkip                      // Skip next rune
)

const (
	fullAlphabet = alphabet + 4
	forP         = alphabet
	forL         = alphabet + 1
	forS         = alphabet + 2
	forD         = alphabet + 3 // default action

	isP = 1 << iota // IsPrint
	isL             // IsLetter
	isD             // IsDigit
	isS             // IsSpace
	isO             // IsOctetDigit
	isH             // IsHexDigit
)

var (
	keywords = map[string]int{
		`break`:    tkBreak,
		`continue`: tkContinue,
		`elif`:     tkElif,
		`else`:     tkElse,
		`false`:    tkFalse,
		`for`:      tkFor,
		`func`:     tkFunc,
		`if`:       tkIf,
		`in`:       tkIn,
		`while`:    tkWhile,
		`return`:   tkReturn,
		`run`:      tkRun,
		`true`:     tkTrue,
		`const`:    tkConst,
		`struct`:   tkStruct,
		`switch`:   tkSwitch,
		`case`:     tkCase,
		`include`:  tkInclude,
		`import`:   tkImport,
		`pub`:      tkPub,
		`fn`:       tkFn,
		`go`:       tkGo,
		`default`:  tkDefault,
	}

	charType [alphabet]int

	preBack     = lexItem{nil, lexBack, nil}
	preFloat    = lexItem{nil, 0, newFloat}
	preError    = lexItem{nil, lexError | ErrLetter, nil}
	preLexTable = map[int][]lexItem{
		lexMain: { // main
			{nil, lexError | ErrLetter, nil},
			{'S', 0, nil},
			{[]rune{'\n', ';'}, 0, newLine},
			{[]rune{'{', '(', ')', '[', ']', ',', '?', '~'}, 0, newSymbol},
			{[]string{`+=`, `++`, `+`}, 0, newOper},
			{[]string{`-=`, `--`, `-`}, 0, newOper},
			{[]string{`*=`, `*`}, 0, newOper},
			{[]string{`...`, `..`, `.`}, 0, newOper},
			{[]string{`==`, `=`}, 0, newOper},
			{[]string{`!=`, `!`}, 0, newOper},
			{[]string{`<<=`, `<=`, `<<`, `<`}, 0, newOper},
			{[]string{`>>=`, `>=`, `>>`, `>`}, 0, newOper},
			{[]string{`||`, `|=`, `|`}, 0, newOper},
			{[]string{`&&`, `&=`, `&`}, 0, newOper},
			{[]string{`%=`, `%`}, 0, newOper},
			{[]string{`^=`, `^`}, 0, newOper},
			{[]string{`#=`, `##`, `#`}, 0, newOper},
			{':', 0, newOper},
			{'$', lexCmd, nil},
			{'}', 0, endExp},
			{[]string{`//`, `/*`, `/=`, `/`}, 0, newDiv},
			{'\'', lexChar, newChar},
			{'`', lexStrQuote | fNewBuf, nil},
			{'"', lexStrDouble | fNewBuf, nil},
			{'L', lexIdent, newIdent},
			{'D', lexInt, newInt},
			{'0', lexOctHex, newInt},
		},
		lexIdent: { // identifier
			{[]rune{'L', 'D', '_', '.'}, 0, nil},
		},
		lexOctHex: { // number
			{[]rune{'x', 'X'}, lexHex, nil},
			{'O', lexOct, nil},
			{[]string{`..`, `.`}, fBack, isFloat},
		},
		lexInt: { // integer
			{[]rune{'L', '_'}, lexError | ErrWord, nil},
			{'D', 0, nil},
			{[]string{`..`, `.`}, fBack, isFloat},
			{[]rune{'e', 'E'}, fBack, isFloatExp},
		},
		lexOct: { // octal integer
			{[]rune{'L', 'D', '_'}, lexError | ErrWord, nil},
			{'O', 0, nil},
		},
		lexHex: { // hex integer
			{[]rune{'L', '_'}, lexError | ErrWord, nil},
			{'H', 0, nil},
		},
		lexFloat: { // float
			{[]rune{'L', '_'}, lexError | ErrWord, nil},
			{[]rune{'e', 'E'}, lexFloatExp, nil},
			{'D', 0, nil},
		},
		lexFloatExp: { // float exp
			{[]rune{'L', '_'}, lexError | ErrWord, nil},
			{'D', 0, nil},
			{[]rune{'+', '-'}, lexFloatSignExp, nil},
		},
		lexFloatSignExp: { // float sign exp
			{[]rune{'L', '_'}, lexError | ErrWord, nil},
			{'D', 0, nil},
		},
		lexChar: { // char
			{nil, 0, nil},
			{'\\', fSkip, nil},
			{'\'', lexBackNext, nil},
		},
		lexStrQuote: {
			{nil, 0, pushBuf},
			{[]string{"``", "`"}, 0, endStr},
			{`%{`, 0, expStr},
			{`${`, lexMustIdent, newExpEnv},
		},
		lexStrDouble: {
			{nil, 0, pushBuf},
			{'\\', 0, backSlash},
			{'"', 0, newStr},
		},
		lexCommentLine: {
			{nil, 0, nil},
			{'\n', lexBack, nil},
		},
		lexComment: {
			{nil, 0, nil},
			{`*/`, lexBackNext, nil},
		},
		lexCmd: {
			{nil, lexError | ErrWord, nil},
			{'L', lexIdent | fBack, newEnv},
			{' ', lexCmdLine | fBack, newCmdLine},
		},
		lexCmdLine: {
			{nil, 0, pushBuf},
			{`%{`, 0, expStr},
			{`${`, lexMustIdent, newExpEnv},
			{'\n', lexBack, nil},
		},
		lexMustIdent: {
			{nil, lexError | ErrEnvName, nil},
			{'L', lexIdent | fBack, nil}, //newExpEnv},
		},
	}
	lexTable = make(map[int][fullAlphabet]*lexItem)
)

func fillCharType() {
	for i := 0; i < alphabet; i++ {
		var val int
		r := rune(i)
		if unicode.IsPrint(r) {
			val = isP
		}
		if unicode.IsLetter(r) {
			val |= isL
			if (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
				val |= isH
			}
		}
		if unicode.IsDigit(r) {
			val |= isD | isH
			if r <= '7' {
				val |= isO
			}
		}
		if unicode.IsSpace(r) {
			val |= isS
		}
		charType[i] = val
	}
}

func setJump(jumps *[fullAlphabet]*lexItem, in rune, ptr *lexItem) {
	for i := 0; i < alphabet; i++ {
		switch in {
		case 'L':
			if charType[i]&isL != 0 {
				(*jumps)[i] = ptr
			}
		case 'S':
			if charType[i]&isS != 0 {
				(*jumps)[i] = ptr
			}
		case 'D':
			if charType[i]&isD != 0 {
				(*jumps)[i] = ptr
			}
		case 'O':
			if charType[i]&isO != 0 {
				(*jumps)[i] = ptr
			}
		case 'H':
			if charType[i]&isH != 0 {
				(*jumps)[i] = ptr
			}
		default:
			if in == rune(i) {
				(*jumps)[in] = ptr
			}
		}
	}
	// for runes which are greater then 127
	switch in {
	case 'P':
		(*jumps)[forP] = ptr
	case 'L':
		(*jumps)[forL] = ptr
	case 'S':
		(*jumps)[forS] = ptr
	}
}

func makeLexTable() {
	fillCharType()
	for state, items := range preLexTable {
		var (
			def   *lexItem
			jumps [fullAlphabet]*lexItem
		)
		if items[0].Pattern == nil {
			def = &items[0]
		} else {
			def = &preBack
		}
		for i := 0; i < fullAlphabet; i++ {
			if i < alphabet && charType[i] == 0 {
				jumps[i] = &preError
			} else {
				jumps[i] = def
			}
		}
		for i, item := range items {
			switch v := item.Pattern.(type) {
			case rune:
				setJump(&jumps, v, &preLexTable[state][i])
			case string:
				setJump(&jumps, rune(v[0]), &preLexTable[state][i])
			case []string:
				setJump(&jumps, rune(v[0][0]), &preLexTable[state][i])
			case []rune:
				for _, r := range v {
					setJump(&jumps, r, &preLexTable[state][i])
				}
			default:
				continue
			}
		}
		lexTable[state] = jumps
	}
}

func newIdent(lex *lexEngine, start, off int) {
	if lex.Callback {
		original := string(lex.Lex.Source[start:off])
		name := strings.TrimRight(original, `.`)
		lex.Off -= len(original) - len(name)
		tokType := tkIdent
		if keyType, ok := keywords[name]; ok {
			tokType = keyType
		}
		lex.Lex.NewToken(tokType, start, lex.Off-start)
	}
}

func newLine(lex *lexEngine, start, off int) {
	if lex.Colon && lex.Lex.Source[off] != ';' {
		lex.Colon = false
		lex.Lex.NewTokens(off, tkRCurly)
	}
	lex.Lex.NewTokens(off, tkLine)
}

func newSymbol(lex *lexEngine, start, off int) {
	lex.Lex.NewTokens(off, oper2tk[string(lex.Lex.Source[off:off+1])])
}

func newOper(lex *lexEngine, start, off int) {
	if lex.Callback {
		return
	}
	oper := string(lex.Lex.Source[start : off+1])
	lex.Lex.NewToken(oper2tk[oper], start, len(oper))
}

func isFloat(lex *lexEngine, start, off int) {
	if lex.Callback {
		return
	}
	oper := string(lex.Lex.Source[start : off+1])
	switch oper2tk[oper] {
	case tkRange:
		lex.State = lexBack
		lex.Off--
	case tkDot:
		lex.State = lexFloat
		lex.Stack[len(lex.Stack)-1].Action = &preFloat
	}
}

func isFloatExp(lex *lexEngine, start, off int) {
	if lex.Callback {
		return
	}
	lex.State = lexFloatExp
	lex.Stack[len(lex.Stack)-1].Action = &preFloat
}

func newFloat(lex *lexEngine, start, off int) {
	if lex.Callback {
		lex.Lex.NewToken(tkFloat, start, off-start)
	}
}

func newInt(lex *lexEngine, start, off int) {
	if lex.Callback {
		lex.Lex.NewToken(tkInt, start, off-start)
	}
}

func newStr(lex *lexEngine, start, off int) {
	start = lex.Stack[len(lex.Stack)-1].Offset
	var index int
	if len(lex.Buf) > 0 {
		index = len(lex.Lex.Strings)
		lex.Lex.Strings = append(lex.Lex.Strings, string(lex.Buf))
	}
	lex.Lex.Tokens = append(lex.Lex.Tokens, core.Token{Type: int32(tkStr), Index: int32(index),
		Offset: start, Length: off - start})
	lex.Buf = lex.Buf[:0]
	if !lex.Callback {
		lex.State = lexBackNext
	}
}

func endStr(lex *lexEngine, start, off int) {
	if off-start == 1 {
		lex.Buf = append(lex.Buf, lex.Lex.Source[off])
	} else {
		newStr(lex, start, off)
	}
}

func expStr(lex *lexEngine, start, off int) {
	if lex.Callback {
		return
	}
	prev := lex.Stack[len(lex.Stack)-1]
	newStr(lex, prev.Offset, off)
	lex.Lex.NewTokens(off, tkStrExp, tkLPar)
	lex.State = lexMain
}

func endExp(lex *lexEngine, start, off int) {
	if len(lex.Stack) != 0 {
		prev := lex.Stack[len(lex.Stack)-1]
		if prev.State == lexStrQuote || prev.State == lexStrDouble || prev.State == lexCmdLine {
			lex.Lex.NewTokens(off, tkRPar, tkStrExp)
			lex.State = lexBackNext
			return
		}
	}
	newSymbol(lex, start, off)
}

func backSlash(lex *lexEngine, start, off int) {
	if lex.Callback {
		return
	}
	if lex.Lex.Source[off+1] == '{' {
		expStr(lex, start, off+1)
		lex.State |= fSkip
	} else {
		lex.Buf = append(lex.Buf, lex.Lex.Source[off], lex.Lex.Source[off+1])
		lex.State = fSkip
	}
}

func pushBuf(lex *lexEngine, start, off int) {
	lex.Buf = append(lex.Buf, lex.Lex.Source[off])
}

func newCmdLine(lex *lexEngine, start, off int) {
	if lex.Callback {
		newStr(lex, start, off)
		lex.Lex.NewTokens(start, tkRPar)
		return
	}
	lex.Lex.NewTokens(start-1, tkIdent, tkLPar)
	lex.Buf = lex.Buf[:0]
}

func newDiv(lex *lexEngine, start, off int) {
	oper := string(lex.Lex.Source[start : off+1])
	token := oper2tk[oper]
	switch token {
	case tkDiv, tkDivEq:
		lex.Lex.NewToken(token, start, len(oper))
	case tkCommentLine:
		lex.State = lexCommentLine
	case tkComment:
		lex.State = lexComment
	}
}

func newEnv(lex *lexEngine, start, off int) {
	if lex.Callback {
		lex.Lex.NewToken(tkEnv, start, off-start)
	}
}

func newExpEnv(lex *lexEngine, start, off int) {
	if lex.Callback {
		prev := lex.Stack[len(lex.Stack)-1]
		newStr(lex, prev.Offset, start)
		lex.Lex.NewTokens(off, tkStrExp)
		if lex.Lex.Source[off] != '}' {
			lex.Error = ErrEnvName
		}
		lex.Lex.NewToken(tkEnv, start+1, off-start-1)
		lex.Lex.NewTokens(off, tkStrExp)
		lex.State = lexBackNext
	}
}

func newChar(lex *lexEngine, start, off int) {
	if lex.Callback {
		lex.Lex.NewToken(tkChar, start, off-start+1)
	}
}
