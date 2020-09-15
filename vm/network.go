// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gentee/gentee/core"
)

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

// HTTPRequest send HTTP request to the specified URL and returns a string result.
func HTTPRequest(urlPath string, method string, params *core.Map, headers *core.Map) (ret string,
	err error) {
	var (
		req         *http.Request
		buf         []byte
		contentType string
		isForm      bool
		body        io.Reader
	)
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
