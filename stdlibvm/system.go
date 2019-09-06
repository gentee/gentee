// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gentee/gentee/core"
)

func splitCmdLine(cmdLine string) (*exec.Cmd, error) {
	var (
		cmds      []string
		offset, i int
		quote, ch rune
	)
	input := []rune(strings.TrimSpace(cmdLine))
	newPar := func(i int) {
		if offset < i {
			end := i
			if quote != 0 && input[offset-1] != quote {
				end++
			}
			cmds = append(cmds, string(input[offset:end]))
		}
		offset = i + 1
	}
	for i, ch = range input {
		if quote != 0 {
			if ch == quote {
				newPar(i)
				quote = 0
			}
			continue
		}
		if ch == '\'' || ch == '"' || ch == '`' {
			quote = ch
			if offset == i {
				offset = i + 1
			}
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
	if cmds[0] == `echo` && runtime.GOOS == "windows" {
		cmds[0] = `cmd.exe`
		cmds = append(cmds[:1], append([]string{`/C`, `echo`}, cmds[1:]...)...)
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
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf(err.Error())
	}
	return err
}

// CommandOutput executes the command line and returns the standard output
func CommandOutput(cmdLine string) (string, error) {
	cmd, err := splitCmdLine(cmdLine)
	if err != nil {
		return ``, err
	}
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(err.Error())
	}
	return string(stdout), err
}

// GetEnv return the value of the environment variable
func GetEnv(name string) string {
	return os.Getenv(name)
}

// SetEnv assign the value to the environment variable
func SetEnv(name string, value interface{}) (string, error) {
	ret := fmt.Sprint(value)
	err := os.Setenv(name, ret)
	return ret, err
}

// SetEnvBool assign the value to the environment variable
func SetEnvBool(name string, value int64) (string, error) {
	ret := StrºBool(value)
	err := os.Setenv(name, ret)
	return ret, err
}
