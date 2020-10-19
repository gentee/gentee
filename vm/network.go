// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gentee/gentee/core"
)

// Download downloads and saves the file by url.
func Download(rt *Runtime, url, filename string) (written int64, err error) {
	var (
		size   int64
		reader io.Reader
		prog   *Progress
	)
	isProgress := rt.Owner.Settings.ProgressHandle != nil

	if rt.Owner.Settings.IsPlayground || isProgress {
		hinfo, err := HeadInfo(rt, url)
		if err != nil {
			return 0, err
		}
		size = hinfo.Values[1].(int64)
		if rt.Owner.Settings.IsPlayground {
			if err = CheckPlaygroundLimits(rt.Owner, filename, size); err != nil {
				return 0, err
			}
		}
	}
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
	if isProgress && size > 0 {
		prog = NewProgress(rt, size, ProgressDownload)
		prog.Start(url, filename)
		reader = NewProgressReader(resp.Body, prog)
	} else {
		isProgress = false
		reader = resp.Body
	}

	if rt.Owner.Settings.IsPlayground {
		written, err = io.CopyN(out, reader, rt.Owner.Settings.Playground.SizeLimit)
		if err != nil && err != io.EOF {
			return written, err
		}
		if written >= rt.Owner.Settings.Playground.SizeLimit {
			return written, fmt.Errorf(`%s [%d MB]`, ErrorText(ErrPlaySize),
				rt.Owner.Settings.Playground.SizeLimit>>20)
		}
		if size == 0 {
			if err = CheckPlaygroundLimits(rt.Owner, filename, written); err != nil {
				return 0, err
			}
		}
	} else {
		written, err = io.Copy(out, reader)
	}
	if isProgress {
		prog.Complete()
	}
	return
}

// HeadInfo function read the header of url
func HeadInfo(rt *Runtime, url string) (*Struct, error) {
	hinfo := NewStruct(rt, &rt.Owner.Exec.Structs[HINFOSTRUCT])
	resp, err := http.Head(url)
	if err != nil {
		return hinfo, err
	}
	hinfo.Values[0] = int64(resp.StatusCode)
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	hinfo.Values[1] = int64(size)
	hinfo.Values[2] = resp.Header.Get("Content-Type")
	return hinfo, nil
}

// HTTPGet issues a GET to the specified URL.
func HTTPGet(rt *Runtime, url string) (buf *core.Buffer, err error) {
	var res *http.Response
	buf = core.NewBuffer()
	res, err = http.Get(url)
	if err == nil {
		if rt.Owner.Settings.IsPlayground {
			out := bytes.NewBuffer(nil)
			written, err := io.CopyN(out, res.Body, rt.Owner.Settings.Playground.SizeLimit)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if written >= rt.Owner.Settings.Playground.SizeLimit {
				return nil, fmt.Errorf(`%s [%d MB]`, ErrorText(ErrPlaySize),
					rt.Owner.Settings.Playground.SizeLimit>>20)
			}
			if err = CheckPlaygroundLimits(rt.Owner, core.RandName(), written); err != nil {
				return nil, err
			}
			buf.Data = out.Bytes()
		} else {
			buf.Data, err = ioutil.ReadAll(res.Body)
		}
		res.Body.Close()
	}
	return
}

// HTTPPage issues a GET to the specified URL and returns a string result.
func HTTPPage(rt *Runtime, url string) (string, error) {
	var (
		ret string
		buf []byte
	)
	res, err := http.Get(url)
	if err == nil {
		if rt.Owner.Settings.IsPlayground {
			out := bytes.NewBuffer(nil)
			written, err := io.CopyN(out, res.Body, rt.Owner.Settings.Playground.SizeLimit)
			if err != nil && err != io.EOF {
				return ``, err
			}
			if written >= rt.Owner.Settings.Playground.SizeLimit {
				return ``, fmt.Errorf(`%s [%d MB]`, ErrorText(ErrPlaySize),
					rt.Owner.Settings.Playground.SizeLimit>>20)
			}
			if err = CheckPlaygroundLimits(rt.Owner, core.RandName(), written); err != nil {
				return ``, err
			}
			ret = out.String()
		} else {
			buf, err = ioutil.ReadAll(res.Body)
		}
		res.Body.Close()
		if err == nil {
			ret = string(buf)
		}
	}
	return ret, err
}

// HTTPRequest send HTTP request to the specified URL and returns a string result.
func HTTPRequest(rt *Runtime, urlPath string, method string, params *core.Map, headers *core.Map) (ret string,
	err error) {
	var (
		req         *http.Request
		buf         []byte
		contentType string
		isForm      bool
		body        io.Reader
	)
	if rt.Owner.Settings.IsPlayground {
		return ``, fmt.Errorf(ErrorText(ErrPlayFunc), `HTTPRequest`)
	}
	for _, key := range headers.Keys {
		if key == `Content-Type` {
			contentType = headers.Data[key].(string)
			break
		}
	}
	if method != `GET` && len(params.Data) > 0 {
		if strings.HasPrefix(contentType, `application/json`) {
			var jsonBody []byte
			jsonBody, err = json.Marshal(params.Data)
			if err != nil {
				return
			}
			body = bytes.NewBuffer(jsonBody)
		} else {
			parVals := url.Values{}
			for _, key := range params.Keys {
				parVals.Add(key, params.Data[key].(string))
			}
			body = strings.NewReader(parVals.Encode())
			isForm = len(contentType) == 0
		}
	}
	if req, err = http.NewRequest(method, urlPath, body); err != nil {
		return
	}

	if method == `GET` && len(params.Data) > 0 {
		q := req.URL.Query()
		for _, key := range params.Keys {
			q.Add(key, params.Data[key].(string))
		}
		req.URL.RawQuery = q.Encode()
	}

	for _, key := range headers.Keys {
		req.Header.Add(key, headers.Data[key].(string))
	}
	//	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	if isForm {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		buf, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err == nil {
			ret = string(buf)
		}
	}
	return
}
