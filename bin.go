package devo

import (
	"os/exec"
	"syscall"

	perrors "github.com/pkg/errors"
)

//Binner build and execs the binary
//go:generate counterfeiter . Binner
type Binner interface {
	Build() error
	Exec() error
}

//BinnerI standard binner
type BinnerI struct {
	config
}

//NewBinner new Binner
func NewBinner(c config) *BinnerI {
	return &BinnerI{c}
}

//Build build the target package
func (b *BinnerI) Build() error {
	cmd := exec.Command("go", "build")
	cmd.Dir = b.target

	err := cmd.Start()
	if err != nil {
		return perrors.Wrap(err, "build not started")
	}

	err = cmd.Wait()
	if err != nil {
		return perrors.Wrap(err, "build failed")
	}

	return nil
}

//Exec the binary
func (b *BinnerI) Exec() {
	syscall.Exec(b.binname, b.binargs, b.environ)
}
