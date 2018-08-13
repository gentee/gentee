// Copyright 2018 The Gentee Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"bitbucket.org/novostrim/go-gentee/workspace"
)

func getWant(v interface{}, want string) error {
	get := fmt.Sprint(v)
	want = strings.Replace(want, `\n`, "\n", -1)
	if get != want {
		return fmt.Errorf("get != want; %s != %s", get, want)
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
		for i, line := range list {
			if !strings.HasPrefix(line, `=====`) {
				source = append(source, line)
				continue
			}
			testErr := func(err error) error {
				return fmt.Errorf(`[%d] of %s  %v`, i, filename, err)
			}

			want := strings.TrimSpace(strings.TrimLeft(line, `=`))
			err = workspace.Compile(strings.Join(source, "\n"))
			source = source[:0]
			if err != nil && err.Error() != strings.TrimSpace(want) {
				return testErr(err)
			}
			if err != nil {
				continue
			}
			result, err := workspace.Run(``)
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
