// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"github.com/gentee/gentee/core"
)

func coReturn(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackReturn, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coReturnBack(cmpl *compiler) error {
	owner := cmpl.curOwner()
	funcObj := cmpl.unit.Objects[len(cmpl.unit.Objects)-1].(*core.FuncObject)
	switch len(owner.Children) {
	case 0:
		if funcObj.Block.Result != nil {
			return cmpl.Error(ErrMustReturn)
		}
	case 1:
		if funcObj.Block.Result == nil {
			return cmpl.Error(ErrReturn)
		}
		if !isEqualTypes(funcObj.Block.Result, owner.Children[0].GetResult()) {
			return cmpl.Error(ErrReturnType)
		}
	default:
		return cmpl.Error(ErrCompiler, `coReturn 1`)
	}
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	return nil
}

func coIf(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackIf, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coIfBack(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID == core.StackIf {
		if len(cmd.Children) == 1 {
			if !isBoolResult(cmd.Children[0]) {
				cmpl.pos = cmd.Children[0].GetToken()
				return cmpl.Error(ErrBoolExp)
			}
			cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
			cmd.Children = append(cmd.Children, &cmdIf)
			cmpl.owners = append(cmpl.owners, &cmdIf)
			cmpl.dynamic = &cmState{tkLCurly, cmLCurly, nil, nil, 0}
		} else {
			cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		}
	} else {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		cmpl.dynamic = &cmState{tkLCurly, cmElseIf, nil, nil, 0}
	}
	return nil
}

func coElse(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID != core.StackIf {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		cmd = cmpl.curOwner()
	}
	cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmd.Children = append(cmd.Children, &cmdIf)
	cmpl.owners = append(cmpl.owners, &cmdIf)
	return nil
}

func coElif(cmpl *compiler) error {
	coExpStart(cmpl)
	return nil
}

func coElifBack(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID == core.StackIf {
		if !isBoolResult(cmd.Children[len(cmd.Children)-1]) {
			cmpl.pos = cmd.Children[len(cmd.Children)-1].GetToken()
			return cmpl.Error(ErrBoolExp)
		}
		cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
			CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
		cmd.Children = append(cmd.Children, &cmdIf)
		cmpl.owners = append(cmpl.owners, &cmdIf)
		cmpl.dynamic = &cmState{tkLCurly, cmLCurly, nil, nil, 0}
	} else {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	}
	return nil
}

func coIfEnd(cmpl *compiler) error {
	if cmpl.curOwner().ID != core.StackIf {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	}
	return nil
}

func coWhile(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackWhile, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coWhileBack(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID == core.StackWhile {
		if len(cmd.Children) == 1 {
			if !isBoolResult(cmd.Children[0]) {
				cmpl.pos = cmd.Children[0].GetToken()
				return cmpl.Error(ErrBoolExp)
			}
			cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
			cmd.Children = append(cmd.Children, &cmdIf)
			cmpl.owners = append(cmpl.owners, &cmdIf)
			cmpl.dynamic = &cmState{tkLCurly, cmLCurly, nil, nil, 0}
		}
	} else {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-2]
	}
	return nil
}

func coFor(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackFor, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	lp := cmpl.unit.Lexeme[0]
	cmpl.newPos = cmpl.pos + 1
	if lp.Tokens[cmpl.newPos].Type != tkIdent {
		return cmpl.ErrorPos(cmpl.newPos, ErrName)
	}
	cmpl.curType = nil
	cmpl.pos = cmpl.newPos
	if err := coVar(cmpl); err != nil {
		return err
	}
	cmpl.newPos++
	if lp.Tokens[cmpl.newPos].Type == tkComma {
		cmpl.newPos++
		if lp.Tokens[cmpl.newPos].Type != tkIdent {
			return cmpl.ErrorPos(cmpl.newPos, ErrName)
		}
		cmpl.pos = cmpl.newPos
		if err := coVar(cmpl); err != nil {
			return err
		}
		cmpl.newPos++
	} else {
		if err := coVarToken(cmpl, randName()); err != nil {
			return err
		}
	}

	if lp.Tokens[cmpl.newPos].Type != tkIn {
		return cmpl.ErrorPos(cmpl.newPos, ErrForIn)
	}
	return nil
}

func coForBack(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID == core.StackFor {
		if len(cmd.Children) == 1 {
			if !isIndexResult(cmd.Children[0]) {
				return cmpl.ErrorPos(cmd.Children[0].GetToken(), ErrSupportIndex,
					cmd.Children[0].GetResult().GetName())
			}
			cmd.Vars[0] = cmd.Children[0].GetResult().IndexOf
			cmd.Vars[1] = cmpl.vm.StdLib().Names[`int`].(*core.TypeObject)
			cmdFor := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
			cmd.Children = append(cmd.Children, &cmdFor)
			cmpl.owners = append(cmpl.owners, &cmdFor)
			cmpl.dynamic = &cmState{tkLCurly, cmLCurly, nil, nil, 0}
		}
	} else {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-2]
	}
	return nil
}

func coSwitch(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackSwitch, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coSwitchBack(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID == core.StackSwitch {
		if len(cmd.Children) == 1 {
			if !isBaseResult(cmd.Children[0]) {
				return cmpl.ErrorPos(cmd.Children[0].GetToken(), ErrSwitchType,
					cmd.Children[0].GetResult().GetName())
			}
			cmpl.dynamic = &cmState{tkCase, cmCaseMust, nil, nil, 0}
		} else {
			cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		}
	}
	return nil
}

func coCase(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackCase, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coCaseBack(cmpl *compiler) error {
	cmd := cmpl.curOwner()
	if cmd.ID == core.StackCase {
		if len(cmd.Children) >= 1 {
			switchType := cmd.Parent.Children[0].GetResult()
			for _, cmdExp := range cmd.Children {
				if switchType != cmdExp.GetResult() {
					return cmpl.ErrorPos(cmdExp.GetToken(), ErrWrongType,
						switchType.GetName())

				}
			}
			cmdIf := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
				CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
			cmd.Children = append(cmd.Children, &cmdIf)
			cmpl.owners = append(cmpl.owners, &cmdIf)
			cmpl.dynamic = &cmState{tkLCurly, cmLCurly, nil, nil, 0}
		}
	} else {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-2]
	}
	return nil
}

func coDefault(cmpl *compiler) error {
	cmd := core.CmdBlock{ID: core.StackDefault, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coDefaultBack(cmpl *compiler) error {
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	return nil
}
