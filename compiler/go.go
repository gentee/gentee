// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"github.com/gentee/gentee/core"
)

type goStack struct {
	Exp        []core.ICmd
	ExpBuf     []ExpBuf
	LatestFunc int
	Name       string
	Params     []core.ICmd
}

func goExpPush(cmpl *compiler) string {
	name := `*` + randName()
	cmpl.goStack = append(cmpl.goStack, goStack{
		Exp:        append(make([]core.ICmd, 0, len(cmpl.exp)), cmpl.exp...),
		ExpBuf:     append(make([]ExpBuf, 0, len(cmpl.expbuf)), cmpl.expbuf...),
		LatestFunc: cmpl.curFunc,
		Name:       name,
	})
	cmpl.exp = cmpl.exp[:0]
	cmpl.expbuf = cmpl.expbuf[:0]
	return name
}

func goExpPop(cmpl *compiler) {
	stack := cmpl.goStack[len(cmpl.goStack)-1]
	cmpl.curFunc = stack.LatestFunc
	cmpl.exp = append(cmpl.exp[:0], stack.Exp...)
	cmpl.expbuf = append(cmpl.expbuf[:0], stack.ExpBuf...)
	cmpl.goStack = cmpl.goStack[:len(cmpl.goStack)-1]
}

func coGo(cmpl *compiler) error {
	newFunc(cmpl, goExpPush(cmpl))
	return nil
}

func coGoBack(cmpl *compiler) error {
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	threadFunc := cmpl.latestFunc()
	params := cmpl.goStack[len(cmpl.goStack)-1].Params
	goExpPop(cmpl)
	cmpl.dynamic = &cmState{tkToken, cmExp, nil, nil, 0}
	*cmpl.states = (*cmpl.states)[:len(*cmpl.states)-1]
	lp := cmpl.getLex()
	nextPos := cmpl.pos + 1
	if len(lp.Tokens) == nextPos || (lp.Tokens[nextPos].Type != tkLine &&
		lp.Tokens[nextPos].Type != tkRCurly) {
		return cmpl.ErrorPos(nextPos, ErrLineRCurly)
	}

	appendExp(cmpl, &core.CmdAnyFunc{CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		Children: params,
		Object:   threadFunc, IsThread: true, Result: cmpl.unit.FindType(`thread`).(*core.TypeObject)})
	return coExpEnd(cmpl)
}

func coGoParams(cmpl *compiler) error {
	cmpl.owners = append(cmpl.owners, cmpl.owners[len(cmpl.owners)-2])
	cmpl.optionals = append(cmpl.optionals, &optInfo{})
	return appendExpBuf(cmpl, tkCallFunc)
}
