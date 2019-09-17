// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"path/filepath"
	"testing"

	ws "github.com/gentee/gentee"
	"github.com/gentee/gentee/vm"
)

func TestNetwork(t *testing.T) {
	workspace := ws.New()

	scriptName := filepath.Join(`scripts`, `network.g`)
	exec, _, err := workspace.CompileFile(scriptName)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := vm.Run(exec, vm.Settings{})
	if err != nil {
		t.Error(err)
		return
	}
	if fmt.Sprint(result) != `OK` {
		t.Errorf(`Wrong result %v`, result)
		return
	}
}
