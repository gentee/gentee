// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"math/rand"
	"strconv"
	"strings"
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

func randName() string {
	alpha := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := len(alpha)
	b := make([]rune, 16)
	for i := range b {
		b[i] = alpha[rand.Intn(length)]
	}
	return string(b)
}

func unNewLine(in string) (string, error) {
	return strconv.Unquote(`"` + strings.Replace(strings.Replace(in, "\n", `\n`, -1),
		"\r", `\r`, -1) + `"`)
}

func (cmpl *compiler) getIntType() *core.TypeObject {
	return cmpl.unit.FindType(`int`).(*core.TypeObject)
}

func (cmpl *compiler) getStrType() *core.TypeObject {
	return cmpl.unit.FindType(`str`).(*core.TypeObject)
}

func (cmpl *compiler) copyNameSpace(srcUnit *core.Unit, imported bool) error {
	for key, item := range srcUnit.NameSpace {
		if (item&core.NSImported) != 0 || (imported && (item&core.NSPub) == 0) {
			continue
		}
		index := cmpl.unit.GetObj(item).GetUnitIndex()
		if _, ok := cmpl.unit.Included[index]; ok {
			continue
		}
		if ind, ok := cmpl.unit.NameSpace[key]; ok {
			return cmpl.Error(ErrDupObject, cmpl.unit.GetObj(ind).GetName())
		}
		if imported {
			item = (item & core.NSIndex) | core.NSImported
		}
		cmpl.unit.NameSpace[key] = item
	}
	for index, itype := range srcUnit.Included {
		if itype {
			continue
		}
		if _, ok := cmpl.unit.Included[index]; !ok {
			cmpl.unit.Included[index] = itype
		}
	}
	return nil
}
