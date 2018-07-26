package devo

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

//SubjectState state of the subjects
type SubjectState map[string]time.Time

//Watcher provides updates on the state
//go:generate counterfeiter . Watcher
type Watcher interface {
	LastUpdated() time.Time
}

//GoFileWatcher watch go files in subdirectories of the subject directories
type GoFileWatcher struct {
	subjects []string
	state    SubjectState
}

//NewGoFileWatcher new GoFileWatcher
func NewGoFileWatcher(subjects []string) *GoFileWatcher {
	g := &GoFileWatcher{
		subjects: subjects,
		state:    SubjectState{},
	}

	return g
}

//LastUpdated attempts to stat all go files and returns true
//if it sees that any of them were updated. Also updates the
//underlying state of this object
func (w *GoFileWatcher) LastUpdated() time.Time {
	var lastUpdated time.Time
	for _, subject := range w.subjects {
		err := filepath.Walk(subject, func(path string, file os.FileInfo, err error) error {
			if err != nil {
				//TODO log
				return nil
			}

			if file.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			//TODO make sure the executable is ALWAYS excluded, even when it happens to end in .go or something weird like that
			//first, learn how t.f. `go build` decides where to put the executable and what to name it. can I control that?

			touched := file.ModTime()

			if touched.After(lastUpdated) {
				lastUpdated = touched
			}
			return nil
		})

		if err != nil {
			//TODO log
		}
	}

	return lastUpdated
}
