// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/apatters/go-logrun"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	localFileExistsTestTable = []existsTestEntry{
		{"Actual file", "/bin/bash", true},
		{"Symbolic link to a file", "/bin/sh", true},
		{"Non-existent file", "/xyzzy", false},
		{"Non-file", "/etc", false},
	}
	localDirExistsTestTable = []existsTestEntry{
		{"Actual directory", "/run/lock", true},
		{"Symbolic link to a directory", "/var/lock", true},
		{"Non-existent directory", "/xyzzy", false},
		{"Non-directory", "/bin/true", false},
	}
	localGlobTestTable = []globTestEntry{
		{"Single file", "/bin/true*", false, []string{"/bin/true"}},
		{"Multiple files", "/etc/passwd*", false, []string{"/etc/passwd", "/etc/passwd-"}},
		{"Failed glob", "xy*zzy", true, []string{}},
	}
	localRsyncTestTable = []rsyncTestEntry{
		{"Directory", "/etc/cron.daily/", false},
		{"Single file", "/bin/true", false},
		{"Non-existent file", "/bin/xyzzy", true},
		{"Non-existent directory", "/bin/xyzzy.", true},
		/*
			{
				"Remote directory",
				fmt.Sprintf("%s@localhost:/etc/cron.daily/", username()),
				false,
			},
		*/
		{
			"Remote single file",
			fmt.Sprintf("%s@localhost:/bin/true", username()),
			false,
		},
		{
			"Remote non-existent directory",
			fmt.Sprintf("%s@localhost:/xyzzy/", username()),
			true,
		},
		{
			"Remote non-existent single file",
			fmt.Sprintf("%s@localhost:/bin/xyzzy", username()),
			true,
		},
	}
)

func username() (name string) {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	return u.Username
}

func runLocalFileExistsTest(t *testing.T, e existsTestEntry) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})

	t.Logf("%s\n", e.Description)
	t.Logf("path = %q", e.Path)
	t.Logf("expected result = %t", e.ExpectedResult)
	exists, _ := l.FileExists(e.Path)
	t.Logf("exists = %t", exists)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.EqualValues(t, e.ExpectedResult, exists)
	assert.EqualValues(t, "/usr/bin/stat --dereference --format %n:%F "+e.Path+"\n", out.String())
	assert.Empty(t, errOut.String())
}

func runLocalDirExistsTest(t *testing.T, e existsTestEntry) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})

	t.Logf("%s\n", e.Description)
	t.Logf("path = %q", e.Path)
	t.Logf("expected result = %t", e.ExpectedResult)
	exists, _ := l.DirExists(e.Path)
	t.Logf("exists = %t", exists)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.EqualValues(t, e.ExpectedResult, exists)
	assert.EqualValues(t, "/usr/bin/stat --dereference --format %n:%F "+e.Path+"\n", out.String())
	assert.Empty(t, errOut.String())
}

func runLocalGlobTest(t *testing.T, e globTestEntry) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})

	t.Log(e.Description)
	t.Logf("Glob = %q", e.Glob)
	t.Logf("ExpectError = %t", e.ExpectError)
	t.Logf("ExpectedPaths = %v", e.ExpectedPaths)

	results, err := l.Glob(e.Glob)
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
	assert.EqualValues(t, "/bin/sh -c \"/bin/ls -1 --directory "+e.Glob+"\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func runLocalRsyncTest(t *testing.T, e rsyncTestEntry) {
	log, out, errOut := newLogger()
	ll := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	l := logrun.NewLocalLogRun(logrun.LocalConfig{})

	t.Log(e.Description)
	t.Logf("SrcPath = %q", e.SrcPath)
	t.Logf("ExpectError = %t", e.ExpectError)

	destDir, err := ioutil.TempDir("", "go-run-")
	t.Logf("err = %v", err)
	require.NoError(t, err)
	err = os.MkdirAll(destDir, 0755)
	t.Logf("err = %v", err)
	require.NoError(t, err)
	defer os.RemoveAll(destDir)

	err = ll.Rsync(e.SrcPath, fmt.Sprintf("%s/", destDir))
	t.Logf("err = %v", err)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.EqualValues(t, "/usr/bin/rsync --rsh ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o GlobalKnownHostsFile=/dev/null --recursive --links --times "+e.SrcPath+" "+destDir+"/\n", out.String())
	assert.Empty(t, errOut.String())
	if e.ExpectError {
		require.Error(t, err)
	} else {
		require.NoError(t, err)
	}
	if e.ExpectError {
		return
	}

	srcGlob := e.SrcPath
	last := strings.LastIndex(e.SrcPath, ":")
	if last > 0 {
		srcGlob = e.SrcPath[last+1:]
	}
	t.Logf("srcGlob = %q", srcGlob)
	if strings.HasSuffix(e.SrcPath, "/") {
		srcGlob = fmt.Sprintf("%s*", srcGlob)
	}
	srcGlobResults, err := l.Glob(srcGlob)
	t.Logf("err = %v", err)
	t.Logf("srcGlob = %q", srcGlob)
	t.Logf("srcGlobResults = %q", srcGlobResults)
	require.NoError(t, err)

	srcGlobBasenameResults := []string{}
	for _, path := range srcGlobResults {
		srcGlobBasenameResults = append(srcGlobBasenameResults, filepath.Base(path))
	}
	t.Logf("srcGlobBasenameResults = %q", srcGlobBasenameResults)

	destGlob := fmt.Sprintf("%s/*", destDir)
	destGlobResults, err := l.Glob(destGlob)
	t.Logf("err = %v", err)
	t.Logf("destGlob = %q", destGlob)
	t.Logf("destGlobResults = %q", destGlobResults)
	require.NoError(t, err)

	destGlobBasenameResults := []string{}
	for _, path := range destGlobResults {
		destGlobBasenameResults = append(destGlobBasenameResults, filepath.Base(path))
	}
	t.Logf("destGlobBasenameResults = %q", destGlobBasenameResults)
	assert.EqualValues(t, srcGlobBasenameResults, destGlobBasenameResults)
}

