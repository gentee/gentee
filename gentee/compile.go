// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"fmt"
	"strconv"
)

const (
	cmdPush   = 1 + iota // push value into stack
	cmdReturn            // return values from the function
)

type Oper struct {
	Token int
	Cmd   *Cmd
}

var (
	opers = map[int]int{
		tkReturn: 0,
	}
)

// Compiler is used for saving compiling information
type Compiler struct {
	Code  *Code // current code
	Opers []*Oper
	Names map[string]*Object
}

// Compile compiles the source code
func (vm *VirtualMachine) Compile(source string) error {
	var (
		state int
	)
	lp, err := LexParsing([]rune(source))
	if err != nil {
		line, column := lp.LineColumn(len(lp.Tokens) - 1)
		return fmt.Errorf(` %d:%d: %s`, line, column, err)
	}
	vm.Lexeme = append(vm.Lexeme, lp)
	vm.Compiler.Code = &vm.Root

	stackState := make([]int, 0, 16)

	for i, token := range lp.Tokens {
		next := compileTable[state][token.Type]
		flag := next.Action & 0xff0000
		if flag&cfSkip != 0 {
			continue
		}
		if flag&cfError != 0 {
			return compileError(lp, next.Action&0xffff, i)
		}
		if next.Func != nil {
			if err := next.Func(vm, i); err != nil {
				return err
			}
		}
		if flag&cfStay != 0 {
			i--
		}
		if flag&cfBack != 0 {
			if len(stackState) == 0 {
				return compileError(lp, ErrCompiler, i, `stackState`)
			}
			state = stackState[len(stackState)-1]
			stackState = stackState[:len(stackState)-1]
			continue
		}
		stackState = append(stackState, state)
		state = next.Action & 0xffff
	}
	return nil
}

func (vm *VirtualMachine) getLex() *Lex {
	return vm.Lexeme[len(vm.Lexeme)-1]
}

func (comp *Compiler) appendCmd(cmd *Cmd) {
	comp.Code.ByteCode = append(comp.Code.ByteCode, *cmd)
}

func coExp(vm *VirtualMachine, cur int) error {
	comp := &vm.Compiler
	for len(comp.Opers) > 0 {
		comp.appendCmd(comp.Opers[len(comp.Opers)-1].Cmd)
		comp.Opers = comp.Opers[:len(comp.Opers)-1]
	}
	return nil
}

func coPush(vm *VirtualMachine, cur int) error {
	lp := vm.getLex()
	switch lp.Tokens[cur].Type {
	case tkInt:
		v, _ := strconv.ParseInt(lp.getToken(cur), 10, 64)
		vm.Compiler.appendCmd(&Cmd{ID: cmdPush, Value: v, TokenID: cur})
		fmt.Println(`INT`, v)
	case tkIntHex, tkIntOct:
		fmt.Println(`INTEGER`)
	}
	return nil
}

func expCmd(comp *Compiler, cmd *Cmd, token int) {
	if len(comp.Opers) == 0 || opers[token] > opers[comp.Opers[len(comp.Opers)-1].Token] {
		comp.Opers = append(comp.Opers, &Oper{Token: token, Cmd: cmd})
	} else {
		comp.appendCmd(comp.Opers[len(comp.Opers)-1].Cmd)
		comp.Opers = comp.Opers[:len(comp.Opers)-1]
	}
}

func coReturn(vm *VirtualMachine, cur int) error {
	expCmd(&vm.Compiler, &Cmd{ID: cmdReturn, Value: 1, TokenID: cur}, tkReturn)
	fmt.Println(`RETURN`)
	return nil
}

func newFunc(vm *VirtualMachine, name string) int {
	code := &Code{
		Owner: vm.Compiler.Code,
		LexID: len(vm.Lexeme) - 1,
	}
	vm.Funcs = append(vm.Funcs, code)
	vm.Compiler.Code.Children = append(vm.Compiler.Code.Children, code)
	vm.Compiler.Code = code
	// !!! TODO insert into Names
	return len(vm.Funcs) - 1
}

func coRun(vm *VirtualMachine, cur int) error {
	if vm.RunID != Undefined {
		return compileError(vm.getLex(), ErrRun, cur)
	}
	vm.RunID = newFunc(vm, `run`)
	return nil
}

func coRetType(vm *VirtualMachine, cur int) error {
	return nil
}
