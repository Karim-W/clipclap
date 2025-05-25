package main

import (
	"fmt"
)

var log Logger

type Logger struct {
	isVerbose bool
}

func (l *Logger) Log(msg string, args ...any) {
	if !l.isVerbose {
		return
	}

	args = append([]any{msg}, args...)

	fmt.Println(args...)
}
