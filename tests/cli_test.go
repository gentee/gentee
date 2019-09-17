// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/gentee/gentee/core"
)

type testItem struct {
	want   string
	params []string
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
	ret := regexp.MustCompile(`GOPATH="?([^"|\n|\r]*)`).FindStringSubmatch(string(stdout))
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
		out := strings.Replace(string(stdout), `\`, `/`, -1)
		if err != nil {
			return getWant(out, want)
		} else if err = getWant(out, want); err != nil {
			return err
		}
		return nil
	}

	testList := []testItem{
		{``, []string{`-t`, `h.g`}},
		{``, []string{`-t`, `ok.g`}},
		{``, []string{`-t`, `struct.g`}},
		{`falsefalseok`, []string{`cmdline.g`}},
		{`1my par º ok`, []string{`cmdline.g`, `my par º ok`}},
		{`my parfalse`, []string{`cmdline.g`, `-p="my par"`, `--flag`}},
		{`my par+second`, []string{`cmdline.g`, `-list`, `my par`, `second`}},
		{`my option+one`, []string{`cmdline.g`, `-o`, `my option`, `-`, `one`}},
		{`-two+'three four'+fivetruefalse10false`,
			[]string{`cmdline.g`, `-i:10`, `-`, `-two`, `'three four'`, `five`}},
		{`один+*.два+.три+Welcome+"-option"+"-"truetrue`,
			[]string{`cmdline.g`, `один`, `*.два`, `.три`, `Welcome`, `"-option"`, `"-"`}},
		{"ok 777\n", []string{`ok.g`}},
		{"test", []string{`runname.g`}},
		{core.Version, []string{`-ver`}},
		{``, []string{`nothing.g`}},
		{core.Version, []string{`const.g`}},
		{"ERROR #3: .../tests/scripts/traceerror.g [2:13] divided by zero\n" +
			".../tests/scripts/traceerror.g [5:5] run -> myfunc\n" +
			".../tests/scripts/traceerror.g [2:13] myfunc -> Div", []string{`traceerror.g`}},
		{"ERROR #300: .../tests/scripts/customerror.g [3:24] Σ custom error №5\n" +
			".../tests/scripts/customerror.g [9:12] run -> myerr\n" +
			".../tests/scripts/customerror.g [3:24] myerr -> error", []string{`customerror.g`}},
		{"ERROR: .../tests/scripts/err-a.g [6:5] duplicate of c_func has been found after include/import",
			[]string{`err-a.g`}},
		{``, []string{`-t`, `a.g`}},
		{``, []string{`-t`, `d.g`}},
		{"ERROR: .../tests/scripts/err-b.g [6:12] function c_func(int) has not been found",
			[]string{`err-b.g`}},
		{``, []string{`-t`, `f.g`}},
		{"ERROR: .../tests/scripts/err-c.g [7:12] function e_func(int, int) has not been found",
			[]string{`err-c.g`}},
		{"ERROR: .../tests/scripts/err-d.g [7:7] function Assign(int, et) has not been found",
			[]string{`err-d.g`}},
		{"ERROR: .../tests/scripts/err-e.g [6:12] unknown identifier EINT",
			[]string{`err-e.g`}},
		{"ERROR: .../tests/scripts/err-f.g [6:5] unknown identifier myf",
			[]string{`err-f.g`}},
		{"ERROR #3: .../tests/scripts/err_thread.g [2:14] divided by zero\n" +
			".../tests/scripts/err_thread.g [7:17] thread -> divZero\n" +
			".../tests/scripts/err_thread.g [2:14] divZero -> Div",
			[]string{`thread1.g`}},
		{"ERROR #1000: .../tests/scripts/err_thread.g [14:9] This is an error message\n" +
			".../tests/scripts/err_thread.g [14:9] thread -> error",
			[]string{`thread2.g`}},
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
