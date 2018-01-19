// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"errors"
)

const (
	// The list of errors

	// ErrLetter returns when an unknown character has been found
	ErrLetter = iota + 1
	// ErrWord returns when a sequence of characters is wrong
	ErrWord
	// ErrDecl returns when the unexpexted token has been found on the top level
	ErrDecl
)

var (
	errText = map[int]string{
		ErrLetter: `unknown character`,
		ErrWord:   `wrong sequence of characters`,
		ErrDecl:   `expected declaration: func, run etc`,
	}
)

func compError(id int) error {
	return errors.New(errText[id])
}
