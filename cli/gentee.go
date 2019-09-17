// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	workspace "github.com/gentee/gentee"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/vm"
)

const (
	errNoFile = iota + 1
	errCompile
	errRun
	errResult
)

func main() {
	var (
		env           string
		testMode, ver bool
		err           error
	)

	flag.StringVar(&env, "env", "", "environment variables")
	flag.BoolVar(&testMode, "t", false, "compare with #result")
	flag.BoolVar(&ver, "ver", false, "compare with #result")
	flag.Parse()

	workspace := workspace.New()
	if ver {
		fmt.Println(workspace.Version())
		return
	}

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Specify Gentee script file: ./gentee yourscript.g")
		os.Exit(errNoFile)
	}

	isError := func(code int) {
		if err != nil {
			fmt.Print(`ERROR`)
			if errTrace, ok := err.(*vm.RuntimeError); ok {
				fmt.Printf(" #%d: %s\n", errTrace.ID, err.Error())
				for _, trace := range errTrace.Trace {
					path := trace.Path
					dirs := strings.Split(filepath.ToSlash(path), `/`)
					if len(dirs) > 3 {
						path = `...` + path[len(path)-len(strings.Join(dirs[len(dirs)-3:], `/`))-1:]
					}
					fmt.Printf("%s [%d:%d] %s -> %s\n", path, trace.Line, trace.Pos, trace.Entry, trace.Func)
				}
				code = errTrace.ID
			} else {
				fmt.Println(`:`, err.Error())
			}
			os.Exit(code)
		}
	}
	script := files[0]
	var (
		result interface{}
		unitID int
		exec   *core.Exec
	)
	exec, unitID, err = workspace.CompileFile(script)
	isError(errCompile)
	result, err = vm.Run(exec, vm.Settings{CmdLine: files[1:]})
	//	workspace.CmdLine(files[1:]...)
	//	result, err = workspace.Run(unitID)
	isError(errRun)
	resultStr := fmt.Sprint(result)
	if testMode {
		for _, line := range strings.Split(workspace.Unit(unitID).Lexeme[0].Header, "\n") {
			ret := regexp.MustCompile(`\s*result\s*=\s*(.*)$`).FindStringSubmatch(strings.TrimSpace(line))
			if len(ret) == 2 {
				if ret[1] == strings.TrimSpace(resultStr) {
					return
				}
			}
		}
		err = fmt.Errorf(`different test result %s`, resultStr)
		isError(errResult)
	}
	if result != nil {
		fmt.Println(resultStr)
	}
}
