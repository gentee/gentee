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

func (rt *Runtime) Run(i int64) (result interface{}, err error) {
	top := Call{}
	code := rt.Owner.Exec.Code
	end := int64(len(code))
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
		case core.LENSTR:
			top.Str--
			rt.SInt[top.Int] = int64(len(rt.SStr[top.Str]))
			top.Int++
		case core.GETVAR:
			//soff := top.Str
			ioff := top.Int
			var count int
			if code[i+2]&0xffff == core.INDEX {
				count = int(code[i+2] >> 16)
				for ind := 0; ind < count; ind++ {
					// for map will be here
					ioff--
				}
			}
			shift := int(code[i]) >> 16
			i++
			typeVar := int(code[i]) >> 16
			index := int32(int(code[i]) & 0xffff)
			blockOff := rt.Calls[int32(len(rt.Calls)-1-shift)]
			switch typeVar & 0xff {
			case core.STACKSTR:
				rt.SStr[top.Str] = rt.SStr[blockOff.Str+index]
				top.Str++
			case core.STACKANY:
				rt.SAny[top.Any] = rt.SAny[blockOff.Any+index]
				top.Any++
			default:
				rt.SInt[top.Int] = rt.SInt[blockOff.Int+index]
				top.Int++
			}
			if count > 0 {
				//				var typeRet int
				i++
				for ind := 0; ind < count; ind++ {
					i++
					typeVar = int(code[i]) >> 16
					//					typeRet = int(code[i]) & 0xffff
					switch typeVar {
					case core.TYPESTR:
						top.Str--
						index := rt.SInt[ioff]
						runes := []rune(rt.SStr[top.Str])
						//						fmt.Println(`INDEX`, index, rt.SStr[top.Str])
						if index < 0 || index >= int64(len(runes)) {
							return nil, runtimeError(rt, i, ErrIndexOut)
						}
						rt.SInt[ioff] = int64(runes[index])
						top.Int = ioff + 1
					}
					//					fmt.Printf("IND %x %x\r\n", typeVar, typeRet)
				}
			}
			/*		case core.SETVAR:
					shift := int(code[i]) >> 16
					i++
					typeVar := int(code[i]) >> 16
					index := int32(int(code[i]) & 0xffff)
					switch typeVar {
					default:
						rt.SInt[rt.Calls[int32(len(rt.Calls)-1-shift)].Int+index] = rt.SInt[top.Int-1]
					}*/
		case core.ADDRESS:
			shift := int(code[i]) >> 16
			i++
			typeVar := int(code[i]) >> 16
			index := int32(int(code[i]) & 0xffff)
			switch typeVar {
			case core.TYPESTR:
				rt.SAny[top.Any] = &rt.SStr[rt.Calls[int32(len(rt.Calls)-1-shift)].Str+index]
			case core.TYPEARR:
				rt.SAny[top.Any] = rt.SAny[rt.Calls[int32(len(rt.Calls)-1-shift)].Any+index]
			default:
				rt.SAny[top.Any] = &rt.SInt[rt.Calls[int32(len(rt.Calls)-1-shift)].Int+index]
			}
			top.Any++
		case core.ASSIGN:
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
			}
		case core.ASSIGNADD:
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
			top.Int++
		case core.DUP:
			rt.SInt[top.Int] = rt.SInt[top.Int-1]
			top.Int++
		case core.CYCLE:
			lenCalls := len(rt.Calls) - 1
			rt.Calls[lenCalls].Cycle--
			//			fmt.Println(`CYCLE`, rt.Calls[lenCalls].Cycle, rt.SInt[:top.Int], rt.SStr[:top.Str], rt.SAny[:top.Any])
			if rt.Calls[lenCalls].Cycle == 0 {
				return nil, runtimeError(rt, i, ErrCycle)
			}
		case core.JMP:
			i += int64(int16(code[i+1]))
			//	top = rt.Calls[len(rt.Calls)-1]
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
					case core.TYPECHAR:
						rt.SInt[top.Int] = int64(' ')
						top.Int++
					case core.TYPESTR:
						rt.SStr[top.Str] = ``
						top.Str++
					case core.TYPEARR:
						rt.SAny[top.Any] = core.NewArray()
						top.Any++
					default:
						rt.SInt[top.Int] = 0
						top.Int++
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
		case core.DELVARS:
			curTop := top
			top = rt.Calls[len(rt.Calls)-1]
			rt.Calls = rt.Calls[:len(rt.Calls)-1]
			for j := top.Any; j < curTop.Any; j++ {
				rt.SAny[j] = nil
			}
		case core.RANGE:
			top.Int -= 2
			rt.SAny[top.Any] = core.Range{From: rt.SInt[top.Int], To: rt.SInt[top.Int+1]}
			top.Any++
		case core.LENGTH:
			var length int64
			if code[i]>>16 == core.TYPESTR {
				length = int64(len([]rune(rt.SStr[top.Str-1])))
			} else {
				switch v := rt.SAny[top.Any-1].(type) {
				case *core.Array:
					length = int64(v.Len())
				case core.Range:
					length = v.To - v.From
					if length < 0 {
						length = -length
					}
					length++
				default:
					fmt.Printf("LENGTH %T\n", v)
				}
			}
			rt.SInt[top.Int] = length
			top.Int++
		case core.FORINDEX:
			typeSrc := code[i] >> 16
			i++
			indDest := int32(code[i] & 0xffff)
			typeDest := code[i] >> 16
			i++
			index := rt.SInt[rt.Calls[int(len(rt.Calls)-1)].Int+int32(code[i]&0xffff)]
			if typeSrc == core.TYPESTR {
				val := []rune(rt.SStr[top.Str-1])
				rt.SInt[rt.Calls[int(len(rt.Calls)-1)].Int+indDest] = int64(val[index])
			} else {
				switch v := rt.SAny[top.Any-1].(type) {
				case *core.Array:
					switch typeDest & 0xff {
					case core.STACKINT:
						rt.SInt[rt.Calls[int(len(rt.Calls)-1)].Int+indDest] = v.Data[index].(int64)
					case core.STACKSTR:
						rt.SStr[rt.Calls[int(len(rt.Calls)-1)].Str+indDest] = v.Data[index].(string)
					}
				case core.Range:
					val := int64(v.From - index)
					if v.From < v.To {
						val = v.From + index
					}
					rt.SInt[rt.Calls[int(len(rt.Calls)-1)].Int+indDest] = val
				default:
					fmt.Println(`FORINDEX`)
				}
			}
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
		case core.SETINDEX:
			top.Any--
			switch v := rt.SAny[top.Any].(type) {
			case *string:
				top.Int--
				index := rt.SInt[top.Int-1]
				runes := []rune(*v)
				if index < 0 || index >= int64(len(runes)) {
					return nil, runtimeError(rt, i, ErrIndexOut)
				}
				runes[index] = rune(rt.SInt[top.Int])
				rt.SInt[top.Int-1] = int64(runes[index])
				*v = string(runes)
			case *core.Array:
				top.Int--
				index := rt.SInt[top.Int]
				if index < 0 || index >= int64(v.Len()) {
					return nil, runtimeError(rt, i, ErrIndexOut)
				}
				switch code[i] >> 16 {
				case core.TYPEINT:
					v.Data[index] = rt.SInt[top.Int-1]
				case core.TYPESTR:
					v.Data[index] = rt.SStr[top.Str-1]
				default:
					v.Data[index] = rt.SAny[top.Any-1]
				}
			case *interface{}:
				top.Int--
				index := rt.SInt[top.Int-1]
				runes := []rune((*v).(string))
				if index < 0 || index >= int64(len(runes)) {
					return nil, runtimeError(rt, i, ErrIndexOut)
				}
				runes[index] = rune(rt.SInt[top.Int])
				rt.SInt[top.Int-1] = int64(runes[index])
				*v = string(runes)

			default:
				fmt.Printf("TYPE=%T %v\r\n", v, v)
			}
		case core.ADDRINDEX:
			switch v := rt.SAny[top.Any-1].(type) {
			case *core.Array:
				top.Int--
				index := rt.SInt[top.Int]
				switch code[i] >> 16 {
				case core.TYPESTR:
					if index < 0 || index >= int64(v.Len()) {
						return nil, runtimeError(rt, i, ErrIndexOut)
					}
					rt.SAny[top.Any-1] = &v.Data[index]
				}
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
