// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package conlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 36
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

// LogLevelFormat is used to set how the log level is displayed in an
// output message.
type LogLevelFormat uint32

const (
	// LogLevelFormatUnknown is used for defensive
	// programming. You should never see this.
	LogLevelFormatUnknown = iota

	// LogLevelFormatNone disables output of the log level.
	LogLevelFormatNone

	// LogLevelFormatShort outputs the abbreviated form of the log
	// level with a trailing space, e.g., "DEBU ".
	LogLevelFormatShort

	// LogLevelFormatLongTitle outputs the full, capitalized, log
	// level with a trailing space, e.g., "Debug ".
	LogLevelFormatLongTitle

	// LogLevelFormatLongLower outputs the full, lower-case,
	// version of the log level with a trailing space, e.g.,
	// "debug ".
	LogLevelFormatLongLower
)

// TimestampType is used to set what type of timestamp is displayed in
// output messages.
type TimestampType uint32

const (
	// TimestampTypeUnknown is used for defensive programming. You
	// should never see this.
	TimestampTypeUnknown = iota

	// TimestampTypeNone disables outputting the timestamp.
	TimestampTypeNone

	// TimestampTypeWall outputs the current wall clock time using
	// the WallclockTimestampFmt format.
	TimestampTypeWall

	// TimestampTypeElapsed outputs the elapsed time in seconds
	// since the start of execution using ElapsedTimestampFmt.
	TimestampTypeElapsed
)

// FormattingOptions are options that control output format.
type FormattingOptions struct {
	// LogLevelFmt is the format used to display the log
	// level. Defaults to LogLevelFormatNone.
	LogLevelFmt LogLevelFormat

	// ShowLogLevelColors controls showing colorized log levels if
	// output is to a TTY. Defaults to false.
	ShowLogLevelColors bool

	// TimestampTypeOutput controls the type of timestamp
	// used. The default is TimestampTypeNone.
	TimestampType TimestampType

	// WallclockTimestampFmt is the time.Format() format used when
	// displaying wall clock timestamps. Defaults to time.RFC3339.
	WallclockTimestampFmt string

	// ElapsedTimestampFmt is the format string used to display
	// elapsed time timestamps. Defaults to "%04d".
	ElapsedTimestampFmt string
}

// NewFormattingOptions is the constructor for Formatting options.
func NewFormattingOptions() *FormattingOptions {
	return &FormattingOptions{
		LogLevelFmt:           LogLevelFormatNone,
		ShowLogLevelColors:    false,
		TimestampType:         TimestampTypeNone,
		WallclockTimestampFmt: DefaultWallclockTimestampFormat,
		ElapsedTimestampFmt:   DefaultElapsedTimestampFormat,
	}
}

// StdFormatter formats logs into text.
type StdFormatter struct {

	// Formatting options used to modify the output.
	Options *FormattingOptions

	// Whether the logger's out is to a terminal.
	isTerminal bool

	sync.Once
}

// NewStdFormatter is the StdFormatter constructor.
func NewStdFormatter() *StdFormatter {
	return &StdFormatter{
		Options: NewFormattingOptions(),
	}
}

func (f *StdFormatter) init(entry *Entry) {
	f.isTerminal = f.checkIfTerminal(entry.Log.out)
}

func (f *StdFormatter) checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

// Format renders a single log entry.
func (f *StdFormatter) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	if entry.Level == printLevel {
		_, err := fmt.Fprint(b, entry.Message)
		if err != nil {
			return []byte{}, nil
		}
		return b.Bytes(), nil
	}

	f.printLeader(b, entry)
	f.printMessage(b, entry)

	return b.Bytes(), nil
}

func (f *StdFormatter) printLeader(w io.Writer, entry *Entry) (n int) {
	var leader string

	switch f.Options.LogLevelFmt {
	case LogLevelFormatShort:
		leader += strings.ToUpper(entry.Level.String())[0:4]
		if f.Options.ShowLogLevelColors && f.isTerminal {
			leader = f.colorOn(entry.Level) + leader + f.colorOff(entry.Level)
		}
	case LogLevelFormatLongTitle:
		leader += strings.Title(entry.Level.String())
		if f.Options.ShowLogLevelColors && f.isTerminal {
			leader = f.colorOn(entry.Level) + leader + f.colorOff(entry.Level)
		}
	case LogLevelFormatLongLower:
		leader += strings.ToLower(entry.Level.String())
		if f.Options.ShowLogLevelColors && f.isTerminal {
			leader = f.colorOn(entry.Level) + leader + f.colorOff(entry.Level)
		}
	default:
	}

	switch f.Options.TimestampType {
	case TimestampTypeWall:
		time := entry.Time
		leader += fmt.Sprintf(
			"[%s]",
			time.Format(f.Options.WallclockTimestampFmt))
	case TimestampTypeElapsed:
		ticks := int(entry.Time.Sub(baseTimestamp) / time.Second)
		leader += fmt.Sprintf(
			"["+f.Options.ElapsedTimestampFmt+"]",
			ticks)
	default:
	}

	if len(leader) == 0 {
		return 0
	}
	n, _ = fmt.Fprintf(w, "%s ", leader)
	return n
}

func (f *StdFormatter) printMessage(w io.Writer, entry *Entry) {
	_, _ = fmt.Fprintf(w, "%s", entry.Message)
}

func (f *StdFormatter) colorOn(level Level) (on string) {
	var levelColor int
	switch level {
	case DebugLevel:
		levelColor = blue
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = green
	}

	on = fmt.Sprintf("\x1b[%dm", levelColor)

	return on
}

func (f *StdFormatter) colorOff(level Level) (off string) {
	return "\x1b[0m"
}
