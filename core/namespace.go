// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

// Prefixes
const (
	npType = `@`

	// NSImported means imported object in NSpace or NCustom
	NSImported = 0x10000000
	// NSPub means public object in NSpace or NCustom
	NSPub = 0x20000000

	NSIndex = 0x0fffffff
)

// GetObj returns the object by its index
func (unit *Unit) GetObj(ind uint32) IObject {
	return unit.VM.Objects[ind&NSIndex]
}

// FindType returns the type object with the specified name
func (unit *Unit) FindType(name string) IObject {
	if ind, ok := unit.NSpace[npType+name]; ok {
		return unit.VM.Objects[ind&NSIndex] //.(*TypeObject)
	}
	return nil
}
