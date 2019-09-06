// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

type Bcode int32

// Bytecode contains bytecode information
type Bytecode struct {
	Code          []Bcode
	Used          map[int32]byte // identifier of used objects
	Init          []int32        // offsets of init funcs
	Strings       map[string]uint16
	StrOffset     []int32 // offsets of PUSHSTR
	Structs       map[string]uint16
	StructsList   []StructInfo
	StructsOffset []int32 // offsets of struct types
	BlockFlags    int16
	Pos           []CodePos
}

type CodePos struct {
	Offset int32  // byte code position
	Path   uint16 // Path index
	Name   uint16 // Name index
	Line   uint16 // Line
	Column uint16 // Column
}

type StructInfo struct {
	Name   string
	Fields []uint16 // types
	Keys   []string
}

type Exec struct {
	Code    []Bcode
	Funcs   map[int32]int32
	Init    []int32  // offsets of init funcs (initializing constants)
	Strings []string // string resources
	Structs []StructInfo
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
type AssignFloatFunc func(*float64, float64) (float64, error)
type AssignStrFunc func(*string, interface{}) (string, error)
type AssignAnyFunc func(interface{}, interface{}) (interface{}, error)

//type SetIndexFunc func(interface{}, interface{}, interface{}) error

const (
	TYPENONE   = 0
	TYPEINT    = 0x011
	TYPEBOOL   = 0x021
	TYPECHAR   = 0x031
	TYPESTR    = 0x042
	TYPEFLOAT  = 0x053
	TYPEARR    = 0x064
	TYPERANGE  = 0x074
	TYPEMAP    = 0x084
	TYPEBUF    = 0x094
	TYPEFUNC   = 0x0a4
	TYPESTRUCT = 0x104

	BlBreak    = 0x0001
	BlContinue = 0x0002
	BlVars     = 0x0004
	BlPars     = 0x0008
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
	PUSHFLOAT // + float64
	PUSHSTR   // & (strid << 16 )
	PUSHFUNC  // + id func
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
	ADDFLOAT  // float + float
	SUBFLOAT  // float - float
	MULFLOAT  // float * float
	DIVFLOAT  // float / float
	SIGNFLOAT // -float
	EQFLOAT   // float == float
	LTFLOAT   // float < float
	GTFLOAT   // float > float
	ADDSTR    // str + str
	EQSTR     // str == str
	LTSTR     // str < str
	GTSTR     // str > str
	GETVAR    // & (block shift<<16) + int16 type + int16 index
	SETVAR    // & (block shift<<16) + int16 type + int16 index + int16 index count + int16 assign
	DUP       // & (type<<16) duplicate top
	POP       // & (type<<16) pop top
	CYCLE     // cycle counter
	JMP       // + int32 jump
	JZE       // + int32 jump if the top value is zero
	JNZ       // + int32 jump if the top value is not zero
	JEQ       // & (type<<16) + int32 jump if equals for case statement
	INITVARS  // & (flags<<16) initializing variables + offset break + offset continue +
	// parcount<<16 | var count +
	DELVARS   // delete variables
	INITOBJ   // & (count<<16) create a new object + int16 type +int16 type item
	RANGE     // create range
	ARRAY     // &(count<<16) create array + int32 types
	LEN       // & (type<<16) length
	FORINC    // & (index<<16) increment counter
	BREAK     // break
	CONTINUE  // continue
	RET       // & (type<<16) return from function
	END       // end of the function
	CONSTBYID // + int32 id of the object
	CALLBYID  // & (par count<<16) + int32 id of the object
	EMBED     // & (embed id << 16) calls embedded func + int32 count for variadic funcs
	// + [variadic types]
	IOTA // & (iota<<16)

	INDEX        // & (int32 count) + {(type input<<16) + result type}
	ASSIGNPTR    // & (int16 type << 16)
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
)
