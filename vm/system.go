// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ClearCarriage deletes output before carriage return character
func ClearCarriage(input string) string {
	var start int
	runes := []rune(string(input))
	out := make([]rune, 0, len(runes))
	for _, char := range []rune(runes) {
		if char == 0xd {
			out = out[:start]
		} else {
			out = append(out, char)
			if char == 0xa {
				start = len(out)
			}
		}
	}
	return string(out)
}

// Command executes the command line
func Command(rt *Runtime, cmdLine string) error {
	if rt.Owner.Settings.IsPlayground {
		return fmt.Errorf(ErrorText(ErrPlayRun))
	}
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
func CommandOutput(rt *Runtime, cmdLine string) (string, error) {
	if rt.Owner.Settings.IsPlayground {
		return ``, fmt.Errorf(ErrorText(ErrPlayRun))
	}
	cmd, err := splitCmdLine(cmdLine)
	if err != nil {
		return ``, err
	}
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(err.Error())
	}
	return ClearCarriage(string(stdout)), err
}

// GetEnv return the value of the environment variable
func GetEnv(name string) string {
	return os.Getenv(name)
}

// SetEnv assign the value to the environment variable
func SetEnv(rt *Runtime, name string, value interface{}) (string, error) {
	if rt.Owner.Settings.IsPlayground {
		return ``, fmt.Errorf(ErrorText(ErrPlayEnv))
	}
	ret := fmt.Sprint(value)
	err := os.Setenv(name, ret)
	return ret, err
}

// SetEnvBool assign the value to the environment variable
func SetEnvBool(rt *Runtime, name string, value int64) (string, error) {
	if rt.Owner.Settings.IsPlayground {
		return ``, fmt.Errorf(ErrorText(ErrPlayEnv))
	}
	ret := strºBool(value)
	err := os.Setenv(name, ret)
	return ret, err
}

// UnsetEnv unsets the environment variable
func UnsetEnv(rt *Runtime, name string) error {
	if rt.Owner.Settings.IsPlayground {
		return fmt.Errorf(ErrorText(ErrPlayEnv))
	}
	return os.Unsetenv(name)
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
		return nil, fmt.Errorf(ErrorText(ErrQuoteCommand))
	}
	if offset < len(input) {
		cmds = append(cmds, string(input[offset:]))
	}
	if len(cmds) == 0 {
		return nil, fmt.Errorf(ErrorText(ErrEmptyCommand))
	}
	if cmds[0] == `echo` && runtime.GOOS == "windows" {
		cmds[0] = `cmd.exe`
		cmds = append(cmds[:1], append([]string{`/C`, `echo`}, cmds[1:]...)...)
	}
	return exec.Command(cmds[0], cmds[1:]...), nil
}
