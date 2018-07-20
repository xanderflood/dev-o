package main

import (
	"fmt"
	"time"

	"github.com/xanderflood/dev-o"
)

const printVal = 9

const waitSeconds = 1

//TODO figure out why `go` isn't in $PATH when using exec.Command
const goExecutable = "/usr/lib/go-1.10/bin/go"

func main() {
	devo.Autoreload(
		devo.WithTarget("github.com/xanderflood/dev-o/cmd"),
		devo.WhileMonitoring("github.com/xanderflood/dev-o"),
	)

	for {
		<-time.NewTimer(waitSeconds * time.Second).C

		fmt.Println(printVal)
	}
}
