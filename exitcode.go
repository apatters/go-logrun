// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun

// Suggested exit code definitions.
const (
	// ExitOK is the exit code indicating that no unrecovered
	// error occurred.
	ExitOK int = iota

	// ExitErrorInternal is the exit code indicating that some
	// internal error occurred where the condition would only make
	// sense to the programmer.
	ExitErrorInternal

	// ExitErrorUsage is the exit code indicating that some sort
	// of command-line usage error occurred, e.g., an invalid
	// command-line option or missing arguments was specified.
	ExitErrorUsage

	// ExitErrorUser is the exit code indicating that some sort of
	// user error occurred, e.g., the program is being used
	// incorrectly.
	ExitErrorUser

	// ExitErrorIO is the exit code indicating that some sort of
	// I/O error occurred, e.g., a truncated write occurred.
	ExitErrorIO

	// ExitErrorPerm is the exit code indicatating that some sort
	// of permission error occurred, e.g, super-user access to a
	// file is needed.
	ExitErrorPerm

	// ExitErrorInvalid is the exit code indicating that something
	// was invalid, e.g., an incorrect parameter was passed.
	ExitErrorInvalid

	// ExitErrorExists is the exit code indicating that something
	// already exists, e.g., the file that you are trying to
	// create already exists.
	ExitErrorExists

	// ExitErrorNotFound is the exit code indicating that
	// something cannot be found, e.g., you are trying to read a
	// file that does not exist.
	ExitErrorNotFound

	// ExitErrorSignal is the exit code indicating that an
	// untrapped or invalid signal occurred. e.g., an unmasked
	// SIGTERM was received.
	ExitErrorSignal

	// ExitErrorRuntime is the exit code indicating that some sort
	// of run-time error occurred. This exit code is mostly a
	// catchall.
	ExitErrorRuntime
)

var (
	// ExitErrorExecute is returned when the Run() or Shell()
	// methods has an issue with executing the command, It is not
	// the non-zero return code of the command itself. An issue is
	// typically either: a permission problem when running the
	// command, running a command that does not exist, using the
	// wrong path to a command, or using missing/incorrect
	// credentials.
	ExitErrorExecute = ExitErrorInternal
)
