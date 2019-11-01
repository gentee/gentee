// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"reflect"
	"strings"

	"github.com/gentee/gentee/core"
)

func checkNewType(cmpl *compiler) (string, error) {
	token := getToken(cmpl.unit.Lexeme, cmpl.pos)
	obj, _ := getType(cmpl)
	if obj != nil {
		return ``, cmpl.Error(ErrTypeExists, token)
	}
	if isCapital(token) {
		return ``, cmpl.Error(ErrCapitalLetters)
	}
	if strings.IndexRune(token, '.') >= 0 {
		return ``, cmpl.Error(ErrIdent)
	}
	return token, nil
}

func coStruct(cmpl *compiler) error {
	token, err := checkNewType(cmpl)
	if err != nil {
		return err
	}
	pType := cmpl.unit.NewType(token, reflect.TypeOf(core.Struct{}), nil).(*core.TypeObject)
	pType.Custom = &core.StructType{
		Fields: make(map[string]int64),
		Types:  make([]*core.TypeObject, 0),
	}
	cmpl.curType = pType
	return nil
}

func coStructEnd(cmpl *compiler) error {
	cmpl.curType = nil
	return nil
}

func coStructLine(cmpl *compiler) error {
	if len(cmpl.curType.Custom.Fields) != len(cmpl.curType.Custom.Types) {
		return cmpl.Error(ErrName)
	}
	return nil
}

func coStructType(cmpl *compiler) error {
	obj, err := getType(cmpl)
	if err != nil {
		return err
	}
	cmpl.curType.Custom.Types = append(cmpl.curType.Custom.Types, obj.(*core.TypeObject))
	return nil
}

func coStructName(cmpl *compiler) error {
	token := getToken(cmpl.unit.Lexeme, cmpl.pos)
	if strings.IndexRune(token, '.') >= 0 {
		return cmpl.Error(ErrIdent)
	}
	if len(cmpl.curType.Custom.Fields) == len(cmpl.curType.Custom.Types) {
		return cmpl.Error(ErrLineRCurly)
	}
	if _, ok := cmpl.curType.Custom.Fields[token]; ok {
		return cmpl.Error(ErrStructField, token)
	}
	if obj, _ := autoType(cmpl, token); obj != nil {
		return cmpl.Error(ErrName)
	}
	cmpl.curType.Custom.Fields[token] = int64(len(cmpl.curType.Custom.Types) - 1)
	return nil
}

func structIndex(cmpl *compiler, typeObj *core.TypeObject, field string) (int64, *core.TypeObject, error) {
	var (
		indField int64
		ok       bool
	)
	if typeObj.Custom == nil {
		return 0, nil, cmpl.ErrorPos(cmpl.pos-1, ErrStructType, typeObj.GetName())
	}
	if indField, ok = typeObj.Custom.Fields[field]; !ok {
		return 0, nil, cmpl.ErrorPos(cmpl.pos-1, ErrStruct, typeObj.GetName(), field)
	}
	return indField, typeObj.Custom.Types[indField], nil
}
