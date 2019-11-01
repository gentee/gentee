// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"github.com/gentee/gentee/core"
)

func coTry(cmpl *compiler) error {
	cmd := &core.CmdBlock{ID: core.StackTry, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, cmd)
	cmpl.owners = append(cmpl.owners, cmd)
	cmdTry := core.CmdBlock{ID: core.StackBlock, Parent: cmd,
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmd.Children = append(cmd.Children, &cmdTry)
	cmpl.owners = append(cmpl.owners, &cmdTry)
	return nil
}

func coTryBack(cmpl *compiler) error {
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	cmd := cmpl.curOwner()
	if cmd.ID == core.StackTry {
		if len(cmd.Children) == 1 {
			cmpl.dynamic = &cmState{tkLCurly, cmCatch, nil, nil, 0}
		} else {
			cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		}
	}
	return nil
}

func coCatch(cmpl *compiler) error {
	token := getToken(cmpl.unit.Lexeme, cmpl.pos)
	if err := checkUsedName(cmpl, token); err != nil {
		return err
	}
	cmd := core.CmdBlock{ID: core.StackBlock, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		ParCount: 1, Vars: []*core.TypeObject{cmpl.unit.FindType(`error`).(*core.TypeObject)},
		VarNames: map[string]int{token: 0}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coCatchBack(cmpl *compiler) error {
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	return nil
}

func isInCatch(cmpl *compiler) bool {
	for _, item := range cmpl.owners {
		if item.GetType() == core.CtStack {
			parent := item.(*core.CmdBlock).Parent
			if parent != nil && parent.ID == core.StackTry && len(parent.Children) == 2 &&
				item == parent.Children[1] {
				return true
			}
		}
	}
	return false
}

func coRecover(cmpl *compiler) error {
	if !isInCatch(cmpl) {
		return cmpl.Error(ErrRecover)
	}
	appendCmd(cmpl, &core.CmdCommand{CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		ID: core.RcRecover})
	return nil
}

func coRetry(cmpl *compiler) error {
	if !isInCatch(cmpl) {
		return cmpl.Error(ErrRetry)
	}
	appendCmd(cmpl, &core.CmdCommand{CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		ID: core.RcRetry})
	return nil
}
