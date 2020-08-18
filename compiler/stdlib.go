// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"hash/crc64"

	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/vm"
)

// InitStdlib appends stdlib types and functions to the virtual machine
func InitStdlib(ws *core.Workspace) {
	stdlib := ws.InitUnit()
	stdlib.Pub = core.PubAll
	ws.Units = append(ws.Units, stdlib)
	ws.UnitNames[core.DefName] = len(ws.Units) - 1
	InitTypes(ws)
	NewStructType(ws, `time`, []string{
		`Year:int`, `Month:int`, `Day:int`,
		`Hour:int`, `Minute:int`, `Second:int`,
		`UTC:bool`,
	})
	NewStructType(ws, `finfo`, []string{
		`Name:str`, `Size:int`, `Mode:int`,
		`Time:time`, `IsDir:bool`, `Dir:str`,
	})
	NewStructType(ws, `trace`, []string{
		`Path:str`, `Entry:str`, `Func:str`, `Line:int`, `Pos:int`,
	})
	InitEmbed(ws)

	ws.IotaID = stdlib.NewConst(core.ConstIota, int64(0), false)
	stdlib.NewConst(core.ConstDepth, int64(1000), true)
	stdlib.NewConst(core.ConstCycle, int64(16000000), true)
	stdlib.NewConst(core.ConstScript, ``, true)
	stdlib.NewConst(core.ConstVersion, core.Version, false)

	// For flag param of ReadDir(str, int, str)
	stdlib.NewConst(core.ConstRecursive, int64(vm.Recursive), false)
	stdlib.NewConst(core.ConstOnlyFiles, int64(vm.OnlyFiles), false)
	stdlib.NewConst(core.ConstRegExp, int64(vm.RegExp), false)

	src := `
	pub	func Run(str cmd, str args...) {
		buf ? stdin &= sysBufNil()
		buf ? stdout &= sysBufNil()
		buf ? stderr &= sysBufNil()
		sysRun(cmd, false, stdin, stdout, stderr, args)
	  }
	  
	pub func Start(str cmd, str args...) {
		buf ? stdin &= sysBufNil()
		buf stdout &= sysBufNil()
		buf stderr &= sysBufNil()
		sysRun(cmd, true, stdin, stdout, stderr, args)
	  }
	`
	unitID, _ := Compile(ws, src, ``)
	ws.Units[0].NameSpace[`?Run`] = ws.Units[unitID].NameSpace[`?Run`]
	ws.Units[0].NameSpace[`?Start`] = ws.Units[unitID].NameSpace[`?Start`]
}

// InitEmbed imports in-line functions
func InitEmbed(ws *core.Workspace) {
	var crc string

	for i, embed := range ws.Embedded {
		crc += fmt.Sprintf("%s(%s)%s", embed.Name, embed.Pars, embed.Ret)
		ws.StdLib().ImportEmbed(embed)
		if i == vm.StdLibCount-1 {
			vm.CRCStdlib = crc64.Checksum([]byte(crc), crc64.MakeTable(crc64.ECMA))
			crc = ``
		}
	}
	if len(crc) > 0 {
		vm.CRCCustom = crc64.Checksum([]byte(crc), crc64.MakeTable(crc64.ECMA))
	}
}
