// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package conlog

import (
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/tevino/abool"
)

// Logger encapsulates basic logging.
type Logger struct {
	// Log messages with level higher than ErrorLevel are
	// `io.Copy`'d to this Writer in a mutex. It is common to set
	// out to a file, or leave it at the default which is
	// `os.Stdout`.
	out io.Writer

	// Log messages with level ErrorLevel and lower are
	// `io.Copy`'d to this in the same mutex as out. It is common
	// to set this to the same as out or a file, or leave it as
	// the default which is `os.Stderr`. Like out, you can use
	// just about any Writer.
	errOut io.Writer

	// All log entries pass through the formatter before being
	// logged to out.
	formatter Formatter

	// If true, suppresses output of the Print*() family of
	// logging functions.
	printEnabled *abool.AtomicBool

	// The logging level the logger should log at. This is
	// typically conlog.InfoLevel, which allows Info(), Warn(),
	// Error() and Fatal() to be logged. The default is InfoLevel
	level Level

	// Used to sync writing to the log. Locking is enabled by
	// Default.
	mu MutexWrap

	// Reusable empty entry
	entryPool sync.Pool
}

// MutexWrap is used to serialize logging output amongst goroutines.
type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

// Lock creates a lock on the logging mutex.
func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

// Unlock removes the lock on the logging mutex. It should generally
// be called in a "defer" statement.
func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

// Disable disables the use of locks when logging which can result in
// a performance increase at the expense of the loss of serialization.
func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

// Exit is used to wrap the exit code when making calls to
// Log.Fatal*() function which uses an internal panic mechanism to
// return the exit code.
type Exit struct {
	Code int
}

// HandleExit is used to call deferred functions before exiting. It
// uses the panic/recover mechanism to defer exiting. This routing
// should be used in your main routine like so:
//
// func main() {
//    defer handleExit()
//    // ready to go
// }
//
// See
// https://stackoverflow.com/questions/27629380/how-to-exit-a-go-program-honoring-deferred-calls
// for details.
func HandleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

// NewLogger is the constructor for Logger.
func NewLogger() *Logger {
	log := &Logger{
		out:          os.Stdout,
		errOut:       os.Stderr,
		formatter:    NewStdFormatter(),
		level:        InfoLevel,
		printEnabled: abool.New(),
	}
	log.printEnabled.Set()

	return log
}

// SetOutput sets the writer used for Print, Debug, and Warning
// messages.
func (log *Logger) SetOutput(w io.Writer) {
	log.mu.Lock()
	log.out = w
	log.mu.Unlock()
}

// SetErrorOutput sets the writer used for Error, Fatal, and Panic
// messages.
func (log *Logger) SetErrorOutput(w io.Writer) {
	log.mu.Lock()
	log.errOut = w
	log.mu.Unlock()
}

// SetLevel sets the logger level.
func (log *Logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&log.level), uint32(level))
}

// GetLevel returns the current logging level.
func (log *Logger) GetLevel() Level {
	return Level(atomic.LoadUint32((*uint32)(&log.level)))
}

// SetPrintEnabled sets the PrintEnabled setting.
func (log *Logger) SetPrintEnabled(enabled bool) {
	log.printEnabled.SetTo(enabled)
}

// GetPrintEnabled returns the PrintEnabled setting.
func (log *Logger) GetPrintEnabled() bool {
	return log.printEnabled.IsSet()
}

// SetFormatter sets the formatter used when printing entries.
func (log *Logger) SetFormatter(formatter Formatter) {
	log.mu.Lock()
	log.formatter = formatter
	log.mu.Unlock()
}

// SetNoLock disables the use of locking. It can be used when the log
// files are opened with appending mode, It is then safe to write
// concurrently to a file (within 4k message on Linux).
func (log *Logger) SetNoLock() {
	log.mu.Disable()
}

// Lock temporarily blocks output.
func (log *Logger) Lock() {
	log.mu.Lock()
}

// Unlock re-enables output that has been blocked with Lock().
func (log *Logger) Unlock() {
	log.mu.Unlock()
}

// newEntry is the constructor for an Entry.
func (log *Logger) newEntry() *Entry {
	entry, ok := log.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(log)
}

// releaseEntry returns an entry to the pool.
func (log *Logger) releaseEntry(entry *Entry) {
	log.entryPool.Put(entry)
}

// Print prints a message to the logger. It ignores logging levels. No
// logging levels, or timestamps are added. No newline is added. The
// equivalent of fmt.Fprint(out, ...).
func (log *Logger) Print(args ...interface{}) {
	if log.GetPrintEnabled() {
		entry := log.newEntry()
		entry.Print(args...)
		log.releaseEntry(entry)
	}
}

