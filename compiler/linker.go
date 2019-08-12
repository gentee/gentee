// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"reflect"

	"github.com/gentee/gentee/core"
)

type Linker struct {
	Blocks []*core.CmdBlock
	Lex    *core.Lex
}

func Link(ws *core.Workspace, unitID int) (*core.Exec, error) {
	var exec *core.Exec
	if unitID < 0 || unitID >= len(ws.Units) {
		return nil, fmt.Errorf(errText[ErrLinkIndex], unitID)
	}
	unit := ws.Units[unitID]
	if unit.RunID == core.Undefined {
		return nil, nil
	}
	bcode := genBytecode(ws, int32(unit.RunID))
	exec = &core.Exec{
		Code:  append([]core.Bcode{}, bcode.Code...),
		Funcs: make(map[int32]int32),
		Pos:   bcode.Pos,
	}
	var (
		ok  bool
		ind uint16
	)
	for ikey := range bcode.Used {
		exec.Funcs[ikey] = int32(len(exec.Code))
		usedCode := ws.Objects[ikey].GetCode()
		shift := int32(len(exec.Code))
		exec.Code = append(exec.Code, usedCode.Code...)

		rebuild := make([]uint16, len(usedCode.Strings))
		for key, curInd := range usedCode.Strings {
			if ind, ok = bcode.Strings[key]; !ok {
				ind = uint16(len(bcode.Strings))
				bcode.Strings[key] = ind
			}
			rebuild[curInd] = ind
		}
		for _, pos := range usedCode.Pos {
			exec.Pos = append(exec.Pos, core.CodePos{
				Offset: pos.Offset + shift,
				Path:   rebuild[pos.Path],
				Name:   rebuild[pos.Name],
				Line:   pos.Line,
				Column: pos.Column,
			})
		}

	}
	exec.Strings = make([]string, len(bcode.Strings))
	for key, ikey := range bcode.Strings {
		exec.Strings[ikey] = key
	}
	//	fmt.Println(`NAMES`, exec.Paths, exec.Names, exec.Pos)
	//	fmt.Println(`USED`, exec.Funcs, exec.Code)
	return exec, nil
}

func copyUsed(src, dest *core.Bytecode) {
	if src.Used == nil {
		return
	}
	if dest.Used == nil {
		dest.Used = make(map[int32]byte)
	}
	for ikey := range src.Used {
		dest.Used[ikey] = 1
	}
}

func type2Code(itype *core.TypeObject) (retType core.Bcode) {
	switch itype.Original {
	case reflect.TypeOf(int64(0)):
		retType = core.TYPEINT
	case reflect.TypeOf(true):
		retType = core.TYPEBOOL
		//				case reflect.TypeOf(float64(0.0)):
		//					retType = core.STACKFLOAT
	case reflect.TypeOf('a'):
		retType = core.TYPECHAR
		//				case reflect.TypeOf(``):
		//					retType = core.STACKSTR
	}
	return retType
}

func genBytecode(ws *core.Workspace, idObj int32) *core.Bytecode {
	bcode := ws.Objects[idObj].GetCode()
	if ws.Objects[idObj].GetType() == core.ObjType {
		return nil
	}
	if bcode.Code != nil {
		return bcode
	}
	bcode.Code = make([]core.Bcode, 0, 64)
	bcode.Strings = make(map[string]uint16)
	cmd2Code(&Linker{Lex: ws.Objects[idObj].GetLex()},
		&ws.Objects[idObj].(*core.FuncObject).Block, bcode)
	bcode.Code = append(bcode.Code, core.END)
	//	fmt.Println(`CODE`, bcode.Code)
	return bcode
}
