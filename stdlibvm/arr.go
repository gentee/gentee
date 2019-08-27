// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"github.com/gentee/gentee/core"
)

// AssignºArrArr copies one array to another one
func AssignºArrArr(arr interface{}, value interface{}) (interface{}, error) {
	core.CopyVar(&arr, value)
	return arr, nil
}

// AssignAddºArrStr appends one string to array
func AssignAddºArrAny(arr interface{}, value interface{}) (interface{}, error) {
	arr.(*core.Array).Data = append(arr.(*core.Array).Data, value)
	return arr, nil
}
