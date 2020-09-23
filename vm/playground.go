// Copyright 2020 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentee/gentee/core"
)

type Playground struct {
	Path         string // path to the temporary folder if it's empty then TempDir is used.
	AllSizeLimit int    // all files size limit. In default, 10MB
	FilesLimit   int    // count of files limit. In default, 1000
	SizeLimit    int    // file size limit. In default, 5MB
}

/*type FSFile struct {
	IsDir bool
	Size  int
}*/

type PlaygroundFS struct {
	Size  int // summary of files size
	Files map[string]int
}

// InitPlayground inits playground settings
func InitPlayground(settings *Settings) (err error) {
	if len(settings.Playground.Path) == 0 {
		settings.Playground.Path = TempDir()
	}
	if settings.Playground.Path, err = filepath.Abs(settings.Playground.Path); err != nil {
		return
	}
	settings.Playground.Path = filepath.Join(settings.Playground.Path, strings.ToLower(core.RandName()))
	if err = os.MkdirAll(settings.Playground.Path, os.ModePerm); err != nil {
		return
	}
	if settings.Playground.AllSizeLimit == 0 {
		settings.Playground.AllSizeLimit = 10 << 20 // 10MB
	}
	if settings.Playground.FilesLimit == 0 {
		settings.Playground.FilesLimit = 1000
	}
	if settings.Playground.SizeLimit == 0 {
		settings.Playground.SizeLimit = 5 << 20 // 5MB
	}
	return os.Chdir(settings.Playground.Path)
}

// DeinitPlayground removes playground files
func DeinitPlayground(vm *VM) {
	os.RemoveAll(vm.Settings.Playground.Path)
}

func PlaygroundAbsPath(vm *VM, fname string) (ret string, err error) {
	ret, err = filepath.Abs(fname)
	if err == nil {
		if !strings.HasPrefix(strings.ToLower(ret), strings.ToLower(vm.Settings.Playground.Path)) {
			return ``, fmt.Errorf(`%s [%s]`, ErrorText(ErrPlayAccess), fname)
		}
	}
	return
}

func CheckPlaygroundLimits(vm *VM, fname string, size int) error {
	var (
		curSize int
		ok      bool
	)
	ret, err := PlaygroundAbsPath(vm, fname)
	if err != nil {
		return err
	}
	name := ret[len(vm.Settings.Playground.Path):]
	if curSize, ok = vm.Playground.Files[name]; !ok {
		vm.Playground.Files[name] = size
	} else {
		vm.Playground.Files[name] += size
	}
	if len(vm.Playground.Files) > vm.Settings.Playground.FilesLimit {
		return fmt.Errorf(`%s [%d]`, ErrorText(ErrPlayCount), vm.Settings.Playground.FilesLimit)
	}
	vm.Playground.Size += size
	if curSize+size > vm.Settings.Playground.SizeLimit {
		return fmt.Errorf(`%s [%d MB]`, ErrorText(ErrPlaySize), vm.Settings.Playground.SizeLimit>>20)
	}
	if vm.Playground.Size > vm.Settings.Playground.AllSizeLimit {
		return fmt.Errorf(`%s [%d MB]`, ErrorText(ErrPlayAllSize), vm.Settings.Playground.AllSizeLimit>>20)
	}
	return nil
}
