package devo

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

//Watcher provides updates on the state
//go:generate counterfeiter . Watcher
type Watcher interface {
	LastUpdated() time.Time
}

//GoFileWatcher watch go files in subdirectories of the subject directories
type GoFileWatcher struct {
	config
}

//NewGoFileWatcher new GoFileWatcher
func NewGoFileWatcher(c config) *GoFileWatcher {
	g := &GoFileWatcher{
		config: c,
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
				w.logger.Errorf("failed checking file %s: %s", path, err)
				return nil
			}

			if file.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			if filepath.Clean(path) == filepath.Clean(w.binpath) {
				return nil
			}

			touched := file.ModTime()

			if touched.After(lastUpdated) {
				lastUpdated = touched
			}
			return nil
		})

		if err != nil {
			w.logger.Error("failed walking files:", err)
		}
	}

	return lastUpdated
}
