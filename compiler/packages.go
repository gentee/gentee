// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

// Code contains a source code
type Code struct {
	CRC int64
}

// Package describes a library
type Package struct {
	Name  string
	Codes []Code
}

var (
	packages map[string]*Package
)

func initPackages() {
	packages = make(map[string]*Package)
}
