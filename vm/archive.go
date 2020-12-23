// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"archive/tar"
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentee/gentee/core"
)

type UnzipFile struct {
	Name   string
	Reader *zip.ReadCloser
}

type ZipFile struct {
	Name   string
	File   *os.File
	Writer *zip.Writer
}

type GzFile struct {
	Name      string
	File      *os.File
	GzWriter  *gzip.Writer
	TarWriter *tar.Writer
}

type Pack interface {
	FileName() string
	Header(finfo os.FileInfo, packname string) (io.Writer, error)
}

func (zf *ZipFile) FileName() string {
	return zf.Name
}

func (zf *ZipFile) Header(finfo os.FileInfo, packname string) (io.Writer, error) {
	header, err := zip.FileInfoHeader(finfo)
	if err != nil {
		return nil, err
	}
	header.Name = packname
	if finfo.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}
	return zf.Writer.CreateHeader(header)
}

func (gzf *GzFile) FileName() string {
	return gzf.Name
}

func (gzf *GzFile) Header(finfo os.FileInfo, packname string) (io.Writer, error) {
	tarhead, err := tar.FileInfoHeader(finfo, ``)
	if err != nil {
		return nil, err
	}
	if finfo.IsDir() {
		packname += "/"
	}
	tarhead.Name = packname
	if err = gzf.TarWriter.WriteHeader(tarhead); err != nil {
		return nil, err
	}
	return gzf.TarWriter, nil
}

func ArchiveName(finfo *Struct, root string) string {
	name := finfo.Values[0].(string)
	dir := strings.TrimPrefix(finfo.Values[5].(string), root)
	if len(root) == 0 || len(dir) == len(finfo.Values[5].(string)) {
		return name
	}
	packname := strings.Trim(strings.ReplaceAll(dir, `\\`, `/`), `/`)
	if len(packname) > 0 {
		packname += `/`
	}
	return packname + name
}

func CompressFile(rt *Runtime, pack Pack, filename, packname string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return err
		}
	}
	finfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, pack.FileName(), finfo.Size()); err != nil {
			return err
		}
	}
	if len(packname) == 0 {
		packname = finfo.Name()
	}
	writer, err := pack.Header(finfo, packname)
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
		prog.Start(filename, pack.FileName())
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

// CloseTarGz closes the created tar.gz file
func CloseTarGz(gzfile *GzFile) (err error) {
	if err = gzfile.TarWriter.Close(); err == nil {
		if err = gzfile.GzWriter.Close(); err == nil {
			err = gzfile.File.Close()
		}
	}
	return
}

// CloseZip closes the created zip file
func CloseZip(zfile *ZipFile) (err error) {
	if err = zfile.Writer.Close(); err == nil {
		err = zfile.File.Close()
	}
	return
}

// CreateZip creates zip file and returns its handle
func CreateZip(rt *Runtime, filename string) (*ZipFile, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return nil, err
		}
	}
	zipfile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	archive := zip.NewWriter(zipfile)
	archive.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	return &ZipFile{Name: filename, File: zipfile, Writer: archive}, nil
}

// CreateTarGz creates tar.gz file and returns its handle
func CreateTarGz(rt *Runtime, filename string) (*GzFile, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, filename, NoLimit); err != nil {
			return nil, err
		}
	}
	gzfile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	var gw *gzip.Writer
	if gw, err = gzip.NewWriterLevel(gzfile, gzip.BestCompression); err != nil {
		return nil, err
	}
	return &GzFile{Name: filename, File: gzfile, GzWriter: gw, TarWriter: tar.NewWriter(gw)}, nil
}

// openUnzip opens zip file and returns its handle
func openUnzip(rt *Runtime, zipfile string) (*UnzipFile, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, zipfile, NoLimit); err != nil {
			return nil, err
		}
	}
	archive, err := zip.OpenReader(zipfile)
	if err != nil {
		return nil, err
	}
	return &UnzipFile{Name: zipfile, Reader: archive}, nil
}

// closeUnzip closes the opened zip file
func closeUnzip(zfile *UnzipFile) (err error) {
	return zfile.Reader.Close()
}

// ReadZip returns the file list of zip
func ReadZip(rt *Runtime, zipfile string) (*core.Array, error) {
	var err error
	zf, err := openUnzip(rt, zipfile)
	ret := core.NewArray()
	for _, finfo := range zf.Reader.File {
		fi := fromFileInfo(finfo.FileInfo(), NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT]))
		fi.Values[0] = strings.TrimRight(finfo.Name, `/`)
		ret.Data = append(ret.Data, fi)
	}
	if err = closeUnzip(zf); err != nil {
		return nil, err
	}
	return ret, nil
}

