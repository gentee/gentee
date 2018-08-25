// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestCli(t *testing.T) {
	var err error

	outputFile := os.ExpandEnv(`${GOPATH}/bin/gentee`)

	cmd := exec.Command(`go`, `build`, `-o`, outputFile, `../cli/gentee.go`)
	if err = cmd.Run(); err != nil {
		t.Error(err)
		return
	}
	cmd = exec.Command(outputFile, `-t`, `scripts/ok.g`)
	stdout, err := cmd.CombinedOutput()
	out := string(stdout)
	if err != nil {
		t.Error(err, out)
		return
	}
	cmd = exec.Command(outputFile, `scripts/ok.g`)
	stdout, err = cmd.CombinedOutput()
	out = string(stdout)
	if err != nil {
		t.Error(err, out)
		return
	}
	if err = getWant(out, "ok 777\n"); err != nil {
		t.Error(err)
		return
	}
	cmd = exec.Command(outputFile, `scripts/runname.g`, `scripts/ok.g`)
	stdout, err = cmd.CombinedOutput()
	out = string(stdout)
	if out != "test\nERROR: 3:1: script ok has already been linked\n" {
		t.Error(strconv.Quote(out))
		return
	}
}
