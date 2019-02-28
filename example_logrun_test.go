package logrun_test

import (
	"fmt"
	"os"

	"github.com/apatters/go-conlog"
	"github.com/apatters/go-logrun"
)

func ExampleLogRun_Run() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println(`Run the "seq 1 3" command locally.`)
	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	stdout, stderr, code := runner.Run("seq", "1", "3")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)
	fmt.Println()

	fmt.Println(`Run the "seq 1 3" command remotely.`)
	runner, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: conlog.Debug,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report the error.
	}
	stdout, stderr, code = runner.Run("seq", "1", "3")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)

	// Output:
	// Run the "seq 1 3" command locally.
	// Debug seq 1 3
	// Stdout = "1\n2\n3\n"
	// Stderr = ""
	// Exit code = 0
	//
	// Run the "seq 1 3" command remotely.
	// Debug ssh buildman@localhost seq 1 3
	// Stdout = "1\n2\n3\n"
	// Stderr = ""
	// Exit code = 0
}

func ExampleLogRun_FormatRun() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println("Get the log output of local command using Run().")
	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	cmdLog := runner.FormatRun("seq", "1", "3")
	fmt.Printf("cmdLog = %q\n", cmdLog)
	fmt.Println()

	fmt.Println("Get the log output of remote command using Run().")
	runner, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: conlog.Debug,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report error
	}
	cmdLog = runner.FormatRun("seq", "1", "3")
	fmt.Printf("cmdLog = %q\n", cmdLog)

	// Output:
	// Get the log output of local command using Run().
	// cmdLog = "seq 1 3"
	//
	// Get the log output of remote command using Run().
	// cmdLog = "ssh buildman@localhost seq 1 3"
}

func ExampleLogRun_Shell() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println(`Run the "seq 1 3 | grep 2" command locally.`)
	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	stdout, stderr, code := runner.Shell("seq 1 3 | grep 2")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)
	fmt.Println()

	fmt.Println(`Run the "seq 1 3 | grep 2" command remotely.`)
	runner, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: conlog.Debug,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report the error.
	}
	stdout, stderr, code = runner.Run("seq 1 3 | grep 2")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)

	// Output:
	// Run the "seq 1 3 | grep 2" command locally.
	// Debug /bin/sh -c "seq 1 3 | grep 2"
	// Stdout = "2\n"
	// Stderr = ""
	// Exit code = 0
	//
	// Run the "seq 1 3 | grep 2" command remotely.
	// Debug ssh buildman@localhost seq 1 3 | grep 2
	// Stdout = "2\n"
	// Stderr = ""
	// Exit code = 0
}

func ExampleLogRun_FormatShell() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println("Get the log output of a local command using Shell().")
	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	cmdLog := runner.FormatShell("seq 1 3")
	fmt.Printf("cmdLog = %q\n", cmdLog)
	fmt.Println()

	fmt.Println("Get the log output of a remote command using Shell().")
	runner, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: conlog.Debug,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report error
	}
	cmdLog = runner.FormatShell("seq 1 3")
	fmt.Printf("cmdLog = %q\n", cmdLog)

	// Output:
	// Get the log output of a local command using Shell().
	// cmdLog = "/bin/sh -c \"seq 1 3\""
	//
	// Get the log output of a remote command using Shell().
	// cmdLog = "ssh buildman@localhost /bin/sh -c \"seq 1 3\""
}

func ExampleLogRun_FileExists() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println(`See if /bin/true exists on the local host.`)
	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	exists, err := runner.FileExists("/bin/true")
	if err != nil {
		// Report the error.
	}
	fmt.Printf("exists = %t\n", exists)
	fmt.Println()

	fmt.Println(`See if /bin/true exists on the remote host.`)
	runner, err = logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: conlog.Debug,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report the error.
	}
	exists, err = runner.FileExists("/bin/true")
	if err != nil {
		// Report the error.
	}
	fmt.Printf("exists = %t\n", exists)

	// Output:
	// See if /bin/true exists on the local host.
	// Debug /usr/bin/stat --dereference --format %n:%F /bin/true
	// exists = true
	//
	// See if /bin/true exists on the remote host.
	// Debug ssh buildman@localhost /usr/bin/stat --dereference --format %n:%F /bin/true
	// exists = true
}

func ExampleLogRun_DirExists() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println(`See if /etc exists on the local host.`)
	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	exists, err := runner.DirExists("/etc")
	if err != nil {
		// Report the error.
	}
	fmt.Printf("exists = %t\n", exists)
	fmt.Println()

	fmt.Println(`See if /etc exists on the remote host.`)
	runner, err = logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: conlog.Debug,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report the error.
	}
	exists, err = runner.DirExists("/etc")
	if err != nil {
		// Report the error.
	}
	fmt.Printf("exists = %t\n", exists)

	// Output:
	// See if /etc exists on the local host.
	// Debug /usr/bin/stat --dereference --format %n:%F /etc
	// exists = true
	//
	// See if /etc exists on the remote host.
	// Debug ssh buildman@localhost /usr/bin/stat --dereference --format %n:%F /etc
	// exists = true
}

func ExampleLogRun_Glob() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println(`Glob local passwd files`)
	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	paths, err := runner.Glob("/etc/passwd*")
	if err != nil {
		// Report the error.
	}
	for _, path := range paths {
		fmt.Println(path)
	}
	fmt.Println()

	fmt.Println(`Glob remote passwd files.`)
	runner, err = logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: conlog.Debug,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report the error.
	}
	paths, err = runner.Glob("/etc/passwd*")
	if err != nil {
		// Report the error.
	}
	for _, path := range paths {
		fmt.Println(path)
	}
	fmt.Println()

	// Output:
	// Glob local passwd files
	// Debug /bin/sh -c "/bin/ls -1 --directory /etc/passwd*"
	// /etc/passwd
	// /etc/passwd-
	//
	// Glob remote passwd files.
	// Debug ssh buildman@localhost /bin/sh -c "/bin/ls -1 --directory /etc/passwd*"
	// /etc/passwd
	// /etc/passwd-
}

func ExampleLogRun_Rsync() {
	// Set up logger to output at debug level and add a log level
	// prefix using title case.
	formatter := conlog.NewStdFormatter()
	formatter.Options.LogLevelFmt = conlog.LogLevelFormatLongTitle
	conlog.SetFormatter(formatter)
	conlog.SetLevel(conlog.DebugLevel)

	fmt.Println(`Copy the contents of directory on a remote host to local temporary directory.`)
	destDir := "/tmp/go-logrun-XXXXX"
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		// Report the error
		return
	}
	defer os.RemoveAll(destDir)

	runner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: conlog.Debug,
	})
	err = runner.Rsync("localhost:/etc/cron.daily/", fmt.Sprintf("%s/", destDir))
	if err != nil {
		// Report the error.
	}

	// Output:
	// Copy the contents of directory on a remote host to local temporary directory.
	// Debug /usr/bin/rsync --rsh ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o GlobalKnownHostsFile=/dev/null --recursive --links --times localhost:/etc/cron.daily/ /tmp/go-logrun-XXXXX/
}
