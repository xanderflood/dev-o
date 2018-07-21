package devo

import (
	"sync"
	"time"
)

//Autoreload asynchronously replace this process with updated versions
//this process is compiled from the go file
func Autoreload(options ...Option) (sync.Locker, error) {
	c, err := configure(options...)
	if err != nil {
		return nil, err
	}

	binner := NewBinner(c)
	watcher := NewGoFileWatcher(c.subjects)
	go AutoreloadDaemon(
		binner,
		watcher,
		c.waitMS,
	)

	// start //
	return binner.lock, nil
}

//TODO much of this is a little inconvenient to test, but
//AutoreloadDaemon can and should be tested be pretty thoroughly

//AutoreloadDaemon AutoreloadDaemon
func AutoreloadDaemon(
	binner Binner,
	watcher Watcher,
	waitMS uint64,
) {
	var updated, built, stateChanged bool
	for {
		if stateChanged {
			//wait for the interval
			<-time.NewTimer(time.Duration(waitMS) * time.Millisecond).C
		} else {
			stateChanged = true
		}

		if built {
			binner.Exec()

			//it didn't exec, try rebuilding
			//send to the error log
			built = false
		} else if updated {
			err := binner.Build()
			if err != nil {
				//TOOD log
				continue
			}

			built = true
			stateChanged = false
		} else {
			if watcher.WasUpdated() {
				updated = true
				stateChanged = true
			}
		}
	}
}
