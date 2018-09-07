// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gentee/gentee/workspace"
)

func getWant(v interface{}, want string) error {
	get := fmt.Sprint(v)
	want = strings.Replace(want, `\n`, "\n", -1)
	if get != want {
		return fmt.Errorf("get != want;\n%s !=\n%s", get, want)
	}
	return nil
}

func TestRun(t *testing.T) {

	workspace := workspace.New()

	testFile := func(filename string) error {
		input, err := ioutil.ReadFile(filename)
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
			unit, err := workspace.Compile(strings.Join(source, "\n"), ``)
			source = source[:0]
			if err != nil && err.Error() != strings.TrimSpace(want) {
				return testErr(err)
			}
			if err != nil {
				continue
			}
			result, err := workspace.Run(unit.Name)
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
	t.Errorf(`OK`)
}
