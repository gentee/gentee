// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	OnlyDirs  = 0x08

	FileCreate   = 0x01
	FileTrunc    = 0x02
	FileReadonly = 0x04
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

func CloseFile(file *core.File) error {
	if file == nil || file.Handle == nil {
		return fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	if err := file.Handle.Close(); err != nil {
		return err
	}
	file.Handle = nil
	return nil
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
	defer func() {
		destFile.Close()
		os.Chtimes(dest, finfo.ModTime(), finfo.ModTime())
	}()
	var (
		prog   *Progress
		reader io.Reader
	)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog = NewProgress(rt, finfo.Size(), ProgressCopy)
		prog.Start(src, dest)
		reader = NewProgressReader(srcFile, prog)
	} else {
		reader = srcFile
	}
	ret, err := io.Copy(destFile, reader)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog.Complete()
	}
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

// FileInfoºFile sets returns the finfo of the file
func FileInfoºFile(rt *Runtime, file *core.File) (*Struct, error) {
	finfo := NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT])
	if file == nil || file.Handle == nil {
		return finfo, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	fileInfo, err := file.Handle.Stat()
	if err != nil {
		return finfo, err
	}
	return fromFileInfo(fileInfo, finfo), nil
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

// IsEmptyDir returns true if the specified folder is empty
func IsEmptyDir(rt *Runtime, path string) (ret int64, err error) {
	if rt.Owner.Settings.IsPlayground {
		if err = CheckPlaygroundLimits(rt.Owner, path, NoLimit); err != nil {
			return
		}
	}
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()
	_, err = f.Readdir(1)

	if err == io.EOF {
		return 1, nil
	}
	return
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

// ObjºFinfo converts finfo structure to object.
func ObjºFinfo(finfo *Struct) *core.Obj {
	obj := core.NewObj()
	val := core.NewMap()
	val.SetIndex(`name`, objºAny(finfo.Values[0]))
	val.SetIndex(`size`, objºAny(finfo.Values[1]))
	val.SetIndex(`mode`, objºAny(finfo.Values[2]))
	val.SetIndex(`time`, objºAny(StrºTime(finfo.Values[3].(*Struct))))
	val.SetIndex(`isdir`, objºAny(finfo.Values[4].(int64) != 0))
	val.SetIndex(`dir`, objºAny(finfo.Values[5]))
	obj.Data = val
	return obj
}

func OpenFileºStr(rt *Runtime, fname string, flags int64) (ret *core.File, err error) {
	if fname, err = filepath.Abs(fname); err != nil {
		return
	}
	if rt.Owner.Settings.IsPlayground {
		var trunc int64
		if (flags & FileTrunc) != 0 {
			trunc = ClearSize
		}
		if err = CheckPlaygroundLimits(rt.Owner, fname, trunc); err != nil {
			return
		}
	}
	ret = core.NewFile()
	var (
		iFlags int
		handle *os.File
	)
	if (flags & FileCreate) != 0 {
		iFlags |= os.O_CREATE
	}
	if (flags & FileReadonly) != 0 {
		iFlags |= os.O_RDONLY
	} else {
		iFlags |= os.O_RDWR
	}
	if handle, err = os.OpenFile(fname, iFlags, 0644); err != nil {
		return
	}
	if (flags & FileTrunc) != 0 {
		if err = handle.Truncate(0); err != nil {
			return
		}
	}
	ret.Name = fname
	ret.Handle = handle
	return
}

func ReadºFileInt(file *core.File, size int64) (*core.Buffer, error) {
	if file == nil || file.Handle == nil {
		return nil, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	buf := core.NewBuffer()
	buf.Data = make([]byte, size)
	n, err := file.Handle.Read(buf.Data)
	if err != nil && err == io.EOF {
		err = nil
	}
	buf.Data = buf.Data[:n]
	return buf, err
}

// ReadDirºStr reads a directory
func ReadDirºStr(rt *Runtime, dirname string) (*core.Array, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, dirname, NoLimit); err != nil {
			return nil, err
		}
	}
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

func readDir(rt *Runtime, ret *core.Array, dirname string, flags int64, patterns *core.Array,
	ignore *core.Array) error {
	fileList, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	isMatch := func(filename, pattern string) (ok int64, err error) {
		if len(pattern) == 0 {
			return 1, nil
		}
		isRegex := flags&RegExp != 0
		if !isRegex {
			if len(pattern) > 2 && pattern[0] == '/' && pattern[len(pattern)-1] == '/' {
				isRegex = true
				pattern = pattern[1 : len(pattern)-1]
			}
		}
		if isRegex {
			if ok, err = MatchºStrStr(filename, pattern); err != nil {
				return
			}
		} else if ok, err = MatchPath(pattern, filename); err != nil {
			return
		}
		return
	}
main:
	for _, fileInfo := range fileList {
		var ok int64
		for _, item := range ignore.Data {
			if pattern := item.(string); len(pattern) > 0 {
				if ok, err = isMatch(fileInfo.Name(), pattern); err != nil {
					return err
				} else if ok != 0 {
					continue main
				}
			}
		}
		if fileInfo.IsDir() {
			if flags&Recursive != 0 {
				err = readDir(rt, ret, filepath.Join(dirname, fileInfo.Name()), flags, patterns, ignore)
				if err != nil {
					return err
				}
			}
			if flags&OnlyFiles != 0 {
				continue
			}
		} else if flags&OnlyDirs != 0 {
			continue
		}
		if len(patterns.Data) > 0 {
			for _, item := range patterns.Data {
				if ok, err = isMatch(fileInfo.Name(), item.(string)); err != nil {
					return err
				}
				if ok != 0 {
					break
				}
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
	patterns := core.NewArray()
	if len(pattern) > 0 {
		patterns.Data = append(patterns.Data, pattern)
	}
	return ReadDirºStrArr(rt, dirname, flags, patterns, core.NewArray())
}

// ReadDirºStrArr reads a directory with additional settings
func ReadDirºStrArr(rt *Runtime, dirname string, flags int64, patterns *core.Array,
	ignore *core.Array) (*core.Array, error) {
	var err error
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, dirname, NoLimit); err != nil {
			return nil, err
		}
	}
	if flags&OnlyFiles != 0 && flags&OnlyDirs != 0 {
		flags &^= OnlyFiles
		flags &^= OnlyDirs
	}
	ret := core.NewArray()
	dirname, err = filepath.Abs(dirname)
	if err != nil {
		return ret, err
	}
	err = readDir(rt, ret, dirname, flags, patterns, ignore)
	return ret, err
}

// ReadFileºStr reads a file
func ReadFileºStr(rt *Runtime, filename string) (string, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return ``, err
		}
	}
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

// ReadFileºStrBuf reads a file to buffer
func ReadFileºStrBuf(rt *Runtime, filename string, buf *core.Buffer) (*core.Buffer, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return nil, err
		}
	}
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return buf, err
	}
	buf.Data = out
	return buf, nil
}

