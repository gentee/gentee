// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"time"
)

func newTime(rt *Runtime) *Struct {
	return NewStruct(rt, &rt.Owner.Exec.Structs[TIMESTRUCT])
}

func toTime(it *Struct) time.Time {
	utc := time.Local
	if it.Values[6].(int64) == 1 {
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
	if in.Location() == time.UTC {
		it.Values[6] = int64(1)
	} else {
		it.Values[6] = int64(0)
	}
	return it
}

func replaceArr(in string, old, new []string) string {
	input := []rune(in)
	out := make([]rune, 0, len(input))
	lin := len(input)
	for i := 0; i < lin; i++ {
		eq := -1
		maxLen := lin - i
		for k, item := range old {
			litem := len([]rune(item))
			if maxLen >= litem && string(input[i:i+litem]) == item {
				eq = k
				break
			}
		}
		if eq >= 0 {
			out = append(out, []rune(new[eq])...)
			i += len([]rune(old[eq])) - 1
		} else {
			out = append(out, input[i])
		}
	}
	return string(out)
}

func layout2go(layout string) string {
	return replaceArr(layout, []string{
		`YYYY`, `YY`, `MMMM`, `MMM`, `MM`, `M`, `DD`, `D`, `dddd`, `ddd`,
		`HH`, `hh`, `h`, `PM`, `pm`, `mm`, `m`, `ss`, `s`, `tz`, `zz`, `z`,
	}, []string{
		`2006`, `06`, `January`, `Jan`, `01`, `1`, `02`, `2`, `Monday`, `Mon`,
		`15`, `03`, `3`, `PM`, `pm`, `04`, `4`, `05`, `5`, `MST`, `-0700`, `-07:00`,
	})
}

// intºTime converts time to Unix time
func intºTime(it *Struct) int64 {
	return toTime(it).Unix()
}

// timeºInt converts Unix time to time
func timeºInt(rt *Runtime, unix int64) *Struct {
	return fromTime(newTime(rt), time.Unix(unix, 0))
}

// AddHoursºTimeInt adds/subtract hours
func AddHoursºTimeInt(rt *Runtime, it *Struct, hours int64) *Struct {
	return fromTime(newTime(rt), toTime(it).Add(time.Duration(hours)*time.Hour))
}

// DateºInts returns time
func DateºInts(rt *Runtime, year, month, day int64) *Struct {
	return DateTimeºInts(rt, year, month, day, 0, 0, 0)
}

// DateTimeºInts returns time
func DateTimeºInts(rt *Runtime, year, month, day, hour, minute, second int64) *Struct {
	return fromTime(newTime(rt), time.Date(int(year), time.Month(month), int(day), int(hour), int(minute),
		int(second), 0, time.Local))
}

// DaysºTime returns the days of the month
func DaysºTime(it *Struct) int64 {
	next := time.Date(int(it.Values[0].(int64)), time.Month(it.Values[1].(int64))+1, 0, 0, 0, 0, 0, time.UTC)
	next.Add(time.Duration(-24 * time.Hour))
	return int64(next.Day())
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

// Now returns the current time
func Now(rt *Runtime) *Struct {
	return fromTime(newTime(rt), time.Now())
}

// UTCºTime converts time to UTC time.
func UTCºTime(rt *Runtime, local *Struct) *Struct {
	return fromTime(newTime(rt), toTime(local).UTC())
}

// WeekdayºTime returns the day of the week specified by t.
func WeekdayºTime(rt *Runtime, t *Struct) int64 {
	return int64(toTime(t).Weekday())
}

// YearDayºTime returns the day of the year specified by t.
func YearDayºTime(t *Struct) int64 {
	return int64(toTime(t).YearDay())
}
