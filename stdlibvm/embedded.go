// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import "github.com/gentee/gentee/core"

var Embedded = []core.Embed{
	{Func: StrºBool, Return: core.TYPESTR, Params: []uint16{core.TYPEBOOL}},
	{Func: StrºInt, Return: core.TYPESTR, Params: []uint16{core.TYPEINT}},
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
	{Func: AddºCharChar, Return: core.TYPESTR, Params: []uint16{core.TYPECHAR, core.TYPECHAR}},
	{Func: AddºCharStr, Return: core.TYPESTR, Params: []uint16{core.TYPECHAR, core.TYPESTR}},
	{Func: AddºStrChar, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPECHAR}},
	{Func: ExpStrºChar, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPECHAR}},
	{Func: GreaterºCharChar, Return: core.TYPEBOOL, Params: []uint16{core.TYPECHAR, core.TYPECHAR}},
	{Func: LessºCharChar, Return: core.TYPEBOOL, Params: []uint16{core.TYPECHAR, core.TYPECHAR}},
	{Func: strºChar, Return: core.TYPESTR, Params: []uint16{core.TYPECHAR}},
	{Func: floatºInt, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEINT}},
	// 20
	{Func: floatºStr, Return: core.TYPEFLOAT, Params: []uint16{core.TYPESTR}, CanError: true},
	{Func: AddºFloatInt, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEFLOAT, core.TYPEINT}},
	{Func: AddºIntFloat, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEINT, core.TYPEFLOAT}},
	{Func: intºFloat, Return: core.TYPEINT, Params: []uint16{core.TYPEFLOAT}},
	{Func: StrºFloat, Return: core.TYPESTR, Params: []uint16{core.TYPEFLOAT}},
	{Func: MulºFloatInt, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEFLOAT, core.TYPEINT}},
	{Func: SubºFloatInt, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEFLOAT, core.TYPEINT}},
	{Func: SubºIntFloat, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEINT, core.TYPEFLOAT}},
	{Func: DivºFloatInt, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEFLOAT, core.TYPEINT},
		CanError: true},
	{Func: DivºIntFloat, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEINT, core.TYPEFLOAT},
		CanError: true},
	// 30
	{Func: MulºIntFloat, Return: core.TYPEFLOAT, Params: []uint16{core.TYPEINT, core.TYPEFLOAT}},
	{Func: ExpStrºFloat, Return: core.TYPESTR, Params: []uint16{core.TYPESTR, core.TYPEFLOAT}},
	{Func: EqualºFloatInt, Return: core.TYPEBOOL, Params: []uint16{core.TYPEFLOAT, core.TYPEINT}},
	{Func: GreaterºFloatInt, Return: core.TYPEBOOL, Params: []uint16{core.TYPEFLOAT, core.TYPEINT}},
	{Func: LessºFloatInt, Return: core.TYPEBOOL, Params: []uint16{core.TYPEFLOAT, core.TYPEINT}},
	{Func: TrimSpaceºStr, Return: core.TYPESTR, Params: []uint16{core.TYPESTR}},
	{Func: LinesºStr, Return: core.TYPEARR, Params: []uint16{core.TYPESTR}},
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

var EmbedFloat = []core.AssignFloatFunc{
	AssignºFloatFloat,
	AssignAddºFloatFloat,
	AssignSubºFloatFloat,
	AssignMulºFloatFloat,
	AssignDivºFloatFloat,
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
	AssignAddºArrAny,
	AssignºMapMap,
	AssignAddºBufChar,
	AssignAddºBufInt,
	AssignAddºBufBuf,
	AssignAddºBufStr,
	AssignºFnFn,
}
