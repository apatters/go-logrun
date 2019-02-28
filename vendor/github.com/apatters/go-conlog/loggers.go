// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package conlog

// Loggers outputs the same message at the same log level to a list of
// loggers. It is often used when you want the same message to go to
// both the console and a log file.
type Loggers struct {
	// A list of Logger(s) to output to.
	Loggers []ConLogger
}

// NewLogs is the constructor for Logs.
func NewLoggers(loggers ...ConLogger) *Loggers {
	return &Loggers{
		Loggers: loggers,
	}
}

// SetLevel sets the logger level for all loggers..
func (logs *Loggers) SetLevel(level Level) {
	for _, logger := range logs.Loggers {
		logger.SetLevel(level)
	}
}

// SetPrintEnabled enables/disables Print*- output for all loggers.
func (logs *Loggers) SetPrintEnabled(enabled bool) {
	for _, logger := range logs.Loggers {
		logger.SetPrintEnabled(enabled)
	}
}

// Print print a message to the loggers. It ignores logging levels. No
// logging levels or timestamps are added. No newline is added. The
// equivalent of fmt.Fprint.
func (logs *Loggers) Print(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Print(args...)
	}
}

// Printf prints a message to the loggers. It ignores logging
// levels. No logging levels or timestamps are added. The equivalent
// of fmt.Fprintf.
func (logs *Loggers) Printf(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Printf(format, args...)
	}
}

// Println prints a message to the logger. It ignores logging
// levels. No logging levels, timestamps, or key files are added. The
// equivalent of fmt.Fprintln.
func (logs *Loggers) Println(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Println(args...)
	}
}

// Debug logs a message at level Debug on all loggers. Arguments are
// handled in the manner of fmt.Print.
func (logs *Loggers) Debug(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Debug(args...)
	}
}

// Debugf logs a message at level Debug on all loggers. Arguments are
// handled in the manner of fmt.Printf.
func (logs *Loggers) Debugf(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Debugf(format, args...)
	}
}

// Debugln logs a message at level Debug on all loggers. Arguments are
// handled in the manner of fmt.Println.
func (logs *Loggers) Debugln(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Debugln(args...)
	}
}

// Info logs a message at level Info on all loggers. Arguments are
// handled in the manner of fmt.Print.
func (logs *Loggers) Info(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Info(args...)
	}
}

// Infof logs a message at level Info on all loggers. Arguments are
// handled in the manner of fmt.Printf.
func (logs *Loggers) Infof(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Infof(format, args...)
	}
}

// Infoln logs a message at level Info on all loggers. Arguments are
// handled in the manner of fmt.Println.
func (logs *Loggers) Infoln(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Infoln(args...)
	}
}

// Warn logs a message at level Warn on all loggers. Arguments are
// handled in the manner of fmt.Print.
func (logs *Loggers) Warn(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Warn(args...)
	}
}

// Warnf logs a message at level Warn on all loggers. Arguments are
// handled in the manner of fmt.Printf.
func (logs *Loggers) Warnf(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Warnf(format, args...)
	}
}

// Warnln logs a message at level Warn on all loggers. Arguments are
// handled in the manner of fmt.Println.
func (logs *Loggers) Warnln(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Warnln(args...)
	}
}

// Warning logs a message at level Warning on all loggers. Arguments
// are handled in the manner of fmt.Print. Warning is an alias for
// Warn.
func (logs *Loggers) Warning(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Warning(args...)
	}
}

// Warningf logs a message at level Warning on all loggers.  Arguments
// are handled in the manner of fmt.Printf. Warningf is an alias for
// Warnf.
func (logs *Loggers) Warningf(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Warningf(format, args...)
	}
}

// Warningln logs a message at level Warning on all loggers. Arguments
// are handled in the manner of fmt.Println. Warningf is an alias for
// Warnln.
func (logs *Loggers) Warningln(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Warnln(args...)
	}
}

// Error logs a message at level Error on all loggers. Arguments are
// handled in the manner of fmt.Print.
func (logs *Loggers) Error(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Error(args...)
	}
}

// Errorf logs a message at level Error on all loggers. Arguments are
// handled in the manner of fmt.Printf.
func (logs *Loggers) Errorf(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Errorf(format, args...)
	}
}

// Errorln logs a message at level Error on all loggers. Arguments are
// handled in the manner of fmt.Println.
func (logs *Loggers) Errorln(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Errorln(args...)
	}
}

// Fatal displays a message at Fatal level to the first logger and
// exits with DefaultExitCode error code. Arguments are handled in the
// manner of fmt.Print.
func (logs *Loggers) Fatal(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Fatal(args...)
	}
}

// Fatalf displays a message at Fatal level to the first logger and
// exits with DefaultExitCode error code. Arguments are handled in the
// manner of fmt.Printf.
func (logs *Loggers) Fatalf(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Fatalf(format, args...)
	}
}

// Fatalln displays a message at Fatal level to the first logger and
// exits with DefaultExitCode error code.  Arguments are handled in
// the manner of fmt.Println.
func (logs *Loggers) Fatalln(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Fatalln(args...)
	}
}

// FatalWithExitCode displays an message at Fatal level to the first
// logger and exits with a given exit code.  Arguments are
// handled in the manner of fmt.Print.
func (logs *Loggers) FatalWithExitCode(code int, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.FatalWithExitCode(code, args...)
	}
}

// FatalfWithExitCode displays an message at Fatal level to the first
// logger and exits with a given exit code.  Arguments are handled in
// the manner of fmt.Printf.
func (logs *Loggers) FatalfWithExitCode(code int, format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.FatalfWithExitCode(code, format, args...)
	}
}

// FatallnWithExitCode displays an message at Fatal level to the first
// logger and exits with a given exit code.  Arguments are handled in
// the manner of fmt.Println.
func (logs *Loggers) FatallnWithExitCode(code int, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.FatallnWithExitCode(code, args...)
	}
}

// FatalIfError displays an message at Fatal level to the first logger
// and exits with a given exit code if err is not nil. Arguments are
// handled in the manner of fmt.Print.
func (logs *Loggers) FatalIfError(err error, code int, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.FatalIfError(err, code, args...)
	}
}

// FatalfIfError displays a message at Fatal level to the first logger
// and exits with a given exit code and if err is not nil.  Arguments
// are handled in the manner of fmt.Printf.
func (logs *Loggers) FatalfIfError(err error, code int, format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.FatalfIfError(err, code, format, args...)
	}
}

// FatallnIfError displays a message at Fatal level to the first
// logger and exits with a given exit code and if err is not nil.
// Arguments are handled in the manner of fmt.Printf.
func (logs *Loggers) FatallnIfError(err error, code int, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.FatallnIfError(err, code, args...)
	}
}

// Panic displays a message to the logger at Panic level and then
// panics. Arguments are handled in the manner of fmt.Print.
func (logs *Loggers) Panic(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Panic(args...)
	}
}

// Panicf displays a message to the logger at Panic level and then
// panics. Arguments are handled in the manner of fmt.Printf.
func (logs *Loggers) Panicf(format string, args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Panicf(format, args...)
	}
}

// Panicln displays a message to the logger at Panic level and then
// panics. Arguments are handled in the manner of fmt.Println.
func (logs *Loggers) Panicln(args ...interface{}) {
	for _, logger := range logs.Loggers {
		logger.Panicln(args...)
	}
}
