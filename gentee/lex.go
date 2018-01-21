// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	//	"fmt"
	"unicode"
)

// Token is a lexical token.
type Token struct {
	Type   int
	Offset int
	Length int
}

// Lex contains the result of the lexical parsing
type Lex struct {
	Source []rune
	Tokens []Token
	Lines  []int // offsets of lines
}

// LineColumn return the line and the column of the ind-th token
func (lp *Lex) LineColumn(ind int) (line int, column int) {
	if len(lp.Tokens) > ind {
		for ; line < len(lp.Lines); line++ {
			if lp.Lines[line] > lp.Tokens[ind].Offset {
				break
			}
		}
		column = lp.Tokens[ind].Offset - lp.Lines[line-1] + 1
	}
	return
}

// LexParsing performs lexical analysis of the input string and returns a sequence of lexical tokens.
func LexParsing(input []rune) (*Lex, error) {
	var (
		off, state, tokOff, line int
	)
	lp := Lex{Source: append(input, ' '), // added stop-character
		Lines: make([]int, 0, 10)}

	newToken := func(tokType int) {
		if tokType == stIdent { // check keywords
			if keyType, ok := keywords[string(input[tokOff:off])]; ok {
				tokType = keyType
			}
		}
		lp.Tokens = append(lp.Tokens, Token{Type: tokType, Offset: tokOff, Length: off - tokOff})
	}
	newLine := func(offset int) {
		lp.Lines = append(lp.Lines, offset)
		line++
	}

	newLine(0)
	length := len(lp.Source)
	lp.Tokens = make([]Token, 0, 32+length/10)

	for off < length {
		ch := lp.Source[off]
		if ch >= 127 {
			if unicode.IsLetter(ch) {
				ch = 127
			} else {
				tokOff = off
				newToken(tkError)
				return &lp, compError(ErrLetter)
			}
		}
		todo := parseTable[state][ch]
		//		fmt.Printf("off %d %x\r\n", off, todo)
		if lp.Source[off] == 0xa {
			newLine(off + 1)
		}
		if todo&fStart != 0 {
			tokOff = off
		}
		if todo&fToken != 0 {
			if state == stMain { // it means one character token
				tokOff = off
			}
			newToken(todo & 0xffff)
			if state != stMain {
				state = stMain
				continue
			}
		} else if todo&fNext == 0 {
			if state = todo & 0xffff; state == stError {
				tokOff = off
				newToken(tkError)
				return &lp, compError(ErrWord)
			}
		}
		off++
	}

	return &lp, nil
}
