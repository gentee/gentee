// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"sort"
	"strings"

	"github.com/gentee/gentee/core"
)

// AssignAddºArrAny appends an item to array
func AssignAddºArrAny(arr interface{}, value interface{}) (interface{}, error) {
	arr.(*core.Array).Data = append(arr.(*core.Array).Data, value)
	return arr, nil
}

// JoinºArrStr concatenates the elements of a to create a single string.
func JoinºArrStr(value *core.Array, sep string) string {
	tmp := make([]string, len(value.Data))
	for i, item := range value.Data {
		tmp[i] = item.(string)
	}
	return strings.Join(tmp, sep)
}

// ReverseºArr reverses an array
func ReverseºArr(arr interface{}) interface{} {
	data := arr.(*core.Array).Data
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	arr.(*core.Array).Data = data
	return arr
}

// SortºArr sorts an array of strings
func SortºArr(value *core.Array) *core.Array {
	sort.Sort(value)
	return value
}
