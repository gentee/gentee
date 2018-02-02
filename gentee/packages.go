// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package gentee

// Source contains a source code
type Source struct {
	CRC int64
}

// Package describes a library
type Package struct {
	Name    string
	Sources []Source
}

var (
	packages map[string]*Package
)

func initPackages() {
	packages = make(map[string]*Package)
}
