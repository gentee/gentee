// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

import (
	"reflect"
)

// ObjectType is used for types of objects
type ObjectType int

const (
	// ObjType is a type
	ObjType ObjectType = iota + 1
	// ObjEmbedded is a embedded golang function
	ObjEmbedded
	// ObjFunc is a gentee function
	ObjFunc
)

// Object contains infromation about any compiled object of the virtual machine
type Object struct {
	Type ObjectType
	Name string
	Next *Object // Next object with the same name
}

// TypeObject contains information about the type
type TypeObject struct {
	Object
	Original reflect.Type // Original golang type
}

// NewType adds a new type to Unit
func (unit *Unit) NewType(name string, original reflect.Type) {
	object := TypeObject{
		Object: Object{
			Type: ObjType,
			Name: name,
		},
		Original: original,
	}
	unit.Names[name] = &object.Object
}
