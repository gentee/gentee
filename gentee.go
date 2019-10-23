// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gentee/gentee/compiler"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/vm"
)

// Exec is a structure with a bytecode that is ready to run
type Exec struct {
	*core.Exec
}

// Unit is a structure describing source code unit
type Unit struct {
	*core.Unit
}

// Settings is a structure with parameters for running bytecode
type Settings struct {
	vm.Settings
}

// Gentee is a common structure for compiling and executing Gentee source code
type Gentee struct {
	*core.Workspace
}

// EmbedItem is a structure for declaration of embedded functions.
type EmbedItem struct {
	Prototype string
	Object    interface{}
}

// Custom is a structure with parameters for compiling and runtime
type Custom struct {
	Embedded []EmbedItem
}

func str2type(in string) (ret uint16) {
	switch in {
	case ``:
		ret = core.TYPENONE
	case `int`, `thread`:
		ret = core.TYPEINT
	case `bool`:
		ret = core.TYPEBOOL
	case `float`:
		ret = core.TYPEFLOAT
	case `char`:
		ret = core.TYPECHAR
	case `str`:
		ret = core.TYPESTR
	case `range`:
		ret = core.TYPERANGE
	case `buf`:
		ret = core.TYPEBUF
	case `fn`:
		ret = core.TYPEFUNC
	case `error`:
		ret = core.TYPEERROR
	case `set`:
		ret = core.TYPESET
	default:
		if in == `arr` || strings.HasPrefix(in, `arr.`) {
			ret = core.TYPEARR
		} else if in == `map` || strings.HasPrefix(in, `map.`) {
			ret = core.TYPEMAP
		} else {
			ret = core.TYPESTRUCT
		}
	}
	return
}

func str2pars(in string) (types []uint16) {
	if len(in) == 0 {
		return
	}
	for _, par := range strings.Split(in, `,`) {
		types = append(types, str2type(strings.TrimSpace(par)))
	}
	return
}

func Customize(custom *Custom) error {
	re, err := regexp.Compile(`^([\wº]+)\(([\w ,\.\*]*)\)\s*([\w\.\*]*)?`)
	if err != nil {
		return err
	}
	for _, v := range custom.Embedded {
		v.Prototype = strings.ReplaceAll(v.Prototype, ` `, ``)
		if len(v.Prototype) == 0 || v.Object == nil {
			return fmt.Errorf("%s %v", vm.ErrorText(vm.ErrCustom), v)
		}
		list := re.FindAllStringSubmatch(v.Prototype, -1)
		vals := list[0]
		if len(vals) < 4 {
			return fmt.Errorf("%s %v", vm.ErrorText(vm.ErrCustom), v)
		}
		t := reflect.TypeOf(v.Object)
		embed := core.Embed{
			Name:     vals[1],
			Pars:     vals[2],
			Ret:      vals[3],
			Code:     uint32(len(vm.EmbedFuncs)),
			Func:     v.Object,
			Return:   str2type(vals[3]),
			Params:   str2pars(vals[2]),
			Variadic: t.IsVariadic(),
			Runtime:  t.NumIn() > 0 && t.In(0) == reflect.TypeOf(&vm.Runtime{}),
			CanError: t.NumOut() >= 1 && t.Out(t.NumOut()-1).String() == `error`,
		}
		vm.EmbedFuncs = append(vm.EmbedFuncs, embed)
	}
	return nil
}

// New creates a new Gentee workspace
func New() *Gentee {
	g := Gentee{
		Workspace: core.NewVM(vm.EmbedFuncs),
	}
	compiler.InitStdlib(g.Workspace)
	return &g
}

// Compile compiles the Gentee source code.
// The function returns bytecode, id of the compiled unit and error code.
func (g *Gentee) Compile(input, path string) (*Exec, int, error) {
	unitID, err := compiler.Compile(g.Workspace, input, path)
	if err != nil {
		return nil, 0, err
	}
	exec, err := compiler.Link(g.Workspace, unitID)
	return &Exec{Exec: exec}, unitID, err
}

// CompileAndRun compiles the specified Gentee source file and run it.
func (g *Gentee) CompileAndRun(filename string) (interface{}, error) {
	exec, _, err := g.CompileFile(filename)
	if err != nil {
		return nil, err
	}
	return exec.Run(Settings{})
}

// CompileFile compiles the specified Gentee source file.
// The function returns bytecode, id of the compiled unit and error code.
func (g *Gentee) CompileFile(filename string) (*Exec, int, error) {
	unitID, err := compiler.CompileFile(g.Workspace, filename)
	if err != nil {
		return nil, 0, err
	}
	exec, err := compiler.Link(g.Workspace, unitID)
	return &Exec{Exec: exec}, unitID, err
}

// Unit returns the unit structure by its index.
func (g *Gentee) Unit(unitID int) Unit {
	return Unit{Unit: g.Units[unitID]}
}

// Run executes the bytecode.
func (exec *Exec) Run(settings Settings) (interface{}, error) {
	return vm.Run(exec.Exec, settings.Settings)
}

// Version returns the current version of the Gentee compiler.
func (g *Gentee) Version() string {
	return core.Version
}
