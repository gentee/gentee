// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"strings"
	"unicode"

	"github.com/gentee/gentee/core"
)

type lexStack struct {
	State  int // lexical state
	Offset int // offset
	Action *lexItem
}

type lexEngine struct {
	Lex      *core.Lex
	Buf      []rune
	State    int
	Stack    []lexStack
	Callback bool
	Colon    bool
	Error    int
}

// LexParsing performs lexical analysis of the input string and returns a sequence of lexical tokens.
func LexParsing(input []rune) (*core.Lex, int) {
	var (
		off, flag, action int
		lex               lexEngine
	)
	lp := core.Lex{Source: append(input, []rune(`   `)...), // added stop-character
		Lines: make([]int, 0, 10), Strings: []string{``}}
	lex.Lex = &lp
	lex.Buf = make([]rune, 0, 4096)
	lex.Stack = make([]lexStack, 0, 16)

	newLine := func(offset int) {
		if len(lp.Lines) == 0 || lp.Lines[len(lp.Lines)-1] != offset {
			lp.Lines = append(lp.Lines, offset)
		}
	}

	newLine(0)
	length := len(lp.Source)
	lp.Tokens = make([]core.Token, 0, 32+length/10)
	// Skip the first lines with # character
	var hashMode bool
	for lp.Source[off] == '#' || hashMode {
		start := off
		for ; off < length && lp.Source[off] != 0xa; off++ {
		}
		if off >= length {
			break
		}
		off++
		line := string(lp.Source[start:off])
		if strings.TrimSpace(line) == `###` {
			hashMode = !hashMode
		} else if start != 0 || lp.Source[1] != '!' {
			if !hashMode {
				line = line[1:]
			}
			lp.Header += line
		}
		newLine(off)
	}
	state := lexMain
main:
	for off < length {
		ch := lp.Source[off]
		start := off
		if ch > 127 {
			if unicode.IsSpace(ch) {
				ch = forS
			} else if unicode.IsLetter(ch) {
				ch = forL
			} else if unicode.IsPrint(ch) {
				ch = forP
			}
		}
		if lp.Source[off] == 0xa {
			newLine(off + 1)
		}
		pLexItem := lexTable[state][ch]
		lex.State = 0
		lex.Callback = false
		flag = pLexItem.Action & 0xffffff00
		if flag&fNewBuf != 0 {
			lex.Buf = lex.Buf[:0]
		}
		action = pLexItem.Action & 0xff
		switch v := pLexItem.Pattern.(type) {
		case string:
			i := 1
			length := len(v)
			for ; i < length; i++ {
				if lp.Source[off+i] != rune(v[i]) {
					pLexDef := lexTable[state][forD]
					if pLexDef.Func != nil {
						pLexDef.Func(&lex, start, off)
					}
					off++
					continue main
				}
			}
			off += length - 1
		case []string:
		pattern:
			for _, pat := range v {
				i := 1
				length := len(pat)
				for ; i < length; i++ {
					if lp.Source[off+i] != rune(pat[i]) {
						continue pattern
					}
				}
				off += length - 1
				break
			}
		}
		if pLexItem.Func != nil {
			pLexItem.Func(&lex, start, off)
			if lex.State != 0 {
				action = lex.State & 0xff
				flag = lex.State & 0xffffff00
			}
		}
		switch action {
		case 0:
		case lexError:
			lex.Error = flag
		case lexBack, lexBackNext:
			lex.Callback = true
			for {
				prev := lex.Stack[len(lex.Stack)-1]
				state = prev.State
				lex.Stack = lex.Stack[:len(lex.Stack)-1]
				if prev.Action.Func != nil {
					prev.Action.Func(&lex, prev.Offset, off)
					if lex.State != 0 {
						action = lex.State
					}
				}
				if prev.Action.Action&fBack == 0 {
					break
				}
			}
			if action == lexBack {
				continue
			}
		default:
			lex.Stack = append(lex.Stack, lexStack{State: state, Offset: off, Action: pLexItem})
			state = action
		}
		if lex.Error != ErrSuccess {
			lp.NewTokens(off, tkError)
			return &lp, lex.Error
		}
		if flag&fSkip != 0 {
			off++
		}
		off++
	}
	if lex.Colon {
		lp.NewTokens(off, tkRCurly)
	}
	return &lp, ErrSuccess
}

func getToken(lp *core.Lex, cur int) string {
	// !!! TODO Added checking out of range
	token := lp.Tokens[cur]
	return string(lp.Source[token.Offset : token.Offset+token.Length])
}
