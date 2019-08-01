// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gentee/gentee/core"
)

// InitNetwork appends stdlib network functions to the virtual machine
func InitNetwork(ws *core.Workspace) {
	for _, item := range []interface{}{
		Download, // Download(str, str) int
		HTTPGet,  // HTTPGet(str) buf
		HTTPPage, // HTTPPage(str) str
	} {
		ws.StdLib().NewEmbed(item)
	}
}

// Download downloads and saves the file by url.
func Download(url, filename string) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	out, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return io.Copy(out, resp.Body)
}

// HTTPGet issues a GET to the specified URL.
func HTTPGet(url string) (buf *core.Buffer, err error) {
	var res *http.Response
	buf = core.NewBuffer()
	res, err = http.Get(url)
	if err == nil {
		buf.Data, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
	}
	return
}

// HTTPPage issues a GET to the specified URL and returns a string result.
func HTTPPage(url string) (string, error) {
	var (
		ret string
		buf []byte
	)
	res, err := http.Get(url)
	if err == nil {
		buf, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err == nil {
			ret = string(buf)
		}
	}
	return ret, err
}
