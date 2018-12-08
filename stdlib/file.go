// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"io/ioutil"
	"os"

	"github.com/gentee/gentee/core"
)

// InitFile appends stdlib int functions to the virtual machine
func InitFile(vm *core.VirtualMachine) {
	for _, item := range []interface{}{
		ChDirºStr,        // ChDir( str )
		CreateDirºStr,    // CreateDir( str )
		GetCurDir,        // GetCurDir( ) str
		ReadFileºStr,     // ReadFile( str ) str
		ReadFileºStrBuf,  // ReadFile( str, buf ) buf
		RemoveºStr,       // Remove( str )
		RemoveDirºStr,    // RemoveDir( str )
		RenameºStrStr,    // Rename( str, str )
		TempDir,          // TempDir()
		TempDirºStrStr,   // TempDir(str, str)
		WriteFileºStrBuf, // WriteFile( str, buf )
		WriteFileºStrStr, // WriteFile( str, str )
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// ChDirºStr change the current directory
func ChDirºStr(dirname string) error {
	return os.Chdir(dirname)
}

// GetCurDir returns the current directory
func GetCurDir() (string, error) {
	return os.Getwd()
}

// CreateDirºStr creates the directory(s)
func CreateDirºStr(dirname string) error {
	return os.MkdirAll(dirname, os.ModePerm)
}

// ReadFileºStr reads a file
func ReadFileºStr(filename string) (string, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

// ReadFileºStrBuf reads a file to buffer
func ReadFileºStrBuf(filename string, buf *core.Buffer) (*core.Buffer, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return buf, err
	}
	buf.Data = out
	return buf, nil
}

// RenameºStrStr renames a file or a directory
func RenameºStrStr(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// RemoveºStr removes a file or an empty directory
func RemoveºStr(filename string) error {
	return os.Remove(filename)
}

// RemoveDirºStr removes a directory
func RemoveDirºStr(dirname string) error {
	return os.RemoveAll(dirname)
}

// TempDir returns the temporary directory
func TempDir() string {
	return os.TempDir()
}

// TempDirºStrStr creates a directory in the temporary directory
func TempDirºStrStr(dir, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}

// WriteFileºStrBuf write a buffer to a file
func WriteFileºStrBuf(filename string, buf *core.Buffer) error {
	err := ioutil.WriteFile(filename, buf.Data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// WriteFileºStrStr write a string to a file
func WriteFileºStrStr(filename, in string) error {
	err := ioutil.WriteFile(filename, []byte(in), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
