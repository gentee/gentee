// Copyright 2018 Alexey Krivonogov. All rights reserved.
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

func TestGentee(t *testing.T) {

	workspace := New()

	testFile := func(filename string) error {
		input, err := ioutil.ReadFile(filepath.Join(`tests`, filename))
		if err != nil {
			return err
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
			testErr := func(err error) error {
				return fmt.Errorf(`[%d] of %s  %v`, i, filename, err)
			}

			want := strings.TrimSpace(strings.TrimLeft(line, `=`))
			unitID, err := workspace.Compile(strings.Join(source, "\n"), ``)
			source = source[:0]
			if err != nil && err.Error() != strings.TrimSpace(want) {
				return testErr(err)
			}
			if err != nil {
				continue
			}
			result, err := workspace.Run(unitID)
			if err == nil {
				if err = getWant(result, want); err != nil {
					return testErr(err)
				}
			} else if err.Error() != strings.TrimSpace(want) {
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
}
