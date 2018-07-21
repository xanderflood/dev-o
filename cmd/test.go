package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/xanderflood/dev-o"
)

var printVals = []int{4, 5, 6}

const waitSeconds = 1

//TODO figure out why `go` isn't in $PATH when using exec.Command
const goExecutable = "/usr/lib/go-1.10/bin/go"

func main() {
	lock, err := devo.Autoreload(
		devo.WithTarget("github.com/xanderflood/dev-o/cmd"),
		devo.WhileMonitoring("github.com/xanderflood/dev-o"),
	)
	if err != nil {
		panic(err)
	}

	for {
		printWithWait(lock)
	}
}

func printWithWait(lock sync.Locker) {
	lock.Lock()
	defer lock.Unlock()

	for i := 0; i < len(printVals); i++ {
		<-time.NewTimer(waitSeconds * time.Second).C

		fmt.Println(printVals[i])
	}
}
