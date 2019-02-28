// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

/*
Package logrun is a Go (golang) package that wraps the standard Go
os/exec and golang.org/x/crypto/ssh packages to run commands either
locally or over ssh while capturing stdout, stderr, and exit
codes. Commands are logged using a specified logging function.
*/
package logrun

import (
	"fmt"
	"strings"

	"github.com/apatters/go-run"
)

var (
	// DefaultLogFunc is the function called when not specified in
	// the LocalConfig or RemoteConfig objects.  The function's
	// type must match LogFunc.
	DefaultLogFunc = DiscardLogFunc

	// FileExistsCmd is the external command used to determine if
	// a file exists or not. This command has been tested on
	// RHEL/CentOS 7 and Ubuntu 18.04.
	FileExistsCmd = "/usr/bin/stat"

	// FileExistsCmdOptions are the command-line options added to
	// FileExistsCmd used to determine if a file exists or
	// not. This command and options has been tested on
	// RHEL/CentOS 7 and Ubuntu 18.04.
	FileExistsCmdOptions = []string{
		"--dereference",
		"--format",
		"%n:%F",
	}

	// DirExistsCmd is the external command used to determine if
	// a directory exists or not. This command has been tested on
	// RHEL/CentOS 7 and Ubuntu 18.04.
	DirExistsCmd = "/usr/bin/stat"

	// DirExistsCmdOptions are the command-line options added to
	// DirExistsCmd used to determine if a file directory or
	// not. This command and options has been tested on
	// RHEL/CentOS 7 and Ubuntu 18.04.
	DirExistsCmdOptions = []string{
		"--dereference",
		"--format",
		"%n:%F",
	}

	// GlobCmd is the external command used to return a list of
	// paths that match a shell glob pattern.  This command has
	// been tested on RHEL/CentOS 7 and Ubuntu 18.04.
	GlobCmd = "/bin/ls"

	// GlobCmdOptions are the command-line options added to
	// GlobCmd used to return a list of paths that match a shell
	// glob pattern. This command and options has been tested on
	// RHEL/CentOS 7 and Ubuntu 18.04.
	GlobCmdOptions = []string{
		"-1",
		"--directory",
	}

	// RsyncCmd is the external command used to copy a directory
	// or file to or from a local or remote destination. This
	// command has been tested on RHEL/CentOS 7 and Ubuntu 18.04.
	RsyncCmd = "/usr/bin/rsync"

	// RsyncCmdOptions are the command-line options added to
	// RsyncCmd used to opy a directory or file to or from a local
	// or remote destination. This command and options has been
	// tested on RHEL/CentOS 7 and Ubuntu 18.04.
	RsyncCmdOptions = []string{
		"--rsh",
		"ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o GlobalKnownHostsFile=/dev/null",
		"--recursive",
		"--links",
		"--times",
	}
)

// LogFunc is the type for the function that will be called to log the
// command.
type LogFunc func(...interface{})

// DiscardLogFunc is used to disable logging. It is the default
// logging function.
func DiscardLogFunc(v ...interface{}) {
}

// LogRun encapsulates a logger used to log and run and either a local
// or remote command.
type LogRun struct {
	Runner  run.Runner
	logFunc LogFunc
	Dryrun  bool
}

// SetLogFunc is used to set the logging function used to log a
// command. The function is typically something like log.Println() or
// logrus.Debug. A custom function of type LogFunc can also be used.
func (r *LogRun) SetLogFunc(f LogFunc) {
	r.logFunc = f
}

// SetDryrun enables/disables the execution of commands. If Dryrun is
// true, the command is only logged.
func (r *LogRun) SetDryrun(dryrun bool) {
	r.Dryrun = dryrun
}

// Run first logs the command and then runs the command. Only logging
// is performed if DryRun is true.
func (r *LogRun) Run(cmd string, args ...string) (string, string, int) {
	msg := r.Runner.FormatRun(cmd, args...)
	r.logFunc(msg)
	if r.Dryrun {
		return "", "", ExitOK
	}

	return r.run(cmd, args...)
}

