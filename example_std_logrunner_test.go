package logrun_test

// This example illustrates how to use the standard log runner to run
// local commands without using a constructor. It is is useful when
// you do not need to customize behavior.

import (
	"fmt"
	"os"

	"github.com/apatters/go-logrun"
)

// Custom logging function.
func MyLogger(v ...interface{}) {
	fmt.Printf("Command: ")
	fmt.Println(v...)
}

func Example() {
	// Set our custom logger.
	logrun.SetLogFunc(MyLogger)
	defer logrun.SetLogFunc(logrun.DefaultLogFunc)

	fmt.Println("Execute a command using Run().")
	stdout, stderr, code := logrun.Run("/usr/bin/seq", "1", "3")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("code = %d\n", code)
	fmt.Println()

	fmt.Println("Execute a command using Shell().")
	stdout, stderr, code = logrun.Shell("seq 1 3 | grep 2")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("code = %d\n", code)
	fmt.Println()

	fmt.Println("See if a file exists.")
	exists, err := logrun.FileExists("/bin/true")
	if err != nil {
		// report error.
	}
	fmt.Printf("/bin/true exists: %t\n", exists)
	exists, err = logrun.FileExists("/bin/xyzzy")
	if err != nil {
		// report error.
	}
	fmt.Printf("/bin/xyzzy exists: %t\n", exists)
	fmt.Println()

	fmt.Println("See if a directory exists.")
	exists, err = logrun.DirExists("/etc")
	if err != nil {
		// report error.
	}
	fmt.Printf("/bin/etc exists: %t\n", exists)
	exists, err = logrun.FileExists("/xyzzy")
	if err != nil {
		// report error.
	}
	fmt.Printf("/xyzzy exists: %t\n", exists)
	fmt.Println()

	fmt.Println("List files using a shell glob pattern.")
	paths, err := logrun.Glob("/etc/passwd*")
	if err != nil {
		// Report the error.
	}
	for _, path := range paths {
		fmt.Println(path)
	}
	fmt.Println()

	fmt.Println("Copy the contents of a remote directory to a local temporary directory.")
	destDir := "/tmp/go-logrun-XXXXX"
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		// Report the error
		return
	}
	defer os.RemoveAll(destDir)
	err = logrun.Rsync("localhost:/etc/cron.daily/", fmt.Sprintf("%s/", destDir))
	if err != nil {
		// Report the error.
	}
	fmt.Println()

	fmt.Println("Log commands but do not execute them.")
	logrun.SetDryrun(true)
	stdout, stderr, code = logrun.Run("/usr/bin/seq", "1", "3")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("code = %d\n", code)
	stdout, stderr, code = logrun.Run("/bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("code = %d\n", code)
	fmt.Println()

	fmt.Println("Disable logging")
	logrun.SetLogFunc(logrun.DiscardLogFunc)
	stdout, stderr, code = logrun.Run("/usr/bin/seq", "1", "3")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("code = %d\n", code)

	// Output:
	// Execute a command using Run().
	// Command: /usr/bin/seq 1 3
	// stdout = "1\n2\n3\n"
	// stderr = ""
	// code = 0
	//
	// Execute a command using Shell().
	// Command: /bin/sh -c "seq 1 3 | grep 2"
	// stdout = "2\n"
	// stderr = ""
	// code = 0
	//
	// See if a file exists.
	// Command: /usr/bin/stat --dereference --format %n:%F /bin/true
	// /bin/true exists: true
	// Command: /usr/bin/stat --dereference --format %n:%F /bin/xyzzy
	// /bin/xyzzy exists: false
	//
	// See if a directory exists.
	// Command: /usr/bin/stat --dereference --format %n:%F /etc
	// /bin/etc exists: true
	// Command: /usr/bin/stat --dereference --format %n:%F /xyzzy
	// /xyzzy exists: false
	//
	// List files using a shell glob pattern.
	// Command: /bin/sh -c "/bin/ls -1 --directory /etc/passwd*"
	// /etc/passwd
	// /etc/passwd-
	//
	// Copy the contents of a remote directory to a local temporary directory.
	// Command: /usr/bin/rsync --rsh ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o GlobalKnownHostsFile=/dev/null --recursive --links --times localhost:/etc/cron.daily/ /tmp/go-logrun-XXXXX/
	//
	// Log commands but do not execute them.
	// Command: /usr/bin/seq 1 3
	// stdout = ""
	// stderr = ""
	// code = 0
	// Command: /bin/false
	// stdout = ""
	// stderr = ""
	// code = 0
	//
	// Disable logging
	// stdout = ""
	// stderr = ""
	// code = 0
}
