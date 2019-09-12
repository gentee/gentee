// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"fmt"
	"time"

	"github.com/gentee/gentee/core"
)

// InitTime appends stdlib time functions to the virtual machine
func InitTime(ws *core.Workspace) {
	NewStructType(ws, `time`, []string{
		`Year:int`, `Month:int`, `Day:int`,
		`Hour:int`, `Minute:int`, `Second:int`,
		`UTC:bool`,
	})

	for _, item := range []embedInfo{
		{intºTime, `time`, `int`},                                // int( time )
		{timeºInt, `int`, `time`},                                // time( int, time )
		{AddHoursºTimeInt, `time,int`, `time`},                   // AddHours(time,int) time
		{DateºInts, `int,int,int`, `time`},                       // Date(day, month, year)
		{DateTimeºInts, `int,int,int,int,int,int`, `time`},       // DateTime()
		{DaysºTime, `time`, `int`},                               // Days(time)
		{EqualºTimeTime, `time,time`, `bool`},                    // binary ==
		{FormatºTimeStr, `str,time`, `str`},                      // Format(time,str)
		{ParseTimeºStrStr, `str,str`, `time`},                    // ParseTime(str,str) time
		{GreaterºTimeTime, `time,time`, `bool`},                  // binary >
		{LessºTimeTime, `time,time`, `bool`},                     // binary <
		{Now, ``, `time`},                                        // Now()
		{core.Link{sleepºInt, 1008<<16 | core.EMBED}, `int`, ``}, // sleep(int)
		{UTCºTime, `time`, `time`},                               // UTC()
		{WeekdayºTime, `time`, `int`},                            // Weekday(time)
		{YearDayºTime, `time`, `int`},                            // YearDay(time) int
	} {
		ws.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
	}
}

func newTime(rt *core.RunTime) *core.Struct {
	return core.NewStructObj(rt, `time`)
}

func fromTime(it *core.Struct, in time.Time) *core.Struct {
	it.Values[0] = int64(in.Year())
	it.Values[1] = int64(in.Month())
	it.Values[2] = int64(in.Day())
	it.Values[3] = int64(in.Hour())
	it.Values[4] = int64(in.Minute())
	it.Values[5] = int64(in.Second())
	it.Values[6] = in.Location() == time.UTC
	return it
}

func toTime(it *core.Struct) time.Time {
	utc := time.Local
	if it.Values[6].(bool) {
		utc = time.UTC
	}
	return time.Date(int(it.Values[0].(int64)), time.Month(it.Values[1].(int64)),
		int(it.Values[2].(int64)), int(it.Values[3].(int64)), int(it.Values[4].(int64)),
		int(it.Values[5].(int64)), 0, utc)
}

// intºTime converts time to Unix time
func intºTime(it *core.Struct) int64 {
	return toTime(it).Unix()
}

// timeºInt converts Unix time to time
func timeºInt(rt *core.RunTime, unix int64) *core.Struct {
	return fromTime(newTime(rt), time.Unix(unix, 0))
}

// AddHoursºTimeInt adds/subtract hours
func AddHoursºTimeInt(rt *core.RunTime, it *core.Struct, hours int64) *core.Struct {
	return fromTime(newTime(rt), toTime(it).Add(time.Duration(hours)*time.Hour))
}

// DateºInts returns time
func DateºInts(rt *core.RunTime, year, month, day int64) *core.Struct {
	return DateTimeºInts(rt, year, month, day, 0, 0, 0)
}

// DateTimeºInts returns time
func DateTimeºInts(rt *core.RunTime, year, month, day, hour, minute, second int64) *core.Struct {
	return fromTime(newTime(rt), time.Date(int(year), time.Month(month), int(day), int(hour), int(minute),
		int(second), 0, time.Local))
}

// DaysºTime returns the days of the month
func DaysºTime(it *core.Struct) int64 {
	next := time.Date(int(it.Values[0].(int64)), time.Month(it.Values[1].(int64))+1, 0, 0, 0, 0, 0, time.UTC)
	next.Add(time.Duration(-24 * time.Hour))
	return int64(next.Day())
}

// EqualºTimeTime returns true if time structures are equal
func EqualºTimeTime(left, right *core.Struct) bool {
	return toTime(left).Equal(toTime(right))
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

// FormatºTimeStr formats the time
func FormatºTimeStr(layout string, t *core.Struct) string {
	return toTime(t).Format(layout2go(layout))
}

// GreaterºTimeTime returns true if left time structures are greater than right
func GreaterºTimeTime(left, right *core.Struct) bool {
	return toTime(left).After(toTime(right))
}

// LessºTimeTime returns true if left time structures are less than right
func LessºTimeTime(left, right *core.Struct) bool {
	return toTime(left).Before(toTime(right))
}

// Now returns the current time
func Now(rt *core.RunTime) *core.Struct {
	return fromTime(newTime(rt), time.Now())
}

// ParseTimeºStrStr parses a formatted string and returns the time value it represents
func ParseTimeºStrStr(rt *core.RunTime, layout, value string) (*core.Struct, error) {
	ret := newTime(rt)
	t, err := time.Parse(layout2go(layout), value)
	if err != nil {
		return ret, err
	}
	fmt.Println(`T`, t.Location(), time.UTC, t)
	return fromTime(ret, t.Local()), nil
}

// sleepºInt pauses the current script for at least the specified duration in milliseconds.
func sleepºInt(rt *core.RunTime, d int64) {
	rt.Thread.Sleep = d
}

// UTCºTime converts time to UTC time.
func UTCºTime(rt *core.RunTime, local *core.Struct) *core.Struct {
	return fromTime(newTime(rt), toTime(local).UTC())
}

// WeekdayºTime returns the day of the week specified by t.
func WeekdayºTime(rt *core.RunTime, t *core.Struct) int64 {
	return int64(toTime(t).Weekday())
}

// YearDayºTime returns the day of the year specified by t.
func YearDayºTime(t *core.Struct) int64 {
	return int64(toTime(t).YearDay())
}
