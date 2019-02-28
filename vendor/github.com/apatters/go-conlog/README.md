go-conlog
=========

Conlog is a golang logging package heavily leveraged from the popular
[logrus](https://github.com/sirupsen/logrus) package. It differs from
logrus in that it is optimized for console output rather than for use
with log files.  Like logrus, conlog is completely API compatible with
the standard library logger.

[![Build Status](https://travis-ci.org/apatters/go-conlog.svg)](https://travis-ci.org/apatters/go-conlog) [![GoDoc](https://godoc.org/github.com/apatters/go-conlog?status.svg)](https://godoc.org/github.com/apatters/go-conlog)


Features
-------

* Level logging -- only log messages at at or below one of the following levels: Panic, Fatal, Error, Warning, Info, or Debug.
* Optionally display log levels in the log message.
* Optionally display wallclock or elapsed time in log messages.
* Optional colorized log level when output is to a TTY.
* Print*-style message logging that ignores the log level which can be optionally suppressed for verbose/non-verbose output.
* Log a message to multiple logs with one call.


Documentation
-------------

Documentation can be found at [GoDoc](https://godoc.org/github.com/apatters/go-conlog)


Installation
------------

Install conlog using the "go get" command:

```bash
$ go get -u github.com/apatters/go-conlog
```

The Go distribution is conlog's only dependency.


Examples
--------

### Basic operation

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/apatters/go-conlog"
)

func main() {
	// Messages are output/not output depending on log level. The
	// log level (in ascending order are:
	//
	// PanicLevel
	// FatalLevel
	// ErrorLevel
	// WarnLevel
	// InfoLevel
	// DebugLevel
	log.SetLevel(log.DebugLevel)
	formatter := log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	log.SetFormatter(formatter)
	// log.Panic("This is a panic message.")
	// log.Fatal("This is a fatal message.")
	log.Error("This is an error message.")
	log.Warn("This is a warning message.")
	log.Warning("This is also a warning message.")
	log.Info("This is an info message.")
	log.Debug("This is a debug message.")

	// The default log level is Level.Info. You can set a
	// different log level using SetLevel.
	log.SetLevel(log.WarnLevel)
	formatter = log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatShort
	log.SetFormatter(formatter)
	log.Info("This message is above the log level so is not output.")
	log.Warning("This message is at the log level so it is output.")
	log.Error("This message is below the log level so it too is output.")

	// There are three forms of output functions for each log
	// level corresponding to the fmt.Print* functions. So the
	// log.Error* output functions have a log.Error(),
	// Log.Errorln, and log.Errorf variations corresponding to
	// fmt.Print(), fmt.Println, and fmt.Printf(). They process
	// parameters in the same way as their fmt counterparts,
	// except that a newline is always output.
	log.SetLevel(log.DebugLevel)
	log.Infoln("Print a number with a newline:", 4)
	log.Info("Print a number also with a newline (note we have to add a space): ", 4)
	log.Infof("Print a formatted number with a newline: %d", 4)

	// Output is send to stderr for log levels of PanicLevel,
	// FatalLevel, and ErrorLevel. Output is sent to stdout for
	// log levels above ErrorLevel. You can change this behavior
	// by setting the Writers in the logger. For
	// example, if we want all output to go to stdout, use:
	log.SetErrorOutput(os.Stdout)
	log.Info("This message is going to stdout.")
	log.Error("This message is now also going to stdout.")

	// We can send output to any Writer. For example, to send
	// output for all levels to a file, we can use:
	logFile, _ := ioutil.TempFile("", "mylogfile-")
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetErrorOutput(logFile)
	log.Info("This message is going to the logfile.")
	log.Error("This message is also going to the logfile.")
	// Dump to stdout so we can see the results.
	contents, _ := ioutil.ReadFile(logFile.Name())
	fmt.Printf("%s", contents)

	// There are a set of Print* output methods that behave
	// exactly like the corresponding fmt.Print* functions except
	// that output goes to the log.Out writer and output can be
	// suppressed using the SetPrintEnabled() method (enabled by
	// default). Note, Print methods are not in any way governed
	// by the log level.
	//
	// Command-line programs should generally use these Print*
	// functions in favor of the corresponding fmt versions if
	// using this logging system.
	//
	// Unlike the log level output methods, the Print and Printf
	// methods do not output a trailing newline.
	log.SetOutput(os.Stdout)       // Restore default Output writer.
	log.SetErrorOutput(os.Stderr)  // Restore default ErrorOutput writer.
	log.Println("Print a number with an implicit newline:", 4.0)
	log.Print("Print a number with an added newline (note we have to add a space): ", 4.0, "\n")
	log.Printf("Print a formatted number with an added newline: %f\n", 4.0)

	// We can suppress Print* output using the SetPrintEnabled
	// method. This feature can be useful for programs that have a
	// --verbose option. You can disable Print* output by default
	// and then enable it when the verbose flag it set.
	log.SetPrintEnabled(false)
	fmt.Printf("Printing enabled: %t\n", log.GetPrintEnabled())
	log.Printf("This print message is suppressed.\n")
	log.SetPrintEnabled(true)
	fmt.Printf("Printing enabled: %t\n", log.GetPrintEnabled())
	log.Printf("Print message are no longer suppressed.\n")
}
```

The above example program outputs:

```
Error This is an error message.
Warning This is a warning message.
Warning This is also a warning message.
Info This is an info message.
Debug This is a debug message.
WARN This message is at the log level so it is output.
ERRO This message is below the log level so it too is output.
INFO Print a number with a newline: 4
INFO Print a number also with a newline (note we have to add a space): 4
INFO Print a formatted number with a newline: 4
INFO This message is going to stdout.
ERRO This message is now also going to stdout.
INFO This message is going to the logfile.
ERRO This message is also going to the logfile.
Print a number with an implicit newline: 4
Print a number with an added newline (note we have to add a space): 4
Print a formatted number with an added newline: 4.000000
Printing enabled: false
Printing enabled: true
Print message are no longer suppressed.
```

### Using the formatter

The formatter is used to modify the output of the message. It can be
used to add a leader to the message which might include various forms
of the the log level, and/or a wall or elapsed time stamp.

```golang
package main

import (
	"time"

	log "github.com/apatters/go-conlog"
)

func main() {
	// Basic logging using the default formatter.
	log.Info("This is an info message without a leader.")
	log.Warning("This is a warning message without a leader.")
	log.Error("This is an error message without a leader.")

	// Basic logging showing log level leaders.
	formatter := log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	log.SetFormatter(formatter)
	log.Info("This is an info message with a leader.")
	log.Warning("This is a warning message with a leader.")
	log.Error("This is an error message with a leader.")

	// The leader can be colorized if going to a tty. If output is
	// going to a file or pipe, no color is used.
	log.SetLevel(log.DebugLevel)
	formatter = log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	formatter.Options.ShowLogLevelColors = true
	log.SetFormatter(formatter)
	log.Debug("Debug messages are blue.")
	log.Info("Info messages are green.")
	log.Warn("Warning messages are yellow.")
	log.Error("Error messages are red.")

	// You can show the traditional logrus log levels using the
	// formatter LogLevelFormatShort option:
	formatter = log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatShort
	log.SetFormatter(formatter)
	log.Debug("Debug message with a short leader.")
	log.Info("Info message with a short leader.")
	log.Warn("Warning message with a short leader.")
	log.Error("Error message with a short leader.")

	// You can show long-form log levels in lower-case using the
	// formatter LogLevelFormatLongLower option:
	formatter = log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongLower
	log.SetFormatter(formatter)
	log.Debug("Debug message with a long, lowercase leader.")
	log.Info("Info message with a long, lowercase leader.")
	log.Warn("Warning message with a long, lowercase leader.")
	log.Error("Error message with a long, lowercase leader.")

	// You can show time stamps in wall clock time with various
	// formats.
	formatter = log.NewStdFormatter()
	formatter.Options.TimestampType = log.TimestampTypeWall
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	log.SetFormatter(formatter)
	log.Info("Info message with wall clock time (default RFC3339 format).")
	formatter = log.NewStdFormatter()
	formatter.Options.TimestampType = log.TimestampTypeWall
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	formatter.Options.WallclockTimestampFmt = time.ANSIC
	log.SetFormatter(formatter)
	log.Info("Info message with wall clock time (ANSIC format).")
	formatter = log.NewStdFormatter()
	formatter.Options.TimestampType = log.TimestampTypeWall
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	formatter.Options.WallclockTimestampFmt = "Jan _2 15:04:05"
	log.SetFormatter(formatter)
	log.Info("Info message with wall clock time (custom format).")

	// You can show time stamps with elapsed time and in various
	// formats.
	formatter = log.NewStdFormatter()
	formatter.Options.TimestampType = log.TimestampTypeElapsed
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	log.SetFormatter(formatter)
	log.Info("Info message with elapsed time (start).")
	time.Sleep(time.Second)
	log.Info("Info message with elapsed time (wait one second).")
	formatter = log.NewStdFormatter()
	formatter.Options.TimestampType = log.TimestampTypeElapsed
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongTitle
	formatter.Options.ElapsedTimestampFmt = "%02d"
	log.SetFormatter(formatter)
	log.Info("Info message with elapsed time with custom format.")
}
```

The above example program outputs:

![Formatter example](images/formatter.png)


### Logging fatal errors and panics

The Fatal* and Panic* methods are used to make logging the output of
fatal conditions consistent with the rest of the conlog logging
system. These methods should be only be used in a "main" package and
never be buried in other packages (you should be returning errors
there instead).

```golang
package main

import (
	log "github.com/apatters/go-conlog"
)

// The *Fatal* routines use a panic/recover mechanism to exit with a
// specific exit code. This mechanism requires that we call the
// HandleExit function as the last function in main() if any Fatal*
// method is used. This is usually best done by creating a "deferred"
// call at the beginning of main() before any other deferred calls.

// We define an exit "wrapper" to be called in main if we want to
// explicitly exit. We do not have to call this sort of function to
// fall-through exit (with exit code 0) out of main(). You should not
// call this sort of function outside the main package.
func exit(code int) {
	panic(log.Exit{code})
}

func main() {
	// The exit routines use a panic/recover mechanism to exit
	// with a specific exit code. We need to call the recovery
	// routine in a defer as the first defer in main() so that it
	// gets called last.
	defer log.HandleExit()

	var err error

	// The Fatalln, Fatal, and Fatalf methods output a message and
	// exit with exit code 1. These calls append a newline to the
	// message.
	if err != nil {
		log.Fatalln("Fatal message exiting with exit code (note we have added a space): ", 1)
	}
	if err != nil {
		log.Fatal("Fatal message exiting with exit code:", 1)
	}
	if err != nil {
		log.Fatalf("Fatal message exiting with exit code %d", 1)
	}

	// The Fatal*WithExitCode methods work like Fatalln, Fatal,
	// and Fatalf except that you can specify the exit code to
	// use.
	if err != nil {
		log.FatallnWithExitCode(2, "Fatal message exiting with exit code (note we have added a space): ", 2)
	}
	if err != nil {
		log.FatalWithExitCode(2, "Fatal message exiting with exit code:", 2)
	}
	if err != nil {
		log.FatalfWithExitCode(2, "Fatal message exiting with exit code %d", 2)
	}

	// The Fatal*IfError methods output a fatal message and exit
	// with a specified exit code if the the err parameter is not
	// nil. They are used as a short-cut for the common:
	//
	//  if err != nil {
	//      log.FatalWithExitCode(...)
	//  }
	//
	// golang programming construct.
	log.FatallnIfError(err, 3, "Fatal message if err != nil exiting with exit code: ", 3)
	log.FatalIfError(err, 3, "Fatal message if err != nil exiting with exit code:", 3)
	log.FatalfIfError(err, 3, "Fatal message if err != nil exiting with exit code %d, err = %s", 3, err)

	// The Panicln, Panic, and Panicf methods output a panic level
	// message and then call the standard builtin panic() function
	// with the message. The panic message goes to the the ErrorOutput
	// Writer, but the actual panic output will always go to
	// stderr.
	var impossibleCond bool
	if impossibleCond {
		log.Panicln("Panic message to log with panic output to stderr: ", impossibleCond)
	}
	if impossibleCond {
		log.Panic("Panic message to log with panic output to stderr:", impossibleCond)
	}
	if impossibleCond {
		log.Panicf("Panic message to log with panic output to stderr: %t", impossibleCond)
	}

	// We can exit explicitly with a specified exit code using the
	// exit() function we defined earlier. We don't need to use
	// this call to exit with exit code 0 when falling-through the
	// end of main.
	exit(5)
}
```

### Logging to multiple logs

You can log to any number of logs with one call using the Loggers
object.


```golang
package main

// This example illustrates how to use the Loggers object to send
// output to multiple logs. The Loggers object is typically used in a
// command-line program to send output to both the TTY/console and to
// a log file with one call.

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/apatters/go-conlog"
)

func main() {
	// Initialize the log going to the TTY/console.
	ttyLog := log.NewLogger()
	ttyLog.SetLevel(log.InfoLevel)
	formatter := log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatLongLower
	formatter.Options.ShowLogLevelColors = true
	ttyLog.SetFormatter(formatter)

	// Initialize the log going to a file.
	logFile, _ := ioutil.TempFile("", "mylogfile-")
	defer os.Remove(logFile.Name()) // You normally wouldn't delete the file.
	fileLog := log.NewLogger()
	fileLog.SetLevel(log.DebugLevel)
	fileLog.SetOutput(logFile)
	fileLog.SetErrorOutput(logFile)
	formatter = log.NewStdFormatter()
	formatter.Options.LogLevelFmt = log.LogLevelFormatShort
	formatter.Options.TimestampType = log.TimestampTypeWall
	formatter.Options.WallclockTimestampFmt = "15:04:05"
	fileLog.SetFormatter(formatter)

	// Initialize the multi-logger. We will use this one when we
	// want output to go to both logs.
	bothLogs := log.NewLoggers(ttyLog, fileLog)

	// We can send messages to individual logs or to both logs.
	ttyLog.Info("This info message only goes to the TTY.")
	fileLog.Info("This info message only goes to the log file.")
	bothLogs.Info("This info message goes to both the TTY and the log file.")

	// Individual log levels are honored.
	bothLogs.Info("This message goes to both logs because they both have log levels >= InfoLevel.")
	bothLogs.Debug("This message only goes the log file because its log level is DebugLevel.")

	// We can use the logs.SetLevel method to set the log level for all logs.
	bothLogs.SetLevel(log.DebugLevel)
	bothLogs.Debug("This debug message goes to both the TTY and the log file now that they are DebugLevel.")

	// We can enable/disable Print* methods for all logs using SetPrintEnabled.
	bothLogs.SetPrintEnabled(false)
	bothLogs.Print("This message is suppressed on both the TTY and the log file.")
	bothLogs.SetPrintEnabled(true)
	bothLogs.Println("This message goes to both TTY and the log file now that Print is re-enabled.")

	// We can use Fatal* and Panic* methods, but the messages will
	// only go the first log as the program will terminate before
	// getting to subsequent logs.
	var err error
	if err != nil {
		bothLogs.Fatal("This fatal message only goes to the TTY.")
	}
	var impossibleCond bool
	if impossibleCond {
		bothLogs.Panic("This panic message only goes to the TTY. The panic() output always goes to stderr.")
	}

	// Dump the log file to stdout so we can see the results.
	contents, _ := ioutil.ReadFile(logFile.Name())
	defer logFile.Close()
	fmt.Printf("%s\n", contents)
}
```

The above example program outputs:

```
info This info message only goes to the TTY.
info This info message goes to both the TTY and the log file.
info This message goes to both logs because they both have log levels >= InfoLevel.
debug This debug message goes to both the TTY and the log file now that they are DebugLevel.
This message goes to both TTY and the log file now that Print is re-enabled.
INFO[20:41:16] This info message only goes to the log file.
INFO[20:41:16] This info message goes to both the TTY and the log file.
INFO[20:41:16] This message goes to both logs because they both have log levels >= InfoLevel.
DEBU[20:41:16] This message only goes the log file because its log level is DebugLevel.
DEBU[20:41:16] This debug message goes to both the TTY and the log file now that they are DebugLevel.
This message goes to both TTY and the log file now that Print is re-enabled.
```

License
-------

The go-conlog package is available under the [MITLicense](https://mit-license.org/).


Thanks
------

Thanks to [Secure64](https://secure64.com/company/) for
contributing this code.




