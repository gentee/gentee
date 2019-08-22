// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"reflect"

	"github.com/gentee/gentee/core"
	stdlib "github.com/gentee/gentee/stdlibvm"
)

type indexInfo struct {
	IntValue  int64
	StrValue  string
	PtrStr    *string
	IndStr    int
	Original  interface{}
	OrigIndex interface{}
	OrigType  int
	Index     interface{}
	RetType   int
	VarType   int
}

func newValue(vtype int) interface{} {
	switch vtype {
	case core.TYPEINT:
		return int64(0)
	case core.TYPECHAR:
		return int64(' ')
	case core.TYPESTR:
		return ``
	case core.TYPEARR:
		return core.NewArray()
	case core.TYPEMAP:
		return core.NewMap()
	}
	return nil
}

func (rt *Runtime) Run(i int64) (result interface{}, err error) {
	var iInfo indexInfo

	top := Call{}
	code := rt.Owner.Exec.Code
	end := int64(len(code))

	pushIndex := func(ptr interface{}, retType int) {
		switch retType & 0xff {
		case core.STACKSTR:
			rt.SStr[top.Str] = *ptr.(*string)
			top.Str++
		case core.STACKINT:
			rt.SInt[top.Int] = *ptr.(*int64)
			top.Int++
		case core.STACKANY:
			rt.SAny[top.Any] = ptr
			top.Any++
		default:
			fmt.Printf("GET TYPE %T\n", ptr)
		}
	}
	getIndex := func(set bool) (ptr interface{}, err error) {
		var count int
		if code[i+2]&0xffff == core.INDEX {
			count = int(code[i+2] >> 16)
		}
		iInfo.Original = nil
		iInfo.Index = nil
		iInfo.PtrStr = nil
		blockOff := rt.Calls[int32(len(rt.Calls)-1-int(code[i])>>16)]
		i++
		typeVar := int(code[i]) >> 16
		typeRet := typeVar
		varIndex := int32(int(code[i]) & 0xffff)
		switch typeVar & 0xff {
		case core.STACKSTR:
			ptr = &rt.SStr[blockOff.Str+varIndex]
		case core.STACKANY:
			ptr = rt.SAny[blockOff.Any+varIndex]
		default:
			ptr = &rt.SInt[blockOff.Int+varIndex]
		}
		if count > 0 {
			i++
			for ind := 0; ind < count; ind++ {
				i++
				typeVar = int(code[i]) >> 16
				typeRet = int(code[i]) & 0x0fff

				if int(code[i])&0x8000 != 0 {
					top.Str--
					iInfo.Index = rt.SStr[top.Str]
				} else {
					top.Int--
					iInfo.Index = rt.SInt[top.Int]
					if iInfo.Index.(int64) < 0 {
						err = runtimeError(rt, i, ErrIndexOut)
						return
					}
				}
				if typeVar&0xff == core.STACKANY && typeRet&0xff != core.STACKANY &&
					iInfo.Original == nil {
					iInfo.Original = ptr
					iInfo.OrigIndex = iInfo.Index
					iInfo.OrigType = typeRet
				}
				switch typeVar {
				case core.TYPESTR:
					if len(*ptr.(*string)) <= int(iInfo.Index.(int64)) {
						err = runtimeError(rt, i, ErrIndexOut)
						return
					}
					iInfo.PtrStr = ptr.(*string)
					iInfo.IndStr = int(iInfo.Index.(int64))
					runes := []rune(*iInfo.PtrStr)
					iInfo.IntValue = int64(runes[iInfo.Index.(int64)])
					ptr = &iInfo.IntValue
				case core.TYPEARR:
					if ptr.(*core.Array).Len() <= int(iInfo.Index.(int64)) {
						err = runtimeError(rt, i, ErrIndexOut)
						return
					}
					value := ptr.(*core.Array).Data[iInfo.Index.(int64)]
					switch typeRet & 0xff {
					case core.STACKINT:
						iInfo.IntValue = value.(int64)
						ptr = &iInfo.IntValue
					case core.STACKSTR:
						iInfo.StrValue = value.(string)
						ptr = &iInfo.StrValue
					default:
						ptr = value
					}
				case core.TYPERANGE:
					rangeVal := ptr.(*core.Range)
					iInfo.IntValue = rangeVal.From - iInfo.Index.(int64)
					if rangeVal.From < rangeVal.To {
						iInfo.IntValue = rangeVal.From + iInfo.Index.(int64)
					}
					ptr = &iInfo.IntValue
				case core.TYPEMAP:
					var value interface{}

					if key, ok := iInfo.Index.(string); ok {
						if value, ok = ptr.(*core.Map).Data[key]; !ok {
							if !set {
								err = runtimeError(rt, i, ErrMapIndex, key)
								return
							}
							value = newValue(typeRet)
							ptr.(*core.Map).Data[key] = value
							ptr.(*core.Map).Keys = append(ptr.(*core.Map).Keys, key)

						}
					} else {
						value = ptr.(*core.Map).Data[ptr.(*core.Map).Keys[iInfo.Index.(int64)]]
					}
					switch typeRet & 0xff {
					case core.STACKINT:
						iInfo.IntValue = value.(int64)
						ptr = &iInfo.IntValue
					case core.STACKSTR:
						iInfo.StrValue = value.(string)
						ptr = &iInfo.StrValue
					default:
						ptr = value
					}
				default:
					fmt.Printf("INDEX ANY %x %x\n", typeVar, typeRet)
				}
				//fmt.Printf("IND %x %x %v\r\n", typeVar, typeRet, iInfo)
			}
		}
		iInfo.RetType = typeRet
		iInfo.VarType = typeVar
		return
	}

main:
	for i < end {
		switch code[i] & 0xffff {
		case core.PUSH32:
			i++
			rt.SInt[top.Int] = int64(code[i])
			top.Int++
		case core.PUSH64:
			i += 2
			rt.SInt[top.Int] = int64((uint64(code[i-1]) << 32) | (uint64(code[i]) & 0xffffffff))
			top.Int++
		case core.PUSHSTR:
			rt.SStr[top.Str] = rt.Owner.Exec.Strings[(code[i])>>16]
			top.Str++
		case core.ADD:
			top.Int--
			rt.SInt[top.Int-1] += rt.SInt[top.Int]
		case core.SUB:
			top.Int--
			rt.SInt[top.Int-1] -= rt.SInt[top.Int]
		case core.MUL:
			top.Int--
			rt.SInt[top.Int-1] *= rt.SInt[top.Int]
		case core.DIV:
			top.Int--
			if rt.SInt[top.Int] == 0 {
				return nil, runtimeError(rt, i, ErrDivZero)
			}
			rt.SInt[top.Int-1] /= rt.SInt[top.Int]
		case core.MOD:
			top.Int--
			if rt.SInt[top.Int] == 0 {
				return nil, runtimeError(rt, i, ErrDivZero)
			}
			rt.SInt[top.Int-1] %= rt.SInt[top.Int]
		case core.BITOR:
			top.Int--
			rt.SInt[top.Int-1] |= rt.SInt[top.Int]
		case core.BITXOR:
			top.Int--
			rt.SInt[top.Int-1] ^= rt.SInt[top.Int]
		case core.BITAND:
			top.Int--
			rt.SInt[top.Int-1] &= rt.SInt[top.Int]
		case core.LSHIFT:
			top.Int--
			if rt.SInt[top.Int] < 0 {
				return nil, runtimeError(rt, i, ErrShift)
			}
			rt.SInt[top.Int-1] <<= uint32(rt.SInt[top.Int])
		case core.RSHIFT:
			top.Int--
			if rt.SInt[top.Int] < 0 {
				return nil, runtimeError(rt, i, ErrShift)
			}
			rt.SInt[top.Int-1] >>= uint32(rt.SInt[top.Int])
		case core.BITNOT:
			rt.SInt[top.Int-1] = ^rt.SInt[top.Int-1]
		case core.SIGN:
			rt.SInt[top.Int-1] = -rt.SInt[top.Int-1]
		case core.EQ:
			top.Int--
			if rt.SInt[top.Int-1] == rt.SInt[top.Int] {
				rt.SInt[top.Int-1] = 1
			} else {
				rt.SInt[top.Int-1] = 0
			}
		case core.LT:
			top.Int--
			if rt.SInt[top.Int-1] < rt.SInt[top.Int] {
				rt.SInt[top.Int-1] = 1
			} else {
				rt.SInt[top.Int-1] = 0
			}
		case core.GT:
			top.Int--
			if rt.SInt[top.Int-1] > rt.SInt[top.Int] {
				rt.SInt[top.Int-1] = 1
			} else {
				rt.SInt[top.Int-1] = 0
			}
		case core.NOT:
			if rt.SInt[top.Int-1] == 0 {
				rt.SInt[top.Int-1] = 1
			} else {
				rt.SInt[top.Int-1] = 0
			}
		case core.ADDSTR:
			top.Str--
			rt.SStr[top.Str-1] += rt.SStr[top.Str]
		case core.EQSTR:
			top.Str -= 2
			if rt.SStr[top.Str] == rt.SStr[top.Str+1] {
				rt.SInt[top.Int] = 1
			} else {
				rt.SInt[top.Int] = 0
			}
			top.Int++
		case core.LTSTR:
			top.Str -= 2
			if rt.SStr[top.Str] < rt.SStr[top.Str+1] {
				rt.SInt[top.Int] = 1
			} else {
				rt.SInt[top.Int] = 0
			}
			top.Int++
		case core.GTSTR:
			top.Str -= 2
			if rt.SStr[top.Str] > rt.SStr[top.Str+1] {
				rt.SInt[top.Int] = 1
			} else {
				rt.SInt[top.Int] = 0
			}
			top.Int++
		case core.GETVAR:
			ptr, err := getIndex(false)
			if err != nil {
				return nil, err
			}
			//			fmt.Println(`GETVAR index`, ptr, iInfo)
			/*if iInfo.Index != nil {
			switch iInfo.VarType {
			case core.TYPESTR:
				runes := []rune(*ptr.(*string))
				iInfo.IntValue = int64(runes[iInfo.Index.(int64)])
			case core.TYPERANGE:
				iInfo.IntValue = int64(ptr.(*core.Range).From - iInfo.Index.(int64))
				if ptr.(*core.Range).From < ptr.(*core.Range).To {
					iInfo.IntValue = ptr.(*core.Range).From + iInfo.Index.(int64)
				}
			default:
				fmt.Printf("GET INDEX TYPE %x\n", iInfo.VarType)
			}*/
			/*				switch iInfo.RetType & 0xff {
							case core.STACKINT:
								ptr = &iInfo.IntValue
								//				case core.STACKSTR:
								//					ptr = &iInfo.StrValue
							default:
								fmt.Println(`GET RET`, iInfo.RetType)
							}*/
			//}
			pushIndex(ptr, iInfo.RetType)
			//			fmt.Println("GETVAR", rt.SInt[:top.Int], rt.SStr[:top.Str])
		case core.SETVAR:
			ptr, err := getIndex(true)
			if err != nil {
				return nil, err
			}
			i++
			assign := code[i] & 0xffff
			rightType := code[i] >> 16
			//fmt.Printf("ASSIGN %x %d %x %x\n", iInfo.RetType, assign, iInfo.VarType, rightType)
			if iInfo.Original == nil && assign == core.ASSIGN &&
				core.Bcode(iInfo.VarType) == rightType {
				switch v := ptr.(type) {
				case *int64:
					*v = rt.SInt[top.Int-1]
				case *string:
					*v = rt.SStr[top.Str-1]
				default:
					fmt.Println(`Assign`, rightType)
				}
				i++
				continue
			}
			var iValue interface{}
			switch rightType & 0xff {
			case core.STACKINT, core.STACKNONE: // STACKNONE is for inc dec
				top.Int--
				iValue = rt.SInt[top.Int]
			case core.STACKSTR:
				top.Str--
				iValue = rt.SStr[top.Str]
			case core.STACKANY:
				top.Any--
				iValue = rt.SAny[top.Any]
			default:
				fmt.Printf("iValue %x\n", rightType)
			}
			/*if iInfo.Index != nil {
								switch v := ptr.(type) {
								case *string:
									runes := []rune(*v)
									iInfo.IntValue = int64(runes[iInfo.Index.(int64)])
								default:
									fmt.Printf("SET INDEX TYPE %T\n", v)
								}
								switch iInfo.RetType & 0xff {
								case core.STACKINT:
									ptr = &iInfo.IntValue
								}
			}*/
			switch v := ptr.(type) {
			case *int64:
				rt.SInt[top.Int], err = stdlib.EmbedInt[assign-core.ASSIGN](
					v, iValue.(int64))
				top.Int++
			case *string:
				rt.SStr[top.Str], err = stdlib.EmbedStr[assign-core.ASSIGN](
					v, iValue)
				top.Str++
			case *core.Array, *core.Map:
				rt.SAny[top.Any], err = stdlib.EmbedAny[assign-core.ASSIGN](
					ptr, iValue)
				top.Any++
			default:
				fmt.Println(`Embed Assign`, rightType)
			}
			if err != nil {
				return nil, runtimeError(rt, i, err)
			}
			if iInfo.PtrStr != nil {
				runes := []rune(*iInfo.PtrStr)
				runes[iInfo.IndStr] = rune(iInfo.IntValue)
				*iInfo.PtrStr = string(runes)
			}
			if iInfo.Original != nil {
				var iValue interface{}
				switch iInfo.OrigType & 0xff {
				case core.STACKINT:
					iValue = iInfo.IntValue
				case core.STACKSTR:
					iValue = iInfo.StrValue
				}
				switch v := iInfo.Original.(type) {
				case *core.Array:
					v.Data[iInfo.OrigIndex.(int64)] = iValue
				case *core.Map:
					v.Data[iInfo.OrigIndex.(string)] = iValue
				}
			}
			/*		case core.ASSIGN:
					typeVar := int(code[i]) >> 16
					top.Any--
					switch typeVar {
					case core.TYPESTR:
						*(rt.SAny[top.Any].(*string)) = rt.SStr[top.Str-1]
					case core.TYPEARR:
						//				*(rt.SAny[top.Any].(*core.Array)) = rt.SAny[top.Any-1]
						*(rt.SAny[top.Any].(*interface{})) = rt.SAny[top.Any-1]
					default:
						*(rt.SAny[top.Any].(*int64)) = rt.SInt[top.Int-1]
					}*/
			/*		case core.ASSIGNADD:
								typeVar := int(code[i]) >> 16
								top.Any--
								switch typeVar {
								case core.TYPESTR:
									*(rt.SAny[top.Any].(*string)) += rt.SStr[top.Str-1]
									rt.SStr[top.Str-1] = *(rt.SAny[top.Any].(*string))
								default:
									*(rt.SAny[top.Any].(*int64)) += rt.SInt[top.Int-1]
									rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
								}
					case core.ASSIGNSUB:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							*(rt.SAny[top.Any].(*int64)) -= rt.SInt[top.Int-1]
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNMUL:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							*(rt.SAny[top.Any].(*int64)) *= rt.SInt[top.Int-1]
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNDIV:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							if rt.SInt[top.Int-1] == 0 {
								return nil, runtimeError(rt, i, ErrDivZero)
							}
							*(rt.SAny[top.Any].(*int64)) /= rt.SInt[top.Int-1]
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNMOD:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							if rt.SInt[top.Int-1] == 0 {
								return nil, runtimeError(rt, i, ErrDivZero)
							}
							*(rt.SAny[top.Any].(*int64)) %= rt.SInt[top.Int-1]
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNBITOR:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							*(rt.SAny[top.Any].(*int64)) |= rt.SInt[top.Int-1]
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNBITXOR:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							*(rt.SAny[top.Any].(*int64)) ^= rt.SInt[top.Int-1]
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNBITAND:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							*(rt.SAny[top.Any].(*int64)) &= rt.SInt[top.Int-1]
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNLSHIFT:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							*(rt.SAny[top.Any].(*int64)) <<= uint32(rt.SInt[top.Int-1])
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.ASSIGNRSHIFT:
						typeVar := int(code[i]) >> 16
						top.Any--
						switch typeVar {
						default:
							*(rt.SAny[top.Any].(*int64)) >>= uint32(rt.SInt[top.Int-1])
							rt.SInt[top.Int-1] = *(rt.SAny[top.Any].(*int64))
						}
					case core.INC:
						top.Any--
						val := *(rt.SAny[top.Any].(*int64)) + 1
						rt.SInt[top.Int] = val - (int64(code[i]) >> 16)
						*(rt.SAny[top.Any].(*int64)) = val
						top.Int++
					case core.DEC:
						top.Any--
						val := *(rt.SAny[top.Any].(*int64)) - 1
						rt.SInt[top.Int] = val + (int64(code[i]) >> 16)
						*(rt.SAny[top.Any].(*int64)) = val
						top.Int++*/
		case core.DUP:
			switch code[i] >> 16 & 0xff {
			case core.STACKSTR:
				rt.SStr[top.Str] = rt.SStr[top.Str-1]
				top.Str++
			case core.STACKANY:
				rt.SAny[top.Any] = rt.SAny[top.Any-1]
				top.Any++
			default:
				rt.SInt[top.Int] = rt.SInt[top.Int-1]
				top.Int++
			}
		case core.POP:
			switch code[i] >> 16 & 0xff {
			case core.STACKSTR:
				top.Str--
			case core.STACKFLOAT:
				top.Any--
			case core.STACKANY:
				top.Any--
			default:
				top.Int--
			}
		case core.CYCLE:
			lenCalls := len(rt.Calls) - 1
			rt.Calls[lenCalls].Cycle--
			//			fmt.Println(`CYCLE`, rt.Calls[lenCalls].Cycle, rt.SInt[:top.Int], rt.SStr[:top.Str], rt.SAny[:top.Any])
			if rt.Calls[lenCalls].Cycle == 0 {
				return nil, runtimeError(rt, i, ErrCycle)
			}
		case core.JMP:
			i += int64(int16(code[i+1]))
			//top = rt.Calls[len(rt.Calls)-1]
			continue
		case core.JZE:
			top.Int--
			if rt.SInt[top.Int] == 0 {
				i += int64(int16(code[i+1])) + 2
				continue
			}
			i++
		case core.JNZ:
			top.Int--
			if rt.SInt[top.Int] != 0 {
				i += int64(int16(code[i+1])) + 2
				continue
			}
			i++
		case core.INITVARS:
			//			parCount := code[i] >> 16
			i++
			var prevTop Call
			curTop := top
			if rt.ParCount > 0 {
				prevTop = rt.Calls[len(rt.Calls)-1] //top
			}
			varCount := int32(code[i])
			for k := int32(0); k < varCount; k++ {
				i++
				varType := int(code[i])
				if rt.ParCount > k {
					switch varType & 0xff {
					case core.STACKSTR:
						prevTop.Str--
						curTop.Str--
					case core.STACKANY:
						prevTop.Any--
						curTop.Any--
					default:
						prevTop.Int--
						curTop.Int--
					}
				} else {
					switch varType {
					case core.TYPEINT, core.TYPEBOOL:
						rt.SInt[top.Int] = 0
						top.Int++
					case core.TYPECHAR:
						rt.SInt[top.Int] = int64(' ')
						top.Int++
					case core.TYPESTR:
						rt.SStr[top.Str] = ``
						top.Str++
					case core.TYPEARR:
						rt.SAny[top.Any] = core.NewArray()
						top.Any++
					case core.TYPEMAP:
						rt.SAny[top.Any] = core.NewMap()
						top.Any++
					default:
						fmt.Println(`INIT ANY`, varType)
					}
				}
			}
			if rt.ParCount > 0 {
				rt.Calls[len(rt.Calls)-1] = prevTop
				rt.ParCount = 0
			}
			rt.Calls = append(rt.Calls, Call{
				IsFunc: false,
				Cycle:  rt.Owner.Settings.Cycle,
				Offset: int32(i),
				Int:    curTop.Int,
				Float:  curTop.Float,
				Str:    curTop.Str,
				Any:    curTop.Any,
			})
			//			fmt.Println(`INITVARS`, rt.SInt[:top.Int])
		case core.DELVARS:
			curTop := top
			top = rt.Calls[len(rt.Calls)-1]
			rt.Calls = rt.Calls[:len(rt.Calls)-1]
			for j := top.Any; j < curTop.Any; j++ {
				rt.SAny[j] = nil
			}
		case core.RANGE:
			top.Int -= 2
			rt.SAny[top.Any] = &core.Range{From: rt.SInt[top.Int], To: rt.SInt[top.Int+1]}
			top.Any++
		case core.LEN:
			var length int64
			if code[i]>>16 == core.TYPESTR {
				top.Str--
				length = int64(len([]rune(rt.SStr[top.Str])))
			} else {
				top.Any--
				switch v := rt.SAny[top.Any].(type) {
				case *core.Array:
					length = int64(v.Len())
				case *core.Range:
					length = v.To - v.From
					if length < 0 {
						length = -length
					}
					length++
				case *core.Map:
					length = int64(len(v.Keys))
				default:
					fmt.Printf("LENGTH %T\n", v)
				}
			}
			rt.SInt[top.Int] = length
			top.Int++
		case core.FORINC:
			rt.SInt[rt.Calls[int(len(rt.Calls)-1)].Int+int32(code[i]>>16)]++
		case core.RET:
			retType := code[i] >> 16
			k := len(rt.Calls) - 1
			for ; k >= 0; k-- {
				if rt.Calls[k].IsFunc {
					break
				}
			}
			//			fmt.Println(`RET`, k, rt.SInt[:top.Int])

			rt.Calls = rt.Calls[:k+1]
			if len(rt.Calls) == 0 { // return from run function
				switch retType {
				case core.TYPEINT:
					result = rt.SInt[top.Int-1]
				case core.TYPEBOOL:
					if rt.SInt[top.Int-1] == 0 {
						result = false
					} else {
						result = true
					}
				case core.TYPECHAR:
					result = rune(rt.SInt[top.Int-1])
				case core.TYPESTR:
					result = rt.SStr[top.Str-1]
				default:
					result = rt.SAny[top.Any-1]
				}
				break main
			}
			curTop := top
			top = rt.Calls[k]
			rt.Calls = rt.Calls[:k]
			switch retType & 0xff {
			case core.STACKNONE:
			case core.STACKSTR:
				rt.SStr[top.Str] = rt.SStr[curTop.Str-1]
				top.Str++
			case core.STACKANY:
				rt.SAny[top.Any] = rt.SAny[curTop.Any-1]
				top.Any++
			default:
				rt.SInt[top.Int] = rt.SInt[curTop.Int-1]
				top.Int++
			}
			for j := top.Any; j < curTop.Any; j++ {
				rt.SAny[j] = nil
			}
			i = int64(top.Offset)
		case core.END:
			if len(rt.Calls) == 0 {
				break main
			}
			top = rt.Calls[len(rt.Calls)-1]
			rt.Calls = rt.Calls[:len(rt.Calls)-1]
			i = int64(top.Offset)
		case core.INDEX:
			top.Any--
			switch v := rt.SAny[top.Any].(type) {
			case *string:
				index := rt.SInt[top.Int-1]
				runes := []rune(*v)
				if index < 0 || index >= int64(len(runes)) {
					return nil, runtimeError(rt, i, ErrIndexOut)
				}
				rt.SInt[top.Int-1] = int64(runes[index])
			case *interface{}:
				index := rt.SInt[top.Int-1]
				runes := []rune((*v).(string))
				if index < 0 || index >= int64(len(runes)) {
					return nil, runtimeError(rt, i, ErrIndexOut)
				}
				rt.SInt[top.Int-1] = int64(runes[index])
			case *core.Array:
				top.Int--
				index := rt.SInt[top.Int]
				if index < 0 || index >= int64(v.Len()) {
					return nil, runtimeError(rt, i, ErrIndexOut)
				}
				switch code[i] >> 16 {
				case core.TYPEINT:
					rt.SInt[top.Int] = v.Data[index].(int64)
					top.Int++
				case core.TYPESTR:
					rt.SStr[top.Str] = v.Data[index].(string)
					top.Str++
				default:
					rt.SAny[top.Any] = v.Data[index]
					top.Any++
				}
			default:
				fmt.Printf("TYPE=%T %v\r\n", v, v)
			}
		case core.CONSTBYID:
			i++
			v := rt.Owner.Consts[int32(code[i])]
			switch v.Type {
			case core.TYPEINT, core.TYPEBOOL, core.TYPECHAR:
				rt.SInt[top.Int] = v.Value.(int64)
				top.Int++
			case core.TYPESTR:
				rt.SStr[top.Str] = v.Value.(string)
				top.Str++
			}
		case core.CALLBYID:
			rt.ParCount = int32(code[i]) >> 16
			i++
			rt.Calls = append(rt.Calls, Call{
				IsFunc: true,
				Offset: int32(i),
				Int:    top.Int,
				Float:  top.Float,
				Str:    top.Str,
				Any:    top.Any,
			})
			if uint32(len(rt.Calls)) >= rt.Owner.Settings.Depth {
				return nil, runtimeError(rt, i, ErrDepth)
			}
			i = int64(rt.Owner.Exec.Funcs[int32(code[i])])
			continue
			//			rt.Run(int64(rt.Owner.Exec.Funcs[int32(code[i])]))
			//			top = rt.States[len(rt.States)-1]
			//			rt.States = rt.States[:len(rt.States)-1]
		case core.EMBED:
			var vCount int
			embed := stdlib.Embedded[uint16(code[i]>>16)]
			count := len(embed.Params)
			if embed.Variadic {
				i++
				vCount = int(code[i])
				count--
			}
			pars := make([]reflect.Value, count+vCount)
			if vCount > 0 {
				for i := vCount - 1; i >= 0; i-- {
					i++
					switch code[i] & 0xff {
					case core.STACKSTR:
						top.Str--
						pars[count+i] = reflect.ValueOf(rt.SStr[top.Str])
					case core.STACKANY:
						top.Any--
						pars[count+i] = reflect.ValueOf(rt.SAny[top.Any])
					default:
						top.Int--
						pars[count+i] = reflect.ValueOf(rt.SInt[top.Int])
					}
				}
			}
			for i := count - 1; i >= 0; i-- {
				switch embed.Params[i] & 0xff {
				case core.STACKSTR:
					top.Str--
					pars[i] = reflect.ValueOf(rt.SStr[top.Str])
				case core.STACKANY:
					top.Any--
					pars[i] = reflect.ValueOf(rt.SAny[top.Any])
				default:
					top.Int--
					pars[i] = reflect.ValueOf(rt.SInt[top.Int])
				}
			}
			/*			if obj.Runtime {
						pars = append(pars, reflect.ValueOf(rt))
					}*/
			/*			for i := lenStack; i < len(rt.Stack); i++ {
							pars = append(pars, reflect.ValueOf(rt.Stack[i]))
						}
						rt.Stack = rt.Stack[:lenStack]*/
			result := reflect.ValueOf(embed.Func).Call(pars)
			if len(result) > 0 {
				last := result[len(result)-1].Interface()
				if last != nil {
					if _, isError := last.(error); isError {
						return nil, runtimeError(rt, i, result[len(result)-1].Interface().(error))
					}
				}
				switch embed.Return & 0xff {
				case core.STACKNONE:
				case core.STACKSTR:
					rt.SStr[top.Str] = result[0].Interface().(string)
					top.Str++
				case core.STACKANY:
					rt.SAny[top.Any] = result[0].Interface()
					top.Any++
				default:
					rt.SInt[top.Int] = result[0].Interface().(int64)
					top.Int++
				}
			}
		case core.IOTA:
			rt.Owner.Consts[rt.Owner.Exec.Init[0]] = Const{
				Type:  core.TYPEINT,
				Value: int64((int32(code[i]) >> 16) - 1),
			}
		}
		i++
	}
	return
}
