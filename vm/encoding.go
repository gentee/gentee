// Copyright 2020 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gentee/gentee/core"
)

func IfaceToObj(val interface{}) (*core.Obj, error) {
	var err error
	ret := core.NewObj()
	switch v := val.(type) {
	case bool:
		ret.Data = v
	case string:
		ret.Data = v
	case int:
		ret.Data = v
	case int64:
		ret.Data = v
	case float64:
		ret.Data = v
	case json.Number:
		if ret.Data, err = v.Int64(); err != nil {
			if ret.Data, err = v.Float64(); err != nil {
				ret.Data = v.String()
			}
		}
	case []string:
		data := core.NewArray()
		data.Data = make([]interface{}, len(v))
		for i, item := range v {
			iobj := core.NewObj()
			iobj.Data = item
			data.Data[i] = iobj
		}
		ret.Data = data
	case []interface{}:
		data := core.NewArray()
		data.Data = make([]interface{}, 0, len(v))
		for _, item := range v {
			if item == nil {
				continue
			}
			iobj, err := IfaceToObj(item)
			if err != nil {
				return nil, err
			}
			data.Data = append(data.Data, iobj)
		}
		ret.Data = data
	case map[string]interface{}:
		data := core.NewMap()
		data.Keys = make([]string, 0, len(v))
		for key, vi := range v {
			if vi == nil {
				continue
			}
			data.Keys = append(data.Keys, key)
			iobj, err := IfaceToObj(vi)
			if err != nil {
				return nil, err
			}
			data.Data[key] = iobj
		}
		ret.Data = data
	case map[interface{}]interface{}:
		data := core.NewMap()
		data.Keys = make([]string, 0, len(v))
		for key, vi := range v {
			if vi == nil {
				continue
			}
			ikey := fmt.Sprint(key)
			data.Keys = append(data.Keys, ikey)
			iobj, err := IfaceToObj(vi)
			if err != nil {
				return nil, err
			}
			data.Data[ikey] = iobj
		}
		ret.Data = data
	default:
		return nil, fmt.Errorf(ErrorText(ErrObjType))
	}
	return ret, nil
}

// JsonToObj converts json to object
func JsonToObj(input string) (ret *core.Obj, err error) {
	d := json.NewDecoder(strings.NewReader(input))
	d.UseNumber()
	var v interface{}
	if err = d.Decode(&v); err != nil {
		return
	}
	return IfaceToObj(v)
}

func objToIface(obj *core.Obj) (ret interface{}) {
	switch v := obj.Data.(type) {
	case int64, float64, string, bool:
		ret = obj.Data
	case *core.Array:
		data := make([]interface{}, len(v.Data))
		for i, item := range v.Data {
			data[i] = objToIface(item.(*core.Obj))
		}
		ret = data
	case *core.Map:
		data := make(map[string]interface{})
		for _, key := range v.Keys {
			data[key] = objToIface(v.Data[key].(*core.Obj))
		}
		ret = data
	case *core.Obj:
		ret = objToIface(v)
	}
	return
}

// Json converts object to json
func Json(obj *core.Obj) (ret string, err error) {
	var out []byte
	if out, err = json.Marshal(objToIface(obj)); err == nil {
		ret = string(out)
	}
	return
}
