// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gentee/gentee/core"
)

// InitPath appends stdlib filepath functions to the virtual machine
func InitPath(ws *core.Workspace) {
	for _, item := range []interface{}{
		AbsPath,   // AbsPath(str) str
		BaseName,  // BaseName(str) str
		Dir,       // Dir(str) str
		Ext,       // Ext(str) str
		JoinPath,  // JoinPath(str pars...) str
		MatchPath, // MatchPath(str, str) bool
	} {
		ws.StdLib().NewEmbed(item)
	}
}

// AbsPath returns an absolute representation of path.
func AbsPath(fname string) (string, error) {
	return filepath.Abs(fname)
}

// BaseName returns the last element of path.
func BaseName(fname string) string {
	return filepath.Base(fname)
}

// Dir returns all but the last element of path.
func Dir(fname string) string {
	return filepath.Dir(fname)
}

// Ext returns the file name extension used by path.
func Ext(fname string) string {
	return strings.TrimLeft(filepath.Ext(fname), `.`)
}

// JoinPath joins any number of path elements into a single path.
func JoinPath(pars ...interface{}) string {
	names := make([]string, len(pars))
	for i, name := range pars {
		names[i] = fmt.Sprint(name)
	}
	return filepath.Join(names...)
}

// MatchPath reports whether name matches the specified file name pattern.
func MatchPath(pattern, fname string) (bool, error) {
	return filepath.Match(pattern, fname)
}
