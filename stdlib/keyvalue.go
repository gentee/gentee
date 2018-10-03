// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitKeyValue appends stdlib range functions to the virtual machine
func InitKeyValue(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		NewKeyValueºStr,  // binary :
		NewKeyValueºInt,  // binary :
		NewKeyValueºBool, // binary :
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// NewKeyValueºStr adds key-value structure
func NewKeyValueºStr(left string, right string) core.KeyValue {
	return core.KeyValue{Key: left, Value: right}
}

// NewKeyValueºInt adds key-value structure
func NewKeyValueºInt(left string, right int64) core.KeyValue {
	return core.KeyValue{Key: left, Value: right}
}

// NewKeyValueºBool adds key-value structure
func NewKeyValueºBool(left string, right bool) core.KeyValue {
	return core.KeyValue{Key: left, Value: right}
}
