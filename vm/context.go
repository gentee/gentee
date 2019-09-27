// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"strings"

	//	"strings"
	"unicode/utf8"
	//	"github.com/gentee/gentee/core"

	stdlib "github.com/gentee/gentee/stdlibvm"
)

const (
	// CtxChar is a context boundery character
	CtxChar = '#'
	// CtxLength is the max length of the context key
	CtxLength = 64
	// CtxDeep is the max deep in replace function
	CtxDeep = 16
)

var (
	// ErrCtxLength is returned when the key is too long
	ErrCtxLength = `key length is longer than %d characters`
	// ErrCtxLoop is returned when there is a loop in values of context
	ErrCtxLoop = `%s key refers to itself`
	// ErrCtxDeep is returned if the maximum depth reached
	ErrCtxDeep = `maximum depth reached`
)

// CtxIsºStr returns true if a context key exists
func CtxIsºStr(rt *Runtime, key string) int64 {
	rt.Owner.CtxMutex.RLock()
	defer rt.Owner.CtxMutex.RUnlock()
	if _, ok := rt.Owner.Context[key]; ok {
		return 1
	}
	return 0
}

// CtxSetºStrStr sets a context value
func CtxSetºStrStr(rt *Runtime, key, value string) (string, error) {
	if utf8.RuneCountInString(key) > CtxLength {
		return ``, fmt.Errorf(ErrCtxLength, CtxLength)
	}
	rt.Owner.CtxMutex.Lock()
	rt.Owner.Context[key] = value
	rt.Owner.CtxMutex.Unlock()
	return value, nil
}

// CtxSetºStrFloat assign a float to a context key
func CtxSetºStrFloat(rt *Runtime, key string, value float64) (string, error) {
	return CtxSetºStrStr(rt, key, strºFloat(value))
}

// CtxSetºStrBool assign a bool to a context key
func CtxSetºStrBool(rt *Runtime, key string, value int64) (string, error) {
	return CtxSetºStrStr(rt, key, stdlib.StrºBool(value))
}

// CtxSetºStrInt assign an integer to a context key
func CtxSetºStrInt(rt *Runtime, key string, value int64) (string, error) {
	return CtxSetºStrStr(rt, key, strºInt(value))
}

// CtxValueºStr returns a context value
func CtxValueºStr(rt *Runtime, key string) string {
	rt.Owner.CtxMutex.RLock()
	defer rt.Owner.CtxMutex.RUnlock()
	return rt.Owner.Context[key]
}

// CtxºStr replaces context values in a string
func CtxºStr(rt *Runtime, input string) (string, error) {
	stack := make([]string, 0)
	ret, err := replace(rt, []rune(input), &stack)
	return string(ret), err
}

// CtxGetºStr replaces context values in the value of the key
func CtxGetºStr(rt *Runtime, key string) (string, error) {
	return CtxºStr(rt, CtxValueºStr(rt, key))
}

func replace(rt *Runtime, input []rune, stack *[]string) ([]rune, error) {
	if len(input) == 0 || strings.IndexRune(string(input), CtxChar) == -1 {
		return input, nil
	}
	var (
		err        error
		isName, ok bool
		value      string
		tmp        []rune
	)
	result := make([]rune, 0, len(input))
	name := make([]rune, 0, CtxLength+1)

	for i := 0; i < len(input); i++ {
		r := input[i]
		if r != CtxChar {
			if isName {
				name = append(name, r)
				if len(name) > CtxChar {
					result = append(append(result, CtxChar), name...)
					isName = false
					name = name[:0]
				}
			} else {
				result = append(result, r)
			}
			continue
		}
		if isName {
			rt.Owner.CtxMutex.RLock()
			value, ok = rt.Owner.Context[string(name)]
			rt.Owner.CtxMutex.RUnlock()
			if ok {
				if len(*stack) < CtxDeep {
					for _, item := range *stack {
						if item == string(name) {
							return result, fmt.Errorf(ErrCtxLoop, item)
						}
					}
				} else {
					return result, fmt.Errorf(ErrCtxDeep)
				}
				*stack = append(*stack, string(name))
				if tmp, err = replace(rt, []rune(value), stack); err != nil {
					return result, err
				}
				*stack = (*stack)[:len(*stack)-1]
				result = append(result, tmp...)
			} else {
				result = append(append(result, CtxChar), name...)
				i--
			}
			name = name[:0]
		}
		isName = !isName
	}
	if isName {
		result = append(append(result, CtxChar), name...)
	}
	return result, nil
}
