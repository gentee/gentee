// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

// InitProcess appends stdlib process functions to the virtual machine
func InitProcess(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		OpenºStr,     // Open( str )
		OpenWithºStr, // OpenWith( str, str )
	} {
		vm.StdLib().NewEmbed(item)
	}
	for _, item := range []embedInfo{
		{ArgCount, ``, `int`},          // ArgCount() int
		{ArgºStr, `str`, `str`},        // Arg(str) str
		{ArgºStrStr, `str,str`, `str`}, // Arg(str, str) str
		{ArgºStrInt, `str,int`, `int`}, // Arg(str, int) int
		{Args, ``, `arr.str`},          // Args() arr.str
		{ArgsºStr, `str`, `arr.str`},   // Args(str) arr.str
		{ArgsTail, ``, `arr.str`},      // ArgsTail() arr.str
		{IsArgºStr, `str`, `bool`},     // IsArg(str) bool
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// Args returns the command-line parameters
func Args(rt *core.RunTime) *core.Array {
	out := core.NewArray()
	for _, par := range rt.CmdLine {
		out.Data = append(out.Data, par)
	}
	return out
}

// ArgCount returns the count of command-line parameters
func ArgCount(rt *core.RunTime) int64 {
	return int64(len(rt.CmdLine))
}

// ArgºStr returns the value of the command-line option
func ArgºStr(rt *core.RunTime, flag string) string {
	return ArgºStrStr(rt, flag, ``)
}

func args(cmdLine []string, flag string) (bool, []string) {
	var ret []string
	lenf := len(flag)
	lenCmd := len(cmdLine)
	if lenf == 0 || len(strings.Trim(flag, `-`)) == 0 {
		var tail int
		for i := 0; i < lenCmd; i++ {
			if strings.HasPrefix(cmdLine[i], `-`) {
				tail = i + 1
				if len(strings.Trim(cmdLine[i], `-`)) == 0 {
					break
				}
			}
		}
		if tail < lenCmd {
			return true, cmdLine[tail:]
		}
		return false, nil
	}
	if flag[0] != '-' {
		flag = `-` + flag
		lenf++
	}
	for i, arg := range cmdLine {
		if arg == `-` {
			break
		}
		if strings.HasPrefix(arg, flag) {
			if len(arg) == lenf {
				for k := i + 1; k < lenCmd && cmdLine[k][0] != '-'; k++ {
					ret = append(ret, cmdLine[k])
				}
				return true, ret
			}
			if arg[lenf] == '=' || arg[lenf] == ':' {
				if len(arg) > lenf+1 {
					val := arg[lenf+1:]
					if len(val) > 0 && (val[0] == val[len(val)-1] && (val[0] == '"' || val[0] == '\'')) {
						val = val[1 : len(val)-1]
					}
					ret = []string{val}
				}
				return true, ret
			}
		}
	}
	return false, nil
}

// ArgºStrStr returns the value of the command-line option or the default value
func ArgºStrStr(rt *core.RunTime, flag, def string) string {
	found, list := args(rt.CmdLine, flag)
	if !found {
		return def
	}
	if len(list) == 0 {
		return ``
	}
	return list[0]
}

// ArgºStrInt returns the number value of the command-line option or the default value
func ArgºStrInt(rt *core.RunTime, flag string, def int64) (int64, error) {
	return strconv.ParseInt(ArgºStrStr(rt, flag, strconv.FormatInt(def, 10)), 10, 64)
}

// ArgsºStr returns the value list of command-line option
func ArgsºStr(rt *core.RunTime, flag string) *core.Array {
	out := core.NewArray()
	_, list := args(rt.CmdLine, flag)
	for _, item := range list {
		out.Data = append(out.Data, item)
	}
	return out
}

// ArgsTail returns the list of command-line parameters
func ArgsTail(rt *core.RunTime) *core.Array {
	return ArgsºStr(rt, ``)
}

// OpenºStr runs corresponding application with the specified file.
func OpenºStr(fname string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", fname).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", fname).Start()
	case "darwin":
		err = exec.Command("open", fname).Start()
	default:
		err = fmt.Errorf(core.ErrorText(core.ErrPlatform))
	}
	return err
}

// OpenWithºStr runs the application with the specified file.
func OpenWithºStr(app, fname string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command(app, fname).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", app, strings.Replace(fname, "&", "^&", -1)).Start()
	case "darwin":
		err = exec.Command("open", "-a", app, fname).Start()
	default:
		err = fmt.Errorf(core.ErrorText(core.ErrPlatform))
	}
	return err
}

// IsArgºStr returns true if the options is present
func IsArgºStr(rt *core.RunTime, flag string) bool {
	found, _ := args(rt.CmdLine, flag)
	return found
}
