// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitRange appends stdlib range functions to the virtual machine
func InitRange(ws *core.Workspace) {
	for _, item := range []interface{}{
		NewRange, // binary ..
	} {
		ws.StdLib().NewEmbed(item)
	}
}

// NewRange adds two rune values
func NewRange(left, right int64) core.Range {
	return core.Range{From: left, To: right}
}
