package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gentee/gentee/vm"
)

const (
	errNoFile = iota + 1
	errCompile
	errRun
	errResult
	errPlaceholder5
	errPlaceholder6
	errPlaceholder7
	errUndefined
)

type CodedError struct {
	Code int
	Err  error
}

func (c *CodedError) Error() string {
	if c.Err != nil {
		return c.Error()
	}
	return fmt.Sprintf("Unspecified Error with Code: %d", c.Code)
}

func (c *CodedError) Unwrap() error {
	return c.Err
}

func codedError(err error, code int) error {
	if err != nil {
		fmt.Print(`ERROR`)
		if errTrace, ok := err.(*vm.RuntimeError); ok {
			fmt.Printf(" #%d: %s\n", errTrace.ID, err.Error())
			for _, trace := range errTrace.Trace {
				path := trace.Path
				dirs := strings.Split(filepath.ToSlash(path), `/`)
				if len(dirs) > 3 {
					path = `...` + path[len(path)-len(strings.Join(dirs[len(dirs)-3:], `/`))-1:]
				}
				fmt.Printf("%s [%d:%d] %s -> %s\n", path, trace.Line, trace.Pos, trace.Entry, trace.Func)
			}
			code = errTrace.ID
		} else {
			fmt.Println(`:`, err.Error())
		}
		return &CodedError{Err: err, Code: code}
	}
	return nil
}
