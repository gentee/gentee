// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gentee/gentee"
)

// To run: go run tests/manual/manual.go

func main() {
	workspace := gentee.New()

	exec, _, err := workspace.CompileFile(`tests/manual/readinput.g`)
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	result, err := exec.Run(gentee.Settings{})
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	fmt.Println(`Result:`, result)
}
