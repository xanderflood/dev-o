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
	watcher := NewGoFileWatcher(c)
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

//AutoreloadDaemon AutoreloadDaemon
func AutoreloadDaemon(
	binner Binner,
	watcher Watcher,
	logger Logger,
	waitMS uint64,
) {
	var lastUpdated, lastBuilt time.Time

	//TODO better way to set this initally?
	//leaving it at zero means we immediately rebuild :/
	//could use an envar to track the build time even
	//across restarts
	lastBuilt = time.Now()
	for {
		//wait for the interval
		<-time.NewTimer(time.Duration(waitMS) * time.Millisecond).C

		lastUpdated = watcher.LastUpdated()
		if lastBuilt.After(lastUpdated) {
			continue
		}

		logger.Info("Rebuilding")
		lastBuilt = time.Now()
		err := binner.Build()
		if err != nil {
			logger.Infof("--- failed:\n%s", err.Error())
			continue
		}

		logger.Info("Restarting")
		binner.Exec()

		logger.Error("--- didn't restart")
	}
}
