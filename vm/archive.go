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

type UntargzFile struct {
	Name      string
	File      *os.File
	GzReader  *gzip.Reader
	TarReader *tar.Reader
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

// openUntargz opens tar.gz file and returns its handle
func openUntargz(rt *Runtime, gzfile string) (*UntargzFile, error) {
	if rt.Owner.Settings.IsPlayground {
		if err := CheckPlaygroundLimits(rt.Owner, gzfile, NoLimit); err != nil {
			return nil, err
		}
	}
	file, err := os.Open(gzfile)
	if err != nil {
		return nil, err
	}
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	return &UntargzFile{Name: gzfile, File: file, GzReader: gzReader,
		TarReader: tar.NewReader(gzReader)}, nil
}

// closeUntargz closes the opened gz file
func closeUntargz(gzfile *UntargzFile) (err error) {
	if err = gzfile.GzReader.Close(); err == nil {
		err = gzfile.File.Close()
	}
	return
}

// ReadTarGz gets the list of files in the .tar.gz file
func ReadTarGz(rt *Runtime, filename string) (*core.Array, error) {
	gzfile, err := openUntargz(rt, filename)
	if err != nil {
		return nil, err
	}
	ret := core.NewArray()
	for true {
		header, err := gzfile.TarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch header.Typeflag {
		case tar.TypeDir, tar.TypeReg:
			fi := fromFileInfo(header.FileInfo(), NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT]))
			fi.Values[0] = strings.TrimRight(header.Name, `/`)
			ret.Data = append(ret.Data, fi)
		default:
			return nil, fmt.Errorf("ReadTarGz: uknown type: %d in %s", header.Typeflag, header.Name)
		}
	}
	return ret, closeUntargz(gzfile)
}

// ReadZip returns the file list of zip
func ReadZip(rt *Runtime, zipfile string) (*core.Array, error) {
	var err error
	zf, err := openUnzip(rt, zipfile)
	if err != nil {
		return nil, err
	}
	ret := core.NewArray()
	for _, finfo := range zf.Reader.File {
		fi := fromFileInfo(finfo.FileInfo(), NewStruct(rt, &rt.Owner.Exec.Structs[FINFOSTRUCT]))
		fi.Values[0] = strings.TrimRight(finfo.Name, `/`)
		ret.Data = append(ret.Data, fi)
	}
	return ret, closeUnzip(zf)
}

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

// UnpackTarGz unpacks a .tar.gz file to the specified folder
func UnpackTarGz(rt *Runtime, filename string, dir string) error {
	empty := core.NewArray()
	return UnpackTarGzºStr(rt, filename, dir, empty, empty)
}

// UnpackTarGzºStr unpacks a .tar.gz file to the specified folder
func UnpackTarGzºStr(rt *Runtime, filename string, dir string, patterns *core.Array,
	ignore *core.Array) error {
	gzfile, err := openUntargz(rt, filename)
	if err != nil {
		return err
	}
	created := make(map[string]bool)
	for true {
		header, err := gzfile.TarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		var path string
		name := header.Name
		if ok, err := matchName(name, patterns, ignore); err != nil {
			return err
		} else if !ok {
			continue
		}
		if path, err = prepareDecompress(rt, name, header.FileInfo(), dir, created); err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir, tar.TypeReg:
			if err = unpackFile(rt, header.FileInfo(), gzfile.TarReader, path); err != nil {
				return err
			}
		default:
			return fmt.Errorf("UnpackTarGz: uknown type: %d in %s", header.Typeflag, header.Name)
		}
	}
	return closeUntargz(gzfile)
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
		if err = unpackFile(rt, fhead.FileInfo(), reader, path); err != nil {
			return err
		}
		if rt.Owner.Settings.ProgressHandle != nil {
			prog.Increment(1)
		}
		reader.Close()
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
			if ok, err = MatchPath(pattern, filename); err != nil {
				return false, err
			} else if ok != 0 {
				return false, nil
			}
		}
	}
	if len(patterns.Data) > 0 {
		for _, item := range patterns.Data {
			if ok, err = MatchPath(item.(string), filename); err != nil {
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
