// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/shlex"

	"github.com/gentee/gentee/core"
)

// ArgCount returns the count of command-line parameters
func ArgCount(rt *Runtime) int64 {
	return int64(len(rt.Owner.Settings.CmdLine))
}

// ArgºStr returns the value of the command-line option
func ArgºStr(rt *Runtime, flag string) string {
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
func ArgºStrStr(rt *Runtime, flag, def string) string {
	found, list := args(rt.Owner.Settings.CmdLine, flag)
	if !found {
		return def
	}
	if len(list) == 0 {
		return ``
	}
	return list[0]
}

// ArgºStrInt returns the number value of the command-line option or the default value
func ArgºStrInt(rt *Runtime, flag string, def int64) (int64, error) {
	return strconv.ParseInt(ArgºStrStr(rt, flag, strconv.FormatInt(def, 10)), 10, 64)
}

// Args returns the command-line parameters
func Args(rt *Runtime) *core.Array {
	out := core.NewArray()
	for _, par := range rt.Owner.Settings.CmdLine {
		out.Data = append(out.Data, par)
	}
	return out
}

// ArgsºStr returns the value list of command-line option
func ArgsºStr(rt *Runtime, flag string) *core.Array {
	out := core.NewArray()
	_, list := args(rt.Owner.Settings.CmdLine, flag)
	for _, item := range list {
		out.Data = append(out.Data, item)
	}
	return out
}

// ArgsTail returns the list of command-line parameters
func ArgsTail(rt *Runtime) *core.Array {
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
		err = fmt.Errorf(ErrorText(ErrPlatform))
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
		err = fmt.Errorf(ErrorText(ErrPlatform))
	}
	return err
}

// IsArgºStr returns true if the options is present
func IsArgºStr(rt *Runtime, flag string) int64 {
	found, _ := args(rt.Owner.Settings.CmdLine, flag)
	if found {
		return 1
	}
	return 0
}

// SplitCmdLine splits the command line parameters to the array of strings
func SplitCmdLine(cmdline string) (ret *core.Array, err error) {
	var pars []string
	ret = core.NewArray()
	pars, err = shlex.Split(cmdline)
	if err != nil {
		return
	}
	for _, par := range pars {
		ret.Data = append(ret.Data, par)
	}
	return
}

// sysRun executes the process.
func sysRun(cmd string, start int64, stdin *core.Buffer, stdout *core.Buffer, stderr *core.Buffer,
	args *core.Array) error {
	var (
		pars                  []string
		bufOut, bufIn, bufErr bytes.Buffer
	)
	for _, arg := range args.Data {
		pars = append(pars, fmt.Sprint(arg))
	}
	command := exec.Command(cmd, pars...)
	if stdin.Data == nil {
		command.Stdin = os.Stdin
	} else {
		bufIn = bytes.Buffer{}
		bufIn.Write(stdin.Data)
		command.Stdin = &bufIn
	}
	if stdout.Data == nil {
		command.Stdout = os.Stdout
	} else {
		bufOut = bytes.Buffer{}
		command.Stdout = &bufOut
	}
	if stderr.Data == nil {
		command.Stderr = os.Stderr
	} else {
		bufErr = bytes.Buffer{}
		command.Stderr = &bufErr
	}
	if start == 1 {
		return command.Start()
	}
	err := command.Run()
	if stdout.Data != nil {
		stdout.Data = bufOut.Bytes()
	}
	if stderr.Data != nil {
		stderr.Data = bufErr.Bytes()
	}
	return err
}
