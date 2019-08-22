// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"sort"
	"strings"

	"github.com/gentee/gentee/core"
)

type embedInfo struct {
	Func    interface{}
	InTypes string
	OutType string
}

// InitArray appends stdlib array functions to the virtual machine
func InitArray(ws *core.Workspace) {
	for _, item := range []embedInfo{
		{core.Link{AssignAddºArrStr /*23<<16 | core.EMBED*/, core.ASSIGN + 1},
			`arr.str,str`, `arr.str`}, // arr += str
		{core.Link{LenºArr /*22<<16 | core.EMBED*/, core.Bcode(core.TYPEARR<<16) | core.LEN},
			`arr*`, `int`}, // the length of array
		{core.Link{AssignAddºArrInt /*24<<16 | core.EMBED*/, core.ASSIGN + 1},
			`arr.int,int`, `arr.int`}, // arr += int
		{core.Link{AssignAddºArrArr, core.ASSIGN + 1}, `arr.arr*,arr*`, `arr.arr*`},  // arr.arr += arr
		{core.Link{AssignAddºArrBool, core.ASSIGN + 1}, `arr.bool,bool`, `arr.bool`}, // arr += bool
		{AssignAddºArrMap, `arr.map*,map*`, `arr.map*`},                              // arr.map += map
		{core.Link{AssignºArrArr, core.ASSIGN}, `arr*,arr*`, `arr*`},                 // arr = arr
		{AssignBitAndºArrArr, `arr*,arr*`, `arr*`},                                   // arr &= arr
		{JoinºArrStr, `arr.str,str`, `str`},                                          // Join( arr.str, str )
		{SortºArr, `arr.str`, `arr.str`},                                             // Sort( arr.str )
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// LenºArr returns the length of the array
func LenºArr(arr *core.Array) int64 {
	return int64(len(arr.Data))
}

// AssignºArrArr copies one array to another one
func AssignºArrArr(ptr *interface{}, value *core.Array) *core.Array {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Array)
}

// AssignAddºArrArr appends one array to another one
func AssignAddºArrArr(ptr *interface{}, value *core.Array) *core.Array {
	(*ptr).(*core.Array).Data = append((*ptr).(*core.Array).Data, value)
	return (*ptr).(*core.Array)
}

// AssignAddºArrMap appends a map to array
func AssignAddºArrMap(ptr *interface{}, value *core.Map) *core.Array {
	(*ptr).(*core.Array).Data = append((*ptr).(*core.Array).Data, value)
	return (*ptr).(*core.Array)
}

// AssignAddºArrStr appends one string to array
func AssignAddºArrStr(ptr *interface{}, value string) *core.Array {
	(*ptr).(*core.Array).Data = append((*ptr).(*core.Array).Data, value)
	return (*ptr).(*core.Array)
}

// AssignAddºArrInt appends one integer to array
func AssignAddºArrInt(ptr *interface{}, value int64) *core.Array {
	(*ptr).(*core.Array).Data = append((*ptr).(*core.Array).Data, value)
	return (*ptr).(*core.Array)
}

// AssignAddºArrBool appends one boolean value to array
func AssignAddºArrBool(ptr *interface{}, value bool) *core.Array {
	(*ptr).(*core.Array).Data = append((*ptr).(*core.Array).Data, value)
	return (*ptr).(*core.Array)
}

// AssignBitAndºArrArr assigns a pointer to the array
func AssignBitAndºArrArr(ptr *interface{}, value *core.Array) *core.Array {
	*ptr = value
	return (*ptr).(*core.Array)
}

// JoinºArrStr concatenates the elements of a to create a single string.
func JoinºArrStr(value *core.Array, sep string) string {
	tmp := make([]string, len(value.Data))
	for i, item := range value.Data {
		tmp[i] = item.(string)
	}
	return strings.Join(tmp, sep)
}

// SortºArr sorts an array of strings
func SortºArr(value *core.Array) *core.Array {
	sort.Sort(value)
	return value
}
