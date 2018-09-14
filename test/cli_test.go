// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gentee/gentee/core"
)

type testItem struct {
	want   string
	params []string
}

func TestCli(t *testing.T) {
	var err error

	outputFile := os.ExpandEnv(`${GOPATH}/bin/gentee`)

	call := func(want string, params ...string) error {
		cmd := exec.Command(outputFile, params...)
		stdout, err := cmd.CombinedOutput()
		out := string(stdout)
		if err != nil {
			return getWant(out, want)
		} else if err = getWant(out, want); err != nil {
			return err
		}
		return nil
	}

	cmd := exec.Command(`go`, `build`, `-o`, outputFile, `../cli/gentee.go`)
	if err = cmd.Run(); err != nil {
		t.Error(err)
		return
	}
	testList := []testItem{
		{``, []string{`-t`, `ok.g`}},
		{"ok 777\n", []string{`ok.g`}},
		{"test\nERROR: 3:1: script ok has already been linked", []string{`runname.g`, `ok.g`}},
		{core.Version, []string{`-ver`}},
		{``, []string{`nothing.g`}},
		{core.Version, []string{`const.g`}},
	}
	for _, item := range testList {
		for i, v := range item.params {
			if strings.HasSuffix(v, `.g`) {
				item.params[i] = `scripts/` + v
			}
		}
		if len(item.want) > 0 {
			item.want += "\n"
		}
		if err = call(item.want, item.params...); err != nil {
			t.Error(err)
			return
		}
	}
}
