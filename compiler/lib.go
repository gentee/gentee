// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"math/rand"
	"time"
	"unicode"

	"github.com/gentee/gentee/core"
)

func isBoolResult(cmd core.ICmd) bool {
	return cmd.GetResult().GetName() == `bool`
}

func isIntResult(cmd core.ICmd) bool {
	return cmd.GetResult().GetName() == `int`
}

func isBaseResult(cmd core.ICmd) bool {
	switch cmd.GetResult().GetName() {
	case `int`, `bool`, `float`, `char`, `str`:
		return true
	}
	return false
}

func isCase(cmpl *compiler) bool {
	parent := cmpl.owners[len(cmpl.owners)-1]
	return parent.GetType() == core.CtStack && parent.(*core.CmdBlock).ID == core.StackCase
}

func isIndexResult(cmd core.ICmd) bool {
	return cmd.GetResult().IndexOf != nil
}

func isCapital(name string) bool {
	for _, ch := range name {
		if !unicode.IsUpper(ch) && ch != '_' && !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func isVariadic(obj core.IObject) bool {
	return (obj.GetType() == core.ObjFunc && obj.(*core.FuncObject).Block.Variadic) ||
		(obj.GetType() == core.ObjEmbedded && obj.(*core.EmbedObject).Variadic)
}

func randName() string {
	alpha := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := len(alpha)
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = alpha[rand.Intn(length)]
	}
	return string(b)
}
