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

func Link(ws *core.Workspace, unitID int) (*core.Exec, error) {
	var exec *core.Exec
	if unitID < 0 || unitID >= len(ws.Units) {
		return nil, fmt.Errorf(errText[ErrLinkIndex], unitID)
	}
	unit := ws.Units[unitID]
	if unit.RunID == core.Undefined {
		return nil, nil
	}
	bcode := genBytecode(ws, uint16(unit.RunID))
	exec = &core.Exec{
		Code:  append([]uint16{}, bcode.Code...),
		Funcs: make(map[uint16]uint32),
	}
	for ikey := range bcode.Used {
		exec.Funcs[ikey] = uint32(len(exec.Code))
		exec.Code = append(exec.Code, ws.Objects[ikey].GetCode().Code...)
	}
	fmt.Println(`USED`, exec.Funcs, exec.Code)
	return exec, nil
}

func copyUsed(src, dest *core.Bytecode) {
	if src.Used == nil {
		return
	}
	if dest.Used == nil {
		dest.Used = make(map[uint16]byte)
	}
	for ikey := range src.Used {
		dest.Used[ikey] = 1
	}
}

func cmd2Code(cmd core.ICmd, out *core.Bytecode) {

	var cmds []core.Bytecode

	save := func(icmd core.ICmd) int {
		code := core.Bytecode{
			Code: make([]uint16, 0, 16),
		}
		cmd2Code(icmd, &code)
		copyUsed(&code, out)
		cmds = append(cmds, code)
		return len(code.Code)
	}

	push := func(pars ...uint16) {
		out.Code = append(out.Code, pars...)
	}
	callFunc := func() {
		obj := cmd.GetObject()
		switch obj.GetType() {
		case core.ObjEmbedded:
			if obj.(*core.EmbedObject).BCode.Code != nil {
				push(obj.(*core.EmbedObject).BCode.Code...)
			}
		case core.ObjFunc:
			id := obj.(*core.FuncObject).ObjID
			push(core.CALLBYID, id)
			if out.Used == nil {
				out.Used = make(map[uint16]byte)
			}
			if out.Used[id] == 0 {
				genBytecode(obj.(*core.FuncObject).Unit.VM, id)
				copyUsed(&obj.(*core.FuncObject).BCode, out)
				out.Used[id] = 1
			}
		}
	}
	switch cmd.GetType() {
	/*	case CtCommand:
		v := cmd.(*CmdCommand).ID
		switch v {
		case RcRecover, RcRetry:
			rt.Catch = v
			rt.Command = RcBreak
		default:
			rt.Command = v
		}*/
	case core.CtFunc:
		anyFunc := cmd.(*core.CmdAnyFunc)
		for _, param := range anyFunc.Children {
			cmd2Code(param, out)
		}
		callFunc()
	case core.CtBinary:
		cmd2Code(cmd.(*core.CmdBinary).Left, out)
		cmd2Code(cmd.(*core.CmdBinary).Right, out)
		callFunc()
	case core.CtUnary:
		cmd2Code(cmd.(*core.CmdUnary).Operand, out)
		callFunc()
	case core.CtValue:
		switch v := cmd.(*core.CmdValue).Value.(type) {
		case int64:
			if v <= math.MaxInt16 && v >= math.MinInt16 {
				push(core.PUSH16, uint16(v))
			} else if v <= math.MaxInt32 && v >= math.MinInt32 {
				u32 := uint32(v)
				push(core.PUSH32, uint16(u32>>16), uint16(u32&0xffff))
			} else {
				u64 := uint64(v)
				push(core.PUSH64, uint16(u64>>48), uint16((u64>>32)&0xffff),
					uint16((u64>>16)&0xffff), uint16(u64&0xffff))
			}
		case bool:
			var val uint16
			if v {
				val = 1
			}
			push(core.PUSH16, val)
		case rune:
			push(core.PUSH16, uint16(v))
		}
		//rt.Stack = append(rt.Stack, cmd.(*CmdValue).Value)
		/*			case CtConst:
						if err = rt.GetConst(cmd); err != nil {
							return err
						}
					case CtVar:
						if err = getVar(rt, cmd.(*CmdVar)); err != nil {
							return err
						}*/
	case core.CtStack:
		cmdStack := cmd.(*core.CmdBlock)
		//		lenStack := len(rt.Stack)
		switch cmd.(*core.CmdBlock).ID {
		/*		case StackSwitch:
					if err = rt.runCmd(cmdStack.Children[0]); err != nil {
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
					}
				case StackNew:
					switch cmd.(*CmdBlock).Result.Original {
					case reflect.TypeOf(Array{}):
						parr := NewArray()
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
						}
					case reflect.TypeOf(Set{}):
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
						rt.Stack[lenStack] = pbuf
					case reflect.TypeOf(Map{}):
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
						}
					default:
						return runtimeError(rt, cmd, ErrRuntime, `init arr`)
					}
					lenStack++
				case StackInit:
					cmdVar := cmdStack.Children[0].(*CmdVar)
					if vars, err = rt.getVars(cmdVar.Block); err != nil {
						return err
					}
					if err = rt.runCmd(cmdStack.Children[1]); err != nil {
						return err
					}
					vars[cmdVar.Index] = rt.Stack[len(rt.Stack)-1]
				case StackQuestion:
					if err = rt.runCmd(cmdStack.Children[0]); err != nil {
						return err
					}
					iExp := 2
					if rt.Stack[len(rt.Stack)-1].(bool) {
						iExp = 1
					}
					if err = rt.runCmd(cmdStack.Children[iExp]); err != nil {
						return err
					}
					rt.Stack[lenStack] = rt.Stack[len(rt.Stack)-1]
					lenStack++
				case StackAnd, StackOr:
					if err = rt.runCmd(cmdStack.Children[0]); err != nil {
						return err
					}
					if (rt.Stack[len(rt.Stack)-1].(bool) && cmd.(*CmdBlock).ID == StackAnd) ||
						(!rt.Stack[len(rt.Stack)-1].(bool) && cmd.(*CmdBlock).ID == StackOr) {
						if err = rt.runCmd(cmdStack.Children[1]); err != nil {
							return err
						}
					}
					rt.Stack[lenStack] = rt.Stack[len(rt.Stack)-1]
					lenStack++
				case StackIncDec:
					cmdVar := cmdStack.Children[0].(*CmdVar)
					if _, err = rt.getVars(cmdVar.Block); err != nil {
						return err
					}
					var post bool
					if err = getVar(rt, cmdVar); err != nil {
						return err
					}
					val := rt.Stack[len(rt.Stack)-1].(int64)
					rt.Stack = rt.Stack[:len(rt.Stack)-1]
					shift := int64(cmdStack.ParCount)
					if (shift & 0x1) == 0 {
						post = true
						shift /= 2
					}
					if err = setVar(rt, &CmdBlock{Children: []ICmd{
						cmdVar, &CmdValue{Value: val + shift},
					}, Object: rt.VM.StdLib().FindObj(DefAssignIntInt)}); err != nil {
						return err
					}
					rt.Stack = rt.Stack[:len(rt.Stack)-1]
					if !post {
						val += shift
					}
					rt.Stack = append(rt.Stack, val)
					lenStack++
				case StackAssign:
					if err = setVar(rt, cmdStack); err != nil {
						return err
					}
					lenStack++*/
		case core.StackIf:
			var k, size int
			lenIf := len(cmdStack.Children) >> 1
			for k = 0; k < lenIf; k++ {
				size += save(cmdStack.Children[k<<1]) + 2
				size += save(cmdStack.Children[(k<<1)+1]) + 2
			}
			// Calling else
			if k == lenIf && len(cmdStack.Children)&1 == 1 {
				size += save(cmdStack.Children[len(cmdStack.Children)-1]) + 2
			}
			for k, code := range cmds {
				size -= len(code.Code) - 2
				out.Code = append(out.Code, code.Code...)
				if k&1 == 0 && k < len(cmds)-1 {
					push(core.JZE, uint16(len(cmds[k+1].Code)+2))
				} else if size > 0 {
					push(core.JMP, uint16(size))
				}
			}
			cmds = cmds[:]
			/*				case StackWhile:
								cycle := rt.Cycle
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
								}
							case StackFor:
								if err = rt.runCmd(cmdStack.Children[0]); err != nil {
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
			for _, item := range cmdStack.Children {
				cmd2Code(item, out)
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
			//deleteVars(rt)
		case core.StackReturn:
			if cmdStack.Children != nil {
				cmd2Code(cmdStack.Children[0], out)
				retType := uint16(core.STACKANY)
				switch cmdStack.Children[0].GetResult().Original {
				case reflect.TypeOf(int64(0)):
					retType = core.STACKINT
				case reflect.TypeOf(true):
					retType = core.STACKBOOL
				case reflect.TypeOf(float64(0.0)):
					retType = core.STACKFLOAT
				case reflect.TypeOf('a'):
					retType = core.STACKCHAR
				case reflect.TypeOf(``):
					retType = core.STACKSTR
				}
				push(core.RET, retType)
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

func genBytecode(ws *core.Workspace, idObj uint16) *core.Bytecode {
	bcode := ws.Objects[idObj].GetCode()
	if ws.Objects[idObj].GetType() == core.ObjType {
		return nil
	}
	if bcode.Code != nil {
		return bcode
	}
	bcode.Code = make([]uint16, 0, 64)
	cmd2Code(&ws.Objects[idObj].(*core.FuncObject).Block, bcode)
	bcode.Code = append(bcode.Code, core.END)
	//	fmt.Println(`CODE`, bcode.Code)
	return bcode
}
