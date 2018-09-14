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
	GetNext() IObject
	Result() *TypeObject
	GetLex() *Lex
	GetParams() []*TypeObject
	GetType() ObjectType
	SetNext(IObject)
}

// Object contains infromation about any compiled object of the virtual machine
type Object struct {
	//	Type  ObjectType
	Name  string
	Next  IObject // Next object with the same name
	LexID int     // the identifier of source code in Lexeme of Unit
	Unit  *Unit
}

// TypeObject contains information about the type
type TypeObject struct {
	Object
	Original reflect.Type // Original golang type
	IndexOf  *TypeObject  // consists of elements
}

// EmbedObject contains information about the golang function
type EmbedObject struct {
	Object
	Func   interface{}   // golang function
	Return *TypeObject   // the type of the result
	Params []*TypeObject // the types of parameters
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
func (typeObj *TypeObject) GetName() string {
	return typeObj.Name
}

// GetLex returns the lex structure of the object
func (typeObj *TypeObject) GetLex() *Lex {
	return getLex(&typeObj.Object)
}

// GetType returns ObjType
func (typeObj *TypeObject) GetType() ObjectType {
	return ObjType
}

// SetNext sets the next with the same name
func (typeObj *TypeObject) SetNext(next IObject) {
	if typeObj.Next == nil {
		typeObj.Next = next
	} else {
		typeObj.Next.SetNext(next)
	}

}

// GetNext returns the next object with the same name
func (typeObj *TypeObject) GetNext() IObject {
	return typeObj.Next
}

// Result returns the type of the result
func (typeObj *TypeObject) Result() *TypeObject {
	return nil
}

// GetParams returns the slice of parameters
func (typeObj *TypeObject) GetParams() []*TypeObject {
	return nil
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

// SetNext sets the next with the same name
func (funcObj *FuncObject) SetNext(next IObject) {
	if funcObj.Next == nil {
		funcObj.Next = next
	} else {
		funcObj.Next.SetNext(next)
	}
}

// GetNext returns the next object with the same name
func (funcObj *FuncObject) GetNext() IObject {
	return funcObj.Next
}

// Result returns the type of the result
func (funcObj *FuncObject) Result() *TypeObject {
	return funcObj.Block.Result
}

// GetParams returns the slice of parameters
func (funcObj *FuncObject) GetParams() []*TypeObject {
	return funcObj.Block.Vars[:funcObj.Block.ParCount]
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

// SetNext sets the next with the same name
func (embedObj *EmbedObject) SetNext(next IObject) {
	if embedObj.Next == nil {
		embedObj.Next = next
	} else {
		embedObj.Next.SetNext(next)
	}
}

// GetNext returns the next object with the same name
func (embedObj *EmbedObject) GetNext() IObject {
	return embedObj.Next
}

// Result returns the type of the result
func (embedObj *EmbedObject) Result() *TypeObject {
	return embedObj.Return
}

// GetParams returns the slice of parameters
func (embedObj *EmbedObject) GetParams() []*TypeObject {
	return embedObj.Params
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

// SetNext sets the next with the same name
func (constObj *ConstObject) SetNext(next IObject) {
	if constObj.Next == nil {
		constObj.Next = next
	} else {
		constObj.Next.SetNext(next)
	}

}

// GetNext returns the next object with the same name
func (constObj *ConstObject) GetNext() IObject {
	return constObj.Next
}

// Result returns the type of the result
func (constObj *ConstObject) Result() *TypeObject {
	return constObj.Return
}

// GetParams returns the slice of parameters
func (constObj *ConstObject) GetParams() []*TypeObject {
	return nil
}

// NewObject adds a new IObject to Unit
func (unit *Unit) NewObject(obj IObject) {
	name := obj.GetName()
	if curName := unit.Names[name]; curName == nil {
		unit.Names[name] = obj
	} else {
		curName.SetNext(obj)
	}
}

// NewType adds a new type to Unit
func (unit *Unit) NewType(name string, original reflect.Type, indexOf string) {
	typeObject := TypeObject{
		Object: Object{
			Name: name,
		},
		Original: original,
	}
	if len(indexOf) > 0 {
		typeObject.IndexOf = unit.Names[indexOf].(*TypeObject)
	}
	unit.NewObject(&typeObject)
}