//=====================================

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

func untarFile(rt *Runtime, tr *tar.Reader, dest string, header *tar.Header, gzfile string) error {
	target, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		header.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer func() {
		target.Close()
		os.Chtimes(dest, header.FileInfo().ModTime(), header.FileInfo().ModTime())
	}()
	var (
		prog   *Progress
		writer io.Writer
	)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog = NewProgress(rt, header.Size, ProgressDecompress)
		prog.Start(gzfile, dest)
		writer = NewProgressWriter(target, prog)
	} else {
		writer = target
	}
	_, err = io.Copy(writer, tr)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog.Complete()
	}
	return err
}

// UnTarGz unpacks a .tar.gz file to the specified folder
func UnTarGz(rt *Runtime, gzfile string, dir string) error {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, gzfile, NoLimit); err != nil {
			return err
		}
	}
	file, err := os.Open(gzfile)
	if err != nil {
		return err
	}
	defer file.Close()
	gzReader, err := gzip.NewReader(file)
	defer gzReader.Close()
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzReader)
	var prevDir string

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		folder := filepath.Dir(strings.TrimRight(header.Name, `/`))
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
		dest := filepath.Join(path, filepath.Base(strings.TrimRight(header.Name, `/`)))
		if rt.Owner.Settings.IsPlayground {
			if err := CheckPlaygroundLimits(rt.Owner, dest, header.Size); err != nil {
				return err
			}
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(dest, header.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeReg:
			if err = untarFile(rt, tarReader, dest, header, gzfile); err != nil {
				return err
			}
		default:
			return fmt.Errorf("UnTarGz: uknown type: %d in %s", header.Typeflag, header.Name)
		}
	}
	return nil
}

//================================

// TarGz packs a file or directory to tar.gz
func TarGz(rt *Runtime, targzfile string, path string) error {
	var (
		err  error
		list *core.Array
	)
	if path, err = filepath.Abs(path); err != nil {
		return err
	}
	if list, err = archiveList(rt, path); err != nil {
		return err
	}
	gzfile, err := CreateTarGz(rt, targzfile)
	if err != nil {
		return err
	}
	if err = packFiles(rt, gzfile, list, path); err != nil {
		return err
	}
	return CloseTarGz(gzfile)
}

// UnpackZip unpacks a zip file to the specified folder
func UnpackZip(rt *Runtime, zipfile string, dir string) error {
	empty := core.NewArray()
	return UnpackZipºStr(rt, zipfile, dir, empty, empty)
}

// UnpackZipºStr unpacks a zip file to the specified folder
func UnpackZipºStr(rt *Runtime, zipfile string, dir string, patterns *core.Array,
	ignore *core.Array) error {
	var (
		err error
	)
	if rt.Owner.Settings.IsPlayground {
		if err = CheckPlaygroundLimits(rt.Owner, dir, NoLimit); err != nil {
			return err
		}
	}
	if dir, err = filepath.Abs(dir); err != nil {
		return err
	}
	zfile, err := openUnzip(rt, zipfile)
	if err != nil {
		return err
	}
	var (
		prog *Progress
	)
	created := make(map[string]bool)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog = NewProgress(rt, int64(len(zfile.Reader.File)), ProgressDecompressCounter)
		prog.Start(zipfile, ``)
	}
	for _, fhead := range zfile.Reader.File {
		var path string
		name := fhead.Name
		if ok, err := matchName(name, patterns, ignore); err != nil {
			return err
		} else if !ok {
			continue
		}
		if path, err = prepareDecompress(rt, name, fhead.FileInfo(), dir, created); err != nil {
			return err
		}
		reader, err := fhead.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		if err = unpackFile(rt, fhead.FileInfo(), reader, path); err != nil {
			return err
		}
		if rt.Owner.Settings.ProgressHandle != nil {
			prog.Increment(1)
		}
	}
	if rt.Owner.Settings.ProgressHandle != nil {
		prog.Complete()
	}
	return closeUnzip(zfile)
}

