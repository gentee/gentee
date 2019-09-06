// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"math"
	"reflect"

	"github.com/gentee/gentee/core"
	stdlib "github.com/gentee/gentee/stdlibvm"
)

func cmd2Code(linker *Linker, cmd core.ICmd, out *core.Bytecode) {

	var (
		cmds []core.Bytecode
	)

	save := func(icmd core.ICmd) int {
		code := core.Bytecode{
			Code:       make([]core.Bcode, 0, 16),
			Strings:    out.Strings,
			BlockFlags: out.BlockFlags,
		}
		cmd2Code(linker, icmd, &code)
		copyUsed(&code, out)
		cmds = append(cmds, code)
		out.BlockFlags = 0
		return len(code.Code)
	}
	push := func(pars ...core.Bcode) {
		out.Code = append(out.Code, pars...)
	}
	getIndex := func(cmdVar *core.CmdVar, command core.Bcode) {
		var shift int
		block := cmdVar.Block
		for shift = len(linker.Blocks) - 1; shift >= 0; shift-- {
			if linker.Blocks[shift].Block == block {
				break
			}
		}
		for i := len(cmdVar.Indexes) - 1; i >= 0; i-- {
			cmd2Code(linker, cmdVar.Indexes[i].Cmd, out)
		}
		inType := int(type2Code(block.Vars[cmdVar.Index], out))
		push(core.Bcode((len(linker.Blocks)-1-shift)<<16)|command,
			core.Bcode(inType<<16|linker.Blocks[shift].Vars[cmdVar.Index]))
		if inType >= core.TYPESTRUCT {
			structOffset(out, -len(out.Code)+1)
		}
		//		getPos(linker, cmdVar, out)
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
			if embed.BCode.Code != nil {
				code := embed.BCode.Code[0]
				push(code) //...)
				if (code & 0xffff) == core.EMBED && (code>>16 < 1000) {
					embed := stdlib.Embedded[code>>16]
					if embed.Variadic {
						push(core.Bcode(len(ptypes)))
						if len(ptypes) > 0 {
							push(ptypes...)
						}
					}
				}
			}
			canError = embed.CanError
		case core.ObjFunc:
			id := obj.(*core.FuncObject).ObjID
			push(core.Bcode(count<<16)|core.CALLBYID, core.Bcode(id))
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
			//			fmt.Println(`POS`, obj.GetName(), out.Pos, out.Code)
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
			/*	case RcRecover, RcRetry:
					rt.Catch = v
					rt.Command = RcBreak
				default:
					rt.Command = v*/
		}
	case core.CtFunc:
		anyFunc := cmd.(*core.CmdAnyFunc)
		for _, param := range anyFunc.Children {
			cmd2Code(linker, param, out)
		}
		obj := cmd.GetObject()
		if obj == nil {
			cmd2Code(linker, anyFunc.FnVar, out)
			push(core.Bcode(len(anyFunc.Children)<<16)|core.CALLBYID, 0)
			getPos(linker, cmd, out)
		} else if obj.GetType() == core.ObjEmbedded && obj.(*core.EmbedObject).Variadic {
			count := len(anyFunc.Children) + 1 - len(obj.(*core.EmbedObject).Params)
			ptypes := make([]core.Bcode, count)
			for i := 0; i < count; i++ {
				ptypes[i] = type2Code(anyFunc.Children[len(obj.(*core.EmbedObject).Params)-1+i].GetResult(), out)
				// we don't need call structOffset here because it doesn't matter type of struct
			}
			callFunc(len(anyFunc.Children), ptypes...)
		} else if obj.GetType() == core.ObjFunc && obj.(*core.FuncObject).Block.Variadic {
			block := obj.(*core.FuncObject).Block
			count := len(anyFunc.Children) - block.ParCount
			push(core.Bcode(count<<16 | core.ARRAY))
			for j := count - 1; j >= 0; j-- {
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
			callFunc(len(anyFunc.Children))
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
		/*if typeValue == nil {
			return runtimeError(rt, cmdVar, ErrRuntime, `getVar.typeValue`)
		}*/
		//			cmds = cmds[:0]
		/*			typeValue := block.Vars[cmdVar.Index]
					if typeValue.Original == reflect.TypeOf(core.Struct{}) {
						//	typeValue = custom.Type.Custom.Types[index]
					} else {
						typeValue = typeValue.IndexOf
					}
				}*/
		/*		if err = getVar(rt, cmd.(*CmdVar)); err != nil {
				return err
			}*/
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
					out.BlockFlags = core.BlBreak
					push(core.JMP, core.Bcode(save(
						caseStack.Children[len(caseStack.Children)-1])+4))
					for _, icase := range cases {
						out.Code[icase+1] = core.Bcode(len(out.Code) - icase)
					}
					offsets = append(offsets, len(out.Code))
					out.Code = append(out.Code, cmds[0].Code...)
					offsets = append(offsets, len(out.Code))
					push(core.JMP, 0)
					cmds = cmds[:0]
				}
				//cmd2Code(linker, caseStack.Children[len(caseStack.Children)-1], out)
			}
			for _, ioff := range offsets {
				out.Code[ioff+1] = core.Bcode(len(out.Code) - ioff)
			}
			/*				if err = rt.runCmd(cmdStack.Children[0]); err != nil {
								return err
							}
							original := rt.Stack[len(rt.Stack)-1]
							rt.Stack = rt.Stack[:len(rt.Stack)-1]
							var (
								done bool
								def  ICmd
							)
							for i := 1; i < len(cmdStack.Children); i++ {
								caseStack := cmdStack.Children[i].(*CmdBlock)
								if caseStack.ID == StackDefault {
									def = caseStack
									break
								}
								for j := 0; j < len(caseStack.Children)-1; j++ {
									if err = rt.runCmd(caseStack.Children[j]); err != nil {
										return err
									}
									val := rt.Stack[len(rt.Stack)-1]
									rt.Stack = rt.Stack[:len(rt.Stack)-1]
									var equal bool
									switch v := original.(type) {
									case int64:
										equal = v == val.(int64)
									case rune:
										equal = v == val.(rune)
									case bool:
										equal = v == val.(bool)
									case string:
										equal = v == val.(string)
									case float64:
										equal = v == val.(float64)
									}
									if equal {
										if err = rt.runCmd(caseStack.Children[len(caseStack.Children)-1]); err != nil {
											return err
										}
										done = true
										if rt.Command == RcBreak && rt.Catch == 0 {
											rt.Command = 0
										}
										break
									}
								}
								if done {
									break
								}
							}
							if !done && def != nil {
								if err = rt.runCmd(def); err != nil {
									return err
								}
								if rt.Command == RcBreak && rt.Catch == 0 {
									rt.Command = 0
								}
							}*/
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
				/*				parr := NewArray()
								for _, icmd := range cmdStack.Children {
									if err = rt.runCmd(icmd); err != nil {
										return err
									}
									var ptr interface{}
									CopyVar(&ptr, rt.Stack[len(rt.Stack)-1])
									parr.Data = append(parr.Data, ptr)
								}
												if lenStack >= len(rt.Stack) {
													rt.Stack = append(rt.Stack, parr)
												} else {
													rt.Stack[lenStack] = parr
												}*/
				/*					case reflect.TypeOf(Set{}):
										pset := NewSet()
										for _, icmd := range cmdStack.Children {
											if err = rt.runCmd(icmd); err != nil {
												return err
											}
											switch v := rt.Stack[len(rt.Stack)-1].(type) {
											case int64:
												pset.Set(v, true)
											default:
												return runtimeError(rt, icmd, ErrRuntime, `init set`)
											}
										}
										rt.Stack[lenStack] = pset
									case reflect.TypeOf(Buffer{}):
										pbuf := NewBuffer()
										for _, icmd := range cmdStack.Children {
											if err = rt.runCmd(icmd); err != nil {
												return err
											}
											switch v := rt.Stack[len(rt.Stack)-1].(type) {
											case int64:
												if uint64(v) > 255 {
													return runtimeError(rt, icmd, ErrByteOut)
												}
												pbuf.Data = append(pbuf.Data, byte(v))
											case string:
												pbuf.Data = append(pbuf.Data, []byte(v)...)
											case rune:
												pbuf.Data = append(pbuf.Data, []byte(string([]rune{v}))...)
											case *Buffer:
												pbuf.Data = append(pbuf.Data, v.Data...)
											default:
												return runtimeError(rt, icmd, ErrRuntime, `init buf`)
											}
										}
										rt.Stack[lenStack] = pbuf*/
				/*									case reflect.TypeOf(Map{}):
														pmap := NewMap()
														for _, icmd := range cmdStack.Children {
															if err = rt.runCmd(icmd); err != nil {
																return err
															}
															var ptr interface{}
															CopyVar(&ptr, rt.Stack[len(rt.Stack)-1])
															keyValue := ptr.(KeyValue)
															pmap.Data[keyValue.Key.(string)] = keyValue.Value
															pmap.Keys = append(pmap.Keys, keyValue.Key.(string))
														}
														if lenStack >= len(rt.Stack) {
															rt.Stack = append(rt.Stack, pmap)
														} else {
															rt.Stack[lenStack] = pmap
														}
													case reflect.TypeOf(Struct{}):
														pstruct := NewStruct(cmd.(*CmdBlock).Result)
														for _, icmd := range cmdStack.Children {
															if err = rt.runCmd(icmd); err != nil {
																return err
															}
															var ptr interface{}
															CopyVar(&ptr, rt.Stack[len(rt.Stack)-1])
															keyValue := ptr.(KeyValue)
															pstruct.Values[keyValue.Key.(int64)] = keyValue.Value
														}
														if lenStack >= len(rt.Stack) {
															rt.Stack = append(rt.Stack, pstruct)
														} else {
															rt.Stack[lenStack] = pstruct
														}*/
			default:
				fmt.Println(`init arr`)
				//return runtimeError(rt, cmd, ErrRuntime, `init arr`)
			}
			//lenStack++
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

			/*					cmdVar := cmdStack.Children[0].(*CmdVar)
								if vars, err = rt.getVars(cmdVar.Block); err != nil {
									return err
								}
								if err = rt.runCmd(cmdStack.Children[1]); err != nil {
									return err
								}
								vars[cmdVar.Index] = rt.Stack[len(rt.Stack)-1]*/
		case core.StackQuestion:
			cmd2Code(linker, cmdStack.Children[0], out)
			push(core.JZE, core.Bcode(save(cmdStack.Children[1])+2))
			out.Code = append(out.Code, cmds[0].Code...)
			push(core.JMP, core.Bcode(save(cmdStack.Children[2])+2))
			out.Code = append(out.Code, cmds[1].Code...)
			cmds = cmds[:0]
		case core.StackAnd, core.StackOr:
			cmd2Code(linker, cmdStack.Children[0], out)
			size := save(cmdStack.Children[1])
			logic := core.Bcode(core.JZE)
			if cmd.(*core.CmdBlock).ID == core.StackOr {
				logic = core.JNZ
			}
			push(core.DUP, logic, core.Bcode(size))
			out.Code = append(out.Code, cmds[0].Code...)
			cmds = cmds[:0]
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
					if embed.BCode.Code != nil {
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
			/*						if len(cmdVar.Indexes) > 0 {
										for i, ival := range cmdVar.Indexes {
															if typeValue == nil {
															return runtimeError(rt, cmdVar, ErrRuntime, `getVar.typeValue`)
														}
											cmd2Code(linker, ival.Cmd, out)
											if typeValue.Original == reflect.TypeOf(core.Struct{}) {
												//	typeValue = custom.Type.Custom.Types[index]
											} else {
												typeValue = typeValue.IndexOf
											}
										}
									}
									}*/
			/*if err = setVar(rt, cmdStack); err != nil {
				return err
			}
			lenStack++*/
		case core.StackIf:
			var (
				k, size int
				isElse  bool
			)
			lenIf := len(cmdStack.Children) >> 1
			for k = 0; k < lenIf; k++ {
				size += save(cmdStack.Children[k<<1]) + 2
				size += save(cmdStack.Children[(k<<1)+1]) + 2
			}
			// Calling else
			if /*k == lenIf &&*/ len(cmdStack.Children)&1 == 1 {
				size += save(cmdStack.Children[len(cmdStack.Children)-1]) + 2
				isElse = true
			}
			for k, code := range cmds {
				size -= len(code.Code) + 2
				out.Code = append(out.Code, code.Code...)
				if k&1 == 0 && k < len(cmds)-1 {
					more := 2
					if !isElse && k == len(cmds)-2 {
						more = 0
					}
					push(core.JZE, core.Bcode(len(cmds[k+1].Code)+more))
				} else if size > 0 {
					push(core.JMP, core.Bcode(size))
				}
			}
			cmds = cmds[:0]
		case core.StackWhile:
			//			push(core.Bcode((core.BlBreak|core.BlContinue)<<16|core.BLOCK), 0, 0)
			pos := len(out.Code)
			push(core.CYCLE)
			getPos(linker, cmdStack, out)
			cmd2Code(linker, cmdStack.Children[0], out)
			out.BlockFlags = core.BlContinue | core.BlBreak
			push(core.JZE, core.Bcode(save(cmdStack.Children[1])+2))
			blockStart := len(out.Code)
			out.Code = append(out.Code, cmds[0].Code...)
			push(core.JMP, core.Bcode(pos-len(out.Code)))
			out.Code[blockStart+1] = core.Bcode(len(out.Code) - blockStart) // set break of BLOCK
			out.Code[blockStart+2] = core.Bcode(pos - blockStart)           // set continue of BLOCK
			cmds = cmds[:0]
			//			out.Code[pos-1] = core.Bcode(len(out.Code) - pos) // set size of BLOCK
			//			linker.Blocks = linker.Blocks[:len(linker.Blocks)-1]
			/*					cycle := rt.Cycle
								for rt.Result == nil {
									if err = rt.runCmd(cmdStack.Children[0]); err != nil {
										return err
									}
									if rt.Stack[len(rt.Stack)-1].(bool) {
										rt.Stack = rt.Stack[:len(rt.Stack)-1]
										if err = rt.runCmd(cmdStack.Children[1]); err != nil {
											return err
										}
										if rt.Command == RcBreak {
											if rt.Catch == 0 {
												rt.Command = 0
											}
											break
										}
										if rt.Command == RcContinue {
											rt.Command = 0
										}
										cycle--
										if cycle == 0 {
											return runtimeError(rt, cmdStack, ErrCycle)
										}
										continue
									}
									break
								}*/
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
			//			push(core.Bcode((core.BlBreak|core.BlContinue)<<16|core.BLOCK), 0, 0)
			pos := len(out.Code)
			push(core.CYCLE)
			getPos(linker, cmdStack, out)
			push(core.GETVAR, core.Bcode(int(core.TYPEINT)<<16|bInfo.Vars[1]),
				(srcType<<16)|core.DUP, (srcType<<16)|core.LEN, core.LT)
			out.BlockFlags = core.BlContinue | core.BlBreak
			push(core.JZE, core.Bcode(save(cmdStack.Children[1])+12))
			push(core.GETVAR, core.Bcode(int(core.TYPEINT)<<16|bInfo.Vars[1])) // set index
			push(core.GETVAR, core.Bcode(int(srcType)<<16|indcur),             // get cur value
				core.Bcode(1<<16|core.INDEX), core.Bcode(int(srcType)<<16)|curType)
			push(core.SETVAR, core.Bcode(int(curType)<<16|bInfo.Vars[0]),
				core.Bcode(int(curType)<<16|core.ASSIGNPTR), core.Bcode(int(curType)<<16|core.POP))
			blockStart := len(out.Code)
			out.Code = append(out.Code, cmds[0].Code...)
			out.Code[blockStart+2] = core.Bcode(len(out.Code) - blockStart) // set continue of BLOCK
			push(core.Bcode(bInfo.Vars[1]<<16) | core.FORINC)
			push(core.JMP, core.Bcode(pos-len(out.Code)))
			cmds = cmds[:0]
			out.Code[blockStart+1] = core.Bcode(len(out.Code) - blockStart) // set break of BLOCK
			push(core.DELVARS)
			linker.Blocks = linker.Blocks[:len(linker.Blocks)-1]

			/*		if err = rt.runCmd(cmdStack.Children[0]); err != nil {
						return err
					}
					value := rt.Stack[len(rt.Stack)-1]
					rt.Stack = rt.Stack[:len(rt.Stack)-1]
					var index int64
					length := getLength(value)
					lenStack -= initVars(rt, cmdStack)
					if vars, err = rt.getVars(cmdStack); err != nil {
						return err
					}
					cycle := rt.Cycle
					for ; index < length; index++ {
						vars[0] = getIndex(value, index)
						vars[1] = index
						if err = rt.runCmd(cmdStack.Children[1]); err != nil {
							return err
						}
						if rt.Result != nil {
							break
						}
						if rt.Command == RcBreak {
							if rt.Catch == 0 {
								rt.Command = 0
							}
							break
						}
						if rt.Command == RcContinue {
							rt.Command = 0
						}
						length = getLength(value)
						if index > cycle {
							return runtimeError(rt, cmdStack, ErrCycle)
						}
					}
					deleteVars(rt)*/
		case core.StackBlock, core.StackDefault:
			/*			rt.Result = nil
						lenStack -= initVars(rt, cmdStack)*/
			initBlock(linker, cmdStack, out)
			for _, item := range cmdStack.Children {
				cmd2Code(linker, item, out)
				/*					if err = rt.runCmd(item); err != nil {
									return err
								}*/
				/*		if rt.Result != nil {
												if cmdStack.Result != nil && cmdStack.Parent == nil {
													rt.Stack = rt.Stack[:lenStack]
													rt.Stack = append(rt.Stack, rt.Result)
													lenStack++
													rt.Result = nil
												}
						    					break
												}
														if rt.Command != 0 {
															break
														}*/
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
			/*			if cmdStack.Children != nil {
							if err = rt.runCmd(cmdStack.Children[0]); err != nil {
								return err
							}
							if rt.Result == nil {
								rt.Result = rt.Stack[len(rt.Stack)-1]
							}
						} else { // return from the function without result value
							rt.Result = true
						}*/
			/*		case StackOptional: // assigns value if the variable has not yet been assigned as optional
						block := rt.Blocks[len(rt.Blocks)-1]
						var defined bool
						for _, v := range block.Optional {
							if cmdStack.ParCount == v {
								defined = true
								break
							}
						}
						if !defined {
							if err = rt.runCmd(cmdStack.Children[0]); err != nil {
								return err
							}
						}
					case StackLocal:
					case StackCallLocal:
						for i := 1; i < len(cmdStack.Children); i++ {
							if err = rt.runCmd(cmdStack.Children[i]); err != nil {
								return err
							}
						}
						if err = rt.runCmd(cmdStack.Children[0]); err != nil {
							return err
						}
						if rt.Command == RcLocal {
							rt.Stack = rt.Stack[:lenStack]
							rt.Stack = append(rt.Stack, rt.Result)
							lenStack++
							rt.Result = nil
							rt.Command = 0
						}
					case StackLocret:
						if cmdStack.Children != nil {
							if err = rt.runCmd(cmdStack.Children[0]); err != nil {
								return err
							}
							rt.Result = rt.Stack[len(rt.Stack)-1]
						} else { // return from the function without result value
							rt.Result = true
						}
						rt.Command = RcLocal
					case StackTry:
						for {
							if err = rt.runCmd(cmdStack.Children[0]); err != nil {
								if _, ok := err.(*RuntimeError); !ok {
									err = runtimeError(rt, cmdStack.Children[0], err)
								}
								rt.Stack = append(rt.Stack, err)
								if errCatch := rt.runCmd(cmdStack.Children[1]); errCatch != nil {
									err = errCatch
								}
								if rt.Catch == RcRecover || rt.Catch == RcRetry {
									rt.Command = 0
									err = nil
									if rt.Catch == RcRetry {
										rt.Catch = 0
										continue
									}
									rt.Catch = 0
								}
							}
							break
						}*/
		}
		//		rt.Stack = rt.Stack[:lenStack]
	}
}
