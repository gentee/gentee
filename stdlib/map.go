// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/core"
)

// InitMap appends stdlib map functions to the virtual machine
func InitMap(ws *core.Workspace) {
	for _, item := range []embedInfo{
		{LenºMap, `map*`, `int`}, // the length of map
		{core.Link{AssignºMapMap, core.ASSIGN + 2}, `map*,map*`, `map*`}, // map = map
		{AssignBitAndºMapMap, `map*,map*`, `map*`},                       // map &= map
		{DelºMapStrAuto, `map*,str`, `map*`},                             // Delete(map, str)
		{IsKeyºMapStrAuto, `map*,str`, `bool`},                           // IsKey(map, str)
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// LenºMap returns the length of the map
func LenºMap(pmap *core.Map) int64 {
	return int64(len(pmap.Data))
}

// AssignºMapMap copies one array to another one
func AssignºMapMap(ptr *interface{}, value *core.Map) *core.Map {
	core.CopyVar(ptr, value)
	return (*ptr).(*core.Map)
}

// AssignBitAndºMapMap assigns a pointer to the map
func AssignBitAndºMapMap(ptr *interface{}, value *core.Map) *core.Map {
	*ptr = value
	return (*ptr).(*core.Map)
}

// DelºMapStrAuto deletes key and value from the map
func DelºMapStrAuto(pmap *core.Map, key string) *core.Map {
	delete(pmap.Data, key)
	for i, ikey := range pmap.Keys {
		if ikey == key {
			pmap.Keys = append(pmap.Keys[:i], pmap.Keys[i+1:]...)
			break
		}
	}
	return pmap
}

// IsKeyºMapStrAuto returns true if there is teh key in the map
func IsKeyºMapStrAuto(pmap *core.Map, key string) bool {
	_, ok := pmap.Data[key]
	return ok
}
