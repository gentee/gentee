// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// AssignAddºObj appends an item to object array
func AssignAddºObj(obj interface{}, value interface{}) (interface{}, error) {
	var err error
	if obj == nil {
		return nil, fmt.Errorf(ErrorText(ErrObjNil))
	}
	val := obj.(*core.Obj)
	switch v := val.Data.(type) {
	case *core.Array:
		v.Data = append(v.Data, value)
	default:
		err = fmt.Errorf(ErrorText(ErrObjValue))
	}
	return obj, err
}

// AssignºObjAny assigns integer, float, str, arr, map to obj
func AssignºObjAny(ptr interface{}, value interface{}) (interface{}, error) {
	return objType(value)
}

// AssignºObjBool assigns bool to obj
func AssignºObjBool(ptr interface{}, value interface{}) (interface{}, error) {
	ptr = objºBool(value.(int64))
	return ptr, nil
}

// arrºObj returns array of the objects
func arrºObj(obj *core.Obj) (*core.Array, error) {
	if obj.Data == nil {
		return nil, fmt.Errorf(ErrorText(ErrObjNil))
	}
	ret, ok := obj.Data.(*core.Array)
	if !ok {
		return nil, fmt.Errorf(ErrorText(ErrObjArr))
	}
	return ret, nil
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

// boolºObjDef converts object to boolean value
func boolºObjDef(val *core.Obj, def int64) (int64, error) {
	if val == nil || val.Data == nil {
		return def, nil
	}
	return boolºObj(val)
}

// ExpStrºObj adds string and obj in string expression
func ExpStrºObj(left string, right *core.Obj) string {
	return left + strºObj(right)
}

// floatºObj converts object to float
func floatºObj(val *core.Obj) (ret float64, err error) {
	switch v := val.Data.(type) {
	case int64:
		ret = floatºInt(v)
	case float64:
		ret = v
	case string:
		ret, err = floatºStr(v)
	default:
		err = fmt.Errorf(ErrorText(ErrObjValue))
	}
	return
}

// floatºObjDef converts object to float
func floatºObjDef(val *core.Obj, def float64) (float64, error) {
	if val == nil || val.Data == nil {
		return def, nil
	}
	return floatºObj(val)
}

// intºObj converts object to integer
func intºObj(val *core.Obj) (ret int64, err error) {
	switch v := val.Data.(type) {
	case int64:
		ret = v
	case bool:
		if v {
			ret = 1
		}
	case float64:
		ret = intºFloat(v)
	case string:
		ret, err = intºStr(v)
	default:
		err = fmt.Errorf(ErrorText(ErrObjValue))
	}
	return
}

// intºObjDef converts object to integer
func intºObjDef(val *core.Obj, def int64) (int64, error) {
	if val == nil || val.Data == nil {
		return def, nil
	}
	return intºObj(val)
}

// IsArrayºObj returns true if the object is an array
func IsArrayºObj(obj *core.Obj) int64 {
	if Type(obj) == `arr.obj` {
		return 1
	}
	return 0
}

// IsMapºObj returns true if the object is a map
func IsMapºObj(obj *core.Obj) int64 {
	if Type(obj) == `map.obj` {
		return 1
	}
	return 0
}

// IsNil returns true if the object is undefined
func IsNil(val *core.Obj) int64 {
	if val.Data == nil {
		return 1
	}
	return 0
}

// ItemºObjInt returns an item from array object
func ItemºObjInt(val *core.Obj, ind int64) (ret *core.Obj, err error) {
	if val == nil || val.Data == nil {
		return
	}
	switch v := val.Data.(type) {
	case *core.Array:
		if ind >= 0 && ind < int64(len(v.Data)) {
			ret = v.Data[ind].(*core.Obj)
		}
	default:
		err = fmt.Errorf(ErrorText(ErrObjValue))
	}
	return
}

// ItemºObjStr returns an item from map object
func ItemºObjStr(val *core.Obj, key string) (ret *core.Obj, err error) {
	if val == nil || val.Data == nil {
		return
	}
	switch v := val.Data.(type) {
	case *core.Map:
		if item, ok := v.Data[key]; ok {
			ret = item.(*core.Obj)
		}
	default:
		err = fmt.Errorf(ErrorText(ErrObjValue))
	}
	return
}

// mapºObj returns map of the objects
func mapºObj(obj *core.Obj) (*core.Map, error) {
	if obj.Data == nil {
		return nil, fmt.Errorf(ErrorText(ErrObjNil))
	}
	ret, ok := obj.Data.(*core.Map)
	if !ok {
		return nil, fmt.Errorf(ErrorText(ErrObjMap))
	}
	return ret, nil
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
		return `map.obj`
	}
	return `str`
}

// strºObj converts object value to string
func strºObj(val *core.Obj) string {
	return fmt.Sprint(val.Data)
}

// strºObjDef converts object value to string
func strºObjDef(val *core.Obj, def string) string {
	if val == nil || val.Data == nil {
		return def
	}
	return strºObj(val)
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