// Printf prints a message to the logger. Ignores logging levels. No
// logging levels, or timestamps are added. The equivalent of
// fmt.Fprintf(out, ...).
func (log *Logger) Printf(format string, args ...interface{}) {
	if log.GetPrintEnabled() {
		entry := log.newEntry()
		entry.Printf(format, args...)
		log.releaseEntry(entry)
	}
}

// Println a message to the logger. It ignores logging levels. No
// logging levels, or timestamps, are added. The equivalent of
// fmt.Fprintln(out, ...).
func (log *Logger) Println(args ...interface{}) {
	if log.GetPrintEnabled() {
		entry := log.newEntry()
		entry.Println(args...)
		log.releaseEntry(entry)
	}
}

// Debug logs a message at level Debug on the logger.
func (log *Logger) Debug(args ...interface{}) {
	if log.GetLevel() >= DebugLevel {
		entry := log.newEntry()
		entry.Debug(args...)
		log.releaseEntry(entry)
	}
}

// Debugf logs a message at level Debug on the logger. Arguments are
// handled in the manner of fmt.Printf.
func (log *Logger) Debugf(format string, args ...interface{}) {
	if log.GetLevel() >= DebugLevel {
		entry := log.newEntry()
		entry.Debugf(format, args...)
		log.releaseEntry(entry)
	}
}

// Debugln logs a message at level Debug on the logger.  It is
// equivalent to Debug().
func (log *Logger) Debugln(args ...interface{}) {
	if log.GetLevel() >= DebugLevel {
		entry := log.newEntry()
		entry.Debugln(args...)
		log.releaseEntry(entry)
	}
}

// Info logs a message at level Info on the logger.
func (log *Logger) Info(args ...interface{}) {
	if log.GetLevel() >= InfoLevel {
		entry := log.newEntry()
		entry.Info(args...)
		log.releaseEntry(entry)
	}
}

// Infof logs a message at level Info on the logger. Arguments are
// handled in the manner of fmt.Printf.
func (log *Logger) Infof(format string, args ...interface{}) {
	if log.GetLevel() >= InfoLevel {
		entry := log.newEntry()
		entry.Infof(format, args...)
		log.releaseEntry(entry)
	}
}

// Infoln logs a message at level Info on logger. It is equivalent to
// Info().
func (log *Logger) Infoln(args ...interface{}) {
	if log.GetLevel() >= InfoLevel {
		entry := log.newEntry()
		entry.Infoln(args...)
		log.releaseEntry(entry)
	}
}

// Warn logs a message at level Warn on the logger.
func (log *Logger) Warn(args ...interface{}) {
	if log.GetLevel() >= WarnLevel {
		entry := log.newEntry()
		entry.Warn(args...)
		log.releaseEntry(entry)
	}
}

// Warnf logs a message at level Warn on the logger. Arguments are
// handled in the manner of fmt.Printf.
func (log *Logger) Warnf(format string, args ...interface{}) {
	if log.GetLevel() >= WarnLevel {
		entry := log.newEntry()
		entry.Warnf(format, args...)
		log.releaseEntry(entry)
	}
}

// Warnln logs a message at level Warn on the logger. It is equivlent
// to Warn().
func (log *Logger) Warnln(args ...interface{}) {
	if log.GetLevel() >= WarnLevel {
		entry := log.newEntry()
		entry.Warnln(args...)
		log.releaseEntry(entry)
	}
}

// Warning logs a message at level Warn on the logger. Warning is an
// alias for Warn.
func (log *Logger) Warning(args ...interface{}) {
	if log.GetLevel() >= WarnLevel {
		entry := log.newEntry()
		entry.Warn(args...)
		log.releaseEntry(entry)
	}
}

// Warningf logs a message at level Warn on the logger. Arguments are
// handled in the manner of fmt.Printf. Warningf is an an alias for
// Warnf.
func (log *Logger) Warningf(format string, args ...interface{}) {
	if log.GetLevel() >= WarnLevel {
		entry := log.newEntry()
		entry.Warnf(format, args...)
		log.releaseEntry(entry)
	}
}

// Warningln logs a message at level Warn on the logger. It is
// equivlent to Warning(). Warningln is an alias for Warnln.
func (log *Logger) Warningln(args ...interface{}) {
	if log.GetLevel() >= WarnLevel {
		entry := log.newEntry()
		entry.Warnln(args...)
		log.releaseEntry(entry)
	}
}

