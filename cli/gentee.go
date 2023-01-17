// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	gentee "github.com/gentee/gentee"
)

type Cli struct {
	workspace *gentee.Gentee
	args      CommandArgs
}

func (c *Cli) Init() *Cli {
	c.workspace = gentee.New()
	c.args.Parse()
	return c
}

type CommandArgs struct {
	Env      string
	TestMode bool
	Ver      bool

	Execute string
	Stdin   bool
}

func (c *CommandArgs) Parse() *CommandArgs {
	flag.StringVar(&c.Env, "env", "", "environment variables")
	flag.BoolVar(&c.TestMode, "t", false, "compare with #result")
	flag.BoolVar(&c.Ver, "ver", false, "print version")
	flag.StringVar(&c.Execute, "e", "", "Execute the string")
	flag.BoolVar(&c.Stdin, "p", false, "read from stdin")
	flag.Parse()
	return c
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
	switch {
	case c.args.Ver:
		c.exec_Ver()
	case c.args.Execute != "":
		return c.exec_RunString(w, c.args.Execute)
	case c.args.Stdin:
		return c.exec_RunStdin(w)
	default:
		return c.exec_RunFile(w)
	}
	return nil
}

func (c *Cli) exec_Ver() {
	fmt.Println(gentee.Version())
}

func (c *Cli) exec_RunString(w io.Writer, str string) error {
	params := flag.Args()
	var (
		result   interface{}
		unitID   int
		exec     *gentee.Exec
		settings gentee.Settings
		err      error
	)
	exec, unitID, err = c.workspace.Compile(str, "stdin")
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
		return codedError(fmt.Errorf(`different test result %s`, resultStr), errResult)
	}
	if result != nil {
		fmt.Fprintf(w, resultStr)
	}
	return nil

}
func (c *Cli) exec_RunStdin(w io.Writer) error {
	bts, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	return c.exec_RunString(w, string(bts))
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
		return codedError(fmt.Errorf(`different test result %s`, resultStr), errResult)
	}
	if result != nil {
		fmt.Fprintf(w, resultStr)
	}
	return nil
}

func main() {
	cli := new(Cli).Init()
	cli.Exec()
}
