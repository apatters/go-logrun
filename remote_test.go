// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun_test

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/apatters/go-logrun"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	remoteFileExistsTestTable = []existsTestEntry{
		{"Actual file", "/bin/bash", true},
		{"Symbolic link to a file", "/bin/sh", true},
		{"Non-existent file", "/xyzzy", false},
		{"Non-file", "/etc", false},
	}
	remoteDirExistsTestTable = []existsTestEntry{
		{"Actual directory", "/run/lock", true},
		{"Symbolic link to a directory", "/var/lock", true},
		{"Non-existent directory", "/xyzzy", false},
		{"Non-directory", "/bin/true", false},
	}
	remoteGlobTestTable = []globTestEntry{
		{"Single file", "/bin/true*", false, []string{"/bin/true"}},
		{"Multiple files", "/etc/passwd*", false, []string{"/etc/passwd", "/etc/passwd-"}},
		{"Failed glob", "xy*zzy", true, []string{}},
	}
)

func runRemoteFileExistsTest(t *testing.T, e existsTestEntry) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)

	t.Logf("%s\n", e.Description)
	t.Logf("path = %q", e.Path)
	t.Logf("expected result = %t", e.ExpectedResult)
	exists, _ := r.FileExists(e.Path)

	t.Logf("exists = %t", exists)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.EqualValues(t, e.ExpectedResult, exists)
	re := fmt.Sprintf("ssh .*@.* /usr/bin/stat --dereference --format %%n:%%F %s\n", e.Path)
	assert.Regexp(
		t,
		regexp.MustCompile(re),
		out.String())
	assert.Empty(t, errOut.String())
}

func runRemoteDirExistsTest(t *testing.T, e existsTestEntry) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)

	t.Logf("%s\n", e.Description)
	t.Logf("path = %q", e.Path)
	t.Logf("expected result = %t", e.ExpectedResult)
	exists, _ := r.DirExists(e.Path)

	t.Logf("exists = %t", exists)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.EqualValues(t, e.ExpectedResult, exists)
	re := fmt.Sprintf("ssh .*@.* /usr/bin/stat --dereference --format %%n:%%F %s\n", e.Path)
	assert.Regexp(
		t,
		regexp.MustCompile(re),
		out.String())
	assert.Empty(t, errOut.String())
}