// Error logs a message at level Error on the logger.
func (log *Logger) Error(args ...interface{}) {
	if log.GetLevel() >= ErrorLevel {
		entry := log.newEntry()
		entry.Error(args...)
		log.releaseEntry(entry)
	}
}

// Errorf logs a message at level Error on the logger. Arguments are
// handled in the manner of fmt.Printf.
func (log *Logger) Errorf(format string, args ...interface{}) {
	if log.GetLevel() >= ErrorLevel {
		entry := log.newEntry()
		entry.Errorf(format, args...)
		log.releaseEntry(entry)
	}
}

// Errorln logs a message at level Error on the logger. It is
// equivalent to Error().
func (log *Logger) Errorln(args ...interface{}) {
	if log.GetLevel() >= ErrorLevel {
		entry := log.newEntry()
		entry.Errorln(args...)
		log.releaseEntry(entry)
	}
}

// Fatal logs a message at level Fatal on the logger and exits with
// the DefaultExitCode.
func (log *Logger) Fatal(args ...interface{}) {
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatal(args...)
		log.releaseEntry(entry)
	}
	panic(Exit{DefaultExitCode})
}

// Fatalf logs a message at level Fatal on the logger and exits with
// the DefaultExitCode. Arguments are handled in the manner of
// fmt.Printf.
func (log *Logger) Fatalf(format string, args ...interface{}) {
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatalf(format, args...)
		log.releaseEntry(entry)
	}
	panic(Exit{DefaultExitCode})
}

// Fatalln logs a message at level Fatal on the logger and exits with
// the DefaultExitCode. It is equivalent to Fatal().
func (log *Logger) Fatalln(args ...interface{}) {
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatalln(args...)
		log.releaseEntry(entry)
	}
	panic(Exit{DefaultExitCode})
}

// FatalWithExitCode logs a message at level Fatal on the logger and
// exits with the specified code.
func (log *Logger) FatalWithExitCode(code int, args ...interface{}) {
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatal(args...)
		log.releaseEntry(entry)
	}
	panic(Exit{code})
}

// FatalfWithExitCode logs a message at level Fatal on the logger and
// exits with the specified code. Arguments are handled in the manner
// of fmt.Printf.
func (log *Logger) FatalfWithExitCode(code int, format string, args ...interface{}) {
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatalf(format, args...)
		log.releaseEntry(entry)
	}
	panic(Exit{code})
}

// FatallnWithExitCode logs a message at level Fatal on the logger and
// exits with the specified exit code.
func (log *Logger) FatallnWithExitCode(code int, args ...interface{}) {
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatalln(args...)
		log.releaseEntry(entry)
	}
	panic(Exit{code})
}

// FatalIfError logs a message to the logger and exits with the
// specified code if err is not nil.
func (log *Logger) FatalIfError(err error, code int, args ...interface{}) {
	if err == nil {
		return
	}
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatal(args...)
		log.releaseEntry(entry)
	}
	panic(Exit{code})
}

// FatalfIfError logs a message to the logger and exits with the
// specified code if err is not nil.
func (log *Logger) FatalfIfError(err error, code int, format string, args ...interface{}) {
	if err == nil {
		return
	}
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatalf(format, args...)
		log.releaseEntry(entry)
	}
	panic(Exit{code})

}

// FatallnIfError logs a message to the logger and exits with the
// specified code if err is not nil.
func (log *Logger) FatallnIfError(err error, code int, args ...interface{}) {
	if err == nil {
		return
	}
	if log.GetLevel() >= FatalLevel {
		entry := log.newEntry()
		entry.Fatalln(args...)
		log.releaseEntry(entry)
	}
	panic(Exit{code})
}

// Panic logs a message at level Panic on the logger and then panics.
func (log *Logger) Panic(args ...interface{}) {
	if log.GetLevel() >= PanicLevel {
		entry := log.newEntry()
		entry.Panic(args...)
		log.releaseEntry(entry)
	}
}

// Panicf logs a message at level Panic on the logger and then
// panics. Arguments are handled in the manner of fmt.Printf.
func (log *Logger) Panicf(format string, args ...interface{}) {
	if log.GetLevel() >= PanicLevel {
		entry := log.newEntry()
		entry.Panicf(format, args...)
		log.releaseEntry(entry)
	}
}

// Panicln logs a message at level Panic on the logger. It is
// equivalent to Panic().
func (log *Logger) Panicln(args ...interface{}) {
	if log.GetLevel() >= PanicLevel {
		entry := log.newEntry()
		entry.Panicln(args...)
		log.releaseEntry(entry)
	}
}
