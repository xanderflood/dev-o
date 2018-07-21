package devo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

//TODO add a logger that can be
//configured to be a named pipe
//insead of stdout

func defaultConfig() (config, error) {
	var c config

	c.gopath = filepath.Join(strings.TrimSpace(os.Getenv("GOPATH")), "src")
	if len(c.gopath) == 0 {
		return config{}, errors.New("GOPATH is not set")
	}

	c.waitMS = 500

	c.binname = os.Args[0]
	if c.binname[0] != '/' {
		wd, err := os.Getwd()
		if err != nil {
			return config{}, err
		}

		c.binname = filepath.Join(wd, c.binname)
	}

	c.binargs = os.Args[1:]

	c.environ = os.Environ()

	return c, nil
}

func configure(options ...Option) (config, error) {
	c, err := defaultConfig()
	if err != nil {
		return config{}, err
	}

	for _, option := range options {
		err := option(&c)
		if err != nil {
			return config{}, err
		}
	}

	return c, nil
}
