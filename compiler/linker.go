// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/vm"
)

// BlockInfo describes block in the linker
type BlockInfo struct {
	Block   *core.CmdBlock
	Vars    []int
	IsLocal bool
}

// Linker is the main structure of the linker
type Linker struct {
	Blocks []BlockInfo
	Lex    *core.Lex
}

// Int32Slice is a slice of int32
type Int32Slice []int32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Link creates a bytecode
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
		Code:    append([]core.Bcode{}, bcode.Code...),
		Funcs:   make(map[int32]int32),
		Init:    bcode.Init,
		Pos:     bcode.Pos,
		Structs: bcode.StructsList,
		Path:    unit.Lexeme.Path,

		CRCStdlib: vm.CRCStdlib,
		CRCCustom: vm.CRCCustom,
	}
	if len(exec.Path) == 0 {
		exec.Path = unit.Name
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
		var retype []uint16
		rebuild := make([]uint16, len(usedCode.Strings))
		if len(usedCode.Structs) > 0 {
			retype = make([]uint16, len(usedCode.Structs))
		}
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
		countStructs := len(exec.Structs)
		for key, curInd := range usedCode.Structs {
			if ind, ok = bcode.Structs[key]; !ok {
				ind = uint16(len(exec.Structs))
				bcode.Structs[key] = ind
				exec.Structs = append(exec.Structs, usedCode.StructsList[curInd])
			}
			retype[curInd] = ind
		}
		for ; countStructs < len(exec.Structs); countStructs++ {
			sinfo := &exec.Structs[countStructs]
			for i, field := range sinfo.Fields {
				if field >= core.TYPESTRUCT {
					sinfo.Fields[i] = retype[(field-core.TYPESTRUCT)>>8]<<8 + core.TYPESTRUCT
				}
			}
		}
		for _, off := range usedCode.StructsOffset {
			var isleft bool
			if off < 0 {
				isleft = true
				off = -off
			}
			cur := exec.Code[shift+off]
			left := int32(cur >> 16)
			right := int32(cur & 0xffff)
			if isleft {
				left = int32(retype[(left-core.TYPESTRUCT)>>8])<<8 + core.TYPESTRUCT
			} else {
				right = int32(retype[(right-core.TYPESTRUCT)>>8])<<8 + core.TYPESTRUCT
			}
			exec.Code[shift+off] = core.Bcode(left<<16 | right)
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
	//fmt.Println(`Structs`, exec.Structs)
	//fmt.Println(`NAMES`, exec.Pos)
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

func structOffset(out *core.Bytecode, shift int) {
	out.StructsOffset = append(out.StructsOffset, int32(shift))
}

func initBlock(linker *Linker, cmd *core.CmdBlock, out *core.Bytecode) (BlockInfo, []core.Bcode) {
	push := func(pars ...core.Bcode) {
		out.Code = append(out.Code, pars...)
	}
	bInfo := BlockInfo{
		Block: cmd,
	}
	flags := int32(out.BlockFlags)
	out.BlockFlags = 0
	if len(cmd.Vars) > 0 {
		flags |= core.BlVars
	}
	if cmd.ParCount > 0 {
		flags |= core.BlPars
	}
	//	push(core.Bcode(cmd.ParCount<<16)|core.INITVARS, core.Bcode(len(cmd.Vars)))
	push(core.Bcode(flags<<16) | core.INITVARS)
	if flags&core.BlBreak != 0 {
		push(0)
	}
	if flags&core.BlContinue != 0 {
		push(0)
	}
	if flags&core.BlTry != 0 {
		push(0)
	}
	if flags&core.BlRecover != 0 {
		push(0)
	}
	if flags&core.BlRetry != 0 {
		push(0)
	}
	if cmd.ParCount > 0 || flags&core.BlVars != 0 {
		push(core.Bcode(cmd.ParCount<<16 | len(cmd.Vars)))
	}
	var types []core.Bcode
	if len(cmd.Vars) > 0 {
		types = make([]core.Bcode, len(cmd.Vars))
		var sInt, sStr, sFloat, sAny int
		bInfo.Vars = make([]int, len(cmd.Vars))
		lenCode := len(out.Code)
		for i, ivar := range cmd.Vars {
			types[i] = type2Code(ivar, out)
			switch types[i] & 0xf {
			case core.STACKSTR:
				bInfo.Vars[i] = sStr
				sStr++
			case core.STACKFLOAT:
				bInfo.Vars[i] = sFloat
				sFloat++
			case core.STACKANY:
				bInfo.Vars[i] = sAny
				sAny++
				if types[i] >= core.TYPESTRUCT {
					structOffset(out, lenCode+i)
				}
			default:
				bInfo.Vars[i] = sInt
				sInt++
			}
		}
		push(types...)
	}
	linker.Blocks = append(linker.Blocks, bInfo)
	return bInfo, types
}

func type2Code(itype *core.TypeObject, out *core.Bytecode) (retType core.Bcode) {
	switch itype.Original {
	case reflect.TypeOf(int64(0)):
		retType = core.TYPEINT
	case reflect.TypeOf(true):
		retType = core.TYPEBOOL
	case reflect.TypeOf(float64(0.0)):
		retType = core.TYPEFLOAT
	case reflect.TypeOf('a'):
		retType = core.TYPECHAR
	case reflect.TypeOf(``):
		retType = core.TYPESTR
	case reflect.TypeOf(core.Array{}):
		retType = core.TYPEARR
	case reflect.TypeOf(core.Range{}):
		retType = core.TYPERANGE
	case reflect.TypeOf(core.Map{}):
		retType = core.TYPEMAP
	case reflect.TypeOf(core.Buffer{}):
		retType = core.TYPEBUF
	case reflect.TypeOf(core.Fn{}):
		retType = core.TYPEFUNC
	case reflect.TypeOf(core.RuntimeError{}):
		retType = core.TYPEERROR
	case reflect.TypeOf(core.Set{}):
		retType = core.TYPESET
	case reflect.TypeOf(core.Obj{}):
		retType = core.TYPEOBJ
	case reflect.TypeOf(core.Struct{}):
		typeName := itype.GetName()
		var (
			ind uint16
			ok  bool
		)
		if ind, ok = out.Structs[typeName]; !ok {
			sInfo := core.StructInfo{
				Name:   typeName,
				Fields: make([]uint16, len(itype.Custom.Types)),
				Keys:   make([]string, len(itype.Custom.Types)),
			}
			var self []int
			for i, item := range itype.Custom.Types {
				if item == itype {
					self = append(self, i)
				} else {
					sInfo.Fields[i] = uint16(type2Code(item, out))
				}
			}
			for name, i := range itype.Custom.Fields {
				sInfo.Keys[i] = name
			}
			ind = uint16(len(out.StructsList))
			if ind == 0 {
				out.Structs = make(map[string]uint16)
			}
			for _, iself := range self {
				sInfo.Fields[iself] = core.TYPESTRUCT + (ind << 8)
			}
			out.Structs[typeName] = ind
			out.StructsList = append(out.StructsList, sInfo)
		}
		retType = core.Bcode(core.TYPESTRUCT + (ind << 8))
		//		fmt.Printf("STRUCT%s %x %d %v\n", typeName, retType, ind, out.StructsList)

		//	case reflect.TypeOf(core.KeyValue{}):
		//		retType = core.TYPEKEYVALUE
	default:
		fmt.Printf("type2Code %v\n", itype.Original)
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
	type2Code(ws.StdLib().FindType(`trace`).(*core.TypeObject), bcode)
	type2Code(ws.StdLib().FindType(`time`).(*core.TypeObject), bcode)
	type2Code(ws.StdLib().FindType(`finfo`).(*core.TypeObject), bcode)
	type2Code(ws.StdLib().FindType(`hinfo`).(*core.TypeObject), bcode)

	cmd2Code(&Linker{Lex: ws.Objects[idObj].GetLex()}, block, bcode)
	if isConst {
		resType := type2Code(block.GetResult(), bcode)
		bcode.Code = append(bcode.Code, (resType<<16)|core.RET)
		if resType >= core.TYPESTRUCT {
			structOffset(bcode, -len(bcode.Code)+1)
		}
	} else {
		bcode.Code = append(bcode.Code, core.END)
	}
	//	fmt.Println(`CODE`, bcode.Code)
	return bcode
}
