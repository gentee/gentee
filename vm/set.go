// Copyright 2019 Alexey Krivonogov. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/gentee/gentee/core"
)

// arrºSet converts set to array of integers
func arrºSet(set *core.Set) *core.Array {
	ret := core.NewArray()
	for i, v := range set.Data {
		for pos := uint64(0); pos < 64; pos++ {
			if v&(1<<pos) != 0 {
				ret.Data = append(ret.Data, int64(i<<6)+int64(pos))
			}
		}
	}
	return ret
}

// AssignAddºSetSet appends set to set
func AssignAddºSetSet(set interface{}, value interface{}) (interface{}, error) {
	for i, v := range value.(*core.Set).Data {
		for pos := uint64(0); pos < 64; pos++ {
			if v&(1<<pos) != 0 {
				set.(*core.Set).Set(int64(i<<6)+int64(pos), true)
			}
		}
	}
	return set, nil
}

// BitAndºSetSet equals set & set
func BitAndºSetSet(left *core.Set, right *core.Set) *core.Set {
	return bitSet(left, right, true)
}

// BitNotºSet changes boolean value of set
func BitNotºSet(set *core.Set) *core.Set {
	ret := core.NewSet()
	ret.Data = make([]uint64, len(set.Data))
	for i, v := range set.Data {
		ret.Data[i] = ^v
	}
	return ret
}

// BitOrºSetSet equals set & set
func BitOrºSetSet(left *core.Set, right *core.Set) *core.Set {
	return bitSet(left, right, false)
}

func bitSet(left *core.Set, right *core.Set, and bool) *core.Set {
	ret := core.NewSet()
	if len(left.Data) < len(right.Data) {
		left, right = right, left
	}
	ret.Data = make([]uint64, len(left.Data))
	for i, v := range left.Data {
		if i < len(right.Data) {
			if and {
				v &= right.Data[i]
			} else {
				v |= right.Data[i]
			}
		}
		ret.Data[i] = v
	}
	return ret
}

func checkIndex(set *core.Set, index int64) error {
	if index < 0 || index >= core.MaxSet {
		return fmt.Errorf(ErrorText(ErrIndexOut))
	}
	return nil
}

// IsSet returns the value of set[index]
func IsSet(set *core.Set, index int64) int64 {
	shift := int(index >> 6)
	pos := uint64(index % 64)
	if len(set.Data) <= shift || set.Data[shift]&(1<<pos) == 0 {
		return 0
	}
	return 1
}

// setºArr converts array of integers to set
func setºArr(arr *core.Array) (set *core.Set, err error) {
	var ind int64
	set = core.NewSet()
	for _, v := range arr.Data {
		ind = v.(int64)
		if err = checkIndex(set, ind); err == nil {
			set.Set(ind, true)
		}
	}
	return
}

// SetºSet sets the item in the set
func SetºSet(set *core.Set, index int64) (*core.Set, error) {
	var err error
	if err = checkIndex(set, index); err == nil {
		set.Set(index, true)
	}
	return set, err
}

// setºStr converts string to set
func setºStr(value string) (*core.Set, error) {
	s := core.NewSet()
	for i, ch := range value {
		switch ch {
		case '0':
		case '1':
			s.Set(int64(i), true)
		default:
			return nil, fmt.Errorf(ErrorText(ErrInvalidParam))
		}
	}
	return s, nil
}

// strºSet converts set to string
func strºSet(set *core.Set) string {
	return set.String()
}

// ToggleºSetInt changes the value of the set
func ToggleºSetInt(set *core.Set, index int64) (prev int64, err error) {
	if err = checkIndex(set, index); err == nil {
		prev = IsSet(set, index)
		set.Set(index, prev == 0)
	}
	return
}

// UnSetºSet sets the item in the set
func UnSetºSet(set *core.Set, index int64) (*core.Set, error) {
	var err error
	if err = checkIndex(set, index); err == nil {
		set.Set(index, false)
	}
	return set, err
}
