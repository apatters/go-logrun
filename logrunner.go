// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun

// LogRunner is the interface for both LocalLogRun and RemoteLogRun.
type LogRunner interface {
	SetLogFunc(f LogFunc)
	SetDryrun(dryrun bool)
	Run(cmd string, args ...string) (string, string, int)
	FormatRun(cmd string, args ...string) string
	Shell(cmd string) (string, string, int)
	FormatShell(cmd string) string
	FileExists(filename string) (bool, error)
	DirExists(dirname string) (bool, error)
	Glob(pattern string) ([]string, error)
	Rsync(src string, dest string) error
}
