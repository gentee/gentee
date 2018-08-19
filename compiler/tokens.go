// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

const (
	tkIdent  = iota + 1 // identifier
	tkLine              // a new line
	tkInt               // integer number (10-base)
	tkIntHex            // integer number (16-base)
	tkIntOct            // integer number (8-base)
	tkType              // type name
	tkStr               // string
	tkError             // tkError can be only the last tken
)

// Operators
const (
	tkAdd          = iota + 16 // +
	tkSub                      // -
	tkMul                      // *
	tkDiv                      // /
	tkAssign                   // =
	tkComma                    // ,
	tkLPar                     // (
	tkRPar                     // )
	tkLCurly                   // {
	tkRCurly                   // }
	tkEqual                    // ==
	tkNot                      // !
	tkNotEqual                 // !=
	tkLess                     // <
	tkLessEqual                // <=
	tkGreater                  // >
	tkGreaterEqual             // >=
	tkAnd                      // &&
	tkOr                       // ||
	tkQuestion                 // ?
	tkBitAnd                   // &
	tkBitOr                    // |
	tkBitXor                   // ^
	tkMod                      // %
	tkLShift                   // <<
	tkRShift                   // >>
	tkBitNot                   // ~
	tkAddEq                    // +=
	tkSubEq                    // -=
	tkMulEq                    // *=
	tkDivEq                    // /=
	tkModEq                    // %=
	tkLShiftEq                 // <<=
	tkRShiftEq                 // >>=
	tkBitAndEq                 // &=
	tkBitOrEq                  // |=
	tkBitXorEq                 // ^=
	tkInc                      // ++
	tkDec                      // --
	tkStrExp                   // expression inside the string
)

// Keywords
const (
	tkRun = iota + 64 // run
	tkReturn
	tkFalse
	tkFunc
	tkTrue
	tkIf
	tkElif
	tkElse
	tkWhile
	tkConst
	tkDefault // is used for preCompileTable
)

// Flags
const (
	tkUnary    = 0x10000 // flag for unary commands
	tkCallFunc = 0x20000 // temporary calling function
	tkPost     = 0x40000 // flag for post unary commands
)
