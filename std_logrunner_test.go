// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/apatters/go-logrun"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStdRunLogger_SetLogFunc(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	stdout, stderr, code := logrun.Run("/bin/true")
	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/bin/true\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestStdRunLogger_SetDryrun(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	t.Logf("Dryrun = %t", false)
	stdout, stderr, code := logrun.Run("/usr/bin/seq", "1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.EqualValues(t, "1\n", stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/usr/bin/seq 1\n", out.String())
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()

	t.Logf("Dryrun = %t", true)
	logrun.SetDryrun(true)
	defer logrun.SetDryrun(false)
	stdout, stderr, code = logrun.Run("/usr/bin/seq", "1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/usr/bin/seq 1\n", out.String())
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()
}

func TestStdRunLogger_Run(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	stdout, stderr, code := logrun.Run("/usr/bin/seq", "1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.EqualValues(t, "1\n", stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/usr/bin/seq 1\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestStdRunLogger_Shell(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	stdout, stderr, code := logrun.Shell("/usr/bin/seq 1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.EqualValues(t, "1\n", stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/bin/sh -c \"/usr/bin/seq 1\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestStdRunLogger_FileExists(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	path := "/bin/true"
	t.Logf("path = %q", path)
	exists, _ := logrun.FileExists(path)
	t.Logf("exists = %t", exists)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.True(t, exists)
	assert.EqualValues(t, "/usr/bin/stat --dereference --format %n:%F "+path+"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestStdRunLogger_DirExists(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	path := "/etc"
	t.Logf("path = %q", path)
	exists, _ := logrun.DirExists(path)
	t.Logf("exists = %t", exists)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.True(t, exists)
	assert.EqualValues(t, "/usr/bin/stat --dereference --format %n:%F "+path+"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestStdRunLogger_Glob(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	glob := "/etc/passwd*"
	t.Logf("Glob = %q", glob)
	expectedPaths := []string{"/etc/passwd", "/etc/passwd-"}
	t.Logf("ExpectedPaths = %v", expectedPaths)

	results, err := logrun.Glob(glob)

	t.Logf("err = %v", err)
	t.Logf("results = %q", results)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.NoError(t, err)
	assert.EqualValues(t, results, expectedPaths)
	assert.EqualValues(t, "/bin/sh -c \"/bin/ls -1 --directory "+glob+"\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestStdRunLogger_Rysnc(t *testing.T) {
	log, out, errOut := newLogger()
	logrun.SetLogFunc(log.Println)
	defer logrun.SetLogFunc(logrun.DiscardLogFunc)

	srcPath := "/bin/true"
	t.Logf("srcPath = %q", srcPath)

	destDir, err := ioutil.TempDir("", "go-run-")
	t.Logf("err = %v", err)
	require.NoError(t, err)
	err = os.MkdirAll(destDir, 0755)
	t.Logf("err = %v", err)
	require.NoError(t, err)
	defer os.RemoveAll(destDir)

	err = logrun.Rsync(srcPath, fmt.Sprintf("%s/", destDir))
	t.Logf("err = %v", err)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.EqualValues(t, "/usr/bin/rsync --rsh ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o GlobalKnownHostsFile=/dev/null --recursive --links --times "+srcPath+" "+destDir+"/\n", out.String())
	assert.Empty(t, errOut.String())
	assert.NoError(t, err)
}
