// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gentee/gentee/core"
)

// InitSystem appends stdlib system functions to the virtual machine
func InitSystem(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		Command,       // $( str )
		CommandOutput, // $( str )
		GetEnv,        // get environment variable
	} {
		vm.StdLib().NewEmbed(item)
	}
}

func splitCmdLine(cmdLine string) (*exec.Cmd, error) {
	var (
		cmds      []string
		offset, i int
		quote, ch rune
	)
	input := []rune(strings.TrimSpace(cmdLine))
	newPar := func(i int) {
		if offset < i {
			cmds = append(cmds, string(input[offset:i]))
		}
		offset = i + 1
	}
	for i, ch = range input {
		if quote != 0 {
			if ch == quote {
				quote = 0
				newPar(i)
			}
			continue
		}
		if ch == '\'' || ch == '"' || ch == '`' {
			quote = ch
			offset = i + 1
			continue
		}
		if ch == ' ' {
			newPar(i)
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf(core.ErrorText(core.ErrQuoteCommand))
	}
	if offset < len(input) {
		cmds = append(cmds, string(input[offset:]))
	}
	if len(cmds) == 0 {
		return nil, fmt.Errorf(core.ErrorText(core.ErrEmptyCommand))
	}
	return exec.Command(cmds[0], cmds[1:]...), nil
}

// Command executes the command line
func Command(cmdLine string) error {
	cmd, err := splitCmdLine(cmdLine)
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// CommandOutput executes the command line and returns the standard output
func CommandOutput(cmdLine string) (string, error) {
	cmd, err := splitCmdLine(cmdLine)
	if err != nil {
		return ``, err
	}
	stdout, err := cmd.CombinedOutput()
	return string(stdout), err
}

// GetEnv return the value of the environment variable
func GetEnv(name string) string {
	return `%` + name
}
