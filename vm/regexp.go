// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"regexp"

	"github.com/gentee/gentee/core"
)

// FindRegExpºStrStr returns an array of all successive matches of the expression
func FindRegExpºStrStr(src, rePattern string) (*core.Array, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return nil, err
	}
	list := re.FindAllStringSubmatch(src, -1)
	out := core.NewArray()
	for _, ilist := range list {
		sub := core.NewArray()
		for _, sublist := range ilist {
			sub.Data = append(sub.Data, sublist)
		}
		out.Data = append(out.Data, sub)
	}
	return out, nil
}

// FindFirstRegExpºStrStr returns an array of the first successive matches of the expression
func FindFirstRegExpºStrStr(src, rePattern string) (*core.Array, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return nil, err
	}
	list := re.FindStringSubmatch(src)
	out := core.NewArray()
	for _, sub := range list {
		out.Data = append(out.Data, sub)
	}
	return out, nil
}

// MatchºStrStr reports whether the string s contains any match of the regular expression
func MatchºStrStr(s string, rePattern string) (int64, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return 0, err
	}
	if re.MatchString(s) {
		return 1, nil
	}
	return 0, nil
}

// RegExpºStrStr returns the first found match of the expression
func RegExpºStrStr(src, rePattern string) (ret string, err error) {
	var re *regexp.Regexp
	re, err = regexp.Compile(rePattern)
	if err != nil {
		return
	}
	list := re.FindStringSubmatch(src)
	if len(list) > 1 {
		ret = list[1]
	}
	return
}

// ReplaceRegExpºStrStr returns a copy of src, replacing matches of the Regexp with the replacement string
func ReplaceRegExpºStrStr(src, rePattern, repl string) (string, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return ``, err
	}
	return re.ReplaceAllString(src, repl), nil
}
