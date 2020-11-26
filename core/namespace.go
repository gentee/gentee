// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"strings"
)

// Prefixes
const (
	npType     = `@`
	npConst    = `$`
	npVariadic = `?`
	npFunc     = `#`

	// NSImported means imported object in NameSpace
	NSImported = 0x10000000
	// NSPub means public object in NameSpace
	NSPub = 0x20000000

	NSIndex = 0x0fffffff
)

// GetObj returns the object by its index
func (unit *Unit) GetObj(ind uint32) IObject {
	return unit.VM.Objects[ind&NSIndex]
}

// FindObj returns the object by its name
func (unit *Unit) FindObj(fullName string) IObject {
	if ind, ok := unit.NameSpace[fullName]; ok {
		return unit.VM.Objects[ind&NSIndex]
	}
	return nil
}

// FindType returns the type object with the specified name
func (unit *Unit) FindType(name string) IObject {
	return unit.FindObj(npType + name)
}

// FindConst returns the constant object with the specified name
func (unit *Unit) FindConst(name string) IObject {
	return unit.FindObj(npConst + name)
}

// FindFunc returns the function with the specified name and parameters
func (unit *Unit) FindFunc(name string, params []*TypeObject) (IObject, bool) {
	var isStruct bool
	key := npFunc + name
	keyAny := key
	for _, v := range params {
		if v == nil {
			return nil, false
		}
		parName := v.GetName()
		key += npFunc + parName
		if strings.HasPrefix(parName, `arr.`) {
			keyAny += npFunc + `arr*`
		} else if strings.HasPrefix(parName, `map.`) {
			keyAny += npFunc + `map*`
		} else {
			keyAny += npFunc + parName
		}
		isStruct = isStruct || v.Custom != nil
	}
	if obj := unit.FindObj(key); obj != nil {
		return obj, false
	}
	if key != keyAny {
		if obj := unit.FindObj(keyAny); obj != nil {
			return obj, false
		}
	}
	if isStruct {
		key = npFunc + name
		for _, v := range params {
			var parName string
			if v.Custom != nil {
				parName = `struct`
			} else {
				parName = v.GetName()
			}
			key += npFunc + parName
		}
		if obj := unit.FindObj(key); obj != nil {
			return obj, false
		}
	}
	return unit.FindObj(npVariadic + name), true
}

// AddConst appends a constant to NameSpace
func (unit *Unit) AddConst(name string) {
	ind := uint32(len(unit.VM.Objects) - 1)
	obj := unit.GetObj(ind).(*ConstObject)
	if obj.Pub {
		ind |= NSPub
	}
	unit.NameSpace[npConst+name] = ind
}

// AddFunc appends func to NameSpace
func (unit *Unit) AddFunc(ind int, obj IObject, pub bool) {
	var key string
	if pub {
		ind |= NSPub
	}
	name := obj.GetName()
	if IsVariadic(obj) {
		key = npVariadic + name
	} else {
		key = npFunc + name
		for _, v := range obj.GetParams() {
			key += npFunc + v.GetName()
		}
	}
	unit.NameSpace[key] = uint32(ind)
}
