// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"github.com/gentee/gentee/core"
)

// AssignºFnFn copies one fn to another one
func AssignºFnFn(ptr interface{}, value interface{}) (interface{}, error) {
	core.CopyVar(&ptr, value)
	return ptr, nil
}
