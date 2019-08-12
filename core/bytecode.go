// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

type Bcode int32

// Bytecode contains bytecode information
type Bytecode struct {
	Code    []Bcode
	Used    map[int32]byte // identifier of used objects
	Strings map[string]uint16
	Pos     []CodePos
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
	Strings []string // string resources
	Pos     []CodePos
}

const (
	TYPENONE = iota
	TYPEINT
	TYPEBOOL
	TYPECHAR

/*	STACKFLOAT
	STACKSTR
	STACKANY*/
)

const (
	NOP      = iota
	PUSH32   // + int32
	PUSH64   // + int64
	ADD      // int + int
	SUB      // int - int
	MUL      // int * int
	DIV      // int / int
	MOD      // int % int
	SIGN     // -int
	EQ       // int == int
	LT       // int < int
	GT       // int > int
	NOT      // logical not 1 => 0, 0 => 1
	GETVAR   // & (block shift<<16) + int16 type + int16 index
	SETVAR   // & (block shift<<16) + int16 type + int16 index
	DUP      // duplicate top int
	CYCLE    // cycle counter
	JMP      // + int32 jump with clearing stack
	JZE      // + int32 jump if the top value is zero
	JNZ      // + int32 jump if the top value is not zero
	INITVARS // initializing variables
	DELVARS  // delete variables
	RET      // & (type<<16) return from function
	END      // end of the function
	CALLBYID // & (par count<<16) + int32 id of the object
)
