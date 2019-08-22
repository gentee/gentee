// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import "github.com/gentee/gentee/core"

var Embedded = []core.Embed{
	{Func: strºBool, Return: core.TYPESTR, Params: []uint16{core.TYPEBOOL}},
	{Func: strºInt, Return: core.TYPESTR, Params: []uint16{core.TYPEINT}},
	{Func: intºStr, Return: core.TYPEINT, Params: []uint16{core.TYPESTR}, CanError: true},
	{Func: boolºInt, Return: core.TYPEBOOL, Params: []uint16{core.TYPEINT}},
	{Func: ExpStrºInt, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPEINT}},
	{Func: ExpStrºBool, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPEBOOL}},
	{Func: Command, Params: []uint16{core.TYPESTR}, CanError: true},
	{Func: CommandOutput, Return: core.TYPESTR, Params: []uint16{core.TYPESTR}, CanError: true},
	{Func: GetEnv, Return: core.TYPESTR, Params: []uint16{core.TYPESTR}},
	{Func: SetEnv, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPESTR},
		CanError: true},
	// 10
	{Func: SetEnv, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPEINT},
		CanError: true},
	{Func: SetEnvBool, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPEBOOL},
		CanError: true},
	/*-*/ {Func: AssignºStrBool, Return: core.TYPESTR, Params: []uint16{core.TYPEPTR, core.TYPEBOOL}},
	/*-*/ {Func: AssignºStrInt, Return: core.TYPESTR, Params: []uint16{core.TYPEPTR, core.TYPEINT}},
	{Func: AddºStrChar, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPECHAR}},
	/*-*/ {Func: AssignAddºStrChar, Return: core.TYPESTR, Params: []uint16{core.TYPEPTR, core.TYPECHAR}},
	{Func: AddºCharChar, Return: core.TYPESTR, Params: []uint16{core.TYPECHAR, core.TYPECHAR}},
	{Func: AddºCharStr, Return: core.TYPESTR, Params: []uint16{core.TYPECHAR, core.TYPESTR}},
	{Func: ExpStrºChar, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPECHAR}},
	{Func: GreaterºCharChar, Return: core.TYPEBOOL, Params: []uint16{core.TYPECHAR, core.TYPECHAR}},
	// 20
	{Func: LessºCharChar, Return: core.TYPEBOOL, Params: []uint16{core.TYPECHAR, core.TYPECHAR}},
	{Func: strºChar, Return: core.TYPESTR, Params: []uint16{core.TYPECHAR}},
	/*-*/ {Func: LenºArr, Return: core.TYPEINT, Params: []uint16{core.TYPEARR}},
	/*-*/ {Func: AssignAddºArrStr, Return: core.TYPEARR, Params: []uint16{core.TYPEPTR, core.TYPESTR}},
	/*-*/ {Func: AssignAddºArrInt, Return: core.TYPEARR, Params: []uint16{core.TYPEPTR, core.TYPEINT}},
}

var EmbedInt = []core.AssignIntFunc{
	AssignºIntInt,
	AssignAddºIntInt,
	AssignSubºIntInt,
	AssignMulºIntInt,
	AssignDivºIntInt,
	AssignModºIntInt,
	AssignBitAndºIntInt,
	AssignBitOrºIntInt,
	AssignBitXorºIntInt,
	AssignLShiftºIntInt,
	AssignRShiftºIntInt,
	IncDecºInt,
}

var EmbedStr = []core.AssignStrFunc{
	AssignºStrStr,
	AssignAddºStrStr,
	AssignºStrBool,
	AssignºStrInt,
	AssignAddºStrChar,
}

var EmbedAny = []core.AssignAnyFunc{
	AssignºArrArr,
	AssignAddºArrStr,
	AssignºMapMap,
}
