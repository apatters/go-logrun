// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package run

var (
	// The standard runner is used to run local commands without
	// the need to explicitly use a constructor.
	std = NewLocal(LocalConfig{})
)

// Run runs a command like glibc's exec() call using the standard
// runner. It returns the standard out, standard error, and exit code
// of the command when it completes.
func Run(cmd string, args ...string) (string, string, int, error) {
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
func Shell(cmd string) (string, string, int, error) {
	return std.Shell(cmd)
}

// FormatShell returns a string representation of the what command
// would be run using the standard runner's Shell() method. Useful
// for logging commands.
func FormatShell(cmd string) string {
	return std.FormatShell(cmd)
}
