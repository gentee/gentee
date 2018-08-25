// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/workspace"
)

const (
	errNoFile = iota + 1
	errCompile
	errRun
	errResult
)

func main() {
	var (
		env      string
		testMode bool
		err      error
	)

	flag.StringVar(&env, "env", "", "environment variables")
	flag.BoolVar(&testMode, "t", false, "compare with #result")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Specify Gentee script file: ./gentee yourscript.g")
		os.Exit(errNoFile)
	}
	workspace := workspace.New()

	isError := func(code int) {
		if err != nil {
			fmt.Println(`ERROR:`, err.Error())
			os.Exit(code)
		}
	}

	for _, script := range files {
		var (
			result interface{}
			unit   *core.Unit
		)
		unit, err = workspace.CompileFile(script)
		isError(errCompile)
		result, err = workspace.Run(unit.Name)
		isError(errRun)
		resultStr := fmt.Sprint(result)
		if testMode {
			for _, line := range strings.Split(unit.Lexeme[0].Header, "\n") {
				ret := regexp.MustCompile(`\s*result\s*=\s*(.*)$`).FindStringSubmatch(line)
				if len(ret) == 2 {
					if ret[1] == resultStr {
						return
					}
				}
			}
			err = fmt.Errorf(`different test result %s`, resultStr)
			isError(errResult)
		}
		fmt.Println(resultStr)
	}
}
