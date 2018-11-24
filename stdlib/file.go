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
		ChdirºStr,        // Chdir( str )
		ReadFileºStr,     // ReadFile( str ) str
		ReadFileºStrBuf,  // ReadFile( str, buf ) buf
		RemoveºStr,       // Remove( str )
		WriteFileºStrBuf, // WriteFile( str, buf )
		WriteFileºStrStr, // WriteFile( str, str )
	} {
		vm.StdLib().NewEmbed(item)
	}
}

// ChdirºStr change the current directory
func ChdirºStr(dirname string) error {
	return os.Chdir(dirname)
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

// RemoveºStr removes a file or an empty directory
func RemoveºStr(filename string) error {
	return os.Remove(filename)
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
