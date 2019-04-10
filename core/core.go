// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package core

// CmdType is used for types of commands
type CmdType uint32

const (
	// CtValue pushes value into stack
	CtValue CmdType = iota + 1
	// CtVar pushes the value of the variable into stack
	CtVar
	// CtConst pushes the value of the constant into stack
	CtConst
	// CtStack is a stack command
	CtStack
	// CtUnary is an unary function
	CtUnary
	// CtBinary is binary function
	CtBinary
	// CtFunc is other functions
	CtFunc
	// CtCommand is a command of vm like break continue
	CtCommand
)
const (
	// Undefined index
	Undefined = -1
)

const (
	// RcBreak means break command
	RcBreak = 1 + iota
	// RcContinue means continue command
	RcContinue
)

const (
	// StackBlock executes function
	StackBlock = 1 + iota
	// StackReturn returns values from the function
	StackReturn
	// StackIf is the condition statement
	StackIf
	// StackWhile is the while statement
	StackWhile
	// StackSwitch is the switch statement
	StackSwitch
	// StackCase is the case statement
	StackCase
	// StackDefault is the default statement of switch
	StackDefault
	// StackAssign is an assign operator
	StackAssign
	// StackAnd is a logical AND
	StackAnd
	// StackOr is a logical OR
	StackOr
	// StackQuestion is ?(condition, exp1, exp2)
	StackQuestion
	// StackIncDec is ++ --
	StackIncDec
	// StackFor is the for statement
	StackFor
	// StackInit inits array and map variables
	StackInit
	// StackNew creates a new array or map
	StackNew
)

// Token is a lexical token.
type Token struct {
	Type   int32
	Index  int32 // index in Lex.Strings if Type is tkStr
	Offset int
	Length int
}

// Lex contains the result of the lexical parsing
type Lex struct {
	Source  []rune
	Tokens  []Token
	Lines   []int    // offsets of lines
	Strings []string // array of constant strings
	Header  string   // # header
	Path    string   // full path to the source
}

// ICmd is an interface for stack commands
type ICmd interface {
	GetType() CmdType
	GetResult() *TypeObject
	GetObject() IObject
	GetToken() int
}

// CmdCommon is a common structure for all commands
type CmdCommon struct {
	TokenID uint32 // the index of the token in lexeme
}

// CmdValue pushes a value into stack
type CmdValue struct {
	CmdCommon
	Value  interface{}
	Result *TypeObject
}

// CmdCommand is a runtime command
type CmdCommand struct {
	CmdCommon
	ID uint32 // id of the command
}

// CmdRet is the command for getting index values
type CmdRet struct {
	Cmd  ICmd        // the value of the index
	Type *TypeObject // the type of the result
}

// CmdVar pushes the value of the variable into stack
type CmdVar struct {
	CmdCommon
	Block   *CmdBlock // pointer to the block of the variable
	Index   int       // the index of the variable in the block
	Indexes []CmdRet  // the indexes list of the variable
}

// CmdConst pushes a value of the constant into stack
type CmdConst struct {
	CmdCommon
	Object IObject
}

// CmdBlock calls a stack command
type CmdBlock struct {
	CmdCommon
	Parent   *CmdBlock
	Object   IObject
	ID       uint32 // cmdType
	Vars     []*TypeObject
	ParCount int // the count of parameters
	Variadic bool
	VarNames map[string]int
	Result   *TypeObject
	Children []ICmd
}

// CmdUnary calls an unary function
type CmdUnary struct {
	CmdCommon
	Object  IObject
	Result  *TypeObject
	Operand ICmd
}

// CmdBinary calls a binary function
type CmdBinary struct {
	CmdCommon
	Object IObject
	Result *TypeObject
	Left   ICmd
	Right  ICmd
}

// CmdAnyFunc calls a function with more than 2 parameters
type CmdAnyFunc struct {
	CmdCommon
	Object   IObject
	Result   *TypeObject
	Children []ICmd
	FnVar    ICmd
	IsThread bool
}

// GetType returns CtValue
func (cmd *CmdValue) GetType() CmdType {
	return CtValue
}

// GetResult returns result
func (cmd *CmdValue) GetResult() *TypeObject {
	return cmd.Result
}

// GetToken returns the index of the token
func (cmd *CmdValue) GetToken() int {
	return int(cmd.TokenID)
}

// GetObject returns nil
func (cmd *CmdValue) GetObject() IObject {
	return nil
}

