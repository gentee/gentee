// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/gentee/gentee/core"
)

type testItem struct {
	want   string
	params []string
}

func TestCli(t *testing.T) {
	var (
		err    error
		gopath string
		stdout []byte
	)

	cmd := exec.Command(`go`, `env`)
	if stdout, err = cmd.CombinedOutput(); err != nil {
		t.Error(err)
		return
	}
	ret := regexp.MustCompile(`GOPATH="(.*)"`).FindStringSubmatch(string(stdout))
	if len(ret) == 2 {
		gopath = ret[1]
	}
	os.Setenv(`GOPATH`, gopath)
	outputFile := os.ExpandEnv(`${GOPATH}/bin/gentee`)
	cmd = exec.Command(`go`, `build`, `-o`, outputFile, `../cli/gentee.go`)
	if err = cmd.Run(); err != nil {
		t.Error(err)
		return
	}

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

	testList := []testItem{
		{``, []string{`-t`, `ok.g`}},
		{"ok 777\n", []string{`ok.g`}},
		{"test\nERROR: .../test/scripts/ok.g [3:1] script ok has already been linked",
			[]string{`runname.g`, `ok.g`}},
		{core.Version, []string{`-ver`}},
		{``, []string{`nothing.g`}},
		{core.Version, []string{`const.g`}},
		{"ERROR 254: .../test/scripts/traceerror.g [2:13] divided by zero\n" +
			".../test/scripts/traceerror.g [5:5] run -> myfunc\n" +
			".../test/scripts/traceerror.g [2:13] myfunc -> Div", []string{`traceerror.g`}},
		{"ERROR 300: .../test/scripts/customerror.g [3:24] Σ custom error №5\n" +
			".../test/scripts/customerror.g [9:12] run -> myerr\n" +
			".../test/scripts/customerror.g [3:24] myerr -> error", []string{`customerror.g`}},
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
