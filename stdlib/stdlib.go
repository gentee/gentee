// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"github.com/gentee/gentee/compiler"
	"github.com/gentee/gentee/core"
)

type embedInfo struct {
	Func    interface{}
	InTypes string
	OutType string
}

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
		`Time:time`, `IsDir:bool`,
	})

	InitEmbed(ws)
	InitRuntime(ws)
	InitRegExp(ws)
	InitContext(ws)
	InitThread(ws)
	InitCrypto(ws)
	InitNetwork(ws)

	ws.IotaID = stdlib.NewConst(core.ConstIota, int64(0), false)
	stdlib.NewConst(core.ConstDepth, int64(1000), true)
	stdlib.NewConst(core.ConstCycle, int64(16000000), true)
	stdlib.NewConst(core.ConstVersion, core.Version, false)

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
	unitID, _ := compiler.Compile(ws, src, ``)
	ws.Units[0].NameSpace[`?Run`] = ws.Units[unitID].NameSpace[`?Run`]
	ws.Units[0].NameSpace[`?Start`] = ws.Units[unitID].NameSpace[`?Start`]
}

func InitEmbed(ws *core.Workspace) {
	for _, embed := range ws.Embedded {
		ws.StdLib().ImportEmbed(embed)
	}
}
