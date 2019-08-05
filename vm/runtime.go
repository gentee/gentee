// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"github.com/gentee/gentee/core"
)

func (rt *Runtime) Run(i int64) (result interface{}, err error) {
	top := rt.States[len(rt.States)-1]
	topInt, topFloat, topStr, topAny := top.Get()
	code := rt.Owner.Exec.Code
	end := int64(len(code))
main:
	for i < end {
		switch code[i] {
		case core.PUSH16:
			i++
			topInt++
			rt.SInt[topInt] = int64(code[i])
		case core.PUSH32:
			i += 2
			topInt++
			rt.SInt[topInt] = int64((uint64(code[i-1]) << 16) | uint64(code[i]))
		case core.PUSH64:
			i += 4
			topInt++
			rt.SInt[topInt] = int64((uint64(code[i-3]) << 48) | (uint64(code[i-2]) << 32) |
				(uint64(code[i-1]) << 16) | uint64(code[i]))
		case core.ADD:
			topInt--
			rt.SInt[topInt] += rt.SInt[topInt+1]
		case core.SUB:
			topInt--
			rt.SInt[topInt] -= rt.SInt[topInt+1]
		case core.MUL:
			topInt--
			rt.SInt[topInt] *= rt.SInt[topInt+1]
		case core.DIV:
			topInt--
			rt.SInt[topInt] /= rt.SInt[topInt+1]
		case core.MOD:
			topInt--
			rt.SInt[topInt] %= rt.SInt[topInt+1]
		case core.SIGN:
			rt.SInt[topInt] = -rt.SInt[topInt]
		case core.JMP:
			i += int64(int16(code[i+1]))
			topInt, topFloat, topStr, topAny = top.Get()
			continue
		case core.JZE:
			topInt--
			if rt.SInt[topInt+1] == 0 {
				i += int64(int16(code[i+1]))
				continue
			}
			i++
		case core.RET:
			i++
			if len(rt.States) == 1 { // return from run function
				switch code[i] {
				case core.STACKINT:
					result = rt.SInt[topInt]
				case core.STACKBOOL:
					if rt.SInt[topInt] == 0 {
						result = false
					} else {
						result = true
					}
				case core.STACKCHAR:
					result = rune(rt.SInt[topInt])
				}
			} else {
				switch code[i] {
				case core.STACKNO:
				default:
					rt.States[len(rt.States)-1].topInt++
					rt.SInt[top.topInt+1] = rt.SInt[topInt]
				}
			}
			break main
		case core.END:
			break main
		case core.CALLBYID:
			i++
			rt.PushState(topInt, topFloat, topStr, topAny)
			rt.Run(int64(rt.Owner.Exec.Funcs[code[i]]))
			topInt, topFloat, topStr, topAny = rt.PopState()
			top = rt.States[len(rt.States)-1]
		}
		i++
	}
	return
}
