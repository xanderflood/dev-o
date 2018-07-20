package devo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

//TODO precedence-sorted options
// type Option interface {
// 	Apply(*config) error
// 	Precedence() uint
// }

//Option dev-o config functors
type Option func(*config) error

//WithTarget target the executable produced by a particular module
func WithTarget(target string) Option {
	return func(c *config) error {
		c.target = filepath.Join(c.gopath, target)
		return nil
	}
}

//WhileMonitoring *.go files in these directories for changes
func WhileMonitoring(subjects ...string) Option {
	return func(c *config) error {
		c.subjects = make([]string, len(subjects))
		for i, subject := range subjects {
			c.subjects[i] = filepath.Join(c.gopath, subject)
		}
		return nil
	}
}

type config struct {
	gopath   string
	target   string
	binname  string
	binargs  []string
	subjects []string
	environ  []string
	waitMS   uint64
}

//Autoreload asynchronously replace this process with updated versions
//this process is compiled from the go file
func Autoreload(options ...Option) (sync.Locker, error) {
	var c config

	// defaults //
	c.gopath = filepath.Join(strings.TrimSpace(os.Getenv("GOPATH")), "src")
	if len(c.gopath) == 0 {
		return nil, errors.New("GOPATH is not set")
	}

	c.waitMS = 500

	c.binname = os.Args[0]
	if c.binname[0] != '/' {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		c.binname = filepath.Join(wd, c.binname)
	}

	c.binargs = os.Args[1:]

	c.environ = os.Environ()

	// overrides //
	for _, option := range options {
		err := option(&c)
		if err != nil {
			return nil, err
		}
	}

	// start //
	err := autoreload(c)
	if err != nil {
		return nil, err
	}

	return &sync.Mutex{}, nil
}

type watcherMessageKind int

const (
	watcherMessageUpdated        watcherMessageKind = iota
	watcherMessageBuildFailed                       = iota
	watcherMessageBuildSucceeded                    = iota
	watcherMessageError                             = iota
)

type watcherMessage struct {
	kind    watcherMessageKind
	subject string
}

type builderState int

const (
	builderStateGood watcherMessageKind = iota
	builderStateBad                     = iota
)

type builderMessage struct {
	state   builderState
	subject string
}

func autoreload(c config) error {
	//first load them all
	oldState, err := loadState(c.subjects)
	if err != nil {
		// TODO handle?
		return err
	} //start the watcher

	go func() {
		//now loop and watch the subject
		var updated, built, shouldWait bool
		for {
			if shouldWait {
				//wait for the interval
				<-time.NewTimer(time.Duration(c.waitMS) * time.Millisecond).C
			} else {
				shouldWait = true
			}

			if built {
				fmt.Println("resetting")

				//TODO: best way to build this absolute path?
				syscall.Exec("/home/xander/go/src/github.com/xanderflood/dev-o/cmd/cmd", nil, os.Environ())
			} else if updated {
				fmt.Println("building")
				ok := rebuild(c.target)
				if ok {
					built = true
					shouldWait = false
				}
			} else {
				fmt.Println("checking")
				newState, err := loadState(c.subjects)
				if err != nil {
					// TODO handle?
					fmt.Printf("failed loading initial state: %v\n", err)
					continue
				}

				for k, v := range newState {
					//if we failed to stat this file
					//TODO this won't correctly handle deletes

					if (v == time.Time{}) {
						newState[k] = oldState[k]
						continue
					}

					if v.After(oldState[k]) {
						updated = true
						shouldWait = true

						break
					}
				}

				oldState = newState
			}
		}
	}()

	return nil
}

func rebuild(target string) bool {
	//TODO shell escape the path

	cmd := exec.Command("go", "build")
	cmd.Dir = target
	err := cmd.Start()
	if err != nil {
		//TODO log?
		fmt.Printf("build not started %v\n", err)
		return false
	}
	err = cmd.Wait()
	if err != nil {
		//TODO log?
		fmt.Printf("build failed %v\n", err)
		return false
	}

	fmt.Printf("successfully built %s\n", time.Now().Format(time.RFC3339))

	//TODO rebuild the target
	return true
}

func loadState(subjects []string) (map[string]time.Time, error) {
	result := map[string]time.Time{}

	for _, subject := range subjects {
		err := filepath.Walk(subject, func(path string, file os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("failed statting file %s: %v\n", path, err)
				return nil
			}

			if file.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			//TODO make sure the executable is ALWAYS excluded, even when it happens to end in .go or something weird like that
			//first, learn how t.f. `go build` decides where to put the executable and what to name it. can I control that?

			result[path] = file.ModTime()
			return nil
		})

		if err != nil {
			fmt.Printf("failed walking subject %s: %v\n", subject, err)
		}
	}

	return result, nil
}
