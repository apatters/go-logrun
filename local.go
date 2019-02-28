// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun

import (
	"io"

	"github.com/apatters/go-run"
)

// LocalConfig is used to set options in the NewLocalLogRun
// constructor.
type LocalConfig struct {
	// LogFunc is used to set the logging function used to log a
	// command. The function is typically something like
	// log.Println() or logrus.Debug. A custom function of type
	// LogFunc can also be used.
	LogFunc LogFunc

	// ShellExecutable is the full path to the shell to be run
	// when executing shell commands.
	ShellExecutable string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	Env []string

	// Dir specifies the working directory of the command.  If Dir
	// is the empty string, Run runs the command in the calling
	// process's current directory.
	Dir string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device
	// (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the
	// command over a pipe. In this case, Wait does not complete
	// until the goroutine stops copying, either because it has
	// reached the end of Stdin (EOF or a read error) or because
	// writing to the pipe returned an error.
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// Dryrun enables/disables the execution of commands. If
	// Dryrun is true, the command is only logged.
	Dryrun bool
}

// NewLocalLogRun is the constructor for LogRun used to log and run a
// local command.
func NewLocalLogRun(config LocalConfig) *LogRun {
	r := new(LogRun)
	r.Runner = run.NewLocal(run.LocalConfig{
		ShellExecutable: config.ShellExecutable,
		Env:             config.Env,
		Dir:             config.Dir,
		Stdin:           config.Stdin,
		Stdout:          config.Stdout,
		Stderr:          config.Stderr,
	})
	if config.LogFunc == nil {
		r.logFunc = DefaultLogFunc
	} else {
		r.logFunc = config.LogFunc
	}
	r.Dryrun = config.Dryrun

	return r
}
