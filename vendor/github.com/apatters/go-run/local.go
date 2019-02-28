// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package run

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// LocalConfig is used to configure the Local constructor.
type LocalConfig struct {
	// ShellExecutable is the path to the shell used to run
	// commands in Shell() methods. See Local for details.
	ShellExecutable string

	// Env specifies the environment of the command to be run.
	// See Local for details.
	Env []string

	// Dir specifies the working directory of the command.  See
	// Local for details. The default is the empty string.
	Dir string

	// Stdin specifies the process's standard input. See Local for
	// details.
	Stdin io.Reader

	// Stdout specifies the process's standard output. See Local
	// for details.
	Stdout io.Writer

	// Stderr specifies the process's standard error. See Local
	// for details.
	Stderr io.Writer
}

// Local wraps os/exec Cmd to make running external commands on the
// local host relatively as easy as when running them in shell script.
type Local struct {
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

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, the descriptor's output is captured and
	// is returned in the Run and Shell functions.
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer
}

// NewLocal is the constuctor for Local. It takes a LocalConfig
// object to configure it. The following configuration options are set
// if the default LocalConfig constructor, LocalConfig{}, is used:
//
//     ShellExecutable = DefaultShellExecutable
//     Env = []string{} // Use existing environment.
//     Dir = nil        // Current working directory.
//     Stdin = nil      // Discard stdin.
//     Stdout = nil     // Capture stdout.
//     Stderr = nil     // Capture stderr,
func NewLocal(config LocalConfig) *Local {
	local := new(Local)
	if len(config.ShellExecutable) == 0 {
		local.ShellExecutable = DefaultShellExecutable
	}
	local.Env = config.Env
	local.Dir = config.Dir
	local.Stdin = config.Stdin
	local.Stdout = config.Stdout
	local.Stderr = config.Stderr

	return local
}

func (l *Local) exec(command string, args ...string) (string, string, int, error) {
	var err error
	code := 0
	cmd := exec.Command(command, args...)
	cmd.Env = l.Env
	cmd.Dir = l.Dir

	// Hook up standard files.
	cmd.Stdin = l.Stdin
	var stdoutPipe io.Reader
	if l.Stdout == nil {
		stdoutPipe, err = cmd.StdoutPipe()
		if err != nil {
			return "", "", 0, err
		}
	} else {
		cmd.Stdout = l.Stdout
	}
	var stderrPipe io.Reader
	if l.Stderr == nil {
		stderrPipe, err = cmd.StderrPipe()
		if err != nil {
			return "", "", 0, err
		}
	} else {
		cmd.Stderr = l.Stderr
	}

	// Run the command.
	err = cmd.Start()
	if err != nil {
		return "", "", 0, err
	}

	// Process the I/O.
	var stdoutBuf []byte
	if l.Stdout == nil {
		stdoutBuf, err = ioutil.ReadAll(stdoutPipe)
		if err != nil {
			return "", "", 0, err
		}
	}
	var stderrBuf []byte
	if l.Stderr == nil {
		stderrBuf, err = ioutil.ReadAll(stderrPipe)
		if err != nil {
			return "", "", 0, err
		}
	}

	// Wait for the command to complete and check for errors.
	if err = cmd.Wait(); err != nil {
		switch err.(type) {
		case *exec.ExitError:
			// Extract exit code from error message.
			re := regexp.MustCompile("^exit status ([0-9]+)$")
			match := re.FindStringSubmatch(err.Error())
			if match != nil {
				code, err = strconv.Atoi(match[1])
				if err != nil {
					return "", "", 0, err
				}
			}
		default:
			return "", "", 0, err
		}
	}

	return string(stdoutBuf), string(stderrBuf), code, err
}

// Run runs a command like glibc's exec() call. It returns the
// standard out, standard error, and exit code of the command when it
// completes.
func (l *Local) Run(cmd string, args ...string) (string, string, int, error) {
	stdout, stderr, code, err := l.exec(cmd, args...)

	return stdout, stderr, code, err
}

// FormatRun returns a string representation of the what command would
// be run using Run(). Useful for logging commands.
func (l *Local) FormatRun(cmd string, args ...string) string {
	return strings.TrimSpace(cmd + " " + strings.Join(args, " "))
}

// Shell runs a command in a shell. The command is passed to the shell
// as the -c option, so just about any shell code that can be used on
// the command-line will be passed to it. It returns the standard out,
// standard error, and exit code of the command when it completes.
func (l *Local) Shell(cmd string) (string, string, int, error) {
	stdout, stderr, code, err := l.exec(l.ShellExecutable, "-c", cmd)

	return stdout, stderr, code, err
}

// FormatShell returns a string representation of the what command
// would be run using Shell(). Useful for logging commands.
func (l *Local) FormatShell(cmd string) string {
	return strings.TrimSpace(fmt.Sprintf(`%s -c "%s"`, l.ShellExecutable, cmd))
}
