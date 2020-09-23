// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"testing"

	gentee "github.com/gentee/gentee"
)

func TestPlayground(t *testing.T) {
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
			settings.Cycle = 1000
			settings.IsPlayground = true
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
	if err := testFile(`playground_test`); err != nil {
		t.Error(err)
		return
	}
}
