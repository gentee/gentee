// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitKeyValue appends stdlib range functions to the virtual machine
func InitKeyValue(ws *core.Workspace) {
	for _, item := range []interface{}{
		//		core.Link{NewKeyValue, core.KEYVALUE}, // binary :
		core.Link{NewKeyValue, core.NOP}, // binary :
	} {
		ws.StdLib().NewEmbed(item)
	}
}

// NewKeyValue adds key-value structure
func NewKeyValue(left interface{}, right interface{}) core.KeyValue {
	return core.KeyValue{Key: left, Value: right}
}
