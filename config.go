package devo

import (
	"path/filepath"
	"sort"
)

type config struct {
	gopath   string
	target   string
	binname  string
	binpath  string
	binargs  []string
	subjects []string
	environ  []string
	waitMS   uint64
	logger   Logger
}

func configure(options ...Option) (config, error) {
	options = append(options, defaults()...)

	sort.Slice(options, func(i, j int) bool {
		if options[i].Precedence == options[j].Precedence {
			return options[i].Default && !options[j].Default
		}

		return options[i].Precedence < options[j].Precedence
	})

	var c config
	for _, option := range options {
		err := option.Apply(&c)
		if err != nil {
			return config{}, err
		}
	}

	if c.binname == "" {
		c.binname = filepath.Base(c.target)
	}

	if c.binpath == "" {
		c.binpath = filepath.Join(c.target, c.binname)
	}

	return c, nil
}