// ReadFileºStrIntInt reads a part of the file to the buffer
func ReadFileºStrIntInt(rt *Runtime, filename string, off int64, length int64) (buf *core.Buffer, err error) {
	var (
		fhandle *os.File
		n       int
	)
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return nil, err
		}
	}
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
func RemoveºStr(rt *Runtime, filename string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, DeleteSize); err != nil {
			return err
		}
	}
	return os.Remove(filename)
}

// RemoveDirºStr removes a directory
func RemoveDirºStr(rt *Runtime, dirname string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, dirname, DeleteAllSize); err != nil {
			return err
		}
	}
	return os.RemoveAll(dirname)
}

// RenameºStrStr renames a file or a directory
func RenameºStrStr(rt *Runtime, oldname, newname string) error {
	if rt.Owner.Settings.IsPlayground {
		size, err := PlaygroundSize(rt.Owner, oldname)
		if err != nil {
			return err
		}
		if err = CheckPlaygroundLimits(rt.Owner, oldname, DeleteSize); err != nil {
			return err
		}
		if err = CheckPlaygroundLimits(rt.Owner, newname, size); err != nil {
			return err
		}
	}
	return os.Rename(oldname, newname)
}

// SetFileTimeºStrTime changes the modification time of the named file
func SetFileTimeºStrTime(rt *Runtime, name string, ftime *Struct) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, name, NoLimit); err != nil {
			return err
		}
	}
	mtime := toTime(ftime)
	return os.Chtimes(name, mtime, mtime)
}

// SetPosºFileIntInt sets the postion in the file
func SetPosºFileIntInt(file *core.File, off int64, whence int64) (int64, error) {
	if file == nil || file.Handle == nil || whence < 0 || whence > 2 {
		return 0, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	return file.Handle.Seek(off, int(whence))
}

// Sha256FileºStr returns sha256 hash of the file as a hex string
func Sha256FileºStr(rt *Runtime, filename string) (string, error) {
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
func TempDirºStrStr(rt *Runtime, dir, prefix string) (string, error) {
	if rt.Owner.Settings.IsPlayground {
		tmp := dir
		if len(tmp) == 0 {
			tmp = TempDir()
		}
		if err := CheckPlaygroundLimits(rt.Owner, filepath.Join(tmp, prefix+`_`), NoLimit); err != nil {
			return ``, err
		}
	}
	return ioutil.TempDir(dir, prefix)
}

// WriteFileºFileBuf writes a buffer to a file
func WriteFileºFileBuf(rt *Runtime, file *core.File, buf *core.Buffer) (*core.File, error) {
	if file == nil || file.Handle == nil {
		return file, fmt.Errorf(ErrorText(ErrInvalidParam))
	}
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, file.Name, int64(len(buf.Data))); err != nil {
			return file, err
		}
	}
	_, err := file.Handle.Write(buf.Data)
	return file, err
}

// WriteFileºStrBuf writes a buffer to a file
func WriteFileºStrBuf(rt *Runtime, filename string, buf *core.Buffer) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, int64(len(buf.Data))); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, buf.Data, os.ModePerm)
}

// WriteFileºStrStr writes a string to a file
func WriteFileºStrStr(rt *Runtime, filename, in string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, int64(len(in))); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, []byte(in), os.ModePerm)
}
