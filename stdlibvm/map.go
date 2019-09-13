// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"github.com/gentee/gentee/core"
)

// AssignºMapMap copies one array to another one
func AssignºMapMap(ptr interface{}, value interface{}) (interface{}, error) {
	core.CopyVar(&ptr, value)
	return ptr, nil
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
func IsKeyºMapStrAuto(pmap *core.Map, key string) int64 {
	_, ok := pmap.Data[key]
	if ok {
		return 1
	}
	return 0
}
