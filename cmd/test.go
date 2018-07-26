package main

import (
	"sync"
	"time"

	"github.com/xanderflood/dev-o"
)

var printVals = []int{4, 1000, 6}

const waitSeconds = 1

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
