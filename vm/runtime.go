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

type indexObj struct {
	Obj   interface{}
	Index interface{}
	Type  int //core.Bcode
}

type indexInfo struct {
	Objects [32]indexObj
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
	var (
		iInfo  indexInfo
		tmpInt int64
		tmpStr string
		count  int
	)

	top := Call{}
	code := rt.Owner.Exec.Code
	end := int64(len(code))

main:
	for i < end {
		switch code[i] & 0x0fff {
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
			var (
				typeRet int
				index   interface{}
			)
			if code[i+2]&0xffff == core.INDEX {
				count = int(code[i+2] >> 16)
			} else {
				count = 0
			}
			blockOff := rt.Calls[int32(len(rt.Calls)-1-(int(code[i])>>16))]
			i++
			typeVar := int(code[i]) >> 16
			root := int64(int(code[i]) & 0xffff)
			switch typeVar & 0xff {
			case core.STACKINT:
				root += int64(blockOff.Int)
			case core.STACKSTR:
				root += int64(blockOff.Str)
			case core.STACKANY:
				root += int64(blockOff.Any)
			default:
				fmt.Println(`root index`, typeVar)
			}
			if count == 0 {
				switch typeVar & 0xff {
				case core.STACKINT:
					rt.SInt[top.Int] = rt.SInt[root]
					top.Int++
				case core.STACKSTR:
					rt.SStr[top.Str] = rt.SStr[root]
					top.Str++
				case core.STACKANY:
					rt.SAny[top.Any] = rt.SAny[root]
					top.Any++
				default:
					fmt.Printf("GET TYPE %T\n", typeVar)
				}
				i++
				continue
			}
			i++
			var ptr, value interface{}
			var ok bool
			if typeVar&0xff == core.STACKSTR {
				ptr = &rt.SStr[root]
			} else {
				ptr = rt.SAny[root]
			}
			if count == 1 {
				if prange, isrange := ptr.(*core.Range); isrange {
					if prange.From < prange.To {
						rt.SInt[top.Int-1] = prange.From + rt.SInt[top.Int-1]
					} else {
						rt.SInt[top.Int-1] = prange.From - rt.SInt[top.Int-1]
					}
					i += 2
					continue main
				}
			}
			for ind := 0; ind < count; ind++ {
				i++
				typeVar = int(code[i]) >> 16
				typeRet = int(code[i]) & 0x0fff
				if int(code[i])&0x8000 != 0 {
					top.Str--
					index = rt.SStr[top.Str]
				} else {
					top.Int--
					index = rt.SInt[top.Int]
				}
				switch typeVar & 0xff {
				case core.STACKSTR:
					runes := []rune(*ptr.(*string))
					if int(index.(int64)) < 0 || len(runes) <= int(index.(int64)) {
						err = runtimeError(rt, i, ErrIndexOut)
						return
					}
					value = int64(runes[index.(int64)])
				case core.STACKANY:
					value, ok = ptr.(core.Indexer).GetIndex(index)
					if !ok {
						if key, ok := index.(string); ok {
							err = runtimeError(rt, i, ErrMapIndex, key)
						} else {
							err = runtimeError(rt, i, ErrIndexOut)
						}
						return
					}
				default:
					fmt.Printf("INDEX ANY %x\n", typeVar)
				}
				switch typeRet & 0xff {
				case core.STACKSTR:
					rt.SStr[top.Str] = value.(string)
					ptr = &rt.SStr[top.Str]
				case core.STACKANY:
					ptr = value
				}
			}
			switch typeRet & 0xff {
			case core.STACKINT:
				rt.SInt[top.Int] = value.(int64)
				top.Int++
			case core.STACKSTR:
				top.Str++
			case core.STACKANY:
				rt.SAny[top.Any] = ptr
				top.Any++
			default:
				fmt.Printf("GET TYPE %T\n", ptr)
			}
		case core.SETVAR:
			var ptr interface{}
			var err error

			if code[i+2]&0xffff == core.INDEX {
				count = int(code[i+2] >> 16)
			} else {
				count = 0
			}

			blockOff := rt.Calls[int32(len(rt.Calls)-1-(int(code[i])>>16))]
			i++
			//lastObj := 0
			typeVar := int(code[i]) >> 16
			typeRet := typeVar
			obj := &iInfo.Objects[0]
			obj.Type = typeVar
			root := int64(int(code[i]) & 0xffff)
			switch typeVar & 0xff {
			case core.STACKINT:
				root += int64(blockOff.Int)
				ptr = &rt.SInt[root]
			case core.STACKSTR:
				root += int64(blockOff.Str)
				ptr = &rt.SStr[root]
			case core.STACKANY:
				root += int64(blockOff.Any)
				ptr = rt.SAny[root]
			default:
				fmt.Println(`root index`, typeVar)
			}
			obj.Index = root
			//fmt.Println(`ROOT`, iInfo.Objects[:iInfo.Count+1], ptr, rt.SAny[:top.Any])
			if count > 0 {
				i++
				for ind := 0; ind < count; ind++ {
					i++
					typeVar = int(code[i]) >> 16
					typeRet = int(code[i]) & 0x0fff
					//					iInfo.Count++
					obj = &iInfo.Objects[ind+1]
					//				fmt.Printf("IND %d %x %x\n", ind, typeVar, typeRet)
					if int(code[i])&0x8000 != 0 {
						top.Str--
						obj.Index = rt.SStr[top.Str]
					} else {
						top.Int--
						obj.Index = rt.SInt[top.Int]
					}
					obj.Obj = ptr
					obj.Type = typeRet
					//		fmt.Println(`OBJ`, iInfo.Objects[:iInfo.Count+1])
					switch typeVar & 0xff {
					case core.STACKSTR:
						runes := []rune(*ptr.(*string))
						if int(obj.Index.(int64)) < 0 || len(runes) <= int(obj.Index.(int64)) {
							return nil, runtimeError(rt, i, ErrIndexOut)
						}
						tmpInt = int64(runes[obj.Index.(int64)])
						ptr = &tmpInt
					case core.STACKANY:
						var (
							ok    bool
							value interface{}
						)
						value, ok = ptr.(core.Indexer).GetIndex(obj.Index)
						if !ok {
							if key, ok := obj.Index.(string); ok {
								value = newValue(typeRet)
								if !ptr.(core.Indexer).SetIndex(key, value) {
									return nil, runtimeError(rt, i, ErrIndexOut)
								}
							} else {
								return nil, runtimeError(rt, i, ErrIndexOut)
							}
						}
						switch typeRet & 0xff {
						case core.STACKINT:
							tmpInt = value.(int64)
							ptr = &tmpInt
						case core.STACKSTR:
							tmpStr = value.(string)
							ptr = &tmpStr
						default:
							ptr = value
						}
					default:
						fmt.Printf("INDEX ANY %x %x\n", typeVar, typeRet)
					}
				}
			}
			i++
			assign := code[i] & 0xffff
			rightType := code[i] >> 16
			if count == 0 && (assign == core.ASSIGN || assign == core.ASSIGNPTR) &&
				core.Bcode(typeVar) == rightType {
				switch rightType & 0xff {
				case core.STACKINT:
					rt.SInt[root] = rt.SInt[top.Int-1]
				case core.STACKSTR:
					rt.SStr[root] = rt.SStr[top.Str-1]
				default:
					if assign == core.ASSIGN {
						core.CopyVar(&rt.SAny[root], rt.SAny[top.Any-1])
					} else {
						rt.SAny[root] = rt.SAny[top.Any-1]
					}
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
			obj = &iInfo.Objects[count]
			if assign == core.ASSIGN || assign == core.ASSIGNPTR &&
				core.Bcode(obj.Type) == rightType {
				switch v := ptr.(type) {
				case *int64:
					*v = iValue.(int64)
				case *string:
					*v = iValue.(string)
				default:
					if assign == core.ASSIGN {
						core.CopyVar(&ptr, iValue)
						iValue = ptr
					}
				}
			} else {
				switch v := ptr.(type) {
				case *int64:
					iValue, err = stdlib.EmbedInt[assign-core.ASSIGN](
						v, iValue.(int64))
				case *string:
					iValue, err = stdlib.EmbedStr[assign-core.ASSIGN](
						v, iValue)
				case *core.Array, *core.Map:
					iValue, err = stdlib.EmbedAny[assign-core.ASSIGN](
						ptr, iValue)
				default:
					fmt.Println(`Embed Assign`, rightType)
				}
				if err != nil {
					return nil, runtimeError(rt, i, err)
				}
			}
			//			typeVar = (int(code[i]) >> 16) & 0xff
			//			fmt.Println(`OBJ`, iInfo.Objects[:iInfo.Count+1], iValue)
			if count > 0 || typeVar&0xff == core.STACKANY {
				if count > 0 && iInfo.Objects[count-1].Type == core.TYPESTR {
					var dest string
					obj = &iInfo.Objects[count-1]
					if obj.Obj == nil {
						dest = rt.SStr[obj.Index.(int64)]
					} else {
						ret, _ := obj.Obj.(core.Indexer).GetIndex(obj.Index)
						dest = ret.(string)
					}
					runes := []rune(dest)
					runes[iInfo.Objects[count].Index.(int64)] = rune(iValue.(int64))
					//					iValue = string(runes)
					dest = string(runes)
					if obj.Obj == nil {
						rt.SStr[obj.Index.(int64)] = dest
					} else {
						obj.Obj.(core.Indexer).SetIndex(obj.Index, dest)
					}
				} else {
					obj = &iInfo.Objects[count]
					if obj.Obj == nil {
						/*switch obj.Type & 0xff {
						case core.STACKINT:
							rt.SInt[obj.Index.(int64)] = iValue.(int64)
						case core.STACKSTR:
							rt.SStr[obj.Index.(int64)] = iValue.(string)
						case core.STACKANY:*/
						rt.SAny[obj.Index.(int64)] = iValue
						//						}
					} else {
						switch obj.Type & 0xff {
						case core.STACKINT:
							iValue = tmpInt
						case core.STACKSTR:
							iValue = tmpStr
						}
						obj.Obj.(core.Indexer).SetIndex(obj.Index, iValue)
					}
				}
			}
			switch iInfo.Objects[count].Type & 0xff {
			case core.STACKINT:
				rt.SInt[top.Int] = iValue.(int64)
				top.Int++
			case core.STACKSTR:
				rt.SStr[top.Str] = iValue.(string)
				top.Str++
			case core.STACKANY:
				rt.SAny[top.Any] = iValue
				top.Any++
			default:
				fmt.Printf("SET TYPE %T\n", ptr)
			}
			//			fmt.Println(`SETVAR`, rt.SInt[:top.Int], rt.SAny[:top.Any])
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
				length = int64(rt.SAny[top.Any].(core.Indexer).Len())
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
			/*		case core.INDEX:
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
					}*/
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
