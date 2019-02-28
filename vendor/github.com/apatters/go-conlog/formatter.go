// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package conlog

import "time"

const (
	// Default time stamp format used when displaying wall clock
	// time. Wall clock time is the current time.Now() value when
	// the message is printed.
	DefaultWallclockTimestampFormat = time.RFC3339

	// Default time stamp format used when displaying elapsed
	// time. Elapsed time is the number of seconds since the
	// program started running is seconds.
	DefaultElapsedTimestampFormat = "%04d"
)

// The Formatter interface is used to implement a custom Formatter. It
// takes an `Entry`.
//
// Format is expected to return an array of bytes which are then
// logged to `log.Out`.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}