// FormatRun returns a string representation of the command that would
// be executed using Run().
func (r *LogRun) FormatRun(cmd string, args ...string) string {
	return r.Runner.FormatRun(cmd, args...)
}

// Shell first logs the command and then runs the command in a
// shell. Only logging is performed if DryRun is true.
func (r *LogRun) Shell(cmd string) (string, string, int) {
	msg := r.Runner.FormatShell(cmd)
	r.logFunc(msg)
	if r.Dryrun {
		return "", "", ExitOK
	}
	return r.shell(cmd)
}

// FormatShell returns a string representation of the command that
// would be executed using Shell().
func (r *LogRun) FormatShell(cmd string) string {
	return r.Runner.FormatShell(cmd)
}

// FileExists returns true if filename exists and is a regular
// file. This function is more suited to run remotely.
func (r *LogRun) FileExists(filename string) (bool, error) {
	cmdArgs := append(FileExistsCmdOptions, filename)
	r.logFunc(r.Runner.FormatRun(FileExistsCmd, cmdArgs...))
	if r.Dryrun {
		return true, nil
	}
	stdout, stderr, code := r.run(FileExistsCmd, cmdArgs...)
	if code != 0 {
		if strings.Contains(stderr, "No such file or directory") {
			return false, nil
		}
		return false, fmt.Errorf("could not access %s: %s", filename, stdout)
	}
	fileType := strings.TrimSpace(strings.Split(stdout, ":")[1])
	if fileType != "regular file" && fileType != "regular empty file" {
		return false, fmt.Errorf("%s is not a regular file", filename)
	}

	return true, nil
}

// DirExists returns true if dirname exists and is a directory. This
// method is more suited to run remotely.
func (r *LogRun) DirExists(dirname string) (bool, error) {
	cmdArgs := append(DirExistsCmdOptions, dirname)
	r.logFunc(r.Runner.FormatRun(DirExistsCmd, cmdArgs...))
	if r.Dryrun {
		return true, nil
	}
	stdout, stderr, code := r.run(DirExistsCmd, cmdArgs...)
	if code != 0 {
		if strings.Contains(stderr, "No such file or directory") {
			return false, nil
		}
		return false, fmt.Errorf("could not access %s: %s", dirname, stdout)
	}
	if strings.TrimSpace(strings.Split(stdout, ":")[1]) != "directory" {
		return false, fmt.Errorf("%s is not a directory", dirname)
	}

	return true, nil
}

// Glob returns a list of files matching a shell glob pattern. This
// method is more suited to run remotely.
func (r *LogRun) Glob(pattern string) ([]string, error) {
	args := []string{GlobCmd}
	args = append(args, GlobCmdOptions...)
	args = append(args, pattern)
	cmd := strings.Join(args, " ")
	r.logFunc(r.Runner.FormatShell(cmd))
	stdout, stderr, code := r.shell(cmd)
	if code != 0 {
		return []string{}, fmt.Errorf("glob '%s' failed: %s", pattern, stderr)
	}
	var results []string
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, line)
		}
	}

	return results, nil
}

// Rsync copies files/directories to or from local and remote
// locations using the rsync command. This method is more suited to
// run locally.
func (r *LogRun) Rsync(src string, dest string) error {
	cmdArgs := RsyncCmdOptions
	cmdArgs = append(cmdArgs, src, dest)
	_, stderr, code := r.Run(RsyncCmd, cmdArgs...)
	if code != 0 {
		return fmt.Errorf("rsync command failed: %s", stderr)
	}

	return nil
}

func (r *LogRun) run(cmd string, args ...string) (string, string, int) {
	stdout, stderr, code, err := r.Runner.Run(cmd, args...)
	if err != nil {
		return "", err.Error(), ExitErrorExecute
	}

	return stdout, stderr, code
}

func (r *LogRun) shell(cmd string) (string, string, int) {
	stdout, stderr, code, err := r.Runner.Shell(cmd)
	if err != nil {
		return "", err.Error(), ExitErrorExecute
	}

	return stdout, stderr, code
}
