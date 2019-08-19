// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// AssignAddºArrInt appends one integer to array
func AssignAddºArrInt(arr *core.Array, value int64) *core.Array {
	arr.Data = append(arr.Data, value)
	return arr
}

// AssignAddºArrStr appends one string to array
func AssignAddºArrStr(arr *core.Array, value string) *core.Array {
	arr.Data = append(arr.Data, value)
	return arr
}

// LenºArr returns the length of the array
func LenºArr(arr *core.Array) int64 {
	fmt.Println(`LENARR`, *arr)
	return int64(len(arr.Data))
}
