package main

import (
	"runtime"
	"statefulset/cmds"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmds.Execute()
}
