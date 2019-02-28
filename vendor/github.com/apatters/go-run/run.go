// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

/*
Package run is a Go (golang) package that wraps the standard Go
os/exec and golang.org/x/crypto/ssh packages to run commands either
locally or over ssh while capturing stdout, stderr, and exit codes.
*/
package run

const (
	// DefaultShellExecutable is the shell that will be run when
	// using Shell() methods.
	DefaultShellExecutable = "/bin/sh"
)
