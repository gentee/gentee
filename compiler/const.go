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
		},
		Redefined: false,
		Exp:       owner.Children[0],
		Return:    owner.Children[0].GetResult(),
		Iota:      core.NotIota,
	}
	if findObj(cmpl, cmpl.curConst, core.ObjConst) {
		return cmpl.ErrorPos(cmpl.pos-1, ErrConstDef, cmpl.curConst)
	}

	cmpl.unit.Objects = append(cmpl.unit.Objects, constObj)
	if curName := cmpl.unit.Names[cmpl.curConst]; curName == nil {
		cmpl.unit.Names[cmpl.curConst] = constObj
	} else {
		curName.SetNext(constObj)
	}
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
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
		},
		Redefined: false,
		Exp:       cmpl.expConst,
		Return:    cmpl.expConst.GetResult(),
		Iota:      cmpl.curIota,
	}
	cmpl.curIota++
	cmpl.unit.Objects = append(cmpl.unit.Objects, constObj)
	if findObj(cmpl, cmpl.curConst, core.ObjConst) {
		return cmpl.Error(ErrConstDef, cmpl.curConst)
	}
	if curName := cmpl.unit.Names[cmpl.curConst]; curName == nil {
		cmpl.unit.Names[cmpl.curConst] = constObj
	} else {
		curName.SetNext(constObj)
	}

	return nil
}
