// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"reflect"

	"github.com/gentee/gentee/core"
)

func coInitStart(cmpl *compiler) error {
	cmpl.inits++
	cmd := core.CmdBlock{ID: core.StackNew, CmdCommon: core.CmdCommon{
		TokenID: uint32(cmpl.pos)}, Result: cmpl.curType}
	appendCmd(cmpl, &cmd)
	cmpl.owners = append(cmpl.owners, &cmd)
	if (cmpl.curType.Original != reflect.TypeOf(core.Array{}) &&
		cmpl.curType.Original != reflect.TypeOf(core.Buffer{}) &&
		cmpl.curType.Original != reflect.TypeOf(core.Map{}) &&
		cmpl.curType.Original != reflect.TypeOf(core.Struct{})) {
		return cmpl.Error(ErrWrongType, cmpl.curType.GetName())
	}
	cmpl.curType = cmpl.curType.IndexOf
	return nil
}

func coInitEnd(cmpl *compiler) error {

	ownerType := cmpl.curOwner().GetResult()
	if ownerType.Original == reflect.TypeOf(core.Map{}) {
		if err := initMapEnd(cmpl); err != nil {
			return err
		}
	}
	if ownerType.Original == reflect.TypeOf(core.Struct{}) {
		if err := initStructEnd(cmpl); err != nil {
			return err
		}
	}
	for _, item := range cmpl.curOwner().Children {
		if ownerType.Original == reflect.TypeOf(core.Array{}) {
			if !isEqualTypes(item.GetResult(), cmpl.curType) {
				return cmpl.ErrorPos(item.GetToken(), ErrWrongType, cmpl.curType.GetName())
			}
		} else if ownerType.Original == reflect.TypeOf(core.Buffer{}) {
			v := map[string]bool{`buf`: true, `char`: true, `int`: true, `str`: true}
			if !v[item.GetResult().GetName()] {
				return cmpl.ErrorPos(item.GetToken(), ErrWrongType, `int, buf, char, str`)
			}
		} else if ownerType.Original == reflect.TypeOf(core.Map{}) ||
			ownerType.Original == reflect.TypeOf(core.Struct{}) {
			if item.GetResult().Original != reflect.TypeOf(core.KeyValue{}) {
				return cmpl.ErrorPos(item.GetToken(), ErrNotKeyValue)
			}
			if item.(*core.CmdBinary).Right == nil {
				return cmpl.ErrorPos(item.GetToken(), ErrValue)
			}
			if ownerType.Original == reflect.TypeOf(core.Map{}) {
				if !isEqualTypes(ownerType.IndexOf, item.(*core.CmdBinary).Right.GetResult()) {
					return cmpl.ErrorPos(item.(*core.CmdBinary).Right.GetToken(),
						ErrWrongType, ownerType.IndexOf.GetName())
				}
			} else {
				ind := item.(*core.CmdBinary).Left.(*core.CmdValue).Value.(int64)
				if ownerType.Custom.Types[ind] != item.(*core.CmdBinary).Right.GetResult() {
					return cmpl.ErrorPos(item.(*core.CmdBinary).Right.GetToken(),
						ErrWrongType, ownerType.Custom.Types[ind].GetName())
				}
			}
		}
	}
	cmpl.curType = cmpl.curOwner().Result
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	if (*cmpl.states)[len(*cmpl.states)-1].State != cmInit {
		cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	}
	cmpl.inits--
	return nil
}

func coInitNext(cmpl *compiler) error {
	/*	if cmpl.curType.IndexOf == nil {
		cmpl.dynamic = &cmState{tkComma, cmExp, nil, nil, cfStay}
		return nil
	}*/
	if len(cmpl.owners[len(cmpl.owners)-1].(*core.CmdBlock).Children) == 0 {
		return cmpl.Error(ErrValue)
	}
	lp := cmpl.getLex()
	if lp.Tokens[cmpl.pos].Type == tkComma {
		for i := cmpl.pos - 1; i > 0; i-- {
			token := lp.Tokens[i]
			if token.Type == tkLine {
				continue
			}
			if token.Type == tkComma {
				return cmpl.Error(ErrValue)
			}
			break
		}
	}
	/*	for i := cmpl.pos + 1; i < len(lp.Tokens); i++ {
		token := lp.Tokens[i]
		if token.Type == tkLine {
			continue
		}
		if token.Type == tkComma {
			return cmpl.ErrorPos(i, ErrValue)
		}
		break
	}*/
	if cmpl.curOwner().GetResult().Original == reflect.TypeOf(core.Map{}) {
		if err := initMapEnd(cmpl); err != nil {
			return err
		}
	}
	if cmpl.curOwner().GetResult().Original == reflect.TypeOf(core.Struct{}) {
		if err := initStructEnd(cmpl); err != nil {
			return err
		}
	}
	return nil
}

