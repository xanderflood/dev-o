package devo

import (
	"fmt"
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
	LoadState() (SubjectState, error)
}

//GoFileWatcher watch go files in subdirectories of the subject directories
type GoFileWatcher struct {
	subjects []string
}

//NewGoFileWatcher new GoFileWatcher
func NewGoFileWatcher(subjects []string) *GoFileWatcher {
	return &GoFileWatcher{
		subjects: subjects,
	}
}

//LoadState get the states of the files
func (w *GoFileWatcher) LoadState() (SubjectState, error) {
	result := map[string]time.Time{}

	for _, subject := range w.subjects {
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
