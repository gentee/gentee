// Copyright 2020 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"io"
	"math/rand"
)

const (
	ProgressCopy = iota
	ProgressDownload
	ProgressCompress
	ProgressDecompress

	ProgressCompressCounter   = 200
	ProgressDecompressCounter = 201

	ProgStatusStart  = 0
	ProgStatusActive = 1
	ProgStatusEnd    = 2
)

type Progress struct {
	ID      uint32
	Type    int32
	Status  int32
	Total   int64
	Current int64
	Source  string
	Dest    string
	Ratio   float64
	Custom  interface{}

	handle ProgressFunc
}

type ProgressFunc func(*Progress) bool

type ProgressReader struct {
	*Progress
	reader io.Reader
}

type ProgressWriter struct {
	*Progress
	writer io.Writer
}

func NewProgress(rt *Runtime, total int64, ptype int64) *Progress {
	return &Progress{
		ID:     rand.Uint32(),
		Total:  total,
		Type:   int32(ptype),
		handle: rt.Owner.Settings.ProgressHandle,
	}
}

func NewProgressReader(reader io.Reader, progress *Progress) *ProgressReader {
	return &ProgressReader{
		Progress: progress,
		reader:   reader,
	}
}

func NewProgressWriter(writer io.Writer, progress *Progress) *ProgressWriter {
	return &ProgressWriter{
		Progress: progress,
		writer:   writer,
	}
}

func (progress *ProgressReader) Read(data []byte) (n int, err error) {
	n, err = progress.reader.Read(data)
	if err == nil && n > 0 {
		progress.Increment(int64(n))
	}
	return n, err
}

func (progress *ProgressWriter) Write(data []byte) (n int, err error) {
	n, err = progress.writer.Write(data)
	if err == nil && n > 0 {
		progress.Increment(int64(n))
	}
	return n, err
}

func (progress *Progress) Increment(inc int64) {
	progress.Status = ProgStatusActive
	progress.Current += int64(inc)
	if progress.Current >= progress.Total {
		progress.Ratio = 1
	} else {
		progress.Ratio = float64(progress.Current) / float64(progress.Total)
	}
	if progress.handle != nil {
		progress.handle(progress)
	}
}

func (progress *Progress) Complete() {
	progress.Ratio = 1
	progress.Status = ProgStatusEnd
	progress.handle(progress)
	return
}

func (progress *Progress) Start(source, dest string) {
	progress.Source = source
	progress.Dest = dest
	progress.Status = ProgStatusStart
	progress.handle(progress)
	return
}

func ProgressStart(rt *Runtime, total int64, ptype int64, src, dest string) int64 {
	if rt.Owner.Unique == nil {
		return 0
	}
	prog := NewProgress(rt, total, ptype)
	rt.Owner.Unique.Store(int64(prog.ID), prog)
	prog.Start(src, dest)
	return int64(prog.ID)
}

func ProgressInc(rt *Runtime, id, inc int64) {
	if rt.Owner.Unique == nil {
		return
	}
	if prog, loaded := rt.Owner.Unique.Load(id); loaded {
		if v, ok := prog.(*Progress); ok {
			v.Increment(inc)
		}
	}
}

func ProgressEnd(rt *Runtime, id int64) {
	if rt.Owner.Unique == nil {
		return
	}
	if prog, loaded := rt.Owner.Unique.LoadAndDelete(id); loaded {
		if v, ok := prog.(*Progress); ok {
			v.Complete()
		}
	}
}
