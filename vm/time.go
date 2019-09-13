// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"time"

	stdlib "github.com/gentee/gentee/stdlibvm"
)

func newTime(rt *Runtime) *Struct {
	return NewStruct(rt, &rt.Owner.Exec.Structs[TIMESTRUCT])
}

func toTime(it *Struct) time.Time {
	utc := time.Local
	if it.Values[6].(bool) {
		utc = time.UTC
	}
	return time.Date(int(it.Values[0].(int64)), time.Month(it.Values[1].(int64)),
		int(it.Values[2].(int64)), int(it.Values[3].(int64)), int(it.Values[4].(int64)),
		int(it.Values[5].(int64)), 0, utc)
}

func fromTime(it *Struct, in time.Time) *Struct {
	it.Values[0] = int64(in.Year())
	it.Values[1] = int64(in.Month())
	it.Values[2] = int64(in.Day())
	it.Values[3] = int64(in.Hour())
	it.Values[4] = int64(in.Minute())
	it.Values[5] = int64(in.Second())
	it.Values[6] = in.Location() == time.UTC
	return it
}

func layout2go(layout string) string {
	return stdlib.ReplaceArr(layout, []string{
		`YYYY`, `YY`, `MMMM`, `MMM`, `MM`, `M`, `DD`, `D`, `dddd`, `ddd`,
		`HH`, `hh`, `h`, `PM`, `pm`, `mm`, `m`, `ss`, `s`, `tz`, `zz`, `z`,
	}, []string{
		`2006`, `06`, `January`, `Jan`, `01`, `1`, `02`, `2`, `Monday`, `Mon`,
		`15`, `03`, `3`, `PM`, `pm`, `04`, `4`, `05`, `5`, `MST`, `-0700`, `-07:00`,
	})
}

// EqualºTimeTime returns true if time structures are equal
func EqualºTimeTime(left, right *Struct) int64 {
	if toTime(left).Equal(toTime(right)) {
		return 1
	}
	return 0
}

// FormatºTimeStr formats the time
func FormatºTimeStr(layout string, t *Struct) string {
	return toTime(t).Format(layout2go(layout))
}

// GreaterºTimeTime returns true if left time structures are greater than right
func GreaterºTimeTime(left, right *Struct) int64 {
	if toTime(left).After(toTime(right)) {
		return 1
	}
	return 0
}

// LessºTimeTime returns true if left time structures are less than right
func LessºTimeTime(left, right *Struct) int64 {
	if toTime(left).Before(toTime(right)) {
		return 1
	}
	return 0
}

// ParseTimeºStrStr parses a formatted string and returns the time value it represents
func ParseTimeºStrStr(rt *Runtime, layout, value string) (*Struct, error) {
	ret := newTime(rt)
	t, err := time.Parse(layout2go(layout), value)
	if err != nil {
		return ret, err
	}
	fmt.Println(`T`, t.Location(), time.UTC, t)
	return fromTime(ret, t.Local()), nil
}
