// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package conlog

import (
	"io"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = NewLogger()
)

// SetOutput sets the writer used for Print, Debug, and Warning
// messages for the standard logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// SetErrorOutput sets the writer used for Error, Fatal, and Panic
// messages for the standard logger.
func SetErrorOutput(w io.Writer) {
	std.SetErrorOutput(w)
}

// SetLevel sets the logger level for the standard logger.
func SetLevel(level Level) {
	std.SetLevel(level)
}

// GetLevel returns the current logging level for the standard logger.
func GetLevel() Level {
	return std.GetLevel()
}

// SetPrintEnabled sets the PrintEnabled setting for the standard
// logger.
func SetPrintEnabled(enabled bool) {
	std.SetPrintEnabled(enabled)
}

// GetPrintEnabled returns the PrintEnabled setting for the standard
// logger.
func GetPrintEnabled() bool {
	return std.GetPrintEnabled()
}

// SetFormatter sets the formatter used when printing entries to the
// standard logger.
func SetFormatter(formatter Formatter) {
	std.SetFormatter(formatter)
}

// Printf prints a message to the standard logger. Ignores logging
// levels. No logging levels, timestamps, or key files are added. The
// equivalent of fmt.Fprintf.
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

// Debugf logs a message at level Debug to the standard
// logger. Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Infof logs a message at level Info to the standard
// logger. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warnf logs a message at level Warn to the standard
// logger. Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Warningf logs a message at level Warn to the standard
// logger. Arguments are handled in the manner of fmt.Printf. Warningf
// is an an alias for Warnf.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Errorf logs a message at level Error to the standard
// logger. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal to the standard logger and
// exits with the DefaultExitCode. Arguments are handled in the manner
// of fmt.Printf.
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// Panicf logs a message at level Panic to the standard logger and
// then panics. Arguments are handled in the manner of fmt.Printf.
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

// Print prints a message to the standard logger. It ignores logging
// levels. No logging levels, timestamps, or key files are added. No
// newline is added. The equivalent of fmt.Fprint.
func Print(args ...interface{}) {
	std.Print(args...)
}

// Debug logs a message at level Debug to the standard logger.
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Info logs a message at level Info to the standard logger.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Warn logs a message at level Warn to the standard logger.
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Warning logs a message at level Warn to the standard
// logger. Warning is an alias for Warn.
func Warning(args ...interface{}) {
	std.Warning(args...)
}

// Error logs a message at level Error to the standard logger.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Fatal logs a message at level Fatal to the standard logger and
// exits with the DefaultExitCode.
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Panic logs a message at level Panic to the standard logger and then
// panics.
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// Println a message to the standard logger. It ignores logging
// levels. No logging levels, timestamps, or key files are added. The
// equivalent of fmt.Fprintln().
func Println(args ...interface{}) {
	std.Println(args...)
}

// Debugln logs a message at level Debug to the standard logger.  It
// is equivalent to Debug().
func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

// Infoln logs a message at level Info to the standard logger. It is
// equivalent to Info().
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

// Warnln logs a message at level Warn to the standard logger. It is
// equivlent to Warn().
func Warnln(args ...interface{}) {
	std.Warnln(args...)
}

// Warningln logs a message at level Warn to the standard logger. It
// is equivlent to Warning(). Warningln is an alias for Warnln.
func Warningln(args ...interface{}) {
	std.Warningln(args...)
}

// Errorln logs a message at level Error to the standard logger. It is
// equivalent to Error().
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

// Fatalln logs a message at level Fatal to the standard logger and
// exits with the DefaultExitCode. It is equivalent to Fatal().
func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}

// Panicln logs a message at level Panic to the standard logger. It is
// equivalent to Panic().
func Panicln(args ...interface{}) {
	std.Panicln(args...)
}

// FatalWithExitCode logs a message at level Fatal to the standard
// logger and exits with the specified code.
func FatalWithExitCode(code int, args ...interface{}) {
	std.FatalWithExitCode(code, args...)
}

// FatalfWithExitCode logs a message at level Fatal to the standard
// logger and exits with the specified code. Arguments are handled in
// the manner of fmt.Printf.
func FatalfWithExitCode(code int, format string, args ...interface{}) {
	std.FatalfWithExitCode(code, format, args...)
}

// FatallnWithExitCode logs a message at level Fatal to the standard
// logger and exits with the specified exit code.
func FatallnWithExitCode(code int, args ...interface{}) {
	std.FatallnWithExitCode(code, args...)
}

// FatalIfError logs a message to the standard logger and exits with
// the specified code if err is not nil.
func FatalIfError(err error, code int, args ...interface{}) {
	std.FatalIfError(err, code, args...)
}

// FatalfIfError logs a message to the standard logger and exits with
// the specified code if err is not nil.
func FatalfIfError(err error, code int, format string, args ...interface{}) {
	std.FatalfIfError(err, code, format, args...)
}

// FatallnIfError logs a message to the standar logger and exits with
// the specified code if err is not nil.
func FatallnIfError(err error, code int, args ...interface{}) {
	std.FatallnIfError(err, code, args...)
}