// ZipºStr packs a file or directory
func ZipºStr(rt *Runtime, zipfile string, path string) error {
	var (
		err  error
		list *core.Array
	)
	if path, err = filepath.Abs(path); err != nil {
		return err
	}
	if list, err = archiveList(rt, path); err != nil {
		return err
	}
	zfile, err := CreateZip(rt, zipfile)
	if err != nil {
		return err
	}
	if err = packFiles(rt, zfile, list, path); err != nil {
		return err
	}
	return CloseZip(zfile)
}

func archiveList(rt *Runtime, path string) (*core.Array, error) {
	var (
		err   error
		finfo os.FileInfo
	)
	if rt.Owner.Settings.IsPlayground {
		if err = CheckPlaygroundLimits(rt.Owner, path, NoLimit); err != nil {
			return nil, err
		}
	}
	list := core.NewArray()
	if path, err = filepath.Abs(path); err != nil {
		return nil, err
	}
	if finfo, err = os.Stat(path); err != nil {
		return nil, err
	}
	if finfo.IsDir() {
		if err = readDir(rt, list, path, Recursive, core.NewArray(), core.NewArray()); err != nil {
			return nil, err
		}
	} else {
		item := fromFileInfo(finfo, NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT]))
		item.Values[5] = filepath.Dir(path)
		list.Data = append(list.Data, item)
	}
	return list, nil
}

func matchName(filename string, patterns *core.Array, ignore *core.Array) (bool, error) {
	var (
		ok  int64
		err error
	)
	for _, item := range ignore.Data {
		if pattern := item.(string); len(pattern) > 0 {
			if ok, err = MatchPathºStrBool(pattern, filename, 0); err != nil {
				return false, err
			} else if ok != 0 {
				return false, nil
			}
		}
	}
	if len(patterns.Data) > 0 {
		for _, item := range patterns.Data {
			if ok, err = MatchPathºStrBool(item.(string), filename, 0); err != nil {
				return false, err
			}
			if ok != 0 {
				break
			}
		}
		if ok == 0 {
			return false, nil
		}
	}
	return true, nil
}

func packFiles(rt *Runtime, pack Pack, list *core.Array, path string) error {
	var err error
	if len(list.Data) != 1 {
		var (
			prog *Progress
		)
		if rt.Owner.Settings.ProgressHandle != nil {
			prog = NewProgress(rt, int64(len(list.Data)), ProgressCompressCounter)
			prog.Start(pack.FileName(), ``)
		}
		for _, item := range list.Data {
			if err = CompressFile(rt, pack, FileInfoToPath(item.(*Struct)),
				ArchiveName(item.(*Struct), path)); err != nil {
				return err
			}
			if rt.Owner.Settings.ProgressHandle != nil {
				prog.Increment(1)
			}
		}
		if rt.Owner.Settings.ProgressHandle != nil {
			prog.Complete()
		}
	} else {
		ifile := list.Data[0].(*Struct)
		name := ifile.Values[0].(string)
		err = CompressFile(rt, pack, filepath.Join(ifile.Values[5].(string), name), name)
	}
	return err
}

func prepareDecompress(rt *Runtime, filename string, finfo os.FileInfo, dir string,
	created map[string]bool) (path string, err error) {
	folder := filepath.Dir(strings.TrimRight(filename, `/`))
	path = dir
	if len(folder) > 0 {
		path = filepath.Join(dir, folder)
		if !created[path] {
			if err = CreateDirºStr(rt, path); err != nil {
				return
			}
			created[path] = true
		}
	}
	path = filepath.Join(path, finfo.Name())
	if rt.Owner.Settings.IsPlayground {
		if err = CheckPlaygroundLimits(rt.Owner, path, finfo.Size()); err != nil {
			return
		}
	}
	return
}

func unpackFile(rt *Runtime, finfo os.FileInfo, reader io.Reader, dest string) error {
	if finfo.IsDir() {
		return os.MkdirAll(dest, finfo.Mode())
	}
	target, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, finfo.Mode())
	if err != nil {
		return err
	}
	defer func() {
		target.Close()
		os.Chtimes(dest, finfo.ModTime(), finfo.ModTime())
	}()
	var (
		prog   *Progress
		writer io.Writer
	)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog = NewProgress(rt, finfo.Size(), ProgressDecompress)
		prog.Start(dest, dest)
		writer = NewProgressWriter(target, prog)
	} else {
		writer = target
	}
	_, err = io.Copy(writer, reader)
	if rt.Owner.Settings.ProgressHandle != nil {
		prog.Complete()
	}
	return err
}
