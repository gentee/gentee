// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

import (
	"errors"
)

var (
	// The list of errors

	// ErrLetter returns when an unknown character has been found
	ErrLetter = errors.New(`unknown character`)
	// ErrWord returns when a sequence of characters is wrong
	ErrWord = errors.New(`wrong sequence of characters`)
)
