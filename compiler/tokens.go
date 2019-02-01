// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

const (
	tkIdent       = iota + 1 // identifier
	tkLine                   // a new line
	tkInt                    // integer number
	tkFloat                  // float number
	tkType                   // type name
	tkChar                   // character
	tkStr                    // string
	tkEnv                    // os environment
	tkCommentLine            // // comment
	tkComment                // /* comment */
	tkError                  // tkError can be only the last tken
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
	tkLSBracket                // [
	tkRSBracket                // ]
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
	tkCmdLine                  // $
	tkColon                    // :
	tkDot                      // .
	tkRange                    // ..
	tkVariadic                 // ...
)

// Keywords
const (
	tkRun = iota + 64 // run
	tkReturn
	tkFalse
	tkFor
	tkFunc
	tkTrue
	tkIf
	tkIn
	tkElif
	tkElse
	tkWhile
	tkConst
	tkStruct
	tkBreak
	tkContinue
	tkSwitch
	tkCase
	tkInclude
	tkImport
	tkPub
	tkDefault
	tkToken // is used for preCompileTable
)

// Flags
const (
	tkUnary    = 0x10000 // flag for unary commands
	tkCallFunc = 0x20000 // temporary calling function
	tkPost     = 0x40000 // flag for post unary commands
	tkIndex    = 0x80000 // temporary index
)

var (
	oper2tk = map[string]int{
		`+`:   tkAdd,
		`-`:   tkSub,
		`*`:   tkMul,
		`/`:   tkDiv,
		`=`:   tkAssign,
		`,`:   tkComma,
		`(`:   tkLPar,
		`)`:   tkRPar,
		`[`:   tkLSBracket,
		`]`:   tkRSBracket,
		`{`:   tkLCurly,
		`}`:   tkRCurly,
		`==`:  tkEqual,
		`!`:   tkNot,
		`!=`:  tkNotEqual,
		`<`:   tkLess,
		`<=`:  tkLessEqual,
		`>`:   tkGreater,
		`>=`:  tkGreaterEqual,
		`&&`:  tkAnd,
		`||`:  tkOr,
		`?`:   tkQuestion,
		`&`:   tkBitAnd,
		`|`:   tkBitOr,
		`^`:   tkBitXor,
		`%`:   tkMod,
		`<<`:  tkLShift,
		`>>`:  tkRShift,
		`+=`:  tkAddEq,
		`-=`:  tkSubEq,
		`*=`:  tkMulEq,
		`/=`:  tkDivEq,
		`%=`:  tkModEq,
		`<<=`: tkLShiftEq,
		`>>=`: tkRShiftEq,
		`&=`:  tkBitAndEq,
		`|=`:  tkBitOrEq,
		`^=`:  tkBitXorEq,
		`++`:  tkInc,
		`--`:  tkDec,
		`$`:   tkCmdLine,
		`//`:  tkCommentLine,
		`/*`:  tkComment,
		`:`:   tkColon,
		`.`:   tkDot,
		`..`:  tkRange,
		`...`: tkVariadic,
	}
)
