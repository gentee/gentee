// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/gentee/gentee/core"
	stdlib "github.com/gentee/gentee/stdlibvm"
)

const (
	// SleepStep is a tick in sleep
	SleepStep = int64(100)
)

type indexObj struct {
	Obj   interface{}
	Index interface{}
	Type  int //core.Bcode
}

type indexInfo struct {
	Objects [32]indexObj
}

func (rt *Runtime) Run(i int64) (result interface{}, err error) {
	var (
		iInfo    indexInfo
		tmpInt   int64
		tmpStr   string
		tmpFloat float64
		count    int
	)

	top := Call{}
	code := rt.Owner.Exec.Code
	end := int64(len(code))

	errHandle := func(pos int64, errPar interface{}, pars ...interface{}) {
		k := len(rt.Calls) - 1
		for ; k > 0; k-- {
			if rt.Calls[k].Flags&core.BlTry != 0 {
				break
			}
		}
		//fmt.Println(`errHandle`, pos, rt.Owner.Exec.Pos)
		err = runtimeError(rt, pos, errPar, pars...)
		if k <= 0 {
			i = end + 1
			return
		}
		top = rt.Calls[k]
		rt.Calls = rt.Calls[:k]
		i = int64(top.Offset + top.Try)
		rt.SAny[top.Any] = err
		rt.Calls[len(rt.Calls)-1].Any++
		top.Any++
		rt.ParCount = 1
	}

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
		case core.PUSHFLOAT:
			i += 2
			rt.SFloat[top.Float] = math.Float64frombits(uint64(code[i-1])<<32 |
				uint64(code[i])&0xffffffff)
			top.Float++
		case core.PUSHSTR:
			rt.SStr[top.Str] = rt.Owner.Exec.Strings[(code[i])>>16]
			top.Str++
		case core.PUSHFUNC:
			i++
			rt.SAny[top.Any] = &Fn{Func: int32(code[i])}
			top.Any++
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
				errHandle(i, ErrDivZero)
				continue
			}
			rt.SInt[top.Int-1] /= rt.SInt[top.Int]
		case core.MOD:
			top.Int--
			if rt.SInt[top.Int] == 0 {
				errHandle(i, ErrDivZero)
				continue
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
				errHandle(i, ErrShift)
				continue
			}
			rt.SInt[top.Int-1] <<= uint32(rt.SInt[top.Int])
		case core.RSHIFT:
			top.Int--
			if rt.SInt[top.Int] < 0 {
				errHandle(i, ErrShift)
				continue
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
		case core.ADDFLOAT:
			top.Float--
			rt.SFloat[top.Float-1] += rt.SFloat[top.Float]
		case core.SUBFLOAT:
			top.Float--
			rt.SFloat[top.Float-1] -= rt.SFloat[top.Float]
		case core.MULFLOAT:
			top.Float--
			rt.SFloat[top.Float-1] *= rt.SFloat[top.Float]
		case core.DIVFLOAT:
			top.Float--
			if rt.SFloat[top.Float] == 0.0 {
				errHandle(i, ErrDivZero)
				continue
			}
			rt.SFloat[top.Float-1] /= rt.SFloat[top.Float]
		case core.SIGNFLOAT:
			rt.SFloat[top.Float-1] = -rt.SFloat[top.Float-1]
		case core.EQFLOAT:
			top.Float -= 2
			if rt.SFloat[top.Float] == rt.SFloat[top.Float+1] {
				rt.SInt[top.Int] = 1
			} else {
				rt.SInt[top.Int] = 0
			}
			top.Int++
		case core.LTFLOAT:
			top.Float -= 2
			if rt.SFloat[top.Float] < rt.SFloat[top.Float+1] {
				rt.SInt[top.Int] = 1
			} else {
				rt.SInt[top.Int] = 0
			}
			top.Int++
		case core.GTFLOAT:
			top.Float -= 2
			if rt.SFloat[top.Float] > rt.SFloat[top.Float+1] {
				rt.SInt[top.Int] = 1
			} else {
				rt.SInt[top.Int] = 0
			}
			top.Int++
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
			switch typeVar & 0xf {
			case core.STACKINT:
				root += int64(blockOff.Int)
			case core.STACKFLOAT:
				root += int64(blockOff.Float)
			case core.STACKSTR:
				root += int64(blockOff.Str)
			case core.STACKANY:
				root += int64(blockOff.Any)
			default:
				fmt.Println(`root index`, typeVar)
			}
			if count == 0 {
				switch typeVar & 0xf {
				case core.STACKINT:
					rt.SInt[top.Int] = rt.SInt[root]
					top.Int++
				case core.STACKFLOAT:
					rt.SFloat[top.Float] = rt.SFloat[root]
					top.Float++
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
			if typeVar&0xf == core.STACKSTR {
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
				typeRet = int(code[i]) & 0x7fff
				if int(code[i])&0x8000 != 0 {
					top.Str--
					index = rt.SStr[top.Str]
				} else {
					top.Int--
					index = rt.SInt[top.Int]
				}
				switch typeVar & 0xf {
				case core.STACKSTR:
					runes := []rune(*ptr.(*string))
					if int(index.(int64)) < 0 || len(runes) <= int(index.(int64)) {
						errHandle(i, ErrIndexOut)
						continue main
						//						err = runtimeError(rt, i, ErrIndexOut)
						//						return
					}
					value = int64(runes[index.(int64)])
				case core.STACKANY:
					value, ok = ptr.(core.Indexer).GetIndex(index)
					if !ok {
						if key, ok := index.(string); ok {
							errHandle(i, ErrMapIndex, key)
							//err = runtimeError(rt, i, ErrMapIndex, key)
						} else {
							errHandle(i, ErrIndexOut)
							//err = runtimeError(rt, i, ErrIndexOut)
						}
						continue main
					}
					if value == nil {
						errHandle(i, ErrUndefined)
						continue main
						//						return nil, runtimeError(rt, i, ErrUndefined)
					}
				default:
					fmt.Printf("INDEX ANY %x\n", typeVar)
				}
				switch typeRet & 0xf {
				case core.STACKSTR:
					rt.SStr[top.Str] = value.(string)
					ptr = &rt.SStr[top.Str]
				case core.STACKANY:
					ptr = value
				}
			}
			switch typeRet & 0xf {
			case core.STACKINT:
				rt.SInt[top.Int] = value.(int64)
				top.Int++
			case core.STACKFLOAT:
				rt.SFloat[top.Float] = value.(float64)
				top.Float++
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
			switch typeVar & 0xf {
			case core.STACKINT:
				root += int64(blockOff.Int)
				ptr = &rt.SInt[root]
			case core.STACKFLOAT:
				root += int64(blockOff.Float)
				ptr = &rt.SFloat[root]
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
					typeRet = int(code[i]) & 0x7fff
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
					switch typeVar & 0xf {
					case core.STACKSTR:
						runes := []rune(*ptr.(*string))
						if int(obj.Index.(int64)) < 0 || len(runes) <= int(obj.Index.(int64)) {
							errHandle(i, ErrIndexOut)
							continue main
							//							return nil, runtimeError(rt, i, ErrIndexOut)
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
								value = newValue(rt, typeRet)
								if ptr.(core.Indexer).SetIndex(key, value) != 0 {
									errHandle(i, ErrIndexOut)
									continue main
									//return nil, runtimeError(rt, i, ErrIndexOut)
								}
							} else {
								errHandle(i, ErrIndexOut)
								continue main
								//return nil, runtimeError(rt, i, ErrIndexOut)
							}
						}
						switch typeRet & 0xf {
						case core.STACKINT:
							tmpInt = value.(int64)
							ptr = &tmpInt
						case core.STACKFLOAT:
							tmpFloat = value.(float64)
							ptr = &tmpFloat
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
			//			fmt.Printf("Assign %d %d %d %d %x %x\n", count, assign, core.ASSIGN, core.ASSIGNPTR,
			//				core.Bcode(typeVar), rightType)
			if count == 0 && (assign == core.ASSIGN || assign == core.ASSIGNPTR) &&
				core.Bcode(typeVar) == rightType {
				switch rightType & 0xf {
				case core.STACKINT:
					rt.SInt[root] = rt.SInt[top.Int-1]
				case core.STACKFLOAT:
					rt.SFloat[root] = rt.SFloat[top.Float-1]
				case core.STACKSTR:
					rt.SStr[root] = rt.SStr[top.Str-1]
				default:
					if rt.SAny[root] == rt.SAny[top.Any-1] {
						errHandle(i, ErrAssignment)
						continue main
						//return nil, runtimeError(rt, i, ErrAssignment)
					}
					if assign == core.ASSIGN {
						CopyVar(rt, &rt.SAny[root], rt.SAny[top.Any-1])
					} else {
						rt.SAny[root] = rt.SAny[top.Any-1]
					}
				}
				i++
				continue
			}
			var iValue interface{}
			switch rightType & 0xf {
			case core.STACKINT, core.STACKNONE: // STACKNONE is for inc dec
				top.Int--
				iValue = rt.SInt[top.Int]
			case core.STACKFLOAT:
				top.Float--
				iValue = rt.SFloat[top.Float]
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
				case *float64:
					*v = iValue.(float64)
				case *string:
					*v = iValue.(string)
				default:
					if ptr == iValue {
						errHandle(i, ErrAssignment)
						continue main
						//return nil, runtimeError(rt, i, ErrAssignment)
					}
					if assign == core.ASSIGN {
						CopyVar(rt, &ptr, iValue)
						iValue = ptr
					}
				}
			} else {
				switch v := ptr.(type) {
				case *int64:
					iValue, err = stdlib.EmbedInt[assign-core.ASSIGN](
						v, iValue.(int64))
				case *float64:
					iValue, err = stdlib.EmbedFloat[assign-core.ASSIGN](
						v, iValue.(float64))
				case *string:
					iValue, err = stdlib.EmbedStr[assign-core.ASSIGN](
						v, iValue)
				default:
					iValue, err = stdlib.EmbedAny[assign-core.ASSIGN](
						ptr, iValue)
					//				default:
					//					fmt.Println(`Embed Assign`, rightType)
				}
				if err != nil {
					errHandle(i, err)
					continue main
					//return nil, runtimeError(rt, i, err)
				}
			}
			//			typeVar = (int(code[i]) >> 16) & 0xf
			//			fmt.Println(`OBJ`, iInfo.Objects[:iInfo.Count+1], iValue)
			if count > 0 || typeVar&0xf == core.STACKANY {
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
						if errID := obj.Obj.(core.Indexer).SetIndex(obj.Index, dest); errID != 0 {
							errHandle(i, errID)
							continue main
							//return nil, runtimeError(rt, i, errID)
						}
					}
				} else {
					obj = &iInfo.Objects[count]
					if obj.Obj == nil {
						/*switch obj.Type & 0xf {
						case core.STACKINT:
							rt.SInt[obj.Index.(int64)] = iValue.(int64)
						case core.STACKSTR:
							rt.SStr[obj.Index.(int64)] = iValue.(string)
						case core.STACKANY:*/
						rt.SAny[obj.Index.(int64)] = iValue
						//						}
					} else {
						switch obj.Type & 0xf {
						case core.STACKINT:
							iValue = tmpInt
						case core.STACKFLOAT:
							iValue = tmpFloat
						case core.STACKSTR:
							iValue = tmpStr
						}
						if obj.Obj == iValue {
							errHandle(i, ErrAssignment)
							continue main
							//return nil, runtimeError(rt, i, ErrAssignment)
						}
						if errID := obj.Obj.(core.Indexer).SetIndex(obj.Index, iValue); errID != 0 {
							errHandle(i, errID)
							continue main
							//return nil, runtimeError(rt, i, errID)
						}
					}
				}
			}
			switch iInfo.Objects[count].Type & 0xf {
			case core.STACKINT:
				rt.SInt[top.Int] = iValue.(int64)
				top.Int++
			case core.STACKFLOAT:
				rt.SFloat[top.Float] = iValue.(float64)
				top.Float++
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
			switch code[i] >> 16 & 0xf {
			case core.STACKFLOAT:
				rt.SFloat[top.Float] = rt.SFloat[top.Float-1]
				top.Float++
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
			switch code[i] >> 16 & 0xf {
			case core.STACKFLOAT:
				top.Float--
			case core.STACKSTR:
				top.Str--
				//			case core.STACKFLOAT:
				//				top.Any--
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
				errHandle(i, ErrCycle)
				continue main
				//return nil, runtimeError(rt, i, ErrCycle)
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
		case core.JEQ:
			switch (code[i] >> 16) & 0xf {
			case core.STACKINT:
				top.Int--
				if rt.SInt[top.Int] == rt.SInt[top.Int-1] {
					i += int64(code[i+1])
					continue
				}
			case core.STACKSTR:
				top.Str--
				if rt.SStr[top.Str] == rt.SStr[top.Str-1] {
					i += int64(code[i+1])
					continue
				}
			case core.STACKFLOAT:
				top.Float--
				if rt.SFloat[top.Float] == rt.SFloat[top.Float-1] {
					i += int64(code[i+1])
					continue
				}
			}
			i++
		case core.JMPOPT:
			idopt := int32(code[i] >> 16)
			for k := len(rt.Calls) - 1; k >= 0; k-- {
				if rt.Calls[k].IsFunc {
					if rt.Calls[k].Optional != nil {
						for _, optVar := range *rt.Calls[k].Optional {
							if idopt == optVar.Var {
								i += int64(code[i+1])
								continue main
							}
						}
					}
					break
				}
			}
			i++
		case core.INITVARS:
			var (
				breakJmp, continueJmp, tryJmp, recoverJmp, retryJmp int32
			)
			flags := int16(code[i] >> 16)
			pos := i
			//			parCount := code[i] >> 16
			if flags&core.BlBreak != 0 {
				i++
				breakJmp = int32(code[i])
			}
			if flags&core.BlContinue != 0 {
				i++
				continueJmp = int32(code[i])
			}
			if flags&core.BlTry != 0 {
				i++
				tryJmp = int32(code[i])
			}
			if flags&core.BlRecover != 0 {
				i++
				recoverJmp = int32(code[i])
			}
			if flags&core.BlRetry != 0 {
				i++
				retryJmp = int32(code[i])
			}
			var prevTop Call
			curTop := top
			if rt.ParCount > 0 {
				prevTop = rt.Calls[len(rt.Calls)-1] //top
			}
			//fmt.Println(`INITVARS`, rt.Optional, rt.ParCount, rt.SInt[:top.Int], rt.SAny[:top.Any])
			if flags&core.BlVars != 0 {
				var optional *[]OptValue
				if len(rt.Calls) > 0 {
					optional = rt.Calls[len(rt.Calls)-1].Optional
				} else {
					optional = rt.Optional
					rt.Optional = nil
				}
				i++
				varCount := int32(code[i] & 0xffff)
				for k := int32(0); k < varCount; k++ {
					i++
					varType := int(code[i])
					//					fmt.Printf("VAR %d %d %x %x \n", i, k, varCount, varType)
					if rt.ParCount > k {
						switch varType & 0xf {
						case core.STACKFLOAT:
							prevTop.Float--
							curTop.Float--
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
						value := newValue(rt, varType)
						if optional != nil {
							for _, optVar := range *optional {
								if k == optVar.Var {
									value = optVar.Value
								}
							}
						}
						switch varType & 0xf {
						case core.STACKINT:
							rt.SInt[top.Int] = value.(int64)
							top.Int++
						case core.STACKSTR:
							rt.SStr[top.Str] = value.(string)
							top.Str++
						case core.STACKFLOAT:
							rt.SFloat[top.Float] = value.(float64)
							top.Float++
						case core.STACKANY:
							rt.SAny[top.Any] = value
							top.Any++
						default:
							fmt.Println(`INIT DEFAULT`, varType)
						}
					}
				}
			}
			if rt.ParCount > 0 {
				rt.Calls[len(rt.Calls)-1] = prevTop
				rt.ParCount = 0
			}
			rt.Calls = append(rt.Calls, Call{
				IsFunc:   false,
				Cycle:    rt.Owner.Settings.Cycle,
				Offset:   int32(pos),
				Int:      curTop.Int,
				Float:    curTop.Float,
				Str:      curTop.Str,
				Any:      curTop.Any,
				Flags:    flags,
				Break:    breakJmp,
				Continue: continueJmp,
				Try:      tryJmp,
				Recover:  recoverJmp,
				Retry:    retryJmp,
			})

			//			fmt.Println(`INIT OK`, rt.SInt[:top.Int], rt.SAny[:top.Any])
			//			fmt.Println(`INITVARS`, rt.Calls)
		case core.DELVARS:
			curTop := top
			top = rt.Calls[len(rt.Calls)-1]
			rt.Calls = rt.Calls[:len(rt.Calls)-1]
			for j := top.Any; j < curTop.Any; j++ {
				rt.SAny[j] = nil
			}
			// fmt.Println(`DELVARS`, rt.Calls)
		case core.OPTPARS:
			count := int64(code[i] >> 16)
			optional := make([]OptValue, count)
			//			rt.Optional = make([]OptValue, count)
			//			rt.IsOptional = true
			for j := count; j > 0; j-- {
				var value interface{}
				itype := int(code[i+j] >> 16)
				switch itype & 0xf {
				case core.STACKINT:
					top.Int--
					value = rt.SInt[top.Int]
				case core.STACKSTR:
					top.Str--
					value = rt.SStr[top.Str]
				case core.STACKFLOAT:
					top.Float--
					value = rt.SFloat[top.Float]
				case core.STACKANY:
					top.Any--
					value = rt.SAny[top.Any]
				default:
					fmt.Println(`OOOPS`, itype)
				}
				optional[j-1] = OptValue{
					Var:   int32(code[i+j] & 0xffff),
					Type:  itype,
					Value: value,
				}
			}
			rt.Optional = &optional
			i += count
		case core.INITOBJ:
			count = int(code[i]) >> 16
			i++
			typeRet := int(code[i]) >> 16
			typeVar := int(code[i]) & 0xffff
			//			fmt.Printf("INITOBJ %d %x %x\n", count, typeVar, typeRet)
			switch typeVar {
			case core.TYPEARR:
				parr := core.NewArray()
				parr.Data = make([]interface{}, count)
				for j := 0; j < count; j++ {
					switch typeRet & 0xf {
					case core.STACKINT:
						top.Int--
						parr.Data[count-j-1] = rt.SInt[top.Int]
					case core.STACKFLOAT:
						top.Float--
						parr.Data[count-j-1] = rt.SFloat[top.Float]
					case core.STACKSTR:
						top.Str--
						parr.Data[count-j-1] = rt.SStr[top.Str]
					case core.STACKANY:
						top.Any--
						/*						var ptr interface{}
												CopyVar(rt, &ptr, rt.SAny[top.Any])
												parr.Data[count-j-1] = ptr*/
						parr.Data[count-j-1] = rt.SAny[top.Any] //ptr
					}
				}
				rt.SAny[top.Any] = parr
			case core.TYPEMAP:
				pmap := core.NewMap()
				pmap.Keys = make([]string, count)
				for j := 0; j < count; j++ {
					var value interface{}
					switch typeRet & 0xf {
					case core.STACKINT:
						top.Int--
						value = rt.SInt[top.Int]
					case core.STACKFLOAT:
						top.Float--
						value = rt.SFloat[top.Float]
					case core.STACKSTR:
						top.Str--
						value = rt.SStr[top.Str]
					case core.STACKANY:
						top.Any--
						//CopyVar(rt, &value, rt.SAny[top.Any])
						value = rt.SAny[top.Any]
					}
					top.Str--
					key := rt.SStr[top.Str]
					pmap.Data[key] = value
					pmap.Keys[count-j-1] = key
				}
				rt.SAny[top.Any] = pmap
			case core.TYPEBUF:
				pbuf := core.NewBuffer()
				tmp := make([][]byte, count)
				for j := 0; j < count; j++ {
					top.Int--
					switch rt.SInt[top.Int] {
					case core.TYPEINT:
						top.Int--
						pos := rt.SInt[top.Int]
						top.Int--
						ind := rt.SInt[top.Int]
						if uint64(ind) > 255 {
							errHandle(pos, ErrByteOut)
							continue main
							//return nil, runtimeError(rt, pos, ErrByteOut)
						}
						tmp[j] = []byte{byte(ind)}
					case core.TYPECHAR:
						top.Int--
						tmp[j] = append(tmp[j], []byte(string([]rune{rune(rt.SInt[top.Int])}))...)
					case core.TYPESTR:
						top.Str--
						tmp[j] = append(tmp[j], []byte(rt.SStr[top.Str])...)
					case core.TYPEBUF:
						top.Any--
						tmp[j] = append(tmp[j], rt.SAny[top.Any].(*core.Buffer).Data...)
					default:
						errHandle(i, ErrRuntime, `init buf`)
						continue main
						//return nil, runtimeError(rt, i, ErrRuntime, `init buf`)
					}
				}
				for j := len(tmp) - 1; j >= 0; j-- {
					pbuf.Data = append(pbuf.Data, tmp[j]...)
				}
				rt.SAny[top.Any] = pbuf
			default:
				if typeVar >= core.TYPESTRUCT {
					pstruct := NewStruct(rt,
						&rt.Owner.Exec.Structs[(typeVar-core.TYPESTRUCT)>>8])
					for j := 0; j < count; j++ {
						top.Int--
						ind := rt.SInt[top.Int]
						var value interface{}
						switch pstruct.Type.Fields[ind] & 0xf {
						case core.STACKINT:
							top.Int--
							value = rt.SInt[top.Int]
						case core.STACKFLOAT:
							top.Float--
							value = rt.SFloat[top.Float]
						case core.STACKSTR:
							top.Str--
							value = rt.SStr[top.Str]
						case core.STACKANY:
							top.Any--
							value = rt.SAny[top.Any]
						}
						pstruct.Values[ind] = value
					}
					rt.SAny[top.Any] = pstruct
				} else {
					fmt.Println(`NEW INITOBJ`, typeVar)
				}
			}
			top.Any++

		case core.RANGE:
			top.Int -= 2
			rt.SAny[top.Any] = &core.Range{From: rt.SInt[top.Int], To: rt.SInt[top.Int+1]}
			top.Any++
		case core.ARRAY:
			count := int(code[i] >> 16)
			ret := core.NewArray()
			ret.Data = make([]interface{}, 0, count)
			for j := 0; j < count; j++ {
				i++
				itype := int(code[i] & 0xffff)
				if int(code[i]>>16) == 1 {
					top.Any--
					for k := len(rt.SAny[top.Any].(*core.Array).Data) - 1; k >= 0; k-- {
						ret.Data = append(ret.Data, rt.SAny[top.Any].(*core.Array).Data[k])
					}
					//					ret.Data = append(ret.Data, rt.SAny[top.Any].(*core.Array).Data...)
				} else {
					var value interface{}
					switch itype & 0xf {
					case core.STACKINT:
						top.Int--
						value = rt.SInt[top.Int]
					case core.STACKSTR:
						top.Str--
						value = rt.SStr[top.Str]
					case core.STACKFLOAT:
						top.Float--
						value = rt.SFloat[top.Float]
					case core.STACKANY:
						top.Any--
						value = rt.SAny[top.Any]
					}
					ret.Data = append(ret.Data, value)
				}
			}
			rt.SAny[top.Any] = stdlib.ReverseºArr(ret)
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
		case core.BREAK:
			k := len(rt.Calls) - 1
			for ; k >= 0; k-- {
				if rt.Calls[k].Flags&core.BlBreak != 0 {
					break
				}
			}
			for j := rt.Calls[k].Any; j < top.Any; j++ {
				rt.SAny[j] = nil
			}
			i = int64(rt.Calls[k].Offset + rt.Calls[k].Break)
			top = rt.Calls[k]
			rt.Calls = rt.Calls[:k]
			continue
		case core.CONTINUE:
			k := len(rt.Calls) - 1
			for ; k >= 0; k-- {
				if rt.Calls[k].Flags&core.BlContinue != 0 {
					break
				}
			}
			for j := rt.Calls[k].Any; j < top.Any; j++ {
				rt.SAny[j] = nil
			}
			i = int64(rt.Calls[k].Offset + rt.Calls[k].Continue)
			top = rt.Calls[k]
			rt.Calls = rt.Calls[:k]
			//			fmt.Println(`RET`, k, rt.SInt[:top.Int])
			continue
		case core.RECOVER, core.RETRY:
			isRecover := code[i]&0xffff == core.RECOVER
			err = nil
			k := len(rt.Calls) - 1
			for ; k >= 0; k-- {
				if rt.Calls[k].Flags&core.BlRecover != 0 { // BlRecover & BlRetry in catch
					break
				}
			}
			for j := rt.Calls[k].Any; j < top.Any; j++ {
				rt.SAny[j] = nil
			}
			i = int64(rt.Calls[k].Offset)
			if isRecover {
				i += int64(rt.Calls[k].Recover)
			} else {
				i += int64(rt.Calls[k].Retry)
			}
			top = rt.Calls[k]
			rt.Calls = rt.Calls[:k]
			continue
		case core.RET:
			retType := code[i] >> 16
			k := len(rt.Calls) - 1
			for ; k >= 0; k-- {
				if rt.Calls[k].IsFunc || rt.Calls[k].IsLocal {
					break
				}
			}
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
				case core.TYPEFLOAT:
					result = rt.SFloat[top.Float-1]
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
			switch retType & 0xf {
			case core.STACKNONE:
			case core.STACKFLOAT:
				rt.SFloat[top.Float] = rt.SFloat[curTop.Float-1]
				top.Float++
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
		case core.CONSTBYID:
			i++
			v := rt.Owner.Consts[int32(code[i])]
			switch v.Type {
			case core.TYPEINT, core.TYPEBOOL, core.TYPECHAR:
				rt.SInt[top.Int] = v.Value.(int64)
				top.Int++
			case core.TYPEFLOAT:
				rt.SFloat[top.Float] = v.Value.(float64)
				top.Float++
			case core.TYPESTR:
				rt.SStr[top.Str] = v.Value.(string)
				top.Str++
			}
		case core.CALLBYID:
			rt.ParCount = int32(code[i]) >> 16
			i++
			id := int32(code[i])
			if id == 0 {
				top.Any--
				id = rt.SAny[top.Any].(*Fn).Func
				if id == 0 {
					errHandle(i, ErrFnEmpty)
					continue main
					//return nil, runtimeError(rt, i, ErrFnEmpty)
				}
			}
			rt.Calls = append(rt.Calls, Call{
				IsFunc:   true,
				Offset:   int32(i),
				Int:      top.Int,
				Float:    top.Float,
				Str:      top.Str,
				Any:      top.Any,
				Optional: rt.Optional,
			})
			rt.Optional = nil
			if uint32(len(rt.Calls)) >= rt.Owner.Settings.Depth {
				errHandle(i, ErrDepth)
				continue main
				//return nil, runtimeError(rt, i, ErrDepth)
			}
			i = int64(rt.Owner.Exec.Funcs[id])
			continue
		case core.GOBYID:
			var pars []int32
			rt.ParCount = int32(code[i]) >> 16
			i++
			id := int32(code[i])
			if rt.ParCount > 0 {
				pars = make([]int32, rt.ParCount)
				for k := int32(0); k < rt.ParCount; k++ {
					i++
					pars[k] = int32(code[i] & 0xffff)
				}
				rt.ParCount = 0
			}
			threadID := rt.GoThread(int64(rt.Owner.Exec.Funcs[id]), pars, &top)
			rt.SInt[top.Int] = threadID
			top.Int++
		case core.EMBED:
			var (
				vCount int
				embed  core.Embed
			)
			idEmbed := uint16(code[i] >> 16)
			if idEmbed < 1000 {
				embed = stdlib.Embedded[idEmbed]
			} else {
				embed = Embedded[idEmbed-1000]
			}
			count := len(embed.Params)
			if embed.Variadic {
				i++
				vCount = int(code[i])
				//				count--
			}
			//			fmt.Println(`CALL`, embed.Variadic, count, vCount)
			pars := make([]reflect.Value, count+vCount)
			if vCount > 0 {
				for j := vCount - 1; j >= 0; j-- {
					switch code[i+int64(j)+1] & 0xf {
					case core.STACKFLOAT:
						top.Float--
						pars[count+j] = reflect.ValueOf(rt.SFloat[top.Float])
					case core.STACKSTR:
						top.Str--
						pars[count+j] = reflect.ValueOf(rt.SStr[top.Str])
					case core.STACKANY:
						top.Any--
						pars[count+j] = reflect.ValueOf(rt.SAny[top.Any])
					default:
						top.Int--
						pars[count+j] = reflect.ValueOf(rt.SInt[top.Int])
					}
				}
				i += int64(vCount)
			}
			for i := count - 1; i >= 0; i-- {
				switch embed.Params[i] & 0xf {
				case core.STACKFLOAT:
					top.Float--
					pars[i] = reflect.ValueOf(rt.SFloat[top.Float])
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
			if embed.Runtime {
				pars = append([]reflect.Value{reflect.ValueOf(rt)}, pars...)
			}
			result := reflect.ValueOf(embed.Func).Call(pars)
			if len(result) > 0 {
				last := result[len(result)-1].Interface()
				if last != nil {
					if _, isError := last.(error); isError {
						errHandle(i, result[len(result)-1].Interface().(error))
						continue
					}
				}
				switch embed.Return & 0xf {
				case core.STACKNONE:
				case core.STACKFLOAT:
					rt.SFloat[top.Float] = result[0].Interface().(float64)
					top.Float++
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
		case core.LOCAL:
			rt.ParCount = int32(code[i]) >> 16
			i++
			shift := int32(code[i])
			rt.Calls = append(rt.Calls, Call{
				IsLocal: true,
				Offset:  int32(i),
				Int:     top.Int,
				Float:   top.Float,
				Str:     top.Str,
				Any:     top.Any,
			})
			if uint32(len(rt.Calls)) >= rt.Owner.Settings.Depth {
				errHandle(i, ErrDepth)
				continue main
				//return nil, runtimeError(rt, i, ErrDepth)
			}
			i += int64(shift)
			continue
		case core.IOTA:
			rt.Owner.Consts[rt.Owner.Exec.Init[0]] = Const{
				Type:  core.TYPEINT,
				Value: int64((int32(code[i]) >> 16) - 1),
			}
		}
		i++
		/*		if i&0x8 != 0x8 {
				continue
			}*/
		step := SleepStep
		check := len(rt.Owner.Runtimes) > 1
		for check || rt.Thread.Status == ThPaused || rt.Thread.Status == ThWait ||
			rt.Thread.Sleep > 0 {
			var x int
			if rt.ThreadID == 0 {
				select {
				case err = <-rt.Owner.ChError:
					return nil, err
				default:
				}
			} else {
				select {
				case x = <-rt.Thread.Chan:
					switch x {
					case ThCmdResume, ThCmdContinue:
						rt.setStatus(ThWork)
					case ThCmdClose:
						rt.setStatus(ThClosed)
					}
				default:
				}
				if rt.Thread.Status == ThClosed {
					errHandle(i, ErrThreadClosed)
					continue main
					//return nil, runtimeError(rt, i, ErrThreadClosed)
				}
			}
			if rt.Thread.Sleep > 0 {
				if step > rt.Thread.Sleep {
					step = rt.Thread.Sleep
				}
				time.Sleep(time.Duration(step) * time.Millisecond)
				rt.Thread.Sleep -= step
			} else if rt.Thread.Status == ThPaused || rt.Thread.Status == ThWait {
				if rt.ThreadID == 0 {
					select {
					case err = <-rt.Owner.ChError:
						return nil, err
					case x = <-rt.Thread.Chan:
						if x == ThCmdContinue {
							rt.setStatus(ThWork)
						}
					}
				} else {
					select {
					case x = <-rt.Thread.Chan:
						switch x {
						case ThCmdResume, ThCmdContinue:
							rt.setStatus(ThWork)
						case ThCmdClose:
							rt.setStatus(ThClosed)
						}
					}
				}
			}
			check = false
		}
	}
	return
}
