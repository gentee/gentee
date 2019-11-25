// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// boolºMap converts map to bool
func boolºMap(pmap *core.Map) int64 {
	if len(pmap.Data) == 0 {
		return 0
	}
	return 1
}

// DelºMapStr deletes key and value from the map
func DelºMapStr(pmap *core.Map, key string) *core.Map {
	delete(pmap.Data, key)
	for i, ikey := range pmap.Keys {
		if ikey == key {
			pmap.Keys = append(pmap.Keys[:i], pmap.Keys[i+1:]...)
			break
		}
	}
	return pmap
}

// IsKeyºMapStr returns true if there is the key in the map
func IsKeyºMapStr(pmap *core.Map, key string) int64 {
	_, ok := pmap.Data[key]
	if ok {
		return 1
	}
	return 0
}

// KeyºMapInt returns the key by the index
func KeyºMapInt(pmap *core.Map, index int64) (string, error) {
	if index >= int64(len(pmap.Keys)) {
		return ``, fmt.Errorf(ErrorText(ErrIndexOut))
	}
	return pmap.Keys[index], nil
}
