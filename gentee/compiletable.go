// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

const (
	// List of compile states
	cmMain     = iota
	cmRun      // run command
	cmBody     // body of the code
	cmExp      // expression
	cmExpStart // begining of the expression
	cmExpValue // expression returns value

	// Flags
	cfSkip  = 0x10000 // stay on the current state
	cfBack  = 0x20000 // go to the previous state
	cfStay  = 0x40000 // stay on the current token
	cfError = 0x80000 // return error
)

type compFunc func() error

type cmState struct {
	Tokens interface{} // can be one token or []token
	Action int
	Func   compFunc
}

var (
	preCompile = map[int][]cmState{
		cmMain: {
			{tkDefault, cfError | ErrDecl, nil},
			{tkRun, cmRun, nil},
		},
		cmRun: {
			{tkLine, cfSkip, nil},
			{tkLCurly, cmBody, nil},
		},
		cmBody: {
			{tkIdent, cfStay | cmExp, nil},
			{tkLine, cfSkip, nil},
			{tkRCurly, cfBack, nil},
			{tkReturn, cmExpValue, nil},
		},
		cmExp: {
			{tkIdent, cmExpStart, nil},
		},
		cmExpStart: {
			{tkAssign, cmExpValue, nil},
		},
		cmExpValue: {
			{[]int{tkInt, tkIntHex, tkIntOct}, cmExpStart, nil},
			//			{tkAssign, cfStay | cmExpStart},
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
