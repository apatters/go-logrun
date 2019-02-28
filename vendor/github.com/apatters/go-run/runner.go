// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package run

// Runner is the interface for both Local and Remote.
type Runner interface {

	// Run runs a command like glibc's exec() call. It returns the
	// standard out, standard error, and exit code of the command
	// when it completes.
	Run(cmd string, args ...string) (string, string, int, error)

	// FormatRun returns a string representation of the what
	// command would be run using Run(). Useful for logging
	// commands.
	FormatRun(cmd string, args ...string) string

	// Shell runs a command in a shell. The command is passed to
	// the shell as the -c option, so just about any shell code
	// that can be used on the command-line will be passed to
	// it. It returns the standard out, standard error, and exit
	// code of the command when it completes
	Shell(cmd string) (string, string, int, error)

	// FormatShell returns a string representation of the what
	// command would be run using Shell(). Useful for logging
	// commands.
	FormatShell(cmd string) string
}
