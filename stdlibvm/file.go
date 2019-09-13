// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/gentee/gentee/core"
)

func appendFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AppendFileºStrBuf appends a buffer to a file
func AppendFileºStrBuf(filename string, buf *core.Buffer) error {
	return appendFile(filename, buf.Data)
}

// AppendFileºStrStr appends a string to a file
func AppendFileºStrStr(filename, s string) error {
	return appendFile(filename, []byte(s))
}

// ChDirºStr change the current directory
func ChDirºStr(dirname string) error {
	return os.Chdir(dirname)
}

// CopyFileºStrStr copies a file
func CopyFileºStrStr(src, dest string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()
	return io.Copy(destFile, srcFile)
}

// CreateDirºStr creates the directory(s)
func CreateDirºStr(dirname string) error {
	return os.MkdirAll(dirname, os.ModePerm)
}

// GetCurDir returns the current directory
func GetCurDir() (string, error) {
	return os.Getwd()
}

// Md5FileºStr returns md5 hash of the file as a hex string
func Md5FileºStr(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ``, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ``, err
	}
	return hex.EncodeToString(hash.Sum(nil)[:]), nil
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

// ReadFileºStrIntInt reads a part of the file to the buffer
func ReadFileºStrIntInt(filename string, off int64, length int64) (buf *core.Buffer, err error) {
	var (
		fhandle *os.File
		n       int
	)
	buf = core.NewBuffer()
	if fhandle, err = os.Open(filename); err != nil {
		return
	}
	defer fhandle.Close()
	fi, err := fhandle.Stat()
	fsize := fi.Size()
	if off < 0 {
		off = fsize + off
	}
	if off < 0 {
		off = 0
	} else if off > fsize-1 {
		return
	}
	if off+length > fsize {
		length = fsize - off
	}
	buf.Data = make([]byte, length)
	n, err = fhandle.ReadAt(buf.Data, off)
	if err != nil && err == io.EOF {
		err = nil
	}
	buf.Data = buf.Data[:n]
	return
}

// RemoveºStr removes a file or an empty directory
func RemoveºStr(filename string) error {
	return os.Remove(filename)
}

// RemoveDirºStr removes a directory
func RemoveDirºStr(dirname string) error {
	return os.RemoveAll(dirname)
}

// RenameºStrStr renames a file or a directory
func RenameºStrStr(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// Sha256FileºStr returns sha256 hash of the file as a hex string
func Sha256FileºStr(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ``, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ``, err
	}
	return hex.EncodeToString(hash.Sum(nil)[:]), nil
}

// TempDir returns the temporary directory
func TempDir() string {
	return os.TempDir()
}

// TempDirºStrStr creates a directory in the temporary directory
func TempDirºStrStr(dir, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}

// WriteFileºStrBuf writes a buffer to a file
func WriteFileºStrBuf(filename string, buf *core.Buffer) error {
	return ioutil.WriteFile(filename, buf.Data, os.ModePerm)
}

// WriteFileºStrStr writes a string to a file
func WriteFileºStrStr(filename, in string) error {
	return ioutil.WriteFile(filename, []byte(in), os.ModePerm)
}
