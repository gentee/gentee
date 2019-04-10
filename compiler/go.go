// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	//	"fmt"
	//	"reflect"
	///	"strings"

	"github.com/gentee/gentee/core"
)

func coGo(cmpl *compiler) error {
	/*	token := getToken(cmpl.getLex(), cmpl.pos)
		if isCapital(token) {
			return cmpl.Error(ErrCapitalLetters)
		}
		if strings.IndexRune(token, '.') >= 0 {
			return cmpl.Error(ErrIdent)
		}*/
	cmd := cmpl.curOwner()
	newFunc(cmpl, randName())
	icmd := &core.CmdAnyFunc{CmdCommon: core.CmdCommon{TokenID: uint32(cmpl.pos)},
		Object: cmpl.latestFunc(), IsThread: true}
	cmd.Children = append(cmd.Children, icmd)
	return nil
}

func coGoBack(cmpl *compiler) error {
	cmpl.owners = cmpl.owners[:len(cmpl.owners)-1]
	return nil
}
