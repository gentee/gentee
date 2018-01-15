// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

const (
	// TokIdent means identifier
	TokIdent = iota + 1
	// TokLine means a new line
	TokLine
	// TokError can be only the last token
	TokError
)

const (
	// List of states
	stMain = iota
	stIdent
	stError // it must be the last state

	// Flags for lexical parser
	fStart = 0x10000 // the beginning of the token
	fToken = 0x20000 // token has been parsed
	fNext  = 0x40000 // stay on the state and get the next character

	alphabet = 128
)

/* Alphabet for preTable
as is: _ 0 . ;

9 0-9
a a-fA-F
n \n
r \r
s space
t \t
x xX
z a-zA-Z and unicode letter
*/

var (
	preTable = map[int]map[string]int{
		stMain: {
			`z`:   fStart | stIdent,
			`srt`: fNext,
			`n;`:  fToken | TokLine,
		},
		stIdent: {
			`z9_`: fNext,
			``:    fToken | TokIdent,
		},
	}

	parseTable [][alphabet]int
)

func init() {
	makeParseTable()
}

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
			ok  bool
		)
		if def, ok = items[``]; !ok {
			def = stError
		}
		for i := 0; i < alphabet; i++ {
			parseTable[state][i] = def
		}
		for key, jump := range items {
			for _, ch := range key {
				switch ch {
				case '9':
					fromto(state, jump, '0', '9')
				case 'n':
					parseTable[state][0xa] = jump
				case 'r':
					parseTable[state][0xd] = jump
				case 's':
					parseTable[state][' '] = jump
				case 't':
					parseTable[state][0x9] = jump
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