func initMapEnd(cmpl *compiler) error {
	obj := cmpl.vm.StdLib().Names[core.DefNewKeyValue]
	block := cmpl.curOwner()
	for i := 0; i < len(block.Children); i++ {
		item := block.Children[i]
		if item.GetResult().Original == reflect.TypeOf(core.KeyValue{}) {
			if item.(*core.CmdBinary).Right == nil {
				if i+1 < len(block.Children) {
					block.Children[i].(*core.CmdBinary).Right = block.Children[i+1]
					block.Children = append(block.Children[:i+1], block.Children[i+2:]...)
					i--
				} else {
					return cmpl.ErrorPos(cmpl.pos, ErrValue)
				}
			}
			continue
		}
		if item.GetResult().Original != reflect.TypeOf(``) {
			return cmpl.ErrorPos(item.GetToken(), ErrWrongType, `str`)
		}
		cmd := &core.CmdBinary{CmdCommon: core.CmdCommon{TokenID: uint32(item.GetToken())},
			Object: obj, Result: obj.Result(), Left: item, Right: nil}
		block.Children[i] = cmd
		if i+1 < len(block.Children) {
			return cmpl.ErrorPos(block.Children[i+1].GetToken(), ErrNotKeyValue)
		}
	}
	return nil
}

func initStructEnd(cmpl *compiler) error {
	var (
		fieldName string
	)
	obj := cmpl.vm.StdLib().Names[core.DefNewKeyValue]
	block := cmpl.curOwner()
	for i := 0; i < len(block.Children); i++ {
		item := block.Children[i]
		if item.GetResult().Original == reflect.TypeOf(core.KeyValue{}) {
			if item.(*core.CmdBinary).Right == nil {
				if i+1 < len(block.Children) {
					block.Children[i].(*core.CmdBinary).Right = block.Children[i+1]
					block.Children = append(block.Children[:i+1], block.Children[i+2:]...)
					i--
				} else {
					return cmpl.ErrorPos(cmpl.pos, ErrValue)
				}
			}
			continue
		}
		if item.GetResult().Original != reflect.TypeOf(``) || item.GetType() != core.CtValue {
			return cmpl.ErrorPos(item.GetToken(), ErrInitField)
		}
		fieldName = item.(*core.CmdValue).Value.(string)
		var (
			ind int64
			ok  bool
		)
		if ind, ok = block.GetResult().Custom.Fields[fieldName]; !ok {
			return cmpl.ErrorPos(cmpl.pos-1, ErrWrongField, fieldName, block.GetResult().GetName())
		}
		item.(*core.CmdValue).Value = ind

		cmd := &core.CmdBinary{CmdCommon: core.CmdCommon{TokenID: uint32(item.GetToken())},
			Object: obj, Result: obj.Result(), Left: item, Right: nil}

		cmpl.curType = block.GetResult().Custom.Types[ind]
		block.Children[i] = cmd
		if i+1 < len(block.Children) {
			return cmpl.ErrorPos(block.Children[i+1].GetToken(), ErrNotKeyValue)
		}
	}
	return nil
}

func coInitKey(cmpl *compiler) error {
	switch cmpl.curOwner().GetResult().Original {
	case reflect.TypeOf(core.Struct{}):
		if err := initStructEnd(cmpl); err != nil {
			return err
		}
	case reflect.TypeOf(core.Map{}):
		if err := initMapEnd(cmpl); err != nil {
			return err
		}
	default:
		return cmpl.ErrorPos(cmpl.pos, ErrKeyValue, cmpl.curType.GetName())
	}
	return nil
}
