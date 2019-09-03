// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"strings"

	"github.com/gentee/gentee/core"
)

// Struct is used for custom struct types
type Struct struct {
	Type   *core.StructInfo
	Values []interface{} // Values of fields
}

// NewStruct creates a new struct object
func NewStruct(rt *Runtime, sInfo *core.StructInfo) *Struct {
	//	ind := (itype - core.TYPESTRUCT) >> 8
	//	sInfo := rt.Owner.Exec.Structs[ind]
	values := make([]interface{}, len(sInfo.Fields))
	for i, v := range sInfo.Fields {
		if v < core.TYPESTRUCT ||
			&rt.Owner.Exec.Structs[(v-core.TYPESTRUCT)>>8] != sInfo {
			values[i] = newValue(rt, int(v))
		}
	}
	return &Struct{
		Type:   sInfo, //rt.Owner.Exec.Structs[ind],
		Values: values,
	}
}

// String interface for Struct
func (pstruct Struct) String() string {
	name := pstruct.Type.Name
	list := make([]string, len(pstruct.Values))
	for i, v := range pstruct.Values {
		list[i] = fmt.Sprintf(`%s:%v`, pstruct.Type.Keys[i], fmt.Sprint(v))
	}
	return name + `[` + strings.Join(list, ` `) + `]`
}

// Len is part of Indexer interface.
func (pstruct *Struct) Len() int {
	return len(pstruct.Values)
}

// GetIndex is part of Indexer interface.
func (pstruct *Struct) GetIndex(index interface{}) (interface{}, bool) {
	sindex := int(index.(int64))
	if sindex < 0 || sindex >= len(pstruct.Values) {
		return nil, false
	}
	return pstruct.Values[sindex], true
}

// SetIndex is part of Indexer interface.
func (pstruct *Struct) SetIndex(index, value interface{}) int {
	sindex := int(index.(int64))
	if sindex < 0 || sindex >= len(pstruct.Values) {
		return core.ErrIndexOut
	}
	pstruct.Values[sindex] = value
	return 0
}
