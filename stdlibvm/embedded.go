// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import "github.com/gentee/gentee/core"

var Embedded = []core.Embed{
	{Func: strºBool, Return: core.TYPESTR, Params: []uint16{core.TYPEBOOL}},
	{Func: strºInt, Return: core.TYPESTR, Params: []uint16{core.TYPEINT}},
	{Func: intºStr, Return: core.TYPEINT, Params: []uint16{core.TYPESTR}, CanError: true},
	{Func: boolºInt, Return: core.TYPEBOOL, Params: []uint16{core.TYPEINT}},
	{Func: ExpStrºInt, Return: core.TYPESTR, Params: []uint16{core.TYPESTR,core.TYPEINT}},
	{Func: ExpStrºBool, Return: core.TYPESTR, Params: []uint16{core.TYPESTR,core.TYPEBOOL}},
}
