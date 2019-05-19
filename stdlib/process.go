// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"os/exec"
	"runtime"
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
