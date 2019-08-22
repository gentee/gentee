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

type AssignIntFunc func(*int64, int64) (int64, error)
type AssignStrFunc func(*string, interface{}) (string, error)
type AssignAnyFunc func(interface{}, interface{}) (interface{}, error)

//type SetIndexFunc func(interface{}, interface{}, interface{}) error

const (
	TYPENONE  = 0
	TYPEINT   = 0x101
	TYPEBOOL  = 0x201
	TYPECHAR  = 0x301
	TYPESTR   = 0x402
	TYPEFLOAT = 0x503
	TYPEARR   = 0x604
	TYPERANGE = 0x704
	TYPEMAP   = 0x804
	TYPEPTR   = 0xf04
)

const (
	STACKNONE = iota
	STACKINT
	STACKSTR
	STACKFLOAT
	STACKANY
)

const (
	NOP       = iota
	PUSH32    // + int32
	PUSH64    // + int64
	PUSHSTR   // & (strid << 16 )
	ADD       // int + int
	SUB       // int - int
	MUL       // int * int
	DIV       // int / int
	MOD       // int % int
	BITOR     // int | int
	BITXOR    // int ^ int
	BITAND    // int & int
	LSHIFT    // int << int
	RSHIFT    // int >> int
	BITNOT    // ^int
	SIGN      // -int
	EQ        // int == int
	LT        // int < int
	GT        // int > int
	NOT       // logical not 1 => 0, 0 => 1
	ADDSTR    // str + str
	EQSTR     // str == str
	LTSTR     // str < str
	GTSTR     // str > str
	GETVAR    // & (block shift<<16) + int16 type + int16 index
	SETVAR    // & (block shift<<16) + int16 type + int16 index + int16 index count + int16 assign
	DUP       // & (type<<16) duplicate top int
	POP       // & (type<<16) pop top
	CYCLE     // cycle counter
	JMP       // + int32 jump with clearing stack
	JZE       // + int32 jump if the top value is zero
	JNZ       // + int32 jump if the top value is not zero
	INITVARS  // initializing variables
	DELVARS   // delete variables
	RANGE     // create range
	LEN       // & (type<<16) length
	FORINC    // & (index<<16) increment counter
	RET       // & (type<<16) return from function
	END       // end of the function
	INDEX     // & (int32 count) + {(type input<<16) + result type}
	CONSTBYID // + int32 id of the object
	CALLBYID  // & (par count<<16) + int32 id of the object
	EMBED     // & (embed id << 16) calls embedded func + int32 count for variadic funcs
	// + [variadic types]
	IOTA // & (iota<<16)

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
	INCDEC

	INC // &( 1 << 16 int++ ) ++int
	DEC // &( 1 << 16 int-- ) --int
)
