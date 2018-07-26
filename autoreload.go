package devo

import (
	"sync"
	"time"
)

//Autoreload Asynchronously replace this process with updated versions
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
		c.logger,
		c.waitMS,
	)

	// start //
	return binner.lock, nil
}

//TODO much of this is a little inconvenient to test, but
//AutoreloadDaemon can and should be tested be pretty thoroughly

//TODO refactor the FSM into something that looks more fun

//AutoreloadDaemon AutoreloadDaemon
func AutoreloadDaemon(
	binner Binner,
	watcher Watcher,
	logger Logger,
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
			logger.Info("Restarting")
			binner.Exec()

			logger.Error("--- didn't restart")
			built = false
		} else if updated {
			logger.Info("Rebuilding")
			err := binner.Build()
			if err != nil {
				logger.Info("--- failed", err.Error())

				//TODO: this parrots the same build error once a second
				//instead, maintain a lastBuiltMoment and a lastModifiedMoment
				//and only rebuild when they're out of whack

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
