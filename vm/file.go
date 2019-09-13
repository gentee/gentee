// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"io/ioutil"
	"os"

	"github.com/gentee/gentee/core"
)

func fromFileInfo(fileInfo os.FileInfo, finfo *Struct) *Struct {
	finfo.Values[0] = fileInfo.Name()
	finfo.Values[1] = fileInfo.Size()
	finfo.Values[2] = fileInfo.Mode()
	fromTime(finfo.Values[3].(*Struct), fileInfo.ModTime())
	finfo.Values[4] = fileInfo.IsDir()
	return finfo
}

// FileInfoºStr returns the finfo describing the named file.
func FileInfoºStr(rt *Runtime, name string) (*Struct, error) {
	finfo := NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT])
	handle, err := os.Open(name)
	if err != nil {
		return finfo, err
	}
	defer handle.Close()
	fileInfo, err := handle.Stat()
	if err != nil {
		return finfo, err
	}
	return fromFileInfo(fileInfo, finfo), nil
}

// ReadDirºStr reads a directory
func ReadDirºStr(rt *Runtime, dirname string) (*core.Array, error) {
	ret := core.NewArray()
	fileList, err := ioutil.ReadDir(dirname)
	if err != nil {
		return ret, err
	}
	for _, fileInfo := range fileList {
		ret.Data = append(ret.Data, fromFileInfo(fileInfo,
			NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT])))
	}
	return ret, nil
}

// SetFileTimeºStrTime changes the modification time of the named file
func SetFileTimeºStrTime(name string, ftime *Struct) error {
	mtime := toTime(ftime)
	return os.Chtimes(name, mtime, mtime)
}
