// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	gentee "github.com/gentee/gentee"
	"github.com/gentee/gentee/core"
	"github.com/gentee/gentee/vm"
)

// Source contains source code and result value
type Source struct {
	Src  string
	Want string
	Line int
}

func loadTest(filename string) (src []Source, err error) {
	var input []byte
	src = make([]Source, 0, 64)
	input, err = ioutil.ReadFile(filename)
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

func noPars() {
}

func retStr() string {
	return `retStr`
}

func parStr(in string) string {
	return in + in
}

func sum(x, y int64) int64 {
	return x + 2*y
}

func Equal(x int64, y string) int64 {
	if fmt.Sprint(x) == y {
		return 1
	}
	return 0
}

func varInt(pars ...int64) int64 {
	var sum int64
	for _, i := range pars {
		sum += i
	}
	return sum
}

func varPar(mul string, pars ...interface{}) string {
	for _, i := range pars {
		mul += fmt.Sprint(i)
	}
	return mul
}

func custErr(s string) (string, error) {
	if len(s) == 1 {
		return s + s, nil
	}
	return ``, fmt.Errorf("string %s is too long", s)
}

func mustErr() error {
	return fmt.Errorf("custom error")
}

func rtStrStack(rt *vm.Runtime, a string, b ...string) (*core.Array, error) {
	ret := core.NewArray()
	for _, s := range rt.SStr[:4] {
		ret.Data = append(ret.Data, s)
	}
	return ret, nil
}

var customLib = []gentee.EmbedItem{
	{Prototype: `nopars()`, Object: noPars},
	{Prototype: `retStr() str`, Object: retStr},
	{Prototype: `ParStr(str) str`, Object: parStr},
	{Prototype: `Sum(int, int) int`, Object: sum},
	{Prototype: `Equal(int, str) bool`, Object: Equal},
	{Prototype: `varInt() int`, Object: varInt},
	{Prototype: `varPar(str) str`, Object: varPar},
	{Prototype: `custErr(str) str`, Object: custErr},
	{Prototype: `mustErr()`, Object: mustErr},
	{Prototype: `rtStrStack(str) arr.str`, Object: rtStrStack},
}

func TestCustom(t *testing.T) {
	err := gentee.Customize(&gentee.Custom{
		Embedded: []gentee.EmbedItem{{Prototype: `myfunc()`, Object: nil},
			{Prototype: ``, Object: nil}},
	})
	if err.Error() != `invalid custom declaration {myfunc() <nil>}` {
		t.Error(err)
		return
	}
	err = gentee.Customize(&gentee.Custom{
		Embedded: customLib,
	})
	if err != nil {
		t.Error(err)
		return
	}

	workspace := gentee.New()

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
			var settings gentee.Settings
			result, err := exec.Run(settings)
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
	if err := testFile(`custom_test`); err != nil {
		t.Error(err)
		return
	}
}
