// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun

var (
	// The standard runner is used to run local commands without
	// the need to explicitly use a constructor.
	std = NewLocalLogRun(LocalConfig{})
)

// SetLogFunc sets the logging function used to log commands by
// calling the standard run logger's SetLogFunc() method.
func SetLogFunc(f LogFunc) {
	std.SetLogFunc(f)
}

// SetDryrun is used to enable/disable whether commands are just
// logged and not executed by calling the run logger's SetDryrun()
// method.
func SetDryrun(dryrun bool) {
	std.SetDryrun(dryrun)
}

// Run runs a command like glibc's exec() call using the standard
// runner. It returns the standard out, standard error, and exit code
// of the command when it completes.
func Run(cmd string, args ...string) (string, string, int) {
	return std.Run(cmd, args...)
}

// FormatRun returns a string representation of the what command would
// be run using the standard runner's Run() method. Useful for logging
// commands.
func FormatRun(cmd string, args ...string) string {
	return std.FormatRun(cmd, args...)
}

// Shell runs a command in a shell using the standard runner. The
// command is passed to the shell as the -c option, so just about any
// shell code that can be used on the command-line will be passed to
// it. It returns the standard out, standard error, and exit code of
// the command when it completes
func Shell(cmd string) (string, string, int) {
	return std.Shell(cmd)
}

// FormatShell returns a string representation of the what command
// would be run using the standard runner's Shell() method. Useful
// for logging commands.
func FormatShell(cmd string) string {
	return std.FormatShell(cmd)
}

// FileExists returns true if filename exists and is a regular file
// using the standard log runner's FileExist() method.
func FileExists(filename string) (bool, error) {
	return std.FileExists(filename)
}

// DirExists returns true if dirname exists and is a directory using
// the standard log runner's DirExists() method..
func DirExists(dirname string) (bool, error) {
	return std.DirExists(dirname)
}

// Glob returns a list of files matching a glob pattern using the
// standard log runner's Glob() method.
func Glob(pattern string) ([]string, error) {
	return std.Glob(pattern)
}

// Rsync copies files/directories using the rsync command by calling
// the standard log runner's Rsync() method().
func Rsync(src string, dest string) error {
	return std.Rsync(src, dest)
}
