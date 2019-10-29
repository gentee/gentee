// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	fmt.Println(`Result:`, result)
	result, err = workspace.CompileAndRun(`tests/scripts/network.g`)
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	if fmt.Sprint(result) != `OK` {
		fmt.Printf(`Wrong result %v`, result)
		return
	}
	err = filepath.Walk(`examples`, func(script string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(script) == ".g" {
			fmt.Println(`FILE`, script)
			exec, unitID, err := workspace.CompileFile(script)
			if err != nil {
				return err
			}
			unit := workspace.Unit(unitID)
			stdout := unit.GetHeader(`stdout`)
			resWant := unit.GetHeader(`result`)
			stdin := unit.GetHeader(`stdin`)
			cycle := unit.GetHeader(`settings.cycle`)
			var (
				rescueStdout, r, w *os.File
			)
			if stdout == `1` {
				rescueStdout = os.Stdout
				r, w, _ = os.Pipe()
				os.Stdout = w
			}
			var settings gentee.Settings
			if len(stdin) > 0 {
				settings.Input = []byte(strings.ReplaceAll(stdin, `\n`, "\n"))
			}
			if len(cycle) > 0 {
				if i, err := strconv.ParseUint(cycle, 10, 64); err == nil {
					settings.Cycle = i
				}
			}
			result, err := exec.Run(settings)
			if stdout == `1` {
				w.Close()
				out, _ := ioutil.ReadAll(r)
				os.Stdout = rescueStdout
				result = string(out)
				if strings.HasPrefix(resWant, `CRC`) {
					result = fmt.Sprintf(`CRC0x%x`, crc64.Checksum([]byte(result.(string)),
						crc64.MakeTable(crc64.ECMA)))
				}
			}
			if err != nil {
				if err.Error() != resWant {
					return fmt.Errorf(`error result %v`, err)
				}
			} else if len(resWant) > 0 && resWant != strings.TrimSpace(fmt.Sprint(result)) {
				return fmt.Errorf(`wrong result %v`, result)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
}
