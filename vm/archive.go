// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentee/gentee/core"
)

type UnzipFile struct {
	Name    string
	Archive *zip.ReadCloser
}

type ZipFile struct {
	Name    string
	File    *os.File
	Archive *zip.Writer
}

// AddFileToZip adds a file to the open zip archive
func AddFileToZip(rt *Runtime, zfile *ZipFile, filename string, zipname string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, -1); err != nil {
			return err
		}
	}
	finfo, err := os.Stat(filename)
	if err != nil {
		return nil
	}
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, zfile.Name, finfo.Size()); err != nil {
			return err
		}
	}
	header, err := zip.FileInfoHeader(finfo)
	if err != nil {
		return err
	}
	if len(zipname) == 0 {
		zipname = filepath.Base(filename)
	}
	header.Name = zipname
	if finfo.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}
	writer, err := zfile.Archive.CreateHeader(header)
	if err != nil {
		return err
	}
	if finfo.IsDir() {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	var (
		prog   *Progress
		reader io.Reader
	)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog = NewProgress(rt, finfo.Size(), ProgressCompress)
		prog.Start(filename, zfile.Name)
		reader = NewProgressReader(file, prog)
	} else {
		reader = file
	}
	_, err = io.Copy(writer, reader)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog.Complete()
	}
	return err
}

// CloseUnzip closes the opened zip file
func CloseUnzip(zfile *UnzipFile) (err error) {
	return zfile.Archive.Close()
}

// CloseZip closes the created zip file
func CloseZip(zfile *ZipFile) (err error) {
	if err = zfile.Archive.Close(); err == nil {
		err = zfile.File.Close()
	}
	return
}

// CreateZip creates zip file and returns its handle
func CreateZip(rt *Runtime, filename string) (*ZipFile, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, -1); err != nil {
			return nil, err
		}
	}
	zipfile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	archive := zip.NewWriter(zipfile)
	return &ZipFile{Name: filename, File: zipfile, Archive: archive}, nil
}

// OpenUnzip opens zip file and returns its handle
func OpenUnzip(rt *Runtime, zipfile string) (*UnzipFile, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, zipfile, NoLimit); err != nil {
			return nil, err
		}
	}
	archive, err := zip.OpenReader(zipfile)
	if err != nil {
		return nil, err
	}
	return &UnzipFile{Name: zipfile, Archive: archive}, nil
}

func fromZipInfo(fileInfo *zip.File, finfo *Struct) *Struct {
	finfo.Values[0] = fileInfo.Name
	finfo.Values[1] = int64(fileInfo.UncompressedSize64)
	finfo.Values[2] = int64(0)
	fromTime(finfo.Values[3].(*Struct), fileInfo.Modified)
	if fileInfo.FileInfo().IsDir() {
		finfo.Values[4] = int64(1)
	} else {
		finfo.Values[4] = int64(0)
	}
	finfo.Values[5] = ``
	return finfo
}

// ReadUnzip returns the file list of zip
func ReadUnzip(rt *Runtime, zfile *UnzipFile) (*core.Array, error) {
	ret := core.NewArray()
	for _, finfo := range zfile.Archive.File {
		ret.Data = append(ret.Data, fromZipInfo(finfo,
			NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT])))
	}
	return ret, nil
}

func unzipByIndex(rt *Runtime, zfile *UnzipFile, index int, dir string) error {
	fhead := zfile.Archive.File[index]
	dest := filepath.Join(dir, filepath.Base(strings.TrimRight(fhead.Name, `/`)))
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, dest, int64(fhead.UncompressedSize64)); err != nil {
			return err
		}
	}
	if fhead.FileInfo().IsDir() {
		os.MkdirAll(dest, fhead.Mode())
		return nil
	}
	rfile, err := fhead.Open()
	if err != nil {
		return err
	}
	defer rfile.Close()

	target, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fhead.Mode())
	if err != nil {
		return err
	}
	defer target.Close()
	var (
		prog   *Progress
		reader io.Reader
	)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog = NewProgress(rt, int64(fhead.CompressedSize64), ProgressDecompress)
		prog.Start(zfile.Name, dest)
		reader = NewProgressReader(rfile, prog)
	} else {
		reader = rfile
	}
	_, err = io.Copy(target, reader)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog.Complete()
	}
	return err
}

// UnzipºStr unzip a zip file to the specified folder
func UnzipºStr(rt *Runtime, zipfile string, dir string) error {
	zfile, err := OpenUnzip(rt, zipfile)
	if err != nil {
		return err
	}
	var (
		prevDir string
		prog    *Progress
	)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog = NewProgress(rt, int64(len(zfile.Archive.File)), ProgressArchiveCounter)
		prog.Start(zipfile, ``)
	}
	for i, fhead := range zfile.Archive.File {
		folder := filepath.Dir(strings.TrimRight(fhead.Name, `/`))
		path := dir
		if len(folder) > 0 {
			path = filepath.Join(dir, folder)
			if prevDir != path {
				if err = CreateDirºStr(rt, path); err != nil {
					return err
				}
				prevDir = path
			}
		}
		if err := unzipByIndex(rt, zfile, i, path); err != nil {
			return err
		}
		if rt.Owner.Settings.ProgressHandle != nil {
			prog.Increment(1)
		}
	}
	if rt.Owner.Settings.ProgressHandle != nil {
		prog.Complete()
	}
	return CloseUnzip(zfile)
}

// UnzipºUnzip unzip the specified file from the open zip archive
func UnzipºUnzip(rt *Runtime, zfile *UnzipFile, filename string, dir string) error {
	for i, finfo := range zfile.Archive.File {
		if filename == finfo.Name {
			if err := unzipByIndex(rt, zfile, i, dir); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
