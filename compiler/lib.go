// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"unicode"

	"bitbucket.org/novostrim/go-gentee/core"
)

func isBoolResult(cmd core.ICmd) bool {
	return cmd.GetResult().GetName() == `bool`
}

func isIntResult(cmd core.ICmd) bool {
	return cmd.GetResult().GetName() == `int`
}

func isCapital(name string) bool {
	for _, ch := range name {
		if !unicode.IsUpper(ch) && ch != '_' && !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}
