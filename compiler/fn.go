// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"reflect"
	///	"strings"

	"github.com/gentee/gentee/core"
)

func coFn(cmpl *compiler) error {
	token, err := checkNewType(cmpl)
	if err != nil {
		return err
	}
	pType := cmpl.unit.NewType(token, reflect.TypeOf(core.Fn{}), nil).(*core.TypeObject)
	pType.Func = &core.FnType{}
	cmpl.curType = pType
	return nil
}

func coFnEnd(cmpl *compiler) error {
	cmpl.curType = nil
	return nil
}

func coFnResult(cmpl *compiler) error {
	obj, err := getType(cmpl)
	if err != nil {
		return err
	}
	cmpl.curType.Func.Result = obj.(*core.TypeObject)
	return nil
}

func coFnType(cmpl *compiler) error {
	obj, err := getType(cmpl)
	if err != nil {
		return err
	}
	cmpl.curType.Func.Params = append(cmpl.curType.Func.Params, obj.(*core.TypeObject))
	return nil
}