// GetType returns CtValue
func (cmd *CmdVar) GetType() CmdType {
	return CtVar
}

// GetResult returns result
func (cmd *CmdVar) GetResult() *TypeObject {
	typeVar := cmd.Block.Vars[cmd.Index]
	if cmd.Indexes != nil {
		typeVar = cmd.Indexes[len(cmd.Indexes)-1].Type
	}
	return typeVar
}

// GetToken returns the index of the token
func (cmd *CmdVar) GetToken() int {
	return int(cmd.TokenID)
}

// GetObject returns nil
func (cmd *CmdVar) GetObject() IObject {
	return nil
}

// GetType returns CtConst
func (cmd *CmdConst) GetType() CmdType {
	return CtConst
}

// GetResult returns result
func (cmd *CmdConst) GetResult() *TypeObject {
	return cmd.Object.Result()
}

// GetToken returns the index of the token
func (cmd *CmdConst) GetToken() int {
	return int(cmd.TokenID)
}

// GetObject returns nil
func (cmd *CmdConst) GetObject() IObject {
	return cmd.Object
}

// GetType returns CtStack
func (cmd *CmdBlock) GetType() CmdType {
	return CtStack
}

// GetResult returns result
func (cmd *CmdBlock) GetResult() *TypeObject {
	return cmd.Result
}

// GetObject returns nil
func (cmd *CmdBlock) GetObject() IObject {
	return cmd.Object
}

// GetToken returns the index of the token
func (cmd *CmdBlock) GetToken() int {
	return int(cmd.TokenID)
}

// GetType returns CtUnary
func (cmd *CmdUnary) GetType() CmdType {
	return CtUnary
}

// GetResult returns the type of the result
func (cmd *CmdUnary) GetResult() *TypeObject {
	return cmd.Result
}

// GetObject returns Object
func (cmd *CmdUnary) GetObject() IObject {
	return cmd.Object
}

// GetToken returns the index of the token
func (cmd *CmdUnary) GetToken() int {
	return int(cmd.TokenID)
}

// GetType returns CtBinary
func (cmd *CmdBinary) GetType() CmdType {
	return CtBinary
}

// GetResult returns the type of the result
func (cmd *CmdBinary) GetResult() *TypeObject {
	return cmd.Result
}

// GetObject returns Object
func (cmd *CmdBinary) GetObject() IObject {
	return cmd.Object
}

// GetToken returns the index of the token
func (cmd *CmdBinary) GetToken() int {
	return int(cmd.TokenID)
}

// GetType returns CtFunc
func (cmd *CmdAnyFunc) GetType() CmdType {
	return CtFunc
}

// GetResult returns the type of the result
func (cmd *CmdAnyFunc) GetResult() *TypeObject {
	return cmd.Result
}

// GetObject returns Object
func (cmd *CmdAnyFunc) GetObject() IObject {
	return cmd.Object
}

// GetToken returns the index of the token
func (cmd *CmdAnyFunc) GetToken() int {
	return int(cmd.TokenID)
}

// GetType returns CtCommand
func (cmd *CmdCommand) GetType() CmdType {
	return CtCommand
}

// GetResult returns result
func (cmd *CmdCommand) GetResult() *TypeObject {
	return nil
}

// GetToken returns the index of the token
func (cmd *CmdCommand) GetToken() int {
	return int(cmd.TokenID)
}

// GetObject returns nil
func (cmd *CmdCommand) GetObject() IObject {
	return nil
}

// LineColumn return the line and the column of the ind-th token
func (lp Lex) LineColumn(ind int) (line int, column int) {
	end := len(lp.Tokens) == ind && ind > 0
	if end {
		ind--
	}
	if len(lp.Tokens) > ind {
		for ; line < len(lp.Lines); line++ {
			if lp.Lines[line] > lp.Tokens[ind].Offset {
				break
			}
		}
		column = lp.Tokens[ind].Offset - lp.Lines[line-1] + 1
		if end {
			column += lp.Tokens[ind].Length
		}
	}
	return
}

// NewToken appends a new token to lexems
func (lp *Lex) NewToken(token, offset, length int) {
	lp.Tokens = append(lp.Tokens, Token{Type: int32(token), Offset: offset, Length: length})
}

// NewTokens appends one-byte new tokens to lexems
func (lp *Lex) NewTokens(offset int, tokens ...int) {
	for _, token := range tokens {
		lp.Tokens = append(lp.Tokens, Token{Type: int32(token), Offset: offset, Length: 1})
	}
}
