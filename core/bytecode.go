// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

// Bytecode contains bytecode information
type Bytecode struct {
	Code []uint16
	Used map[uint16]byte // identifier of used objects
}

type Exec struct {
	Code  []uint16
	Funcs map[uint16]uint32
}

const (
	STACKNO = iota
	STACKINT
	STACKBOOL
	STACKCHAR
	STACKFLOAT
	STACKSTR
	STACKANY
)

const (
	NOP      = iota
	PUSH16   // + int16
	PUSH32   // + int32
	PUSH64   // + int64
	ADD      // int + int
	SUB      // int - int
	MUL      // int * int
	DIV      // int / int
	MOD      // int % int
	SIGN     // -int
	EQINT    // int == int
	LTINT    // int < int
	GTINT    // int > int
	NOT      // logical not 1 => 0, 0 => 1
	JMP      // + int16 jump with clearing stack
	JZE      // + int16 jump if the top value is zero
	RET      // return + int16 stack id
	END      // end of the function
	CALLBYID // + uint16 id of the object
)
