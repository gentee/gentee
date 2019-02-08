// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gentee/gentee/workspace"
)

// To run: go run test/manual/manual.go

func main() {
	workspace := workspace.New()

	unitID, err := workspace.CompileFile(`test/manual/readinput.g`)
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	result, err := workspace.Run(unitID)
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	fmt.Println(`Result:`, result)
}
