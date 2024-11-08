package main

import (
	"fmt"
	"runtime"
)

var (
	Compiler  = runtime.Compiler
	GoVersion = runtime.Version()
	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)
