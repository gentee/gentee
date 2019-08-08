// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	ws "github.com/gentee/gentee"
	"github.com/gentee/gentee/compiler"
)

// Source contains source code and result value
type Source struct {
	Src  string
	Want string
	Line int
}

func getWant(v interface{}, want string) error {
	get := fmt.Sprint(v)
	if runtime.GOOS == `windows` {
		get = strings.Replace(get, "\r", ``, -1)
		get = strings.Replace(get, `\"`, `"`, -1)
	}
	want = strings.Replace(want, `\n`, "\n", -1)
	if get != want {
		return fmt.Errorf("get != want;\n%s !=\n%s", get, want)
	}
	return nil
}

func loadTest(filename string) (src []Source, err error) {
	var input []byte
	src = make([]Source, 0, 64)
	input, err = ioutil.ReadFile(filepath.Join(`../tests`, filename))
	if err != nil {
		return
	}
	list := strings.Split(string(input), "\n")
	source := make([]string, 0, 32)
	on := true
	for i, line := range list {
		if on && strings.HasPrefix(line, `OFF`) {
			on = false
			continue
		}
		if !on {
			if strings.HasPrefix(line, `ON`) {
				on = true
			}
			continue
		}

		if !strings.HasPrefix(line, `=====`) {
			source = append(source, line)
			continue
		}
		src = append(src, Source{
			Src:  strings.Join(source, "\n"),
			Want: strings.TrimSpace(strings.TrimLeft(line, `=`)),
			Line: i,
		})
		source = source[:0]
	}
	return
}

func TestVM(t *testing.T) {
	workspace := ws.New()

	testFile := func(filename string) error {
		src, err := loadTest(filename)
		if err != nil {
			return err
		}
		for i := len(src) - 1; i >= 0; i-- {
			testErr := func(err error) error {
				return fmt.Errorf(`[%d] of %s  %v`, src[i].Line, filename, err)
			}
			unitID, err := workspace.Compile(src[i].Src, ``)
			if err != nil && err.Error() != src[i].Want {
				return testErr(err)
			}
			if err != nil {
				continue
			}
			linked, err := compiler.Link(workspace.Workspace, unitID)
			if err != nil {
				return testErr(err)
			}
			fmt.Println(`LINE START`, src[i].Line)
			result, err := Run(linked, Settings{})
			if err == nil {
				if err = getWant(result, src[i].Want); err != nil {
					fmt.Println(`LINE`, src[i].Line)
					return testErr(err)
				}
			} else if err.Error() != src[i].Want {
				return testErr(err)
			}
		}
		return nil
	}
	/*	var top1, top2, top3, top4 int64
		for i := 0; i < 4000000000; i++ {
			top1++
			top2++
			top3++
			top4++
			top1 = top2 + top3
		}*/
	/*var top State
	for i := 0; i < 4000000000; i++ {
		top.topInt++
		top.topFloat++
		top.topStr++
		top.topAny++
		top.topInt = top.topFloat + top.topStr
	}*/
	//		var top [4]uint16
	//	top := make([]interface{}, 4)
	/*	top := make([]int64, 4)
		top[3] = int64(0)
		for i := 0; i < 400000000; i++ {
			top[0] = top[3]              //++
			top[1] = top[0]              //++
			top[2] = top[1]              //++
			top[3] = top[2]              //++
			top[0] = int64(len(top)) - 1 //top[2]
		}*/
	/*	top := make([]int, 0, 8)
		for i := 0; i < 4000000000; i++ {
			top = append(top, 1)
			top = append(top, 2)
			top = append(top, 3)
			top = append(top, 4)
			//		top[0] = len(top) - 1 //top[2]
			top = top[:0]
		}*/
	//	return
	for _, name := range []string{`run_test`, `err_test`} {
		if err := testFile(name); err != nil {
			t.Error(err)
			return
		}
	}
	files, err := ioutil.ReadDir(filepath.Join("../tests", "stdlib"))
	if err != nil {
		t.Error(err)
		return
	}
	if len(files) < 8 {
		t.Error(`stdlib tests cannot be found`)
		return
	}
	for _, file := range files {
		if err := testFile(filepath.Join(`stdlib`, file.Name())); err != nil {
			t.Error(err)
			return
		}
	}
	if runtime.GOOS == `linux` {
		for _, name := range []string{`linux_test`} {
			if err := testFile(name); err != nil {
				t.Error(err)
				return
			}
		}
	}
	scriptName := filepath.Join(`../tests`, filepath.Join(`scripts`, `const.g`))
	unitID, err := workspace.CompileFile(scriptName)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := workspace.Run(unitID)
	if err != nil {
		t.Error(err)
		return
	}
	if result != workspace.Version() {
		t.Errorf(`Wrong version %v`, result)
		return
	}
	if !strings.HasSuffix(workspace.Unit(unitID).Name, scriptName) {
		t.Errorf(`Wrong unit name %v`, workspace.Unit(unitID).Name)
		return
	}
}
