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

	ProgressStart  = 0
	ProgressActive = 1
	ProgressEnd    = 2
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
	Progress
	reader io.Reader
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
		Progress: *progress,
		reader:   reader,
	}
}

func (progress *ProgressReader) Read(data []byte) (n int, err error) {
	n, err = progress.reader.Read(data)
	if err == nil && n > 0 {
		progress.Status = ProgressActive
		progress.Current += int64(n)
		if progress.Current >= progress.Total {
			progress.Ratio = 1
		} else {
			progress.Ratio = float64(progress.Current) / float64(progress.Total)
		}
		if progress.handle != nil {
			progress.handle(&progress.Progress)
		}
	}
	return n, err
}

func (progress *Progress) Complete() {
	progress.Ratio = 1
	progress.Status = ProgressEnd
	progress.handle(progress)
	return
}

func (progress *Progress) Start(source, dest string) {
	progress.Source = source
	progress.Dest = dest
	progress.Status = ProgressStart
	progress.handle(progress)
	return
}
