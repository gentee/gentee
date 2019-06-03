// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"strings"

	"github.com/gentee/gentee/core"
)

func getLocalBlock(cmpl *compiler) *core.CmdBlock {
	for i := len(cmpl.owners) - 1; i >= 0; i-- {
		block := cmpl.owners[i].(*core.CmdBlock)
		if block.ID == core.StackBlock && block.Parent != nil && block.Parent.ID == core.StackLocal {
			return cmpl.owners[i].(*core.CmdBlock)
		}
	}
	return nil
}

func coLocalBack(cmpl *compiler) error {
	block := cmpl.curOwner()
	if block.Result != nil {
		if len(block.Children) == 0 {
			return cmpl.Error(ErrMustReturn)
		}
		last := block.Children[len(block.Children)-1]
		if last.GetType() != core.CtStack ||
			(last.(*core.CmdBlock).ID != core.StackReturn &&
				last.(*core.CmdBlock).ID != core.StackLocret) {
			return cmpl.Error(ErrMustReturn)
		}
	}
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	return nil
}

func coLocalName(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if isCapital(token) {
		return cmpl.Error(ErrCapitalLetters)
	}
	if strings.IndexRune(token, '.') >= 0 {
		return cmpl.Error(ErrIdent)
	}
	for _, owner := range cmpl.owners {
		block := owner.(*core.CmdBlock)
		if _, ok := block.LocalNames[token]; ok {
			return cmpl.Error(ErrLocalName, token)
		}
	}
	cmd := core.CmdBlock{ID: core.StackLocal, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	ownerBlock := cmpl.curOwner()
	if ownerBlock.LocalNames == nil {
		ownerBlock.LocalNames = map[string]int{}
	}
	cmdBlock := core.CmdBlock{ID: core.StackBlock, Parent: &cmd,
		CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmd.Children = append(cmd.Children, &cmdBlock)

	ownerBlock.LocalNames[token] = len(ownerBlock.Locals)
	ownerBlock.Locals = append(ownerBlock.Locals, &cmdBlock)

	cmpl.owners = append(cmpl.owners, &cmdBlock)

	return nil
}

func coLocalRetType(cmpl *compiler) error {
	obj, err := getType(cmpl)
	if err != nil {
		return err
	}
	cmpl.owners[len(cmpl.owners)-1].(*core.CmdBlock).Result = obj.(*core.TypeObject)
	return coLocalStart(cmpl)
}

func coLocalStart(cmpl *compiler) error {
	block := cmpl.curOwner()
	block.ParCount = len(block.Vars)
	if block.Variadic {
		return cmpl.Error(ErrLocalVariadic)
	}
	return nil
}

func getLocal(cmpl *compiler, name string, params []*core.TypeObject) (cmd core.ICmd) {
	for _, owner := range cmpl.owners {
		block := owner.(*core.CmdBlock)
		if ind, ok := block.LocalNames[name]; ok {
			local := block.Locals[ind].(*core.CmdBlock)
			if len(params) == local.ParCount {
				for i, par := range params {
					if par != local.Vars[i] {
						return nil
					}
				}
				return local
			}
		}
	}
	return nil
}

func coLocret(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackLocret, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coLocretBack(cmpl *compiler) error {
	owner := cmpl.curOwner()
	local := getLocalBlock(cmpl)
	if local == nil {
		return cmpl.Error(ErrLocalRet)
	}
	switch len(owner.Children) {
	case 0:
		if local.Result != nil {
			return cmpl.Error(ErrMustReturn)
		}
	case 1:
		if local.Result == nil {
			return cmpl.Error(ErrReturn)
		}
		if !isEqualTypes(local.Result, owner.Children[0].GetResult()) {
			return cmpl.Error(ErrReturnType)
		}
	default:
		return cmpl.Error(ErrCompiler, `coLocalReturn 1`)
	}
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	return nil
}
