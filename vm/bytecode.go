// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

const (
	NOP    = iota
	PUSH16 // + int16
	PUSH32 // + int32
	PUSH64 // + int64
	ADD    // int + int
	SUB    // int - int
	MUL    // int * int
	DIV    // int / int
	MOD    // int % int
	EQINT  // int == int
	LTINT  // int < int
	GTINT  // int > int
	NOT    // logical not 1 => 0, 0 => 1
)
