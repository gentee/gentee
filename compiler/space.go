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
	for key, indexes := range srcUnit.NCustom {
		var (
			ok   bool
			list []uint32
		)
		if list, ok = cmpl.unit.NCustom[key]; !ok {
			list = make([]uint32, 0, len(indexes))
		}
		for _, item := range indexes {
			if (item&core.NSImported) != 0 || (imported && (item&core.NSPub) == 0) {
				continue
			}
			if imported {
				item = (item & core.NSIndex) | core.NSImported
			}
			/*
				if ind, ok := cmpl.unit.NSpace[key]; ok {
					return cmpl.Error(ErrDupObject, cmpl.unit.GetObj(ind).GetName())
				}*/

			list = append(list, item)
		}
		cmpl.unit.NCustom[key] = list
	}

	return nil
}
