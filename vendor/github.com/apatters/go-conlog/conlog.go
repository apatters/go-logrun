// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

/*
Package conlog provides console logging.

The conlog package is a variant of the popular
https://github.com/Sirupsen/logrus package. It is optimized for
console output and is meant to be used by command-line utilities with
output seen interactively by the user rather than sent to a log file
(although it can be used for that too).

The conlog package can be used as a drop-in replacement for the
standard golang logger.
*/
package conlog

import (
	"fmt"
	"io"
	"strings"
)

// DefaultExitCode is used in the Fatal*()-style functions.
const DefaultExitCode = 1

// Level wraps the enumerated constants used to set and report logging
// levels.
type Level uint32

const (
	// PanicLevel level. The highest level of severity, e.g.,
	// always output. Logs and then calls panic with the message.
	PanicLevel = iota

	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will
	// exit even if the logging level is set to Panic.
	FatalLevel

	// ErrorLevel level. Used for errors that should definitely be
	// noted.
	ErrorLevel

	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel

	// InfoLevel level. General operational entries about what is
	// going on inside the application.
	InfoLevel

	// DebugLevel level. Usually only enabled when debugging. Very
	// verbose logging.
	DebugLevel

	// A pseudo-level. Print-level output is controlled by the
	// PrintEnabled flag.
	printLevel
)

// String converts the Level to a string. E.g. PanicLevel becomes
// "panic".
func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}
	return "unknown"
}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("Not a valid log Level: %q", lvl)
}

// AllLevels is a constant exposing all usable logging levels.
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
}

// The StdLogger interface is compatible with the standard library log package.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})

	SetOutput(io.Writer)
}

// The ConLogger interface extends the standard interface by adding
// level output functions and some utiliities.
type ConLogger interface {
	SetLevel(level Level)
	GetLevel() Level
	SetPrintEnabled(enabled bool)
	GetPrintEnabled() bool

	Printf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Print(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Println(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	FatalWithExitCode(code int, args ...interface{})
	FatalfWithExitCode(code int, format string, args ...interface{})
	FatallnWithExitCode(code int, args ...interface{})
	FatalIfError(err error, code int, args ...interface{})
	FatalfIfError(err error, code int, format string, args ...interface{})
	FatallnIfError(err error, code int, args ...interface{})
}
