package logrun_test

import (
	"bytes"

	"github.com/apatters/go-conlog"
)

type existsTestEntry struct {
	Description    string
	Path           string
	ExpectedResult bool
}

type globTestEntry struct {
	Description   string
	Glob          string
	ExpectError   bool
	ExpectedPaths []string
}

type rsyncTestEntry struct {
	Description string
	SrcPath     string
	ExpectError bool
}

// Create a logger to string.
func newLogger() (*conlog.Logger, *bytes.Buffer, *bytes.Buffer) {
	var outBuf = make([]byte, 0, 256)
	var errOutBuf = make([]byte, 0, 256)
	var out = bytes.NewBuffer(outBuf)
	var errOut = bytes.NewBuffer(errOutBuf)

	log := conlog.NewLogger()
	log.SetOutput(out)
	log.SetErrorOutput(errOut)

	return log, out, errOut
}
