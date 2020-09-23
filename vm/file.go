// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gentee/gentee/core"
)

const (
	Recursive = 0x01
	OnlyFiles = 0x02
	RegExp    = 0x04
)

func appendFile(rt *Runtime, filename string, data []byte) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, int64(len(data))); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// AppendFileºStrBuf appends a buffer to a file
func AppendFileºStrBuf(rt *Runtime, filename string, buf *core.Buffer) error {
	return appendFile(rt, filename, buf.Data)
}

// AppendFileºStrStr appends a string to a file
func AppendFileºStrStr(rt *Runtime, filename, s string) error {
	return appendFile(rt, filename, []byte(s))
}

// ChDirºStr change the current directory
func ChDirºStr(rt *Runtime, dirname string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, dirname, -1); err != nil {
			return err
		}
	}
	return os.Chdir(dirname)
}

// ChModeºStr change the file mode.
func ChModeºStr(rt *Runtime, name string, mode int64) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, name, -1); err != nil {
			return err
		}
	}
	return os.Chmod(name, os.FileMode(mode))
}

// CopyFileºStrStr copies a file
func CopyFileºStrStr(rt *Runtime, src, dest string) (int64, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, src, -1); err != nil {
			return 0, err
		}
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	finfo, err := srcFile.Stat()
	defer srcFile.Close()
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, dest, finfo.Size()); err != nil {
			return 0, err
		}
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()
	ret, err := io.Copy(destFile, srcFile)
	//	if finfo.Size() != ret {
	destFile.Chmod(finfo.Mode())
	return ret, err
}

// CreateDirºStr creates the directory(s)
func CreateDirºStr(rt *Runtime, dirname string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, dirname, 0); err != nil {
			return err
		}
	}
	return os.MkdirAll(dirname, os.ModePerm)
}

// CreateFileºStrBool creates an empty file
func CreateFileºStrBool(rt *Runtime, filename string, always int64) error {
	if rt.Owner.Settings.IsPlayground {
		var trunc int64
		if always != 0 {
			trunc = ClearSize
		}
		if err := CheckPlaygroundLimits(rt.Owner, filename, trunc); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if always != 0 {
		if err = f.Truncate(0); err != nil {
			return err
		}
	}
	f.Close()
	return nil
}

// ExistFile returns true if the file or directory exists
func ExistFile(rt *Runtime, filename string) (int64, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return 0, err
		}
	}
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	return 1, nil
}

func fromFileInfo(fileInfo os.FileInfo, finfo *Struct) *Struct {
	finfo.Values[0] = fileInfo.Name()
	finfo.Values[1] = fileInfo.Size()
	finfo.Values[2] = int64(fileInfo.Mode())
	fromTime(finfo.Values[3].(*Struct), fileInfo.ModTime())
	if fileInfo.IsDir() {
		finfo.Values[4] = int64(1)
	} else {
		finfo.Values[4] = int64(0)
	}
	finfo.Values[5] = ``
	return finfo
}

// FileInfoºStr returns the finfo describing the named file.
func FileInfoºStr(rt *Runtime, name string) (*Struct, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, name, NoLimit); err != nil {
			return nil, err
		}
	}
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

// FileModeºStr returns the file mode.
func FileModeºStr(rt *Runtime, name string) (int64, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, name, NoLimit); err != nil {
			return 0, err
		}
	}
	fStat, err := os.Stat(name)
	if err != nil {
		return 0, err
	}
	return int64(fStat.Mode()), nil
}

// GetCurDir returns the current directory
func GetCurDir() (string, error) {
	return os.Getwd()
}

// Md5FileºStr returns md5 hash of the file as a hex string
func Md5FileºStr(rt *Runtime, filename string) (string, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return ``, err
		}
	}
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

func readDir(rt *Runtime, ret *core.Array, dirname string, flags int64, pattern string) error {
	fileList, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, fileInfo := range fileList {
		if fileInfo.IsDir() {
			if flags&Recursive != 0 {
				err = readDir(rt, ret, filepath.Join(dirname, fileInfo.Name()), flags, pattern)
				if err != nil {
					return err
				}
			}
			if flags&OnlyFiles != 0 {
				continue
			}
		}
		if len(pattern) > 0 {
			var ok int64
			if flags&RegExp != 0 {
				if ok, err = MatchºStrStr(fileInfo.Name(), pattern); err != nil {
					return err
				}
			} else if ok, err = MatchPath(pattern, fileInfo.Name()); err != nil {
				return err
			}
			if ok == 0 {
				continue
			}
		}
		finfo := fromFileInfo(fileInfo, NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT]))
		finfo.Values[5] = dirname
		ret.Data = append(ret.Data, finfo)
	}
	return nil
}

// ReadDirºStrIntStr reads a directory with additional settings
func ReadDirºStrIntStr(rt *Runtime, dirname string, flags int64, pattern string) (*core.Array, error) {
	var err error
	ret := core.NewArray()
	dirname, err = filepath.Abs(dirname)
	if err != nil {
		return ret, err
	}
	err = readDir(rt, ret, dirname, flags, pattern)
	return ret, err
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

// SetFileTimeºStrTime changes the modification time of the named file
func SetFileTimeºStrTime(name string, ftime *Struct) error {
	mtime := toTime(ftime)
	return os.Chtimes(name, mtime, mtime)
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
