// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"time"

	"github.com/gentee/gentee/core"
)

// InitTime appends stdlib time functions to the virtual machine
func InitTime(ws *core.Workspace) {
	for _, item := range []embedInfo{
		//{core.Link{intºTime, 1041<<16 | core.EMBED}, `time`, `int`},              // int( time )
		//{core.Link{timeºInt, 1042<<16 | core.EMBED}, `int`, `time`},              // time( int, time )
		//{core.Link{AddHoursºTimeInt, 1043<<16 | core.EMBED}, `time,int`, `time`}, // AddHours(time,int) time
		//		{core.Link{DateºInts, 1044<<16 | core.EMBED}, `int,int,int`, `time`}, // Date(day, month, year)
		//{core.Link{DateTimeºInts, 1045<<16 | core.EMBED},
		//	`int,int,int,int,int,int`, `time`}, // DateTime()
		//{core.Link{DaysºTime, 1046<<16 | core.EMBED}, `time`, `int`},              // Days(time)
		//{core.Link{EqualºTimeTime, 1019<<16 | core.EMBED}, `time,time`, `bool`},   // binary ==
		//{core.Link{FormatºTimeStr, 1020<<16 | core.EMBED}, `str,time`, `str`},     // Format(time,str)
		//{core.Link{ParseTimeºStrStr, 1015<<16 | core.EMBED}, `str,str`, `time`},   // ParseTime(str,str) time
		//core.Link{GreaterºTimeTime, 1021<<16 | core.EMBED}, `time,time`, `bool`}, // binary >
		//{core.Link{LessºTimeTime, 1022<<16 | core.EMBED}, `time,time`, `bool`}, // binary <
		//{core.Link{Now, 1037<<16 | core.EMBED}, ``, `time`},             // Now()
		{core.Link{sleepºInt, 1008<<16 | core.EMBED}, `int`, ``},        // sleep(int)
		{core.Link{UTCºTime, 1038<<16 | core.EMBED}, `time`, `time`},    // UTC()
		{core.Link{WeekdayºTime, 1039<<16 | core.EMBED}, `time`, `int`}, // Weekday(time)
		{core.Link{YearDayºTime, 1040<<16 | core.EMBED}, `time`, `int`}, // YearDay(time) int
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