func runRemoteGlobTest(t *testing.T, e globTestEntry) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)

	t.Log(e.Description)
	t.Logf("Glob = %q", e.Glob)
	t.Logf("ExpectError = %t", e.ExpectError)
	t.Logf("ExpectedPaths = %v", e.ExpectedPaths)

	results, err := r.Glob(e.Glob)
	t.Logf("err = %v", err)
	t.Logf("results = %q", results)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	if e.ExpectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
	assert.EqualValues(t, results, e.ExpectedPaths)
	re := fmt.Sprintf(
		"ssh .*@.* /bin/sh -c \"/bin/ls -1 --directory %s\"\n",
		regexp.QuoteMeta(e.Glob))
	t.Logf("re = %q", re)
	assert.Regexp(
		t,
		regexp.MustCompile(re),
		out.String())
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_SetLogFunc(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Run("/bin/true")

	t.Log("Use logrun.DiscardLogFunc")
	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Empty(t, out.String())
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()

	t.Log("Use log.Println")
	r.SetLogFunc(log.Println)
	stdout, stderr, code = r.Run("/bin/true")
	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/true\n`),
		out)
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()
}

func TestRemoteLogRun_SetDryrun(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	t.Logf("Dryrun = %t", false)
	stdout, stderr, code := r.Run("/usr/bin/seq", "1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.False(t, r.Dryrun)
	assert.EqualValues(t, "1\n", stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /usr/bin/seq 1\n`),
		out)
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()

	t.Logf("Dryrun = %t", true)
	r.SetDryrun(true)
	stdout, stderr, code = r.Run("/usr/bin/seq", "1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.True(t, r.Dryrun)
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /usr/bin/seq 1\n`),
		out)
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()
}

func TestRemoteLogRun_RunSuccess(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Run("/bin/true")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/true\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_RunFail(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Run("/bin/false")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NotZero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/false\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_RunNonexistent(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Run("/bin/xyzzy")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`.*[nN]o such file or directory.*`),
		stderr)
	assert.NotZero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/xyzzy\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_RunOutput(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)

	stdout, stderr, code := r.Run("/bin/ls", "-1", "/bin/true", "/bin/false", "/xyzzy")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, stdout, "/bin/false\n/bin/true\n")
	assert.Regexp(
		t,
		regexp.MustCompile(".*/bin/ls: cannot access .*xyzzy.*: [nN]o such file or directory.*"),
		stderr)
	assert.NotZero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/ls -1 /bin/true /bin/false /xyzzy\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_RunStdin(t *testing.T) {
	log, out, errOut := newLogger()
	stdinStr := "Hello, World"
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
		Stdin:   strings.NewReader(stdinStr),
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Run(
		"/usr/bin/tr",
		"HW",
		"hw")

	t.Logf("stdin = %q", stdinStr)
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, strings.ToLower(stdinStr), stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /usr/bin/tr HW hw\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_RunStdout(t *testing.T) {
	log, out, errOut := newLogger()
	var b bytes.Buffer
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
		Stdout:  bufio.NewWriter(&b),
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")

	t.Logf("b = %q", b.String())
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, b.String(), "/bin/false\n/bin/true\n")
	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`/bin/ls: cannot access .*xyzzy.*: No such file or directory`),
		stderr)
	assert.NotZero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/ls -1 /bin/true /bin/false /xyzzy\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_RunStderr(t *testing.T) {
	log, out, errOut := newLogger()
	var b bytes.Buffer
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
		Stderr:  bufio.NewWriter(&b),
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")

	t.Logf("b = %q", b.String())
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, stdout, "/bin/false\n/bin/true\n")
	assert.Regexp(
		t,
		regexp.MustCompile(`/bin/ls: cannot access .*xyzzy.*: No such file or directory`),
		b.String())
	assert.Empty(t, stderr)
	assert.NotZero(t, code, 2)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/ls -1 /bin/true /bin/false /xyzzy\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_ShellSuccess(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Shell("exit 0")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "exit 0"\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_ShellFail(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Shell("exit 1")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 1)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "exit 1"\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_ShellNonexistent(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Shell("/bin/xyzzy")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`.*/bin/xyzzy: not found.*`),
		stderr)
	assert.NotZero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "/bin/xyzzy"\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_ShellOutput(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code := r.Shell("cd /bin && ls true false xyzzy")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, stdout, "false\ntrue\n")
	assert.Regexp(
		t,
		regexp.MustCompile(`.*cannot access .*xyzzy.*: [nN]o such file or directory.*`),
		stderr)
	assert.NotZero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "cd /bin && ls true false xyzzy"\n`),
		out)
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_FormatRun(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	msg := r.FormatRun("uname", "-a")

	t.Logf("cmd = %q", "uname -a")
	t.Logf("msg = %q", msg)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* uname -a`),
		msg)
	assert.Empty(t, out.String())
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_FormatShell(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	msg := r.FormatShell("uname -a")

	t.Logf("msg = %q", msg)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "uname -a"`),
		msg)
	assert.Empty(t, out.String())
	assert.Empty(t, errOut.String())
}

func TestRemoteLogRun_FileExists(t *testing.T) {
	for _, entry := range remoteFileExistsTestTable {
		runRemoteFileExistsTest(t, entry)
	}
}

func TestRemoteLogRun_DirExists(t *testing.T) {
	for _, entry := range remoteDirExistsTestTable {
		runRemoteDirExistsTest(t, entry)
	}
}

func TestRemoteLogRun_Glob(t *testing.T) {
	for _, entry := range remoteGlobTestTable {
		runRemoteGlobTest(t, entry)
	}
}

func TestRemoteLogRun_RunDryrun(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
		Dryrun:  true,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)

	stdout, stderr, code := r.Run("/bin/true")
	t.Logf("stdout %q", stdout)
	t.Logf("stderr %q", stderr)
	t.Logf("code %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/true`),
		out)
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()

	stdout, stderr, code = r.Run("/bin/false")
	t.Logf("stdout %q", stdout)
	t.Logf("stderr %q", stderr)
	t.Logf("code %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/false`),
		out)
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()
}

func TestRemoteLogRun_ShellDryrun(t *testing.T) {
	log, out, errOut := newLogger()
	r, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
		Dryrun:  true,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)

	stdout, stderr, code := r.Shell("/bin/true")
	t.Logf("stdout %q", stdout)
	t.Logf("stderr %q", stderr)
	t.Logf("code %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "/bin/true"`),
		out)
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()

	stdout, stderr, code = r.Shell("/bin/false")
	t.Logf("stdout %q", stdout)
	t.Logf("stderr %q", stderr)
	t.Logf("code %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/sh -c "/bin/false"`),
		out)
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()
}
