// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

// Prefixes
const (
	npType   = `@`
	npConst  = `$`
	npCustom = `?`
	npFunc   = `#`

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

// FindObj returns the object by its name
func (unit *Unit) FindObj(fullName string) IObject {
	if ind, ok := unit.NSpace[fullName]; ok {
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
func (unit *Unit) FindFunc(name string, params []*TypeObject) (obj IObject, custom []uint32) {
	key := npFunc + name
	for _, v := range params {
		key += npFunc + v.GetName()
	}
	if obj = unit.FindObj(key); obj != nil {
		return
	}
	custom = unit.NCustom[npCustom+name]
	return
}

// AddConst appends a constant to NSpace
func (unit *Unit) AddConst(name string) {
	ind := uint32(len(unit.VM.Objects) - 1)
	obj := unit.GetObj(ind).(*ConstObject)
	if obj.Pub {
		ind |= NSPub
	}
	unit.NSpace[npConst+name] = ind
}

// AddCustom appends func to NCustom
func (unit *Unit) AddCustom(ind int, name string, pub bool) {
	if pub {
		ind |= NSPub
	}
	key := npCustom + name
	if v, ok := unit.NCustom[key]; !ok {
		unit.NCustom[key] = []uint32{uint32(ind)}
	} else {
		unit.NCustom[key] = append(v, uint32(ind))
	}
}

// AddFunc appends func to NSpace
func (unit *Unit) AddFunc(ind int, name string, params []*TypeObject, pub bool) {
	if pub {
		ind |= NSPub
	}
	key := npFunc + name
	for _, v := range params {
		key += npFunc + v.GetName()
	}
	unit.NSpace[key] = uint32(ind)
}
