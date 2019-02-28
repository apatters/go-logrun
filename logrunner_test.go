// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package logrun_test

import (
	"regexp"
	"testing"

	"github.com/apatters/go-logrun"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogRunnerSuccess(t *testing.T) {
	var runner logrun.LogRunner

	// Create a logger to string.
	log, out, errOut := newLogger()

	runner = logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	stdout, stderr, code := runner.Run("/bin/true")
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

	runner, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
	})
	t.Logf("err = %v", err)
	require.NoError(t, err)
	stdout, stderr, code = runner.Run("/bin/true")
	t.Logf("stdout = %s", stdout)
	t.Logf("stderr = %s", stderr)
	t.Logf("code = %d", code)
	t.Logf("out = %q", out)
	t.Logf("errOut = %q", errOut)
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Zero(t, code)
	assert.Empty(t, errOut.String())
	assert.Regexp(
		t,
		regexp.MustCompile(`ssh .*@.* /bin/true`),
		out.String())
	out.Reset()
	errOut.Reset()
}
