package devo

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"

	perrors "github.com/pkg/errors"
)

//Binner build and execs the binary
//go:generate counterfeiter . Binner
type Binner interface {
	Build() error
	Exec()
}

//BinnerI standard binner
type BinnerI struct {
	config
	lock sync.Locker
}

//NewBinner new Binner
func NewBinner(c config) *BinnerI {
	return &BinnerI{
		config: c,
		lock:   &sync.Mutex{},
	}
}

//Build build the target package
func (b *BinnerI) Build() error {
	cmd := exec.Command("go", "build", "-o", b.binpath)
	cmd.Dir = b.target

	stdout, err := cmd.StderrPipe()
	if err != nil {
		return perrors.Wrap(err, fmt.Sprintf("failed to open build process stdout for readings"))
	}

	doneChan := make(chan struct{})
	var readError error
	var message string
	go func() {
		buf := bytes.NewBuffer([]byte{})
		_, readError = io.Copy(buf, stdout)
		message = buf.String()

		doneChan <- struct{}{}
	}()

	err = cmd.Start()
	if err != nil {
		return perrors.Wrap(err, "build not started")
	}

	err = cmd.Wait()
	if err != nil {
		<-doneChan
		if readError != nil {
			return perrors.Wrap(readError, fmt.Sprintf("build failed (could not read the build output)"))
		}

		return perrors.Wrap(err, message)
	}

	return nil
}

//Exec the binary
func (b *BinnerI) Exec() {
	b.lock.Lock()
	defer b.lock.Unlock()

	syscall.Exec(b.binpath, b.binargs, b.environ)
}
