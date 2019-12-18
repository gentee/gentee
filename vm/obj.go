// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// AssignºObjAny assigns integer, float, str, arr, map to obj
func AssignºObjAny(ptr interface{}, value interface{}) (interface{}, error) {
	return objType(value)
}

// AssignºObjBool assigns bool to obj
func AssignºObjBool(ptr interface{}, value interface{}) (interface{}, error) {
	ptr = objºBool(value.(int64))
	return ptr, nil
}

// boolºObj converts object to boolean value
func boolºObj(val *core.Obj) (ret int64, err error) {
	if val.Data == nil {
		return 0, fmt.Errorf(ErrorText(ErrObjNil))
	}
	switch v := val.Data.(type) {
	case int64:
		ret = boolºInt(v)
	case bool:
		if v {
			ret = 1
		}
	case float64:
		ret = boolºFloat(v)
	case string:
		ret = boolºStr(v)
	case *core.Array:
		ret = boolºArr(v)
	case *core.Map:
		ret = boolºMap(v)
	}
	return
}

// ExpStrºObj adds string and obj in string expression
func ExpStrºObj(left string, right *core.Obj) string {
	return left + strºObj(right)
}

// IsNil returns true if the object is undefined
func IsNil(val *core.Obj) int64 {
	if val.Data == nil {
		return 1
	}
	return 0
}

// objºBool converts boolean value to object
func objºBool(val int64) *core.Obj {
	obj := core.NewObj()
	obj.Data = val != 0
	return obj
}

// objºAny converts int, float, string to object
func objºAny(val interface{}) *core.Obj {
	obj := core.NewObj()
	obj.Data = val
	return obj
}

// Type returns the type of the object's value
func Type(val *core.Obj) string {
	if val.Data == nil {
		return `nil`
	}
	switch val.Data.(type) {
	case bool:
		return `bool`
	case int64:
		return `int`
	case float64:
		return `float`
	case *core.Array:
		return `arr.obj`
	case *core.Map:
		return `arr.map`
	}
	return `str`
}

// strºObj converts object value to string
func strºObj(val *core.Obj) string {
	return fmt.Sprint(val.Data)
}

// objType converts variable to object
func objType(val interface{}) (*core.Obj, error) {
	obj := core.NewObj()
	switch v := val.(type) {
	case int64, float64, string:
		obj.Data = val
	case bool:
		obj.Data = v
	case *core.Array:
		data := core.NewArray()
		for _, item := range v.Data {
			iobj, err := objType(item)
			if err != nil {
				return nil, err
			}
			data.Data = append(data.Data, iobj)
		}
		obj.Data = data
	case *core.Map:
		data := core.NewMap()
		data.Keys = make([]string, len(v.Keys))
		for i, key := range v.Keys {
			data.Keys[i] = key
			iobj, err := objType(v.Data[key])
			if err != nil {
				return nil, err
			}
			data.Data[key] = iobj
		}
		obj.Data = data
	case *core.Obj:
		obj = v
	default:
		return nil, fmt.Errorf(ErrorText(ErrObjType))
	}
	return obj, nil
}

// objºArrMap converts array and map to object
func objºArrMap(val interface{}) (*core.Obj, error) {
	return objType(val)
}
