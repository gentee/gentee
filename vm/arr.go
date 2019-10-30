// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gentee/gentee/core"
)

// AssignAddºArr appends one array to another one
func AssignAddºArr(dest interface{}, src interface{}) (interface{}, error) {
	for _, item := range src.(*core.Array).Data {
		dest.(*core.Array).Data = append(dest.(*core.Array).Data, item)
	}
	return dest, nil
}

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
func ReverseºArr(arr *core.Array) *core.Array {
	for i, j := 0, len(arr.Data)-1; i < j; i, j = i+1, j-1 {
		arr.Data[i], arr.Data[j] = arr.Data[j], arr.Data[i]
	}
	return arr
}

// SliceºArr extracts some consecutive elements from within an array.
func SliceºArr(arr *core.Array, start, end int64) (*core.Array, error) {
	ret := core.NewArray()
	if start < 0 || end > int64(len(arr.Data)) {
		return ret, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	if end == 0 {
		end = int64(len(arr.Data))
	}
	for ; start < end; start++ {
		ret.Data = append(ret.Data, arr.Data[start])
	}
	return ret, nil
}

// SortºArr sorts an array of strings
func SortºArr(value *core.Array) *core.Array {
	sort.Sort(value)
	return value
}
