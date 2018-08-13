// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"strings"

	"github.com/gentee/gentee/core"
)

const (
	// The list of errors

	// ErrSuccess means no error
	ErrSuccess = iota
	// ErrLetter is returned when an unknown character has been found
	ErrLetter
	// ErrWord is returned when a sequence of characters is wrong
	ErrWord
	// ErrDecl is returned when the unexpexted token has been found on the top level
	ErrDecl
	// ErrLCurly is returned when the unexpexted token, expecting {
	ErrLCurly
	// ErrEnd is returned when the unexpexted end of the source, expecting }
	ErrEnd
	// ErrExp is returned when the unexpected token, expecting expression or statement {
	ErrExp
	// ErrName is return when compiler is expecting the identifier
	ErrName
	// ErrValue is returned when the unexpected token, expecting value, identifier or calling func
	ErrValue
	// ErrRun is returned when the compiler has found the second run function.
	ErrRun
	// ErrType is returned when the unexpected token, expecting type name
	ErrType
	// ErrReturn is returned when the function returns a value but it must not return
	ErrReturn
	// ErrMustReturn is returned when the function doesn't return a value but it must return
	ErrMustReturn
	// ErrReturnType is returned when the function returns a wrong type
	ErrReturnType
	// ErrOutOfRange is returned when the number is out of range
	ErrOutOfRange
	// ErrLPar is returned when there is an unclosed left parenthesis
	ErrLPar
	// ErrRPar is returned when extra right parenthesis has been found
	ErrRPar
	// ErrEmptyCode is returned when the source code is empty
	ErrEmptyCode
	// ErrFunction is returned when the compiler could not find a corresponding function
	ErrFunction
	// ErrBoolExp is returned when the compiler expects boolean result but gets different type
	ErrBoolExp
	// ErrFuncExists is returned when the function ahs already been defined
	ErrFuncExists
	// ErrUsedName is returned when the specified name has already been used
	ErrUsedName
	// ErrUnknownIdent is returned when the compiler gets unknown identifier
	ErrUnknownIdent
	// ErrAssign is returned when assign operator is used for different types
	ErrAssign
	// ErrLValue is returned when left operand of assign is not l-value
	ErrLValue
	// ErrOper is return when there is not operator
	ErrOper
	// ErrBoolOper is returned when && or || gets not boolen operands
	ErrBoolOper
	// ErrQuestion is returned when exp1 and exp2 have different types
	ErrQuestion
	// ErrQuestionPars is returned when ?(condition,exp1,exp2) has wrong parameters
	ErrQuestionPars
	// ErrCapitalLetters is returned when the var or func name consists of only capital letters
	ErrCapitalLetters
	// ErrConstName is returned when the name of constant doesn't consist of only capital letters
	ErrConstName
	// ErrMustAssign is returned when teh constant is described without assign
	ErrMustAssign
	// ErrIota is returned when IOTA is used outside const expression
	ErrIota
	// ErrIntOper is returned when ++ or -- gets not int value
	ErrIntOper
	// ErrDoubleQuotes is return when there is a wrong command of backslash in double quotes strings
	ErrDoubleQuotes

	// ErrCompiler error. It means a bug.
	ErrCompiler
)

var (
	errText = map[int]string{
		ErrLetter:         `unknown character`,
		ErrWord:           `wrong sequence of characters`,
		ErrLCurly:         `unexpected token, expecting {`,
		ErrEnd:            `unexpected end of the source`,
		ErrDecl:           `expected declaration: func, run etc`,
		ErrExp:            `unexpected token, expecting expression or statement`,
		ErrName:           `unexpected token, expecting the name of the identifier`,
		ErrRun:            `run function has already been defined`,
		ErrValue:          `unexpected token, expecting value, identifier or calling func`,
		ErrType:           `unexpected token, expecting type`,
		ErrReturn:         `function cannot return any value`,
		ErrMustReturn:     `function must return a value`,
		ErrReturnType:     `function returns wrong type`,
		ErrOutOfRange:     `the number %s is out of range`,
		ErrLPar:           `there is an unclosed left parenthesis`,
		ErrRPar:           `extra right parenthesis`,
		ErrEmptyCode:      `source code is empty`,
		ErrFunction:       `function %s has not been found`,
		ErrBoolExp:        `wrong type of expression, expecting boolean type`,
		ErrFuncExists:     `function %s has already been defined`,
		ErrUsedName:       `"%s" has already been used as the name of the function, type or variable`,
		ErrUnknownIdent:   `unknown identifier %s`,
		ErrAssign:         `assign different types`,
		ErrLValue:         `expecting l-value in the left operand of assign operator`,
		ErrOper:           `unexpected token, expecting operator`,
		ErrBoolOper:       `wrong type of operands, expecting boolean type`,
		ErrQuestion:       `different types of exp1 and exp2 in ?(cond, exp1, exp2)`,
		ErrQuestionPars:   `operator ? must be called as ?(boolean condition, exp1, exp2)`,
		ErrCapitalLetters: `The name of variable or function can't consists of only capital letters`,
		ErrConstName:      `The name of constant must consist of only capital letters`,
		ErrMustAssign:     `unexpected token, expecting =`,
		ErrIota:           `IOTA can be only used in const expression`,
		ErrIntOper:        `wrong type of operands, expecting int type`,
		ErrDoubleQuotes:   `invalid syntax of double quotes string`,

		ErrCompiler: `you have found a compiler bug [%s]. Let us know, please`,
	}
)

func (cmpl *compiler) ErrorPos(pos int, errID int, pars ...interface{}) error {

	line, column := cmpl.getLex().LineColumn(pos)
	return fmt.Errorf(`%d:%d: %s`, line, column, fmt.Sprintf(errText[errID], pars...))
}

func (cmpl *compiler) Error(errID int, pars ...interface{}) error {
	return cmpl.ErrorPos(cmpl.pos, errID, pars...)
}

func (cmpl *compiler) ErrorFunction(errID int, pos int, name string, pars []*core.TypeObject) error {
	var params []string
	for _, par := range pars {
		params = append(params, par.GetName())
	}
	return cmpl.ErrorPos(pos, errID, fmt.Sprintf(`%s(%s)`, name, strings.Join(params, `, `)))
}
