// Copyright 2018-21 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
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
	input, err = ioutil.ReadFile(filepath.Join(`tests`, filename))
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

func TestGentee(t *testing.T) {
	workspace := New()

	testFile := func(filename string) error {
		src, err := loadTest(filename)
		if err != nil {
			return err
		}
		for i := len(src) - 1; i >= 0; i-- {
			testErr := func(err error) error {
				return fmt.Errorf(`[%d] of %s  %v`, src[i].Line, filename, err)
			}
			exec, _, err := workspace.Compile(src[i].Src, ``)
			if err != nil && err.Error() != src[i].Want {
				return testErr(err)
			}
			if err != nil {
				continue
			}
			var settings Settings
			if filename == `err_test` {
				settings.Cycle = 1000000
			}
			//			fmt.Println(`i`, src[i].Line, filename)
			result, err := exec.Run(settings)
			//			result, err := workspace.Run(unitID)
			if err == nil {
				if err = getWant(result, src[i].Want); err != nil {
					return testErr(err)
				}
			} else if err.Error() != src[i].Want {
				return testErr(err)
			}
		}
		return nil
	}
	for _, name := range []string{`run_test`, `err_test`} {
		if err := testFile(name); err != nil {
			t.Error(err)
			return
		}
	}
	files, err := ioutil.ReadDir(filepath.Join("tests", "stdlib"))
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
	scriptName := filepath.Join(`tests`, filepath.Join(`scripts`, `const.g`))
	exec, unitID, err := workspace.CompileFile(scriptName)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := exec.Run(Settings{})
	//result, err := workspace.Run(unitID)
	if err != nil {
		t.Error(err)
		return
	}
	if result != Version() {
		t.Errorf(`Wrong version %v`, result)
		return
	}
	if !strings.HasSuffix(workspace.Unit(unitID).Name, scriptName) {
		t.Errorf(`Wrong unit name %v`, workspace.Unit(unitID).Name)
		return
	}
}
