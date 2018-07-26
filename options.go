package devo

import (
	"os"
	"path/filepath"
	"strings"
)

type Option struct {
	Apply      func(*config) error
	Precedence Precedence
	Default    bool
}

type Precedence uint

const (
	//changing the order of these declartions will change the option precedence
	PrecedenceEnv     Precedence = iota
	PrecedenceDetails Precedence = iota
	PrecedenceGoPath  Precedence = iota
	PrecedencePaths   Precedence = iota
	PrecedenceUtils   Precedence = iota
)

//WithGoPath interpret other paths relative to this
func WithGoPath(gopath string) Option {
	return Option{
		Apply: func(c *config) error {
			c.target = filepath.Join(c.gopath, gopath)
			return nil
		},
		Precedence: PrecedenceGoPath,
	}
}

//WithTarget target the executable produced by a particular module
func WithTarget(target string) Option {
	return Option{
		Apply: func(c *config) error {
			c.target = filepath.Join(c.gopath, target)
			return nil
		},
		Precedence: PrecedencePaths,
	}
}

//WhileMonitoring monitor *.go files in these directories for changes
func WhileMonitoring(subjects ...string) Option {
	return Option{
		Apply: func(c *config) error {
			c.subjects = make([]string, len(subjects))
			for i, subject := range subjects {
				c.subjects[i] = filepath.Join(c.gopath, subject)
			}
			return nil
		},
		Precedence: PrecedencePaths,
	}
}

//WithLogger with a specified logger
func WithLogger(logger Logger) Option {
	return Option{
		Apply: func(c *config) error {
			c.logger = logger
			return nil
		},
		Precedence: PrecedenceUtils,
	}
}

//WithWaitMS with a specified wait time between scans
func WithWaitMS(waitMS uint64) Option {
	return Option{
		Apply: func(c *config) error {
			c.waitMS = waitMS
			return nil
		},
		Precedence: PrecedenceDetails,
	}
}

func defaults() []Option {
	return []Option{
		{
			Apply: func(c *config) error {
				c.gopath = filepath.Join(strings.TrimSpace(os.Getenv("GOPATH")), "src")
				return nil
			},
			Precedence: PrecedenceGoPath,
			Default:    true,
		},
		{
			Apply: func(c *config) error {
				c.binargs = os.Args
				c.environ = os.Environ()
				return nil
			},
			Precedence: PrecedenceEnv,
			Default:    true,
		},
		{
			Apply: func(c *config) error {
				c.logger = NewStandardLogger(LogLevelInfo)
				return nil
			},
			Precedence: PrecedenceUtils,
			Default:    true,
		},
		{
			Apply: func(c *config) error {
				c.waitMS = 500
				return nil
			},
			Precedence: PrecedenceDetails,
			Default:    true,
		},
	}
}
