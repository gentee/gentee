// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

type Bcode int32

// Bytecode contains bytecode information
type Bytecode struct {
	Code      []Bcode
	Used      map[int32]byte // identifier of used objects
	Init      []int32        // offsets of init funcs
	Strings   map[string]uint16
	StrOffset []int32 // offsets of PUSHSTR
	Pos       []CodePos
}

type CodePos struct {
	Offset int32  // byte code position
	Path   uint16 // Path index
	Name   uint16 // Name index
	Line   uint16 // Line
	Column uint16 // Column
}

type Exec struct {
	Code    []Bcode
	Funcs   map[int32]int32
	Init    []int32  // offsets of init funcs (initializing constants)
	Strings []string // string resources
	Pos     []CodePos
}

// Embed contains information about the golang function
type Embed struct {
	Func     interface{} // golang function
	Return   uint16      // the type of the result
	Params   []uint16    // the types of parameters
	Variadic bool        // variadic function
	Runtime  bool        // the first parameter is rt
	CanError bool        // can generate error
}

const (
	TYPENONE = iota
	TYPEINT
	TYPEBOOL
	TYPECHAR
	TYPESTR

/*	STACKFLOAT
	STACKSTR
	STACKANY*/
)

const (
	NOP     = iota
	PUSH32  // + int32
	PUSH64  // + int64
	PUSHSTR // & (strid << 16 )
	ADD     // int + int
	SUB     // int - int
	MUL     // int * int
	DIV     // int / int
	MOD     // int % int
	BITOR   // int | int
	BITXOR  // int ^ int
	BITAND  // int & int
	LSHIFT  // int << int
	RSHIFT  // int >> int
	BITNOT  // ^int
	SIGN    // -int
	EQ      // int == int
	LT      // int < int
	GT      // int > int
	NOT     // logical not 1 => 0, 0 => 1
	ADDSTR  // str + str
	EQSTR   // str == str
	LTSTR   // str < str
	GTSTR   // str > str
	LENSTR  // *str
	GETVAR  // & (block shift<<16) + int16 type + int16 index
	//	SETVAR    // & (block shift<<16) + int16 type + int16 index
	ADDRESS      // & (block shift<<16) + int16 type + int16 index
	ASSIGN       // & (int16 type << 16)
	ASSIGNADD    // int += int  & (int16 type << 16) str += str
	ASSIGNSUB    // int -= int
	ASSIGNMUL    // int *= int
	ASSIGNDIV    // int /= int
	ASSIGNMOD    // int %= int
	ASSIGNBITOR  // int |= int
	ASSIGNBITXOR // int ^= int
	ASSIGNBITAND // int &= int
	ASSIGNLSHIFT // int <<= int
	ASSIGNRSHIFT // int >>= int
	INC          // &( 1 << 16 int++ ) ++int
	DEC          // &( 1 << 16 int-- ) --int
	DUP          // duplicate top int
	CYCLE        // cycle counter
	JMP          // + int32 jump with clearing stack
	JZE          // + int32 jump if the top value is zero
	JNZ          // + int32 jump if the top value is not zero
	INITVARS     // initializing variables
	DELVARS      // delete variables
	RET          // & (type<<16) return from function
	END          // end of the function
	CONSTBYID    // + int32 id of the object
	CALLBYID     // & (par count<<16) + int32 id of the object
	EMBED        // & (embed id << 16) calls embedded func + int32 count for variadic funcs
	// + [variadic types]
	IOTA // & (iota<<16)
)
