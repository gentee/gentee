// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// Fn is used for custom func types
type Fn struct {
	Func int32 // id of function
}

// CopyVar copies one object to another one
func CopyVar(rt *Runtime, ptr *interface{}, value interface{}) {
	switch vItem := value.(type) {
	case *Fn:
		var pfn *Fn
		if ptr == nil || *ptr == nil {
			pfn = &Fn{}
		} else {
			pfn = (*ptr).(*Fn)
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

func newValue(rt *Runtime, vtype int) interface{} {
	switch vtype {
	case core.TYPEINT, core.TYPEBOOL:
		return int64(0)
	case core.TYPECHAR:
		return int64(' ')
	case core.TYPESTR:
		return ``
	case core.TYPEFLOAT:
		return float64(0.0)
	case core.TYPEARR:
		return core.NewArray()
	case core.TYPEMAP:
		return core.NewMap()
	case core.TYPEBUF:
		return core.NewBuffer()
	case core.TYPEFUNC:
		return &Fn{}
	case core.TYPEERROR:
		return &RuntimeError{}
	case core.TYPESET:
		return core.NewSet()
	default:
		if vtype >= core.TYPESTRUCT {
			return NewStruct(rt, &rt.Owner.Exec.Structs[(vtype-core.TYPESTRUCT)>>8])

		} else {
			fmt.Println(`NEW VALUE`, vtype)
		}
	}
	return nil
}
