// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"io/ioutil"
	"os"
	//	"github.com/gentee/gentee/core"
)

func coInclude(cmpl *compiler) error {
	var (
		v   interface{}
		err error
	)
	lp := cmpl.getLex()
	token := getToken(lp, cmpl.pos)
	v = lp.Strings[lp.Tokens[cmpl.pos].Index]
	if token[0] == '"' {
		if v, err = unNewLine(v.(string)); err != nil {
			return cmpl.Error(ErrDoubleQuotes)
		}
	}
	includeFile := os.ExpandEnv(v.(string))
	fmt.Println(`Include`, includeFile)
	_, err = ioutil.ReadFile(includeFile)
	if err != nil {
		return cmpl.Error(ErrIncludeFile, includeFile)
	}

	return nil
}
