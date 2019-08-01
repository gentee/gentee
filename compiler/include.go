// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gentee/gentee/core"
)

// CompileFile compiles the source file
func CompileFile(ws *core.Workspace, filename string) (unitID int, err error) {
	var (
		absname, curDir string
		input           []byte
	)
	if absname, err = filepath.Abs(filename); err != nil {
		return
	}
	if unitID = ws.Linked[absname]; unitID != 0 {
		return
	}
	if curDir, err = os.Getwd(); err != nil {
		return
	}
	if err = os.Chdir(filepath.Dir(absname)); err != nil {
		return
	}
	defer os.Chdir(curDir)
	if input, err = ioutil.ReadFile(absname); err != nil {
		return
	}
	unitID, err = Compile(ws, string(input), absname)
	if err == nil {
		ws.Linked[absname] = unitID
	}
	return
}

func coInclude(cmpl *compiler) error {
	cmpl.isImport = false
	return nil
}

func coImport(cmpl *compiler) error {
	cmpl.isImport = true
	return nil
}

func coPub(cmpl *compiler) error {
	cmpl.unit.Pub = core.PubOne
	return nil
}

func coIncludeImport(cmpl *compiler) error {
	var (
		v      interface{}
		err    error
		unitID int
	)
	lp := cmpl.getLex()
	token := getToken(lp, cmpl.pos)
	v = lp.Strings[lp.Tokens[cmpl.pos].Index]
	if len(lp.Tokens) > cmpl.pos+1 && lp.Tokens[cmpl.pos+1].Type == tkStrExp {
		return cmpl.ErrorPos(cmpl.pos+1, ErrImportStr)
	}
	if token[0] == '"' {
		if v, err = unNewLine(v.(string)); err != nil {
			return cmpl.Error(ErrDoubleQuotes)
		}
	}
	includeFile := os.ExpandEnv(v.(string))
	unitID, err = CompileFile(cmpl.ws, includeFile)
	if err != nil && unitID == 0 {
		return cmpl.Error(ErrIncludeFile, includeFile)
	}
	if err == nil {
		if v, ok := cmpl.unit.Included[uint32(unitID)]; !ok || (v && !cmpl.isImport) {
			err = cmpl.copyNameSpace(cmpl.ws.Units[unitID], cmpl.isImport)
			cmpl.unit.Included[uint32(unitID)] = cmpl.isImport
		}
	}
	return err
}
