// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"github.com/gentee/gentee/core"
)

func (cmpl *compiler) CopyNameSpace(srcUnit *core.Unit, imported bool) error {
	for key, item := range srcUnit.NSpace {
		if (item&core.NSImported) != 0 || (imported && (item&core.NSPub) == 0) {
			continue
		}
		if ind, ok := cmpl.unit.NSpace[key]; ok {
			return cmpl.Error(ErrDupObject, cmpl.unit.GetObj(ind).GetName())
		}
		if imported {
			item = (item & core.NSIndex) | core.NSImported
		}
		cmpl.unit.NSpace[key] = item
	}

	return nil
}
