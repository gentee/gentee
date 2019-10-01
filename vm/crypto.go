// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"crypto/md5"
	"crypto/sha256"

	"github.com/gentee/gentee/core"
)

func md5Hash(in []byte) (out *core.Buffer) {
	out = core.NewBuffer()
	md5Hash := md5.Sum(in)
	out.Data = md5Hash[:]
	return out
}

// Md5ºBuf returns md5 hash of the buffer
func Md5ºBuf(in *core.Buffer) (out *core.Buffer) {
	return md5Hash(in.Data)
}

// Md5ºStr returns md5 hash of the string as hex string
func Md5ºStr(in string) (out *core.Buffer) {
	return md5Hash([]byte(in))
}

func sha256Hash(in []byte) (out *core.Buffer) {
	out = core.NewBuffer()
	shaHash := sha256.Sum256(in)
	out.Data = shaHash[:]
	return out
}

// Sha256ºBuf returns md5 hash of the buffer
func Sha256ºBuf(in *core.Buffer) (out *core.Buffer) {
	return sha256Hash(in.Data)
}

// Sha256ºStr returns md5 hash of the string as hex string
func Sha256ºStr(in string) (out *core.Buffer) {
	return sha256Hash([]byte(in))
}
