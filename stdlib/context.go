// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gentee/gentee/core"
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

// strºInt converts integer value to string
func strºInt(val int64) string {
	return strconv.FormatInt(val, 10)
}

// InitContext appends stdlib context functions to the virtual machine
func InitContext(ws *core.Workspace) {
	for _, item := range []embedInfo{
		{core.Link{CtxSetºStrStr, 1000<<16 | core.EMBED}, `str,str`, `str`}, // CtxSet( str, str )
		{core.Link{CtxSetºStrBool, 1001<<16 | core.EMBED}, `str,bool`, `str`},
		// CtxSet( str, bool )
		{core.Link{CtxSetºStrFloat, 1002<<16 | core.EMBED}, `str,float`, `str`}, // CtxSet( str, float )
		{core.Link{CtxSetºStrInt, 1003<<16 | core.EMBED}, `str,int`, `str`},     // CtxSet( str, int )
		{core.Link{CtxValueºStr, 1004<<16 | core.EMBED}, `str`, `str`},          // CtxValue( str )
		{core.Link{CtxIsºStr, 1005<<16 | core.EMBED}, `str`, `bool`},            // CtxIs( str )
		{core.Link{CtxºStr, 1006<<16 | core.EMBED}, `str`, `str`},               // Ctx( str )
		{core.Link{CtxGetºStr, 1007<<16 | core.EMBED}, `str`, `str`},            // CtxGet( str )
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

// CtxIsºStr returns true if a context key exists
func CtxIsºStr(rt *core.RunTime, key string) bool {
	th := rt.Root.Threads
	th.ConstMutex.RLock()
	defer th.ConstMutex.RUnlock()
	if _, ok := th.Context[key]; ok {
		return true
	}
	return false
}

// CtxSetºStrStr sets a context value
func CtxSetºStrStr(rt *core.RunTime, key, value string) (string, error) {
	if utf8.RuneCountInString(key) > CtxLength {
		return ``, fmt.Errorf(ErrCtxLength, CtxLength)
	}
	th := rt.Root.Threads
	th.ConstMutex.Lock()
	th.Context[key] = value
	th.ConstMutex.Unlock()
	return value, nil
}

// strºFloat converts float value to string
func strºFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

// CtxSetºStrFloat assign a float to a context key
func CtxSetºStrFloat(rt *core.RunTime, key string, value float64) (string, error) {
	return CtxSetºStrStr(rt, key, strºFloat(value))
}

// strºBool converts boolean value to string
func strºBool(val bool) string {
	if val {
		return `true`
	}
	return `false`
}

// CtxSetºStrBool assign a bool to a context key
func CtxSetºStrBool(rt *core.RunTime, key string, value bool) (string, error) {
	return CtxSetºStrStr(rt, key, strºBool(value))
}

// CtxSetºStrInt assign an integer to a context key
func CtxSetºStrInt(rt *core.RunTime, key string, value int64) (string, error) {
	return CtxSetºStrStr(rt, key, strºInt(value))
}

// CtxValueºStr returns a context value
func CtxValueºStr(rt *core.RunTime, key string) string {
	th := rt.Root.Threads
	th.ConstMutex.RLock()
	defer th.ConstMutex.RUnlock()
	return th.Context[key]
}

// CtxºStr replaces context values in a string
func CtxºStr(rt *core.RunTime, input string) (string, error) {
	stack := make([]string, 0)
	ret, err := replace(rt, []rune(input), &stack)
	return string(ret), err
}

// CtxGetºStr replaces context values in the value of the key
func CtxGetºStr(rt *core.RunTime, key string) (string, error) {
	return CtxºStr(rt, CtxValueºStr(rt, key))
}

func replace(rt *core.RunTime, input []rune, stack *[]string) ([]rune, error) {
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
			th := rt.Root.Threads
			th.ConstMutex.RLock()
			value, ok = th.Context[string(name)]
			th.ConstMutex.RUnlock()
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
