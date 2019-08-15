// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/gentee/gentee/core"
)

type Linker struct {
	Blocks []*core.CmdBlock
	Lex    *core.Lex
}

type Int32Slice []int32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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
		Init:  bcode.Init,
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
		for _, id := range usedCode.Init {
			exec.Init = append(exec.Init, id)
		}
		rebuild := make([]uint16, len(usedCode.Strings))
		for key, curInd := range usedCode.Strings {
			if ind, ok = bcode.Strings[key]; !ok {
				ind = uint16(len(bcode.Strings))
				bcode.Strings[key] = ind
			}
			rebuild[curInd] = ind
		}
		for _, off := range usedCode.StrOffset {
			cur := exec.Code[shift+off] >> 16
			exec.Code[shift+off] = core.Bcode((uint32(rebuild[cur]) << 16) | core.PUSHSTR)
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
	sort.Sort(Int32Slice(exec.Init))
	if len(exec.Init) > 0 && exec.Init[0] != ws.IotaID {
		exec.Init = append([]int32{ws.IotaID}, exec.Init...)
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
	case reflect.TypeOf(``):
		retType = core.TYPESTR
	}
	return retType
}

func getPos(linker *Linker, cmd core.ICmd, out *core.Bytecode) {
	var (
		ok   bool
		name string
	)
	if obj := cmd.GetObject(); obj != nil {
		name = obj.GetName()
	}
	line, column := linker.Lex.LineColumn(cmd.GetToken())
	if _, ok = out.Strings[linker.Lex.Path]; !ok {
		out.Strings[linker.Lex.Path] = uint16(len(out.Strings))
	}
	if _, ok = out.Strings[name]; !ok {
		out.Strings[name] = uint16(len(out.Strings))
	}
	out.Pos = append(out.Pos, core.CodePos{
		Offset: int32(len(out.Code) - 1),
		Path:   out.Strings[linker.Lex.Path],
		Name:   out.Strings[name],
		Line:   uint16(line),
		Column: uint16(column),
	})
}

func genBytecode(ws *core.Workspace, idObj int32) *core.Bytecode {
	var (
		block   core.ICmd
		isConst bool
	)
	bcode := ws.Objects[idObj].GetCode()
	if bcode.Code != nil {
		return bcode
	}
	bcode.Code = make([]core.Bcode, 0, 64)
	bcode.Strings = make(map[string]uint16)
	switch ws.Objects[idObj].GetType() {
	case core.ObjType:
		return nil
	case core.ObjFunc:
		block = &ws.Objects[idObj].(*core.FuncObject).Block
	case core.ObjConst:
		constObj := ws.Objects[idObj].(*core.ConstObject)
		block = constObj.Exp
		if constObj.Iota != core.NotIota {
			bcode.Code = append(bcode.Code, core.Bcode((constObj.Iota+1)<<16)|core.IOTA)
		}
		isConst = true
	}
	cmd2Code(&Linker{Lex: ws.Objects[idObj].GetLex()}, block, bcode)
	if isConst {
		bcode.Code = append(bcode.Code, (type2Code(block.GetResult())<<16)|core.RET)
	} else {
		bcode.Code = append(bcode.Code, core.END)
	}
	//	fmt.Println(`CODE`, bcode.Code)
	return bcode
}
