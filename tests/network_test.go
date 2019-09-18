// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package test

import (
	"fmt"
	"path/filepath"
	"testing"

	ws "github.com/gentee/gentee"
)

func TestNetwork(t *testing.T) {
	workspace := ws.New()

	scriptName := filepath.Join(`scripts`, `network.g`)
	result, err := workspace.CompileAndRun(scriptName)
	if err != nil {
		t.Error(err)
		return
	}
	if fmt.Sprint(result) != `OK` {
		t.Errorf(`Wrong result %v`, result)
		return
	}
}
