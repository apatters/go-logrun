// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package conlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	stdoutNames = []string{
		"/dev/stdout",
		"|1",
	}
	bufferPool *sync.Pool
)

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

// write wraps output to the io.Writer. The go test command requires
// output go directly to os.Stdout.
func write(w io.Writer, s []byte) (n int, err error) {
	switch v := w.(type) {
	case *os.File:
		for _, name := range stdoutNames {
			if v.Name() == name {
				return fmt.Printf("%s", string(s)) // nolint: megacheck
			}
		}
	default:
	}

	return w.Write(s)
}

// Entry is the final or intermediate logging entry. It's finally
// logged when Debug, Info, Warn, Error, Fatal or Panic is called on
// it.
type Entry struct {
	Log *Logger

	// Time at which the log entry was created.
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn,
	// Error, Fatal or Panic This field will be set on entry
	// firing and the value will be equal to the one in Logger
	// struct field.
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic.
	Message string

	// When formatter is called in entry.log(), a Buffer may be
	// set to entry.
	Buffer *bytes.Buffer
}

// NewEntry is the Entry constructor.
func NewEntry(log *Logger) *Entry {
	return &Entry{
		Log: log,
	}
}

// String returns the string representation from the reader and
// ultimately the formatter.
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Log.formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// log outputs the message to the Writer after formatting it. It is
// not declared with a pointer value because otherwise race conditions
// will occur when using multiple goroutines.
func (entry Entry) log(level Level, w io.Writer, msg string) {
	var buffer *bytes.Buffer
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer
	serialized, err := entry.Log.formatter.Format(&entry)
	entry.Buffer = nil
	if err != nil {
		entry.Log.mu.Lock()
		_, _ = fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Log.mu.Unlock()
	} else {
		entry.Log.mu.Lock()
		_, err = write(w, serialized)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
		entry.Log.mu.Unlock()
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel {
		panic(&entry)
	}
}

// Print writes a message ala fmt.Print if printing is enabled,
// otherwise it is discarded.
func (entry *Entry) Print(args ...interface{}) {
	if entry.Log.GetPrintEnabled() {
		entry.log(printLevel, entry.Log.out, fmt.Sprint(args...))
	}
}

// Debug writes a message ala fmt.Print.
func (entry *Entry) Debug(args ...interface{}) {
	if entry.Log.GetLevel() >= DebugLevel {
		args = append(args, "\n")
		entry.log(DebugLevel, entry.Log.out, fmt.Sprint(args...))
	}
}

// Info writes a message ala fmt.Print.
func (entry *Entry) Info(args ...interface{}) {
	if entry.Log.GetLevel() >= InfoLevel {
		args = append(args, "\n")
		entry.log(InfoLevel, entry.Log.out, fmt.Sprint(args...))
	}
}

// Warn writes a message ala fmt.Print.
func (entry *Entry) Warn(args ...interface{}) {
	if entry.Log.GetLevel() >= WarnLevel {
		args = append(args, "\n")
		entry.log(WarnLevel, entry.Log.out, fmt.Sprint(args...))
	}
}

// Warning writes a message ala fmt.Print. It is an alias for Warn.
func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

// Error writes a message ala fmt.Print.
func (entry *Entry) Error(args ...interface{}) {
	if entry.Log.GetLevel() >= ErrorLevel {
		args = append(args, "\n")
		entry.log(ErrorLevel, entry.Log.errOut, fmt.Sprint(args...))
	}
}

// Fatal writes a message ala fmt.Print.
func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Log.GetLevel() >= FatalLevel {
		args = append(args, "\n")
		entry.log(FatalLevel, entry.Log.errOut, fmt.Sprint(args...))
	}
}

// Panic writes a message ala fmt.Print and then calls panic.
func (entry *Entry) Panic(args ...interface{}) {
	if entry.Log.GetLevel() >= PanicLevel {
		args = append(args, "\n")
		entry.log(PanicLevel, entry.Log.errOut, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

// Printf writes a message ala fmt.Printf if printing is enabled,
// otherwise it is discarded. No newline is appended.
func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.log(printLevel, entry.Log.out, fmt.Sprintf(format, args...))
}

// Debugf writes a message ala fmt.Printf.
func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Log.GetLevel() >= DebugLevel {
		format += "\n"
		entry.log(DebugLevel, entry.Log.out, fmt.Sprintf(format, args...))
	}
}

// Infof writes a message ala fmt.Printf.
func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Log.GetLevel() >= InfoLevel {
		format += "\n"
		entry.log(InfoLevel, entry.Log.out, fmt.Sprintf(format, args...))
	}
}

// Warnf writes a message ala fmt.Printf.
func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Log.GetLevel() >= WarnLevel {
		format += "\n"
		entry.log(WarnLevel, entry.Log.out, fmt.Sprintf(format, args...))
	}
}

// Warningf writes a message ala fmt.Printf. It is an alias for Warnf.
func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

// Errorf writes a message ala fmt.Printf.
func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Log.GetLevel() >= ErrorLevel {
		format += "\n"
		entry.log(ErrorLevel, entry.Log.errOut, fmt.Sprintf(format, args...))
	}
}

// Fatalf writes a message ala fmt.Printf.
func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Log.GetLevel() >= FatalLevel {
		format += "\n"
		entry.log(FatalLevel, entry.Log.errOut, fmt.Sprintf(format, args...))
	}
}

// Panicf writes a message ala fmt.Print and then calls panicf.
func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Log.GetLevel() >= PanicLevel {
		format += "\n"
		entry.log(PanicLevel, entry.Log.errOut, fmt.Sprintf(format, args...))
	}
}

// Entry Println family functions

// Println writes a message ala fmt.Println if printing is enabled,
// otherwise it is discarded.
func (entry *Entry) Println(args ...interface{}) {
	if entry.Log.GetPrintEnabled() {
		msg := fmt.Sprintln(args...)
		entry.Print(msg)
	}
}

// Debugln writes a message ala fmt.Println.
func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Log.GetLevel() >= DebugLevel {
		entry.log(DebugLevel, entry.Log.out, fmt.Sprintln(args...))
	}
}

// Infoln writes a message ala fmt.Println.
func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Log.GetLevel() >= InfoLevel {
		entry.log(InfoLevel, entry.Log.out, fmt.Sprintln(args...))
	}
}

// Warnln writes a message ala fmt.Println.
func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Log.GetLevel() >= WarnLevel {
		entry.log(WarnLevel, entry.Log.out, fmt.Sprintln(args...))
	}
}

// Warningln writes a message ala fmt.Println. It is an alias for
// Warnln.
func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

// Errorln writes a message ala fmt.Println.
func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Log.GetLevel() >= ErrorLevel {
		entry.log(ErrorLevel, entry.Log.errOut, fmt.Sprintln(args...))
	}
}

// Fatalln writes a message ala fmt.Println.
func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Log.GetLevel() >= FatalLevel {
		entry.log(FatalLevel, entry.Log.errOut, fmt.Sprintln(args...))
	}
}

// Panicln writes a message ala fmt.Println and then calls panic.
func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Log.GetLevel() >= PanicLevel {
		entry.log(PanicLevel, entry.Log.errOut, fmt.Sprintln(args...))
	}
}
