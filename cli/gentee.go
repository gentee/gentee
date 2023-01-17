// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	gentee "github.com/gentee/gentee"
	"github.com/gentee/gentee/vm"
)

const (
	errNoFile = iota + 1
	errCompile
	errRun
	errResult
	errPlaceholder5
	errPlaceholder6
	errPlaceholder7
	errUndefined
)

type CodedError struct {
	Code int
	Err  error
}

func (c *CodedError) Error() string {
	if c.Err != nil {
		return c.Error()
	}
	return fmt.Sprintf("Unspecified Error with Code: %d", c.Code)
}

func (c *CodedError) Unwrap() error {
	return c.Err
}

type CommandArgs struct {
	Env      string
	TestMode bool
	Ver      bool

	Execute string
	Stdin   bool
	Output  string
}

func (c *CommandArgs) Parse() *CommandArgs {
	flag.StringVar(&c.Env, "env", "", "environment variables")
	flag.BoolVar(&c.TestMode, "t", false, "compare with #result")
	flag.BoolVar(&c.Ver, "ver", false, "print version")
	flag.StringVar(&c.Execute, "e", "", "Execute the string")
	flag.BoolVar(&c.Stdin, "p", false, "read from stdin")
	flag.StringVar(&c.Output, "o", "", "output to file (default stdout)")
	flag.Parse()
	return c
}

type Cli struct {
	workspace *gentee.Gentee
	args      CommandArgs
}

func (c *Cli) Init() *Cli {
	c.workspace = gentee.New()
	c.args.Parse()
	return c
}

func (c *Cli) exec_RunFile(w io.Writer) error {
	var params []string
	args := flag.Args()
	switch len(args) {
	case 0:
		fmt.Println("Specify Gentee script file: ./gentee yourscript.g")
		os.Exit(errNoFile)
	case 1:
	default:
		params = args[1:]
	}
	file := flag.Arg(0)
	var (
		result   interface{}
		unitID   int
		exec     *gentee.Exec
		settings gentee.Settings
		err      error
	)
	exec, unitID, err = c.workspace.CompileFile(file)
	if err != nil {
		return codedError(err, errCompile)
	}
	settings.CmdLine = params
	result, err = exec.Run(settings)
	if err != nil {
		return codedError(err, errRun)
	}
	resultStr := fmt.Sprint(result)
	if c.args.TestMode {
		ret := c.workspace.Unit(unitID).GetHeader(`result`)
		if len(ret) > 0 && ret == strings.TrimSpace(resultStr) {
			return nil
		}
		err = fmt.Errorf(`different test result %s`, resultStr)
		return codedError(err, errResult)
	}
	if result != nil {
		fmt.Fprintf(w, resultStr)
	}
	return nil
}

func (c *Cli) exec_Ver() {
	fmt.Println(gentee.Version())
}

func (c *Cli) Exec() {
	stderr := os.Stderr
	err := c.exec()
	if err != nil {
		if coded, ok := err.(*CodedError); ok {
			fmt.Fprintf(stderr, "%s", coded.Unwrap())
			os.Exit(coded.Code)
		}
		fmt.Fprintf(stderr, "%s", err)
		os.Exit(errUndefined)
	}
	os.Exit(0)
}
func (c *Cli) exec() error {
	w := os.Stdout
	if c.args.Output != "" {
		fp, err := os.OpenFile(c.args.Output, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0o640)
		if err != nil {
			return err
		}
		defer fp.Close()
		w = fp
	}
	switch {
	case c.args.Ver:
		c.exec_Ver()
	default:
		return c.exec_RunFile(w)
	}
	return nil
}

func main() {
	cli := new(Cli).Init()
	if cli.args.Ver {
		return
	}
	cli.exec()
}
func codedError(err error, code int) error {
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
		return &CodedError{Err: err, Code: code}
	}
	return nil
}
