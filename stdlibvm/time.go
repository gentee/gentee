// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlibvm

/*
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
// Now returns the current time
func Now(rt *core.RunTime) *core.Struct {
	return fromTime(newTime(rt), time.Now())
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
*/
