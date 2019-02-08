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
	// ObjConst is a constant
	ObjConst
)

// IObject describes interface for all objects
type IObject interface {
	GetName() string
	Result() *TypeObject
	GetLex() *Lex
	GetParams() []*TypeObject
	GetType() ObjectType
	SetPub()
	GetUnitIndex() uint32
}

// Object contains information about any compiled object of the virtual machine
type Object struct {
	Name  string
	LexID int // the identifier of source code in Lexeme of Unit
	Unit  *Unit
	Pub   bool // public object
}

// TypeObject contains information about the type
type TypeObject struct {
	Object
	Original reflect.Type // Original golang type
	IndexOf  *TypeObject  // consists of elements
	Custom   *StructType  // for custom struct type
}

// EmbedObject contains information about the golang function
type EmbedObject struct {
	Object
	Func     interface{}   // golang function
	Return   *TypeObject   // the type of the result
	Params   []*TypeObject // the types of parameters
	Variadic bool          // variadic function
	Runtime  bool          // the first parameter is rt
}

// FuncObject contains information about the function
type FuncObject struct {
	Object
	Block CmdBlock
}

// ConstObject contains information about the constant
type ConstObject struct {
	Object
	Redefined bool
	Exp       ICmd        // expression
	Return    *TypeObject // the type of the constant
	Iota      int64       // iota
}

func getLex(obj *Object) *Lex {
	if obj.LexID < len(obj.Unit.Lexeme) {
		return obj.Unit.Lexeme[obj.LexID]
	}
	return nil
}

// GetName returns the name of the object
func (typeObj *TypeObject) GetName() (ret string) {
	ret = typeObj.Name
	if typeObj.IndexOf != nil && (ret == `arr` || ret == `map`) {
		ret += `.` + typeObj.IndexOf.GetName()
	}
	return
}

// GetLex returns the lex structure of the object
func (typeObj *TypeObject) GetLex() *Lex {
	return getLex(&typeObj.Object)
}

// GetType returns ObjType
func (typeObj *TypeObject) GetType() ObjectType {
	return ObjType
}

// Result returns the type of the result
func (typeObj *TypeObject) Result() *TypeObject {
	return nil
}

// GetParams returns the slice of parameters
func (typeObj *TypeObject) GetParams() []*TypeObject {
	return nil
}

// SetPub set Pub state
func (typeObj *TypeObject) SetPub() {
	typeObj.Pub = true
}

// GetUnitIndex returns the index of the unit of this object
func (typeObj *TypeObject) GetUnitIndex() uint32 {
	return typeObj.Unit.Index
}

// GetName returns the name of the object
func (funcObj *FuncObject) GetName() string {
	return funcObj.Name
}

// GetLex returns the lex structure of the object
func (funcObj *FuncObject) GetLex() *Lex {
	return getLex(&funcObj.Object)
}

// GetType returns ObjFunc
func (funcObj *FuncObject) GetType() ObjectType {
	return ObjFunc
}

// Result returns the type of the result
func (funcObj *FuncObject) Result() *TypeObject {
	return funcObj.Block.Result
}

// GetParams returns the slice of parameters
func (funcObj *FuncObject) GetParams() []*TypeObject {
	return funcObj.Block.Vars[:funcObj.Block.ParCount]
}

// SetPub set Pub state
func (funcObj *FuncObject) SetPub() {
	funcObj.Pub = true
}

// GetUnitIndex returns the index of the unit of this object
func (funcObj *FuncObject) GetUnitIndex() uint32 {
	return funcObj.Unit.Index
}

// GetName returns the name of the object
func (embedObj *EmbedObject) GetName() string {
	return embedObj.Name
}

// GetLex returns the lex structure of the object
func (embedObj *EmbedObject) GetLex() *Lex {
	return getLex(&embedObj.Object)
}

// GetType returns ObjEmbedded
func (embedObj *EmbedObject) GetType() ObjectType {
	return ObjEmbedded
}

// Result returns the type of the result
func (embedObj *EmbedObject) Result() *TypeObject {
	return embedObj.Return
}

// GetParams returns the slice of parameters
func (embedObj *EmbedObject) GetParams() []*TypeObject {
	return embedObj.Params
}

// SetPub set Pub state
func (embedObj *EmbedObject) SetPub() {
	embedObj.Pub = true
}

// GetUnitIndex returns the index of the unit of this object
func (embedObj *EmbedObject) GetUnitIndex() uint32 {
	return embedObj.Unit.Index
}

// GetName returns the name of the object
func (constObj *ConstObject) GetName() string {
	return constObj.Name
}

// GetLex returns the lex structure of the object
func (constObj *ConstObject) GetLex() *Lex {
	return getLex(&constObj.Object)
}

// GetType returns ObjType
func (constObj *ConstObject) GetType() ObjectType {
	return ObjConst
}

// Result returns the type of the result
func (constObj *ConstObject) Result() *TypeObject {
	return constObj.Return
}

// GetParams returns the slice of parameters
func (constObj *ConstObject) GetParams() []*TypeObject {
	return nil
}

// SetPub set Pub state
func (constObj *ConstObject) SetPub() {
	constObj.Pub = true
}

// GetUnitIndex returns the index of the unit of this object
func (constObj *ConstObject) GetUnitIndex() uint32 {
	return constObj.Unit.Index
}

// NewObject adds a new IObject to Unit
func (unit *Unit) NewObject(obj IObject) IObject {
	if unit.Pub > 0 {
		obj.SetPub()
		if unit.Pub == PubOne {
			unit.Pub = 0
		}
	}
	unit.VM.Objects = append(unit.VM.Objects, obj)
	return obj
}

// NewType adds a new type to Unit
func (unit *Unit) NewType(name string, original reflect.Type, indexOf IObject) IObject {
	typeObject := TypeObject{
		Object: Object{
			Name: name,
			Unit: unit,
		},
		Original: original,
	}
	if indexOf != nil {
		typeObject.IndexOf = indexOf.(*TypeObject)
	}
	unit.NewObject(&typeObject)
	ind := uint32(len(unit.VM.Objects) - 1)
	if typeObject.Pub {
		ind |= NSPub
	}
	unit.NameSpace[npType+name] = ind
	return &typeObject
}

// IsVariadic returns true if th efunction is variadic
func IsVariadic(obj IObject) bool {
	return (obj.GetType() == ObjFunc && obj.(*FuncObject).Block.Variadic) ||
		(obj.GetType() == ObjEmbedded && obj.(*EmbedObject).Variadic)
}
