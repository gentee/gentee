// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"unicode"

	"bitbucket.org/novostrim/go-gentee/core"
)

// LexParsing performs lexical analysis of the input string and returns a sequence of lexical tokens.
func LexParsing(input []rune) (*core.Lex, int) {
	var (
		off, state, tokOff, line int
	)
	lp := core.Lex{Source: append(input, ' '), // added stop-character
		Lines: make([]int, 0, 10), Strings: []string{``}}
	buf := make([]rune, 0, 4096)

	newToken := func(tokType int) {
		var index int
		if tokType&fBuf != 0 {
			if len(buf) > 0 {
				index = len(lp.Strings)
				lp.Strings = append(lp.Strings, string(buf))
			}
		}
		tokType &= 0xffff
		if tokType == stIdent { // check keywords
			if keyType, ok := keywords[string(input[tokOff:off])]; ok {
				tokType = keyType
			}
		}
		length := off - tokOff
		if length == 0 { // one-byte token
			length = 1
		}
		lp.Tokens = append(lp.Tokens, core.Token{Type: int32(tokType), Index: int32(index),
			Offset: tokOff, Length: length})
	}
	newLine := func(offset int) {
		if len(lp.Lines) == 0 || lp.Lines[len(lp.Lines)-1] != offset {
			lp.Lines = append(lp.Lines, offset)
			line++
		}
	}

	newLine(0)
	length := len(lp.Source)
	lp.Tokens = make([]core.Token, 0, 32+length/10)

	for off < length {
		ch := lp.Source[off]
		if ch >= 127 {
			if unicode.IsLetter(ch) {
				ch = 127
			} else {
				tokOff = off
				newToken(tkError)
				return &lp, ErrLetter
			}
		}
		todo := parseTable[state][ch]
		if lp.Source[off] == 0xa {
			newLine(off + 1)
		}
		if todo&fStart != 0 {
			tokOff = off
		}
		if todo&fStartBuf != 0 {
			buf = buf[:0]
		}
		if todo&fPushBuf != 0 {
			buf = append(buf, lp.Source[off])
		}
		if todo&fToken != 0 {
			if state == stMain { // it means one character token
				tokOff = off
			} else if todo&fNext != 0 {
				off++
			}
			newToken(todo)
			if state != stMain {
				state = stMain
				continue
			}
		} else if todo&fNext == 0 {
			if state = todo & 0xffff; state == stError {
				tokOff = off
				newToken(tkError)
				return &lp, ErrWord
			}
			if todo&fStay != 0 {
				continue
			}
		}
		off++
	}

	return &lp, ErrSuccess
}

func getToken(lp *core.Lex, cur int) string {
	// !!! TODO Added checking out of range
	token := lp.Tokens[cur]
	return string(lp.Source[token.Offset : token.Offset+token.Length])
}
