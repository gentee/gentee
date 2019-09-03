// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"github.com/gentee/gentee/core"
)

// CopyVar copies one object to another one
func CopyVar(rt *Runtime, ptr *interface{}, value interface{}) {
	switch vItem := value.(type) {
	case *core.Fn:
		var pfn *core.Fn
		if ptr == nil || *ptr == nil {
			pfn = &core.Fn{}
		} else {
			pfn = (*ptr).(*core.Fn)
		}
		pfn.Func = vItem.Func
		*ptr = pfn
	case *Struct:
		var pstruct *Struct
		if ptr == nil || *ptr == nil {
			pstruct = NewStruct(rt, vItem.Type)
		} else {
			pstruct = (*ptr).(*Struct)
		}
		pstruct.Values = make([]interface{}, len(vItem.Values))
		for i, v := range vItem.Values {
			CopyVar(rt, &pstruct.Values[i], v)
		}
		*ptr = pstruct
	case *core.Set:
		var pset *core.Set
		if ptr == nil || *ptr == nil {
			pset = core.NewSet()
		} else {
			pset = (*ptr).(*core.Set)
		}
		pset.Data = make([]uint64, len(vItem.Data))
		copy(pset.Data, vItem.Data)
		*ptr = pset
	case *core.Buffer:
		var pbuf *core.Buffer
		if ptr == nil || *ptr == nil {
			pbuf = core.NewBuffer()
		} else {
			pbuf = (*ptr).(*core.Buffer)
		}
		pbuf.Data = make([]byte, len(vItem.Data))
		copy(pbuf.Data, vItem.Data)
		*ptr = pbuf
	case *core.Array:
		var parr *core.Array
		if ptr == nil || *ptr == nil {
			parr = core.NewArray()
		} else {
			parr = (*ptr).(*core.Array)
		}
		parr.Data = make([]interface{}, len(vItem.Data))
		for i, v := range vItem.Data {
			CopyVar(rt, &parr.Data[i], v)
		}
		*ptr = parr
	case *core.Map:
		var pmap *core.Map
		if ptr == nil || *ptr == nil {
			pmap = core.NewMap()
		} else {
			pmap = (*ptr).(*core.Map)
		}
		pmap.Keys = make([]string, len(vItem.Keys))
		var mapPtr interface{}
		for i, v := range vItem.Keys {
			mapPtr = pmap.Data[v]
			pmap.Keys[i] = v
			CopyVar(rt, &mapPtr, vItem.Data[v])
			pmap.Data[v] = mapPtr
		}
		*ptr = pmap
	default:
		*ptr = value
	}
}
