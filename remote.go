// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun

import (
	"io"

	"github.com/apatters/go-run"
)

// Credentials contains needed credentials to SSH to a host. It can
// use either a password or SSH private key.
type Credentials struct {
	// Hostname is either the hostname or IP of the remote host.
	Hostname string

	// Port is the port used to connect to the ssh server on the
	// remote host.
	Port int

	// Username is the account name used to authenticate on the
	// remote host.
	Username string

	// Password is password used to authenticate on the remote
	// host. Not needed if using PrivateKeyFilename.
	Password string

	// PrivateKeyFilename is the full path the SSH private key
	// used to authenticate with the remote host.  Not used if
	// Password is specified. You must use ssh-agent or something
	// similar to provide the passphrase if the key is passphrase
	// protected.
	PrivateKeyFilename string
}

// RemoteConfig is used to set options in the NewRemoteLoggingRunner
// constructor.
type RemoteConfig struct {
	// LogFunc is used to set the logging function used to log a
	// command. The function is typically something like
	// log.Println() or logrus.Debug. A custom function of type
	// LogFunc can also be used.
	LogFunc LogFunc

	// ShellExecutable is the full path to the shell on the remote
	// host to be run when executing shell commands.
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

	// Credentials are used to authenticate with the remote host.
	Credentials Credentials

	// Dryrun enables/disables the execution of commands. If
	// Dryrun is true, the command is only logged.
	Dryrun bool
}

// NewRemoteLogRun is the constructor for RemoteLogRun used to log and
// run a remote command.
func NewRemoteLogRun(config RemoteConfig) (*LogRun, error) {
	r := new(LogRun)
	creds := run.Credentials{
		Hostname:           config.Credentials.Hostname,
		Port:               config.Credentials.Port,
		Username:           config.Credentials.Username,
		Password:           config.Credentials.Password,
		PrivateKeyFilename: config.Credentials.PrivateKeyFilename,
	}
	remote, err := run.NewRemote(run.RemoteConfig{
		ShellExecutable: config.ShellExecutable,
		Stdin:           config.Stdin,
		Stdout:          config.Stdout,
		Stderr:          config.Stderr,
		Credentials:     creds,
	})
	if err != nil {
		return nil, err
	}

	r.Runner = remote
	if config.LogFunc == nil {
		r.logFunc = DefaultLogFunc
	} else {
		r.logFunc = config.LogFunc
	}
	r.Dryrun = config.Dryrun

	return r, nil
}
