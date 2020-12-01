// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"math"
	"reflect"

	"github.com/gentee/gentee/core"
)

func cmd2Code(linker *Linker, cmd core.ICmd, out *core.Bytecode) {

	var (
		cmds []core.Bytecode
	)
	push := func(pars ...core.Bcode) {
		out.Code = append(out.Code, pars...)
	}
	getIndex := func(cmdVar *core.CmdVar, command core.Bcode) {
		var (
			shift  int
			locOut bool
		)
		block := cmdVar.Block
		for shift = len(linker.Blocks) - 1; shift >= 0; shift-- {
			if linker.Blocks[shift].Block == block {
				break
			}
			if linker.Blocks[shift].IsLocal {
				locOut = true
			}
		}
		for i := len(cmdVar.Indexes) - 1; i >= 0; i-- {
			cmd2Code(linker, cmdVar.Indexes[i].Cmd, out)
		}
		inType := int(type2Code(block.Vars[cmdVar.Index], out))
		blockShift := len(linker.Blocks) - 1 - shift
		if locOut {
			blockShift = 0x0f00 + shift
		}
		push(core.Bcode(blockShift<<16)|command,
			core.Bcode(inType<<16|linker.Blocks[shift].Vars[cmdVar.Index]))
		if inType >= core.TYPESTRUCT {
			structOffset(out, -len(out.Code)+1)
		}
		if len(cmdVar.Indexes) > 0 {
			push(core.Bcode(len(cmdVar.Indexes)<<16) | core.INDEX)
			for _, ival := range cmdVar.Indexes {
				retType := type2Code(ival.Type, out)
				code := core.Bcode(inType<<16) | retType
				if type2Code(ival.Cmd.GetResult(), out) == core.TYPESTR {
					code |= 0x8000
				}
				push(code)
				if inType >= core.TYPESTRUCT {
					structOffset(out, -len(out.Code)+1)
				}
				if retType >= core.TYPESTRUCT {
					structOffset(out, len(out.Code)-1)
				}
				getPos(linker, ival.Cmd, out)
				inType = int(retType)
				/*if typeValue == nil {
					return runtimeError(rt, cmdVar, ErrRuntime, `getVar.typeValue`)
				}*/
			}
		}
	}
	callFunc := func(count int, ptypes ...core.Bcode) {
		canError := true
		obj := cmd.GetObject()
		switch obj.GetType() {
		case core.ObjEmbedded:
			embed := obj.(*core.EmbedObject)
			if ind, ok := embed.Func.(int32); ok {
				push(core.Bcode(ind<<16) | core.EMBED)
				if embed.Variadic {
					push(core.Bcode(len(ptypes)))
					if len(ptypes) > 0 {
						push(ptypes...)
					}
				}
			} else if embed.BCode.Code != nil {
				code := embed.BCode.Code[0]
				if code != core.NOP {
					push(code) //...)
				}
			} else {
				fmt.Println(`EMBED obj`, obj)
			}
			canError = embed.CanError
		case core.ObjFunc:
			id := obj.(*core.FuncObject).ObjID
			anyFunc := cmd.(*core.CmdAnyFunc)
			if anyFunc.IsThread {
				block := obj.(*core.FuncObject).Block
				push(core.Bcode(block.ParCount<<16)|core.GOBYID, core.Bcode(id))
				for k := 0; k < block.ParCount; k++ {
					ptype := type2Code(block.Vars[k], out)
					push(core.Bcode(ptype))
					if ptype >= core.TYPESTRUCT {
						structOffset(out, len(out.Code)-1)
					}
				}
			} else {
				push(core.Bcode(count<<16)|core.CALLBYID, core.Bcode(id))
			}
			if out.Used == nil {
				out.Used = make(map[int32]byte)
			}
			if out.Used[id] == 0 {
				genBytecode(obj.(*core.FuncObject).Unit.VM, id)
				copyUsed(&obj.(*core.FuncObject).BCode, out)
				out.Used[id] = 1
			}
		case core.ObjConst:
			id := obj.(*core.ConstObject).ObjID
			push(core.CONSTBYID, core.Bcode(id))
			if out.Used == nil {
				out.Used = make(map[int32]byte)
			}
			if out.Used[id] == 0 {
				out.Init = append(out.Init, id)
				genBytecode(obj.(*core.ConstObject).Unit.VM, id)
				copyUsed(&obj.(*core.ConstObject).BCode, out)
				out.Used[id] = 1
			}
		}
		if canError {
			getPos(linker, cmd, out)
		}
	}
	switch cmd.GetType() {
	case core.CtCommand:
		v := cmd.(*core.CmdCommand).ID
		switch v {
		case core.RcBreak:
			push(core.BREAK)
		case core.RcContinue:
			push(core.CONTINUE)
		case core.RcRecover:
			push(core.RECOVER)
		case core.RcRetry:
			push(core.RETRY)
		}
	case core.CtFunc:
		anyFunc := cmd.(*core.CmdAnyFunc)
		for _, param := range anyFunc.Children {
			cmd2Code(linker, param, out)
		}
		obj := cmd.GetObject()
		count := len(anyFunc.Children)
		if obj == nil {
			cmd2Code(linker, anyFunc.FnVar, out)
			push(core.Bcode(count<<16)|core.CALLBYID, 0)
			getPos(linker, cmd, out)
		} else if obj.GetType() == core.ObjEmbedded {
			if obj.(*core.EmbedObject).Variadic {
				vcount := count - len(obj.(*core.EmbedObject).Params)
				ptypes := make([]core.Bcode, vcount)
				for i := 0; i < vcount; i++ {
					ptypes[i] = type2Code(anyFunc.Children[len(obj.(*core.EmbedObject).Params)+i].GetResult(), out)
					// we don't need call structOffset here because it doesn't matter type of struct
				}
				callFunc(count, ptypes...)
			} else {
				callFunc(count)
			}
		} else if obj.GetType() == core.ObjFunc {
			var optCount int
			if optCount = len(anyFunc.Optional); optCount > 0 {
				push(core.Bcode(optCount<<16) | core.OPTPARS)
				for _, num := range anyFunc.Optional {
					ptype := type2Code(obj.(*core.FuncObject).Block.Vars[num], out)
					push(ptype<<16 | core.Bcode(num))
					if ptype >= core.TYPESTRUCT {
						structOffset(out, -len(out.Code)+1)
					}
				}
				count -= optCount
			}
			if obj.(*core.FuncObject).Block.Variadic {
				block := obj.(*core.FuncObject).Block
				vcount := count - block.ParCount
				push(core.Bcode(vcount<<16 | core.ARRAY))
				for j := vcount - 1; j >= 0; j-- {
					typeRet := anyFunc.Children[block.ParCount+j].GetResult()
					itype := type2Code(typeRet, out)
					var isarray int32
					if itype == core.TYPEARR && isEqualTypes(block.Vars[block.ParCount], typeRet) {
						isarray = 1
					}
					push(core.Bcode(isarray<<16) | itype)
					if itype >= core.TYPESTRUCT {
						structOffset(out, len(out.Code)-1)
					}
				}
				callFunc(block.ParCount + 1)
			} else {
				callFunc(count)
			}
		}
	case core.CtBinary:
		cmd2Code(linker, cmd.(*core.CmdBinary).Left, out)
		cmd2Code(linker, cmd.(*core.CmdBinary).Right, out)
		callFunc(2)
	case core.CtUnary:
		cmd2Code(linker, cmd.(*core.CmdUnary).Operand, out)
		callFunc(1)
	case core.CtValue:
		switch v := cmd.(*core.CmdValue).Value.(type) {
		case int64:
			if v <= math.MaxInt32 && v >= math.MinInt32 {
				push(core.PUSH32, core.Bcode(v))
			} else {
				u64 := uint64(v)
				push(core.PUSH64, core.Bcode(uint32(u64>>32)),
					core.Bcode(uint32(u64&0xffffffff)))
			}
		case float64:
			u64 := math.Float64bits(v)
			push(core.PUSHFLOAT, core.Bcode(uint32(u64>>32)),
				core.Bcode(uint32(u64&0xffffffff)))
		case bool:
			var val int32
			if v {
				val = 1
			}
			push(core.PUSH32, core.Bcode(val))
		case rune:
			push(core.PUSH32, core.Bcode(v))
		case string:
			var (
				ok bool
				id uint16
			)
			if id, ok = out.Strings[v]; !ok {
				id = uint16(len(out.Strings))
				out.Strings[v] = id
			}
			out.StrOffset = append(out.StrOffset, int32(len(out.Code)))
			push(core.Bcode(uint32(id)<<16) | core.PUSHSTR)
		case *core.Fn:
			id := v.Func.(*core.FuncObject).ObjID
			push(core.PUSHFUNC, core.Bcode(id))
			if out.Used == nil {
				out.Used = make(map[int32]byte)
			}
			if out.Used[id] == 0 {
				genBytecode(v.Func.(*core.FuncObject).Unit.VM, id)
				copyUsed(&v.Func.(*core.FuncObject).BCode, out)
				out.Used[id] = 1
			}
		default:
			fmt.Printf("CmdValue %T %v\n", cmd.(*core.CmdValue).Value, v)
		}
	case core.CtConst:
		callFunc(1)
	case core.CtVar:
		getIndex(cmd.(*core.CmdVar), core.GETVAR)
	case core.CtStack:
		cmdStack := cmd.(*core.CmdBlock)
		switch cmdStack.ID {
		case core.StackSwitch:
			cmd2Code(linker, cmdStack.Children[0], out)
			cmpType := type2Code(cmdStack.Children[0].GetResult(), out)
			offsets := make([]int, 0)
			for i := 1; i < len(cmdStack.Children); i++ {
				caseStack := cmdStack.Children[i].(*core.CmdBlock)
				if caseStack.ID == core.StackDefault {
					out.BlockFlags = core.BlBreak
					offsets = append(offsets, len(out.Code))
					cmd2Code(linker, caseStack, out)
					break
				} else {
					cases := make([]int, 0)
					for j := 0; j < len(caseStack.Children)-1; j++ {
						cmd2Code(linker, caseStack.Children[j], out)
						cases = append(cases, len(out.Code))
						push(core.Bcode(cmpType<<16)|core.JEQ, 0)
					}
					pos := len(out.Code)
					push(core.JMP, 0)
					for _, icase := range cases {
						out.Code[icase+1] = core.Bcode(len(out.Code) - icase)
					}
					offsets = append(offsets, len(out.Code))
					out.BlockFlags = core.BlBreak
					cmd2Code(linker, caseStack.Children[len(caseStack.Children)-1], out)
					offsets = append(offsets, len(out.Code))
					push(core.JMP, 0)
					out.Code[pos+1] = core.Bcode(len(out.Code) - pos)
					cmds = cmds[:0]
				}
			}
			for _, ioff := range offsets {
				out.Code[ioff+1] = core.Bcode(len(out.Code) - ioff)
			}
		case core.StackNew:
			switch cmd.(*core.CmdBlock).Result.Original {
			case reflect.TypeOf(core.Array{}), reflect.TypeOf(core.Map{}):
				for _, icmd := range cmdStack.Children {
					cmd2Code(linker, icmd, out)
				}
				left := type2Code(cmd.(*core.CmdBlock).Result.IndexOf, out)
				right := type2Code(cmd.(*core.CmdBlock).Result, out)
				push(core.Bcode(len(cmdStack.Children)<<16|core.INITOBJ), left<<16|right)
				if left >= core.TYPESTRUCT {
					structOffset(out, -len(out.Code)+1)
				}
				if right >= core.TYPESTRUCT {
					structOffset(out, len(out.Code)-1)
				}
			case reflect.TypeOf(core.Buffer{}):
				for _, icmd := range cmdStack.Children {
					cmd2Code(linker, icmd, out)
					retType := type2Code(icmd.GetResult(), out)
					if retType == core.TYPEINT {
						push(core.PUSH32, core.Bcode(len(out.Code)))
						getPos(linker, icmd, out)
					}
					push(core.PUSH32, retType)
				}
				push(core.Bcode(len(cmdStack.Children)<<16|core.INITOBJ), core.TYPEBUF)
			case reflect.TypeOf(core.Set{}):
				for _, icmd := range cmdStack.Children {
					cmd2Code(linker, icmd, out)
				}
				push(core.Bcode(len(cmdStack.Children)<<16|core.INITOBJ), core.TYPESET)
			case reflect.TypeOf(core.Struct{}):
				for _, item := range cmdStack.Children {
					cmd2Code(linker, item.(*core.CmdBinary).Right, out)
					cmd2Code(linker, item.(*core.CmdBinary).Left, out)
				}
				right := type2Code(cmd.(*core.CmdBlock).Result, out)
				push(core.Bcode(len(cmdStack.Children)<<16|core.INITOBJ), right)
				if right >= core.TYPESTRUCT {
					structOffset(out, len(out.Code)-1)
				}
			default:
				fmt.Println(`init arr`)
			}
		case core.StackInit, core.StackInitPtr:
			cmd2Code(linker, cmdStack.Children[1], out)
			rightType := type2Code(cmdStack.Children[1].GetResult(), out)
			cmdVar := cmdStack.Children[0].(*core.CmdVar)
			getIndex(cmdVar, core.SETVAR)
			cmdID := core.Bcode(core.ASSIGN)
			if cmdStack.ID == core.StackInitPtr {
				cmdID = core.ASSIGNPTR
			}
			push(rightType<<16 | cmdID)
			if rightType >= core.TYPESTRUCT {
				structOffset(out, -len(out.Code)+1)
			}
		case core.StackQuestion:
			cmd2Code(linker, cmdStack.Children[0], out)
			pos := len(out.Code)
			push(core.JZE, 0)
			cmd2Code(linker, cmdStack.Children[1], out)
			out.Code[pos+1] = core.Bcode(len(out.Code) - pos + 2)
			pos = len(out.Code)
			push(core.JMP, 0)
			cmd2Code(linker, cmdStack.Children[2], out)
			out.Code[pos+1] = core.Bcode(len(out.Code) - pos)
		case core.StackAnd, core.StackOr:
			cmd2Code(linker, cmdStack.Children[0], out)
			logic := core.Bcode(core.JZE)
			if cmd.(*core.CmdBlock).ID == core.StackOr {
				logic = core.JNZ
			}
			pos := len(out.Code)
			push(core.DUP, logic, 0) //core.Bcode(size))
			cmd2Code(linker, cmdStack.Children[1], out)
			out.Code[pos+2] = core.Bcode(len(out.Code) - pos - 3 + 2)
		case core.StackAssign, core.StackIncDec:
			rightType := core.Bcode(core.TYPEINT)
			if cmdStack.ID == core.StackAssign {
				cmd2Code(linker, cmdStack.Children[1], out)
				rightType = type2Code(cmdStack.Children[1].GetResult(), out)
			} else {
				push(core.PUSH32, core.Bcode(cmdStack.ParCount))
			}
			cmdVar := cmdStack.Children[0].(*core.CmdVar)
			getIndex(cmdVar, core.SETVAR)
			if cmdStack.ID == core.StackAssign {
				obj := cmd.GetObject()
				if obj.GetType() == core.ObjEmbedded {
					embed := obj.(*core.EmbedObject)
					if ind, ok := embed.Func.(int32); ok {
						push(rightType<<16 | core.Bcode(ind+core.EMBEDFUNC))
						if rightType >= core.TYPESTRUCT {
							structOffset(out, -len(out.Code)+1)
						}
						getPos(linker, cmdStack, out)
					} else if embed.BCode.Code != nil {
						push(rightType<<16 | embed.BCode.Code[0])
						if rightType >= core.TYPESTRUCT {
							structOffset(out, -len(out.Code)+1)
						}
						getPos(linker, cmdStack, out)
					} else {
						fmt.Println(`OOOPS CODE`, embed.BCode.Code, obj.GetName(),
							obj.(*core.EmbedObject).Params[0], obj.(*core.EmbedObject).Params[1])
					}
				} else {
					fmt.Println(`OOOPS`, obj)
				}
			} else {
				push(core.INCDEC)
			}
		case core.StackIf:
			var k int
			lenIf := len(cmdStack.Children) >> 1
			jumps := make([]int, lenIf)
			for k = 0; k < lenIf; k++ {
				cmd2Code(linker, cmdStack.Children[k<<1], out)
				pos := len(out.Code)
				push(core.JZE, 0)
				cmd2Code(linker, cmdStack.Children[(k<<1)+1], out)
				jumps[k] = len(out.Code)
				push(core.JMP, core.Bcode(jumps[k]))
				out.Code[pos+1] = core.Bcode(len(out.Code) - pos)
			}
			if len(cmdStack.Children)&1 == 1 {
				cmd2Code(linker, cmdStack.Children[len(cmdStack.Children)-1], out)
			}
			for _, off := range jumps {
				out.Code[off+1] = core.Bcode(len(out.Code) - off)
			}
		case core.StackWhile:
			pos := len(out.Code)
			push(core.CYCLE)
			getPos(linker, cmdStack, out)
			cmd2Code(linker, cmdStack.Children[0], out)
			push(core.JZE, 0) // core.Bcode(save(cmdStack.Children[1])+2))
			blockStart := len(out.Code)
			out.BlockFlags = core.BlContinue | core.BlBreak
			cmd2Code(linker, cmdStack.Children[1], out)
			out.Code[blockStart-1] = core.Bcode(len(out.Code) - blockStart + 4)
			push(core.JMP, core.Bcode(pos-len(out.Code)))
			out.Code[blockStart+1] = core.Bcode(len(out.Code) - blockStart) // set break of BLOCK
			out.Code[blockStart+2] = core.Bcode(pos - blockStart)           // set continue of BLOCK
		case core.StackFor:
			bInfo, _ := initBlock(linker, cmdStack, out)

			cmd2Code(linker, cmdStack.Children[0], out)
			srcType := type2Code(cmdStack.Children[0].GetResult(), out)
			curType := type2Code(cmdStack.Vars[0], out)
			// we don't need to use structOffset because for doesn't support structs
			indcur := 0
			if curType&0xf == core.STACKANY {
				indcur = 1
			}
			pos := len(out.Code)
			push(core.CYCLE)
			getPos(linker, cmdStack, out)
			push(core.GETVAR, core.Bcode(int(core.TYPEINT)<<16|bInfo.Vars[1]),
				(srcType<<16)|core.DUP, (srcType<<16)|core.LEN, core.LT)
			posJmp := len(out.Code)
			push(core.JZE, 0)
			push(core.GETVAR, core.Bcode(int(core.TYPEINT)<<16|bInfo.Vars[1])) // set index
			push(core.GETVAR, core.Bcode(int(srcType)<<16|indcur),             // get cur value
				core.Bcode(1<<16|core.INDEX), core.Bcode(int(srcType)<<16)|curType)
			push(core.SETVAR, core.Bcode(int(curType)<<16|bInfo.Vars[0]),
				core.Bcode(int(curType)<<16|core.ASSIGNPTR), core.Bcode(int(curType)<<16|core.POP))
			blockStart := len(out.Code)
			out.BlockFlags = core.BlContinue | core.BlBreak
			cmd2Code(linker, cmdStack.Children[1], out)
			out.Code[blockStart+2] = core.Bcode(len(out.Code) - blockStart) // set continue of BLOCK
			push(core.Bcode(bInfo.Vars[1]<<16) | core.FORINC)
			push(core.JMP, core.Bcode(pos-len(out.Code)))
			out.Code[blockStart+1] = core.Bcode(len(out.Code) - blockStart) // set break of BLOCK
			out.Code[posJmp+1] = core.Bcode(len(out.Code) - posJmp)
			push(core.DELVARS)
			linker.Blocks = linker.Blocks[:len(linker.Blocks)-1]

		case core.StackBlock, core.StackDefault, core.StackCatch:
			initBlock(linker, cmdStack, out)
			for _, item := range cmdStack.Children {
				cmd2Code(linker, item, out)
			}
			if cmdStack.ID == core.StackCatch {
				push(core.CATCH)
			}
			push(core.DELVARS)
			linker.Blocks = linker.Blocks[:len(linker.Blocks)-1]
			//deleteVars(rt)
		case core.StackReturn:
			if cmdStack.Children != nil {
				cmd2Code(linker, cmdStack.Children[0], out)
				retType := type2Code(cmdStack.Children[0].GetResult(), out)
				push((retType << 16) | core.RET)
				if retType >= core.TYPESTRUCT {
					structOffset(out, -len(out.Code)+1)
				}
			} else {
				push(core.END)
			}
		case core.StackOptional:
			pos := len(out.Code)
			// assigns value if the variable has not yet been assigned as optional
			push(core.Bcode(cmdStack.ParCount<<16)|core.JMPOPT, 0)
			cmd2Code(linker, cmdStack.Children[0], out)
			out.Code[pos+1] = core.Bcode(len(out.Code) - pos)
		case core.StackLocal:
			linker.Blocks = append(linker.Blocks, BlockInfo{
				Block:   cmdStack.Children[0].(*core.CmdBlock),
				IsLocal: true,
			})
			pos := len(out.Code)
			push(core.JMP, 0) //core.Bcode(save(cmdStack.Children[0]))+3)
			out.Locals = append(out.Locals, core.Local{
				Cmd:    cmdStack.Children[0].(*core.CmdBlock),
				Offset: len(out.Code),
			})
			cmd2Code(linker, cmdStack.Children[0], out)
			push(core.END)
			out.Code[pos+1] = core.Bcode(len(out.Code) - pos)
			linker.Blocks = linker.Blocks[:len(linker.Blocks)-1]
		case core.StackCallLocal:
			offset := -1
			for _, local := range out.Locals {
				if cmdStack.Children[0].(*core.CmdBlock) == local.Cmd {
					offset = local.Offset
					break
				}
			}
			if offset < 0 {
				fmt.Println(`STACK CALL LOCAL ERROR`)
			}
			for i := 1; i < len(cmdStack.Children); i++ {
				cmd2Code(linker, cmdStack.Children[i], out)
			}
			push(core.Bcode((len(cmdStack.Children)-1)<<16)|core.LOCAL,
				core.Bcode(offset-len(out.Code)-1))
		case core.StackLocret:
			cmd2Code(linker, cmdStack.Children[0], out)
			retType := type2Code(cmdStack.Children[0].GetResult(), out)
			push((retType << 16) | core.RET)
			if retType >= core.TYPESTRUCT {
				structOffset(out, -len(out.Code)+1)
			}
		case core.StackTry:
			out.BlockFlags = core.BlTry
			blockTry := len(out.Code)
			cmd2Code(linker, cmdStack.Children[0], out)
			pos := len(out.Code)
			push(core.JMP, 0)
			out.Code[blockTry+1] = core.Bcode(len(out.Code) - blockTry)
			blockCatch := len(out.Code)
			out.BlockFlags = core.BlRecover | core.BlRetry
			cmd2Code(linker, cmdStack.Children[1], out)
			out.Code[pos+1] = core.Bcode(len(out.Code) - pos)
			out.Code[blockCatch+1] = core.Bcode(len(out.Code) - blockCatch) // recover jump
			out.Code[blockCatch+2] = core.Bcode(blockTry - blockCatch)      // retry jump
		}
	}
}
