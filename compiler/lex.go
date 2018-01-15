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
	Line   int
	Column int
}

// Lex performs lexical analysis of the input string and returns a sequence of lexical tokens.
func Lex(input []rune) (*[]Token, error) {
	var (
		off, state, lineOff, tokOff int
	)
	line := 1
	input = append(input, ' ')
	length := len(input)
	tokens := make([]Token, 0, length/4)

	newToken := func(tokType int) {
		tokens = append(tokens, Token{Type: tokType, Line: line, Column: tokOff + 1})
	}

	for off < length {
		ch := input[off]
		if ch >= 127 {
			if unicode.IsLetter(ch) {
				ch = 127
			} else {
				newToken(TokError)
				return &tokens, ErrLexem
			}
		}
		todo := parseTable[state][ch]
		if input[off] == 0xa {
			line++
			lineOff = off + 1
		}
		if todo&fStart != 0 {
			tokOff = off - lineOff
		}
		//		fmt.Printf("%v %x %d %d\r\n", input[off], todo, state, off)
		if todo&fToken != 0 {
			if state == stMain { // it means one character token
				tokOff = off - lineOff
			}
			newToken(todo & 0xffff)
			state = stMain
		} else if todo&fNext == 0 {
			if state = todo & 0xffff; state == stError {
				newToken(TokError)
				return &tokens, ErrLexem
			}
		}
		off++
	}

	return &tokens, nil
}
