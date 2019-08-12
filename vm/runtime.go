// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"github.com/gentee/gentee/core"
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
			rt.SInt[top.Int-1] %= rt.SInt[top.Int]
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
		case core.GETVAR:
			shift := int(code[i]) >> 16
			i++
			typeVar := int(code[i]) >> 16
			index := int32(int(code[i]) & 0xffff)
			switch typeVar {
			default:
				rt.SInt[top.Int] = rt.SInt[rt.Calls[int32(len(rt.Calls)-1-shift)].Int+index]
				top.Int++
			}
		case core.SETVAR:
			shift := int(code[i]) >> 16
			i++
			typeVar := int(code[i]) >> 16
			index := int32(int(code[i]) & 0xffff)
			switch typeVar {
			default:
				rt.SInt[rt.Calls[int32(len(rt.Calls)-1-shift)].Int+index] = rt.SInt[top.Int-1]
			}
		case core.DUP:
			rt.SInt[top.Int] = rt.SInt[top.Int-1]
			top.Int++
		case core.CYCLE:
			lenCalls := len(rt.Calls) - 1
			rt.Calls[lenCalls].Cycle--
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
					switch varType {
					default:
						prevTop.Int--
						curTop.Int--
					}
				} else {
					switch varType {
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
				}
				break main
			}
			curTop := top
			top = rt.Calls[k]
			rt.Calls = rt.Calls[:k]
			switch retType {
			case core.TYPENONE:
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
		}
		i++
	}
	return
}
