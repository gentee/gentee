// Copyright 2018 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package stdlib

import (
	"time"

	"github.com/gentee/gentee/core"
)

// InitTime appends stdlib time functions to the virtual machine
func InitTime(vm *core.VirtualMachine) {
	NewStructType(vm, `time`, []string{
		`Year:int`, `Month:int`, `Day:int`,
		`Hour:int`, `Minute:int`, `Second:int`,
		`UTC:bool`,
	})

	for _, item := range []interface{}{
		SleepºInt, // Sleep( int )
	} {
		vm.StdLib().NewEmbed(item)
	}

	for _, item := range []embedInfo{
		{intºTime, `time`, `int`},                          // int( time )
		{timeºInt, `int`, `time`},                          // time( int, time )
		{DateTimeºInts, `int,int,int,int,int,int`, `time`}, // DateTime()
		{EqualºTimeTime, `time,time`, `bool`},              // binary ==
		{GreaterºTimeTime, `time,time`, `bool`},            // binary >
		{LessºTimeTime, `time,time`, `bool`},               // binary <
		{Now, ``, `time`},                                  // Now()
		{UTCºTime, `time`, `time`},                         // UTC()
	} {
		vm.StdLib().NewEmbedExt(item.Func, item.InTypes, item.OutType)
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

// DateTimeºInts returns time
func DateTimeºInts(rt *core.RunTime, year, month, day, hour, minute, second int64) *core.Struct {
	return fromTime(newTime(rt), time.Date(int(year), time.Month(month), int(day), int(hour), int(minute),
		int(second), 0, time.Local))
}

// EqualºTimeTime returns true if time structures are equal
func EqualºTimeTime(left, right *core.Struct) bool {
	return toTime(left).Equal(toTime(right))
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

// SleepºInt pauses the current script for at least the specified duration in milliseconds.
func SleepºInt(d int64) {
	time.Sleep(time.Duration(d) * time.Millisecond)
}

// UTCºTime converts time to UTC time.
func UTCºTime(rt *core.RunTime, local *core.Struct) *core.Struct {
	return fromTime(newTime(rt), toTime(local).UTC())
}