func TestLocalLogRun_SetLogFunc(t *testing.T) {
	l := logrun.NewLocalLogRun(logrun.LocalConfig{})
	stdout, stderr, code := l.Run("/bin/true")

	t.Log("Use logrun.DiscardLogFunc")
	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)

	t.Log("Use log.Println")
	log, out, errOut := newLogger()
	l.SetLogFunc(log.Println)
	stdout, stderr, code = l.Run("/bin/true")
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
	out.Reset()
	errOut.Reset()
}

func TestLocalLogRun_SetDryrun(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	t.Logf("Dryrun = %t", false)
	stdout, stderr, code := l.Run("/usr/bin/seq", "1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.False(t, l.Dryrun)
	assert.EqualValues(t, "1\n", stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/usr/bin/seq 1\n", out.String())
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()

	t.Logf("Dryrun = %t", true)
	l.SetDryrun(true)
	stdout, stderr, code = l.Run(
		"/usr/bin/seq",
		"1")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.True(t, l.Dryrun)
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/usr/bin/seq 1\n", out.String())
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()
}

func TestLocalLogRun_RunSuccess(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Run("/bin/true")

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

func TestLocalLogRun_RunFail(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Run("/bin/false")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NotZero(t, code)
	assert.EqualValues(t, "/bin/false\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunNonexistent(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Run("/bin/xyzzy")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`.*/bin/xyzzy: no such file or directory.*`),
		stderr)
	assert.NotZero(t, code)
	assert.EqualValues(t, "/bin/xyzzy\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunExit(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Run("/bin/sh", "-c", "exit 6")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 6)
	assert.EqualValues(t, "/bin/sh -c exit 6\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunOutput(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Run("/bin/ls", "-1", "/bin/true", "/bin/false", "/xyzzy")
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
	assert.EqualValues(t, "/bin/ls -1 /bin/true /bin/false /xyzzy\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunStdin(t *testing.T) {
	log, out, errOut := newLogger()
	stdinStr := "Hello, world"
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Stdin:   strings.NewReader(stdinStr),
	})
	stdout, stderr, code := l.Run(
		"/usr/bin/tr",
		"[:upper:]",
		"[:lower:]")
	t.Logf("stdin = %q", stdinStr)
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, strings.ToLower(stdinStr), stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/usr/bin/tr [:upper:] [:lower:]\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunStdout(t *testing.T) {
	log, out, errOut := newLogger()
	var b bytes.Buffer
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Stdout:  bufio.NewWriter(&b),
	})
	stdout, stderr, code := l.Run(
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
	assert.EqualValues(t, "/bin/ls -1 /bin/true /bin/false /xyzzy\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunStderr(t *testing.T) {
	log, out, errOut := newLogger()
	var b bytes.Buffer
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Stderr:  bufio.NewWriter(&b),
	})
	stdout, stderr, code := l.Run(
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
	assert.EqualValues(t, "/bin/ls -1 /bin/true /bin/false /xyzzy\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunEnv(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Env:     []string{"FIRST=1st", "SECOND=2nd"},
	})
	stdout, stderr, code := l.Run("/usr/bin/env")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, "FIRST=1st\nSECOND=2nd\n", stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.EqualValues(t, "/usr/bin/env\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_RunDir(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Dir:     "/",
	})
	savedDir, _ := os.Getwd()
	stdout, stderr, code := l.Run("/bin/pwd")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	curDir, _ := os.Getwd()
	assert.Equal(t, savedDir, curDir, "Working directory not restored.")
	assert.Equal(t, "/\n", stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.EqualValues(t, "/bin/pwd\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_ShellSuccess(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Shell("exit 0")

	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/bin/sh -c \"exit 0\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_ShellFail(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Shell("exit 1")
	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 1)
	assert.EqualValues(t, "/bin/sh -c \"exit 1\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_ShellNonexistent(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Shell("/bin/xyzzy")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Regexp(
		t,
		regexp.MustCompile(`.*/bin/xyzzy'?: ([nN]o such file or directory)|(not found).*`),
		stderr)
	assert.NotZero(t, code)
	assert.EqualValues(t, "/bin/sh -c \"/bin/xyzzy\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_ShellEnv(t *testing.T) {
	log, out, errOut := newLogger()
	envVars := []string{"FIRST=1st", "SECOND=2nd"}
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Env:     envVars,
	})
	stdout, stderr, code := l.Shell("/usr/bin/env")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	stdoutLines := strings.Split(stdout, "\n")
	assert.Subset(t, stdoutLines, envVars)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.EqualValues(t, "/bin/sh -c \"/usr/bin/env\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_ShellDir(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Dir:     "/",
	})
	savedDir, _ := os.Getwd()
	stdout, stderr, code := l.Shell("pwd")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	curDir, _ := os.Getwd()
	assert.Equal(t, savedDir, curDir, "Working directory not restored.")
	assert.Equal(t, "/\n", stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, code, 0)
	assert.EqualValues(t, "/bin/sh -c \"pwd\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_ShellOutput(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := l.Shell("cd /bin && ls true false xyzzy")

	t.Logf("stdout = %q", stdout)
	t.Logf("stderr = %q", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Equal(t, stdout, "false\ntrue\n")
	assert.Regexp(
		t,
		regexp.MustCompile(`.*cannot access '?xyzzy'?: [nN]o such file or directory.*`),
		stderr)
	assert.NotZero(t, code)
	assert.EqualValues(t, "/bin/sh -c \"cd /bin && ls true false xyzzy\"\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_FormatRun(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})

	msg := l.FormatRun("uname")
	t.Logf("cmd = %q", "uname")
	t.Logf("msg = %q", msg)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Equal(t, msg, "uname")
	assert.Empty(t, out.String())
	assert.Empty(t, errOut.String())
	out.Reset()
	errOut.Reset()

	msg = l.FormatRun("uname", "-a")
	t.Logf("cmd = %q", "uname -a")
	t.Logf("msg = %q", msg)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Equal(t, msg, "uname -a")
	assert.Empty(t, out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_FormatShell(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})

	cmd := fmt.Sprintf(`%s -c "%s"`, "/bin/sh", "uname")
	msg := l.FormatShell("uname")
	t.Logf("cmd = %q", cmd)
	t.Logf("msg = %q", msg)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Equal(t, msg, cmd)
	assert.Empty(t, out.String())
	assert.Empty(t, errOut.String())

	cmd = fmt.Sprintf(`%s -c "%s"`, "/bin/sh", "uname -a")
	msg = l.FormatShell("uname -a")
	t.Logf("cmd = %q", cmd)
	t.Logf("msg = %q", msg)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Equal(t, msg, cmd)
	assert.Empty(t, out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_FileExists(t *testing.T) {
	for _, entry := range localFileExistsTestTable {
		runLocalFileExistsTest(t, entry)
	}
}

func TestLocalLogRun_DirExists(t *testing.T) {
	for _, entry := range localDirExistsTestTable {
		runLocalDirExistsTest(t, entry)
	}
}

func TestLocalLogRum_Glob(t *testing.T) {
	for _, entry := range localGlobTestTable {
		runLocalGlobTest(t, entry)
	}
}

func TestLocalLogRun_Rsync(t *testing.T) {
	for _, entry := range localRsyncTestTable {
		runLocalRsyncTest(t, entry)
	}
}

func TestLocalLogRun_RunDryrun(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Dryrun:  true,
	})
	stdout, stderr, code := l.Run("/bin/false")

	t.Logf("stdout %q", stdout)
	t.Logf("stderr %q", stderr)
	t.Logf("code %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/bin/false\n", out.String())
	assert.Empty(t, errOut.String())
}

func TestLocalLogRun_ShellDryrun(t *testing.T) {
	log, out, errOut := newLogger()
	l := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
		Dryrun:  true,
	})
	stdout, stderr, code := l.Shell("/bin/false")

	t.Logf("stdout %q", stdout)
	t.Logf("stderr %q", stderr)
	t.Logf("code %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.EqualValues(t, "/bin/sh -c \"/bin/false\"\n", out.String())
	assert.Empty(t, errOut.String())
}
