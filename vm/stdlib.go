// Code generated by go generate; DO NOT EDIT.
// This file was generated by github.com/gentee/gentee/vm/generate/generate.go at
// 2019/09/27 20:42:46 +05

package vm

import "github.com/gentee/gentee/core"

var EmbedFuncs = []core.Embed{
	{"Add", "buf,buf", "buf", 0, AddºBufBuf, core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPEBUF}, false, false, false},
	{"Add", "char,char", "str", 1, AddºCharChar, core.TYPESTR, []uint16{core.TYPECHAR,core.TYPECHAR}, false, false, false},
	{"Add", "char,str", "str", 2, AddºCharStr, core.TYPESTR, []uint16{core.TYPECHAR,core.TYPESTR}, false, false, false},
	{"Add", "float,float", "float", core.ADDFLOAT, nil, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Add", "float,int", "float", 4, AddºFloatInt, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, false},
	{"Add", "int,float", "float", 5, AddºIntFloat, core.TYPEFLOAT, []uint16{core.TYPEINT,core.TYPEFLOAT}, false, false, false},
	{"Add", "int,int", "int", core.ADD, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Add", "str,char", "str", 7, AddºStrChar, core.TYPESTR, []uint16{core.TYPESTR,core.TYPECHAR}, false, false, false},
	{"Add", "str,str", "str", core.ADDSTR, nil, core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"AddHours", "time,int", "time", 9, AddHoursºTimeInt, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPEINT}, false, true, false},
	{"arr", "set", "arr.int", 10, arrºSet, core.TYPEARR, []uint16{core.TYPESET}, false, false, false},
	{"Assign", "bool,bool", "bool", core.ASSIGN, nil, core.TYPEBOOL, []uint16{core.TYPEBOOL,core.TYPEBOOL}, false, false, false},
	{"Assign", "buf,buf", "buf", core.ASSIGN, nil, core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPEBUF}, false, false, false},
	{"Assign", "char,char", "char", core.ASSIGN, nil, core.TYPECHAR, []uint16{core.TYPECHAR,core.TYPECHAR}, false, false, false},
	{"Assign", "float,float", "float", core.ASSIGN, nil, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Assign", "int,char", "int", core.ASSIGN, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPECHAR}, false, false, false},
	{"Assign", "int,int", "int", core.ASSIGN, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Assign", "set,set", "set", core.ASSIGN, nil, core.TYPESET, []uint16{core.TYPESET,core.TYPESET}, false, false, false},
	{"Assign", "str,bool", "str", 18, core.AssignStrFunc(AssignºStrBool), core.TYPESTR, []uint16{core.TYPESTR,core.TYPEBOOL}, false, false, false},
	{"Assign", "str,int", "str", 19, core.AssignStrFunc(AssignºStrInt), core.TYPESTR, []uint16{core.TYPESTR,core.TYPEINT}, false, false, false},
	{"Assign", "str,str", "str", core.ASSIGN, nil, core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"AssignºArrArr", "arr*,arr*", "arr*", core.ASSIGN, nil, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"AssignºFnFn", "fn,fn", "fn", core.ASSIGN, nil, core.TYPEFUNC, []uint16{core.TYPEFUNC,core.TYPEFUNC}, false, false, false},
	{"AssignºMapMap", "map*,map*", "map*", core.ASSIGN, nil, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"AssignºStructStruct", "struct,struct", "struct", core.ASSIGN, nil, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"AssignAdd", "arr.bool,bool", "arr.bool", 25, core.AssignAnyFunc(AssignAddºArrAny), core.TYPEARR, []uint16{core.TYPEARR,core.TYPEBOOL}, false, false, false},
	{"AssignAdd", "arr.int,int", "arr.int", 26, core.AssignAnyFunc(AssignAddºArrAny), core.TYPEARR, []uint16{core.TYPEARR,core.TYPEINT}, false, false, false},
	{"AssignAdd", "arr.str,str", "arr.str", 27, core.AssignAnyFunc(AssignAddºArrAny), core.TYPEARR, []uint16{core.TYPEARR,core.TYPESTR}, false, false, false},
	{"AssignAdd", "buf,buf", "buf", 28, core.AssignAnyFunc(AssignAddºBufBuf), core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPEBUF}, false, false, false},
	{"AssignAdd", "buf,char", "buf", 29, core.AssignAnyFunc(AssignAddºBufChar), core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPECHAR}, false, false, false},
	{"AssignAdd", "buf,int", "buf", 30, core.AssignAnyFunc(AssignAddºBufInt), core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPEINT}, false, false, true},
	{"AssignAdd", "buf,str", "buf", 31, core.AssignAnyFunc(AssignAddºBufStr), core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPESTR}, false, false, false},
	{"AssignAdd", "float,float", "float", 32, core.AssignFloatFunc(AssignAddºFloatFloat), core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"AssignAdd", "int,int", "int", 33, core.AssignIntFunc(AssignAddºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"AssignAdd", "set,set", "set", 34, core.AssignAnyFunc(AssignAddºSetSet), core.TYPESET, []uint16{core.TYPESET,core.TYPESET}, false, false, false},
	{"AssignAdd", "str,char", "str", 35, core.AssignStrFunc(AssignAddºStrChar), core.TYPESTR, []uint16{core.TYPESTR,core.TYPECHAR}, false, false, false},
	{"AssignAdd", "str,str", "str", 36, core.AssignStrFunc(AssignAddºStrStr), core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"AssignAddºArrArr", "arr.arr*,arr*", "arr.arr*", 37, core.AssignAnyFunc(AssignAddºArrAny), core.TYPEARR, []uint16{core.TYPEARR,core.TYPESTRUCT}, false, false, false},
	{"AssignAddºArrMap", "arr.map*,map*", "arr.map*", 38, core.AssignAnyFunc(AssignAddºArrAny), core.TYPEARR, []uint16{core.TYPEARR,core.TYPESTRUCT}, false, false, false},
	{"AssignBitAnd", "buf,buf", "buf", core.ASSIGNPTR, nil, core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPEBUF}, false, false, false},
	{"AssignBitAnd", "int,int", "int", 40, core.AssignIntFunc(AssignBitAndºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"AssignBitAnd", "set,set", "set", core.ASSIGNPTR, nil, core.TYPESET, []uint16{core.TYPESET,core.TYPESET}, false, false, false},
	{"AssignBitAndºArrArr", "arr*,arr*", "arr*", core.ASSIGNPTR, nil, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"AssignBitAndºMapMap", "map*,map*", "map*", core.ASSIGNPTR, nil, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"AssignBitAndºStructStruct", "struct,struct", "struct", core.ASSIGNPTR, nil, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"AssignBitOr", "int,int", "int", 45, core.AssignIntFunc(AssignBitOrºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"AssignBitXor", "int,int", "int", 46, core.AssignIntFunc(AssignBitXorºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"AssignDiv", "float,float", "float", 47, core.AssignFloatFunc(AssignDivºFloatFloat), core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, true},
	{"AssignDiv", "int,int", "int", 48, core.AssignIntFunc(AssignDivºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"AssignMod", "int,int", "int", 49, core.AssignIntFunc(AssignModºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"AssignLShift", "int,int", "int", 50, core.AssignIntFunc(AssignLShiftºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"AssignMul", "float,float", "float", 51, core.AssignFloatFunc(AssignMulºFloatFloat), core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"AssignMul", "int,int", "int", 52, core.AssignIntFunc(AssignMulºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"AssignRShift", "int,int", "int", 53, core.AssignIntFunc(AssignRShiftºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"AssignSub", "float,float", "float", 54, core.AssignFloatFunc(AssignSubºFloatFloat), core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"AssignSub", "int,int", "int", 55, core.AssignIntFunc(AssignSubºIntInt), core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Base64", "buf", "str", 56, Base64ºBuf, core.TYPESTR, []uint16{core.TYPEBUF}, false, false, false},
	{"BitAnd", "int,int", "int", core.BITAND, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"BitAnd", "set,set", "set", 58, BitAndºSetSet, core.TYPESET, []uint16{core.TYPESET,core.TYPESET}, false, false, false},
	{"BitNot", "int", "int", core.BITNOT, nil, core.TYPEINT, []uint16{core.TYPEINT}, false, false, false},
	{"BitNot", "set", "set", 60, BitNotºSet, core.TYPESET, []uint16{core.TYPESET}, false, false, false},
	{"BitOr", "int,int", "int", core.BITOR, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"BitOr", "set,set", "set", 62, BitOrºSetSet, core.TYPESET, []uint16{core.TYPESET,core.TYPESET}, false, false, false},
	{"BitXor", "int,int", "int", core.BITXOR, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"bol", "float", "bool", 64, boolºFloat, core.TYPEBOOL, []uint16{core.TYPEFLOAT}, false, false, false},
	{"bool", "int", "bool", 65, boolºInt, core.TYPEBOOL, []uint16{core.TYPEINT}, false, false, false},
	{"bool", "str", "bool", 66, boolºStr, core.TYPEBOOL, []uint16{core.TYPESTR}, false, false, false},
	{"buf", "str", "buf", 67, bufºStr, core.TYPEBUF, []uint16{core.TYPESTR}, false, false, false},
	{"Ceil", "float", "int", 68, CeilºFloat, core.TYPEINT, []uint16{core.TYPEFLOAT}, false, false, false},
	{"Command", "str", "", 69, Command, core.TYPENONE, []uint16{core.TYPESTR}, false, false, true},
	{"CommandOutput", "str", "str", 70, CommandOutput, core.TYPESTR, []uint16{core.TYPESTR}, false, false, true},
	{"Date", "int,int,int", "time", 71, DateºInts, core.TYPESTRUCT, []uint16{core.TYPEINT,core.TYPEINT,core.TYPEINT}, false, true, false},
	{"DateTime", "int,int,int,int,int,int", "time", 72, DateTimeºInts, core.TYPESTRUCT, []uint16{core.TYPEINT,core.TYPEINT,core.TYPEINT,core.TYPEINT,core.TYPEINT,core.TYPEINT}, false, true, false},
	{"Days", "time", "int", 73, DaysºTime, core.TYPEINT, []uint16{core.TYPESTRUCT}, false, false, false},
	{"Del", "buf,int,int", "buf", 74, DelºBufIntInt, core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPEINT,core.TYPEINT}, false, false, false},
	{"DelAuto", "map*,str", "map*", 75, DelºMapStr, core.TYPESTRUCT, []uint16{core.TYPESTRUCT,core.TYPESTR}, false, false, false},
	{"Div", "float,float", "float", core.DIVFLOAT, nil, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, true},
	{"Div", "float,int", "float", 77, DivºFloatInt, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, true},
	{"Div", "int,float", "float", 78, DivºIntFloat, core.TYPEFLOAT, []uint16{core.TYPEINT,core.TYPEFLOAT}, false, false, true},
	{"Div", "int,int", "int", core.DIV, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"Equal", "char,char", "bool", core.EQ, nil, core.TYPEBOOL, []uint16{core.TYPECHAR,core.TYPECHAR}, false, false, false},
	{"Equal", "float,float", "bool", core.EQFLOAT, nil, core.TYPEBOOL, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Equal", "float,int", "bool", 82, EqualºFloatInt, core.TYPEBOOL, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, false},
	{"Equal", "int,int", "bool", core.EQ, nil, core.TYPEBOOL, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Equal", "str,str", "bool", core.EQSTR, nil, core.TYPEBOOL, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"Equal", "time,time", "bool", 85, EqualºTimeTime, core.TYPEBOOL, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"ExpStr", "str,bool", "str", 86, ExpStrºBool, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEBOOL}, false, false, false},
	{"ExpStr", "str,char", "str", 87, ExpStrºChar, core.TYPESTR, []uint16{core.TYPESTR,core.TYPECHAR}, false, false, false},
	{"ExpStr", "str,float", "str", 88, ExpStrºFloat, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEFLOAT}, false, false, false},
	{"ExpStr", "str,int", "str", 89, ExpStrºInt, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEINT}, false, false, false},
	{"ExpStr", "str,str", "str", core.ADDSTR, nil, core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"Find", "str,str", "int", 91, FindºStrStr, core.TYPEINT, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"float", "int", "float", 92, floatºInt, core.TYPEFLOAT, []uint16{core.TYPEINT}, false, false, false},
	{"float", "str", "float", 93, floatºStr, core.TYPEFLOAT, []uint16{core.TYPESTR}, false, false, true},
	{"Floor", "float", "int", 94, FloorºFloat, core.TYPEINT, []uint16{core.TYPEFLOAT}, false, false, false},
	{"Format", "str", "str", 95, FormatºStr, core.TYPESTR, []uint16{core.TYPESTR}, true, false, false},
	{"Format", "str,time", "str", 96, FormatºTimeStr, core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTRUCT}, false, false, false},
	{"GetEnv", "str", "str", 97, GetEnv, core.TYPESTR, []uint16{core.TYPESTR}, false, false, false},
	{"Greater", "char,char", "bool", 98, GreaterºCharChar, core.TYPEBOOL, []uint16{core.TYPECHAR,core.TYPECHAR}, false, false, false},
	{"Greater", "float,float", "bool", core.GTFLOAT, nil, core.TYPEBOOL, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Greater", "float,int", "bool", 100, GreaterºFloatInt, core.TYPEBOOL, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, false},
	{"Greater", "int,int", "bool", core.GT, nil, core.TYPEBOOL, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Greater", "str,str", "bool", core.GTSTR, nil, core.TYPEBOOL, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"Greater", "time,time", "bool", 103, GreaterºTimeTime, core.TYPEBOOL, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"HasPrefix", "str,str", "bool", 104, HasPrefixºStrStr, core.TYPEBOOL, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"HasSuffix", "str,str", "bool", 105, HasSuffixºStrStr, core.TYPEBOOL, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"Hex", "buf", "str", 106, HexºBuf, core.TYPESTR, []uint16{core.TYPEBUF}, false, false, false},
	{"Join", "arr.str,str", "str", 107, JoinºArrStr, core.TYPESTR, []uint16{core.TYPEARR,core.TYPESTR}, false, false, false},
	{"Insert", "buf,int,buf", "buf", 108, InsertºBufIntBuf, core.TYPEBUF, []uint16{core.TYPEBUF,core.TYPEINT,core.TYPEBUF}, false, false, false},
	{"int", "bool", "int", core.NOP, nil, core.TYPEINT, []uint16{core.TYPEBOOL}, false, false, false},
	{"int", "char", "int", core.NOP, nil, core.TYPEINT, []uint16{core.TYPECHAR}, false, false, false},
	{"int", "float", "int", 111, intºFloat, core.TYPEINT, []uint16{core.TYPEFLOAT}, false, false, false},
	{"int", "str", "int", 112, intºStr, core.TYPEINT, []uint16{core.TYPESTR}, false, false, true},
	{"int", "time", "int", 113, intºTime, core.TYPEINT, []uint16{core.TYPESTRUCT}, false, false, false},
	{"IsKeyAuto", "map*,str", "bool", 114, IsKeyºMapStr, core.TYPEBOOL, []uint16{core.TYPESTRUCT,core.TYPESTR}, false, false, false},
	{"Left", "str,int", "str", 115, LeftºStrInt, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEINT}, false, false, false},
	{"LenºArr", "arr*", "int", core.TYPESTRUCT<<16 | core.LEN, nil, core.TYPEINT, []uint16{core.TYPESTRUCT}, false, false, false},
	{"Len", "buf", "int", core.TYPEBUF<<16 | core.LEN, nil, core.TYPEINT, []uint16{core.TYPEBUF}, false, false, false},
	{"LenºMap", "map*", "int", core.TYPESTRUCT<<16 | core.LEN, nil, core.TYPEINT, []uint16{core.TYPESTRUCT}, false, false, false},
	{"Len", "set", "int", core.TYPESET<<16 | core.LEN, nil, core.TYPEINT, []uint16{core.TYPESET}, false, false, false},
	{"Len", "str", "int", core.TYPESTR<<16 | core.LEN, nil, core.TYPEINT, []uint16{core.TYPESTR}, false, false, false},
	{"Less", "char,char", "bool", 121, LessºCharChar, core.TYPEBOOL, []uint16{core.TYPECHAR,core.TYPECHAR}, false, false, false},
	{"Less", "float,float", "bool", core.LTFLOAT, nil, core.TYPEBOOL, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Less", "float,int", "bool", 123, LessºFloatInt, core.TYPEBOOL, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, false},
	{"Less", "int,int", "bool", core.LT, nil, core.TYPEBOOL, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Less", "str,str", "bool", core.LTSTR, nil, core.TYPEBOOL, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"Less", "time,time", "bool", 126, LessºTimeTime, core.TYPEBOOL, []uint16{core.TYPESTRUCT,core.TYPESTRUCT}, false, false, false},
	{"Lines", "str", "arr.str", 127, LinesºStr, core.TYPEARR, []uint16{core.TYPESTR}, false, false, false},
	{"Lower", "str", "str", 128, LowerºStr, core.TYPESTR, []uint16{core.TYPESTR}, false, false, false},
	{"LShift", "int,int", "int", core.LSHIFT, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"Max", "float,float", "float", 130, MaxºFloatFloat, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Max", "int,int", "int", 131, MaxºIntInt, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Min", "float,float", "float", 132, MinºFloatFloat, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Min", "int,int", "int", 133, MinºIntInt, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Mod", "int,int", "int", core.MOD, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"Mul", "float,float", "float", core.MULFLOAT, nil, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Mul", "float,int", "float", 136, MulºFloatInt, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, false},
	{"Mul", "int,float", "float", 137, MulºIntFloat, core.TYPEFLOAT, []uint16{core.TYPEINT,core.TYPEFLOAT}, false, false, false},
	{"Mul", "int,int", "int", core.MUL, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"NewKeyValue", "int,int", "keyval", core.NOP, nil, core.TYPESTRUCT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"NewRange", "int,int", "range", core.RANGE, nil, core.TYPERANGE, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Not", "bool", "bool", core.NOT, nil, core.TYPEBOOL, []uint16{core.TYPEBOOL}, false, false, false},
	{"Now", "", "time", 142, Now, core.TYPESTRUCT, nil, false, true, false},
	{"ParseTime", "str,str", "time", 143, ParseTimeºStrStr, core.TYPESTRUCT, []uint16{core.TYPESTR,core.TYPESTR}, false, true, true},
	{"Repeat", "str,int", "str", 144, RepeatºStrInt, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEINT}, false, false, false},
	{"Replace", "str,str,str", "str", 145, ReplaceºStrStrStr, core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTR,core.TYPESTR}, false, false, false},
	{"Round", "float", "int", 146, RoundºFloat, core.TYPEINT, []uint16{core.TYPEFLOAT}, false, false, false},
	{"Round", "float,int", "float", 147, RoundºFloatInt, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, false},
	{"RShift", "int,int", "int", core.RSHIFT, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, true},
	{"set", "arr.int", "set", 149, setºArr, core.TYPESET, []uint16{core.TYPEARR}, false, false, true},
	{"Set", "set,int", "set", 150, SetºSet, core.TYPESET, []uint16{core.TYPESET,core.TYPEINT}, false, false, true},
	{"set", "str", "set", 151, setºStr, core.TYPESET, []uint16{core.TYPESTR}, false, false, true},
	{"SetEnv", "str,str", "str", 152, SetEnv, core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTR}, false, false, true},
	{"SetEnv", "str,int", "str", 153, SetEnv, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEINT}, false, false, true},
	{"SetEnv", "str,bool", "str", 154, SetEnvBool, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEBOOL}, false, false, true},
	{"Shift", "str", "str", 155, ShiftºStr, core.TYPESTR, []uint16{core.TYPESTR}, false, false, false},
	{"Sign", "float", "float", core.SIGNFLOAT, nil, core.TYPEFLOAT, []uint16{core.TYPEFLOAT}, false, false, false},
	{"Sign", "int", "int", core.SIGN, nil, core.TYPEINT, []uint16{core.TYPEINT}, false, false, false},
	{"Sort", "arr.str", "arr.str", 158, SortºArr, core.TYPEARR, []uint16{core.TYPEARR}, false, false, false},
	{"Split", "str,str", "arr.str", 159, SplitºStrStr, core.TYPEARR, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"str", "bool", "str", 160, strºBool, core.TYPESTR, []uint16{core.TYPEBOOL}, false, false, false},
	{"str", "buf", "str", 161, strºBuf, core.TYPESTR, []uint16{core.TYPEBUF}, false, false, false},
	{"str", "char", "str", 162, strºChar, core.TYPESTR, []uint16{core.TYPECHAR}, false, false, false},
	{"str", "float", "str", 163, strºFloat, core.TYPESTR, []uint16{core.TYPEFLOAT}, false, false, false},
	{"str", "int", "str", 164, strºInt, core.TYPESTR, []uint16{core.TYPEINT}, false, false, false},
	{"str", "set", "str", 165, strºSet, core.TYPESTR, []uint16{core.TYPESET}, false, false, false},
	{"Sub", "float,float", "float", core.SUBFLOAT, nil, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEFLOAT}, false, false, false},
	{"Sub", "float,int", "float", 167, SubºFloatInt, core.TYPEFLOAT, []uint16{core.TYPEFLOAT,core.TYPEINT}, false, false, false},
	{"Sub", "int,float", "float", 168, SubºIntFloat, core.TYPEFLOAT, []uint16{core.TYPEINT,core.TYPEFLOAT}, false, false, false},
	{"Sub", "int,int", "int", core.SUB, nil, core.TYPEINT, []uint16{core.TYPEINT,core.TYPEINT}, false, false, false},
	{"Substr", "str,int,int", "str", 170, SubstrºStrIntInt, core.TYPESTR, []uint16{core.TYPESTR,core.TYPEINT,core.TYPEINT}, false, false, true},
	{"sysBufNil", "", "buf", 171, sysBufNil, core.TYPEBUF, nil, false, false, false},
	{"time", "int", "time", 172, timeºInt, core.TYPESTRUCT, []uint16{core.TYPEINT}, false, true, false},
	{"Toggle", "set,int", "bool", 173, ToggleºSetInt, core.TYPEBOOL, []uint16{core.TYPESET,core.TYPEINT}, false, false, false},
	{"TrimRight", "str,str", "str", 174, TrimRightºStr, core.TYPESTR, []uint16{core.TYPESTR,core.TYPESTR}, false, false, false},
	{"TrimSpace", "str", "str", 175, TrimSpaceºStr, core.TYPESTR, []uint16{core.TYPESTR}, false, false, false},
	{"UnBase64", "str", "buf", 176, UnBase64ºStr, core.TYPEBUF, []uint16{core.TYPESTR}, false, false, true},
	{"UnHex", "str", "buf", 177, UnHexºStr, core.TYPEBUF, []uint16{core.TYPESTR}, false, false, true},
	{"UnSet", "set, int", "set", 178, UnSetºSet, core.TYPESET, []uint16{core.TYPESET,core.TYPEINT}, false, false, true},
	{"Upper", "str", "str", 179, UpperºStr, core.TYPESTR, []uint16{core.TYPESTR}, false, false, false},
}