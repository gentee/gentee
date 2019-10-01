// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gentee/gentee"
)

// To run: go run tests/manual/manual.go

func main() {
	workspace := gentee.New()

	result, err := workspace.CompileAndRun(`tests/scripts/readinput.g`)
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	result, err = workspace.CompileAndRun(`tests/scripts/network.g`)
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	if fmt.Sprint(result) != `OK` {
		fmt.Printf(`Wrong result %v`, result)
		return
	}
	fmt.Println(`Result:`, result)
	err = filepath.Walk(`examples`, func(script string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(script) == ".g" {

			fmt.Println(`FILE`, script)
			_, err := workspace.CompileAndRun(script)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
}
