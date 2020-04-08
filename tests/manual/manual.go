// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"bufio"
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

func stdInOut(ws *gentee.Gentee) error {
	var (
		settings   gentee.Settings
		rIn, wIn   *os.File
		rOut, wOut *os.File
		console    *os.File
	)

	chNum := make(chan string, 10)
	exec, _, err := ws.CompileFile(`tests/scripts/stdinout.g`)
	if err != nil {
		return err
	}
	console = os.Stdout
	rIn, wIn, _ = os.Pipe()
	settings.Stdin = rIn
	rOut, wOut, _ = os.Pipe()
	settings.Stdout = wOut
	var (
		got   string
		count int
	)
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := rOut.Read(buf)
			if err != nil {
				got += err.Error()
				break
			}
			console.Write(buf)
			if strings.HasPrefix(string(buf), `Enter `) {
				chNum <- string(buf[6 : n-1])
			} else {
				got += string(buf)
			}
		}
	}()
	go func() {
		var num, buf string
		for {
			num = <-chNum
			buf = fmt.Sprintf("%s\n", num)
			got += `=` + num
			count++
			if count > 5 {
				buf = "100\n"
			}
			wIn.Write([]byte(buf))
		}
	}()

	_, err = exec.Run(settings)

	fmt.Println("\n" + got) //strings.ReplaceAll(got, "\r", "\n"))
	return nil
}

func sysChan(ws *gentee.Gentee) error {
	var (
		settings gentee.Settings
	)

	finished := make(chan error)
	settings.SysChan = make(chan int)
	exec, _, err := ws.CompileFile(`tests/scripts/syschan.g`)
	if err != nil {
		return err
	}
	go func() {
		_, err = exec.Run(settings)
		finished <- err
	}()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(`Press one of the following number and Enter
1 - for suspend
2 - for resume
3 - for terminate
any other - for exit`)

	x := int64(1)
	for x == 1 || x == 2 {
		ret, _ := reader.ReadString('\n')
		x, _ = strconv.ParseInt(strings.TrimSpace(ret), 0, 64)
		if x == 1 || x == 2 || x == 3 {
			settings.SysChan <- int(x)
		}
	}
	result := <-finished
	fmt.Println(`res`, result)
	return nil
}

func main() {
	var (
		err    error
		result interface{}
	)
	workspace := gentee.New()

	if err = stdInOut(workspace); err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	if err = sysChan(workspace); err != nil {
		fmt.Println(`ERROR:`, err)
		return
	}
	result, err = workspace.CompileAndRun(`tests/scripts/readinput.g`)
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
