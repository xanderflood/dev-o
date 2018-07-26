package devo

import (
	"fmt"
	"io"
	"os"
	"time"
)

//Logger logger
//go:generate counterfeiter . Logger
type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Detail(...interface{})
	Detailf(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
}

//LogLevel enum
type LogLevel int

const (
	//LogLevelError LogLevelError
	LogLevelError LogLevel = iota

	//LogLevelInfo LogLevelInfo
	LogLevelInfo LogLevel = iota

	//LogLevelDetail LogLevelDetail
	LogLevelDetail LogLevel = iota

	//LogLevelDebug LogLevelDebug
	LogLevelDebug LogLevel = iota
)

//NewStandardLogger create a new logger that writes to stdout
func NewStandardLogger(level LogLevel) Logger {
	return NewPassLogger(level, os.Stdout)
}

//NewPassLogger create a new logger that writes to an arbitrary target
func NewPassLogger(level LogLevel, target io.Writer) Logger {
	return &LoggerImpl{
		LogLevel: level,
		Target:   target,
	}
}

//LoggerImpl standard logger implementation
type LoggerImpl struct {
	LogLevel LogLevel
	Target   io.Writer
}

func (li *LoggerImpl) print(vals ...interface{}) {
	line := fmt.Sprintln(vals...)
	line = fmt.Sprintf("[DEVO-%s] %s", time.Now().Format(time.RFC3339), line)
	li.Target.Write([]byte(line))
}

func (li *LoggerImpl) printf(format string, vals ...interface{}) {
	li.print(fmt.Sprintf(format, vals...))
}

//Error Error
func (li *LoggerImpl) Error(vals ...interface{}) {
	if li.LogLevel >= LogLevelError {
		li.print(vals...)
	}
}

//Errorf Errorf
func (li *LoggerImpl) Errorf(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelError {
		li.printf(format, vals...)
	}
}

//Info Info
func (li *LoggerImpl) Info(vals ...interface{}) {
	if li.LogLevel >= LogLevelInfo {
		li.print(vals...)
	}
}

//Infof Infof
func (li *LoggerImpl) Infof(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelInfo {
		li.printf(format, vals...)
	}
}

//Detail Detail
func (li *LoggerImpl) Detail(vals ...interface{}) {
	if li.LogLevel >= LogLevelDetail {
		li.print(vals...)
	}
}

//Detailf Detailf
func (li *LoggerImpl) Detailf(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelDetail {
		li.printf(format, vals...)
	}
}

//Debug Debug
func (li *LoggerImpl) Debug(vals ...interface{}) {
	if li.LogLevel >= LogLevelDebug {
		li.print(vals...)
	}
}

//Debugf Debugf
func (li *LoggerImpl) Debugf(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelDebug {
		li.printf(format, vals...)
	}
}
