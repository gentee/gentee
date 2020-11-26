// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"math/rand"

	"github.com/gentee/gentee/core"
	"golang.org/x/crypto/scrypt"
)

const (
	scryptN = 32768 //262144
	scryptr = 8
	scryptp = 1
	AESKey  = 32
	SaltLen = 16
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

func RandomBytes(size int) (ret []byte, err error) {
	ret = make([]byte, size)
	_, err = rand.Read(ret)
	return
}

func DerivePassphrase(passphrase, salt []byte, keyLen int) ([]byte, []byte, error) {
	var (
		err error
		key []byte
	)
	if len(salt) == 0 {
		if salt, err = RandomBytes(SaltLen); err != nil {
			return nil, nil, err
		}
	}
	key, err = scrypt.Key(passphrase, salt, scryptN, scryptr, scryptp, keyLen)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}

func AESEncrypt(passphrase, data []byte) ([]byte, error) {
	key, salt, err := DerivePassphrase(passphrase, nil, AESKey)
	if err != nil {
		return nil, err
	}
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce, err := RandomBytes(gcm.NonceSize())
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return append(ciphertext, salt...), nil
}

func AESDecrypt(passphrase, data []byte) ([]byte, error) {
	key, _, err := DerivePassphrase(passphrase, data[len(data)-SaltLen:], AESKey)
	if err != nil {
		return nil, err
	}
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	data = data[:len(data)-SaltLen]
	return gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
}

// Random(int size) buf
func RandomBuf(size int64) (buf *core.Buffer, err error) {
	var salt []byte
	if salt, err = RandomBytes(int(size)); err != nil {
		return
	}
	buf = core.NewBuffer()
	buf.Data = salt
	return
}

// AESEncrypt(str, buf) buf
func AESEncryptBuf(passphrase string, buf *core.Buffer) (bufout *core.Buffer, err error) {
	bufout = core.NewBuffer()
	bufout.Data, err = AESEncrypt([]byte(passphrase), buf.Data)
	return
}

// AESDecrypt(str, buf) buf
func AESDecryptBuf(passphrase string, buf *core.Buffer) (bufout *core.Buffer, err error) {
	bufout = core.NewBuffer()
	bufout.Data, err = AESDecrypt([]byte(passphrase), buf.Data)
	return
}
