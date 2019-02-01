// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"github.com/gentee/gentee/core"
)

func coConst(cmpl *compiler) error {
	token := getToken(cmpl.getLex(), cmpl.pos)
	if !isCapital(token) {
		return cmpl.Error(ErrConstName)
	}
	cmpl.curConst = token
	return nil
}

func coConstEnum(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackBlock, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmpl.owners = append(cmpl.owners, &cmd)
	cmpl.curIota = 0
	cmpl.expConst = nil
	return nil
}

func coConstEnumBack(cmpl *compiler) error {
	if cmpl.expConst == nil {
		owner := cmpl.curOwner()
		cmpl.expConst = owner.Children[0]
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
		cmpl.dynamic = &cmState{tkToken, cmConstListStart, nil, nil, 0}
	} else { // const finishes
		cmpl.curIota = core.NotIota
	}
	return nil
}

func coConstExp(cmpl *compiler) error {
	coExpStart(cmpl)
	cmd := core.CmdBlock{ID: core.StackBlock, CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)}}
	cmpl.owners = append(cmpl.owners, &cmd)
	return nil
}

func coConstExpBack(cmpl *compiler) error {
	owner := cmpl.curOwner()
	constObj := &core.ConstObject{
		Object: core.Object{
			Name:  cmpl.curConst,
			LexID: len(cmpl.unit.Lexeme) - 1,
			Unit:  cmpl.unit,
			Pub:   cmpl.unit.Pub != 0,
		},
		Redefined: false,
		Exp:       owner.Children[0],
		Return:    owner.Children[0].GetResult(),
		Iota:      core.NotIota,
	}
	if cmpl.unit.FindConst(cmpl.curConst) != nil {
		return cmpl.ErrorPos(cmpl.pos-1, ErrConstDef, cmpl.curConst)
	}

	cmpl.appendObj(constObj)
	cmpl.unit.AddConst(cmpl.curConst)

	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	return nil
}

func coConstBack(cmpl *compiler) error {
	if cmpl.unit.Pub == core.PubOne {
		cmpl.unit.Pub = 0
	}
	return nil
}

func coConstList(cmpl *compiler) error {
	if err := coConst(cmpl); err != nil {
		return err
	}
	constObj := &core.ConstObject{
		Object: core.Object{
			Name:  cmpl.curConst,
			LexID: len(cmpl.unit.Lexeme) - 1,
			Unit:  cmpl.unit,
			Pub:   cmpl.unit.Pub != 0,
		},
		Redefined: false,
		Exp:       cmpl.expConst,
		Return:    cmpl.expConst.GetResult(),
		Iota:      cmpl.curIota,
	}
	cmpl.curIota++
	cmpl.appendObj(constObj)
	if cmpl.unit.FindConst(cmpl.curConst) != nil {
		return cmpl.Error(ErrConstDef, cmpl.curConst)
	}
	cmpl.unit.AddConst(cmpl.curConst)

	return nil
}
