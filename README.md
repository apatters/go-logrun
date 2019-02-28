# Logrun

Logrun is a Go (golang) package that wraps the standard Go os/exec
package and golang.org/x/crypto/ssh package to run commands either
locally or over ssh while capturing stdout, stderr, exit
codes. Commands are optionally logged using a specified logging
function.

[![GoDoc](https://godoc.org/github.com/apatters/go-run?status.svg)](https://godoc.org/github.com/apatters/go-logrun)


Features
-------

* Run commands either locally or remotely over SSH.
* Run commands in a shell or directly ala glibc's exec().
* Capture stdout, stderr, and exit code.
* Output can be redirected to any Writer.
* Commands are logged using a specified logging function.

Documentation
-------------

Documentation can be found at [GoDoc](https://godoc.org/github.com/apatters/go-logrun)


Installation
------------

Install go-logrun using the "go get" command:

```bash
$ go get github.com/apatters/go-logrun
```

Example
-------

``` golang
package main

import (
	"fmt"
	"log"

	"github.com/apatters/go-logrun"
)

func main() {
	locRunner := logrun.NewLocalLogRun(logrun.LocalConfig{
		LogFunc: log.Println,
	})
	rmtRunner, err := logrun.NewRemoteLogRun(logrun.RemoteConfig{
		LogFunc: log.Println,
		// We are faking credentials here to run in a test
		// environment.
		Credentials: logrun.Credentials{
			Hostname: "localhost",
		},
	})
	if err != nil {
		// Report error.
	}

	fmt.Println(`Run the "seq 1 3" command locally.`)
	stdout, stderr, code := locRunner.Run("seq", "1", "3")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)
	fmt.Println()

	fmt.Println(`Run the "seq 1 3" command remotely.`)
	stdout, stderr, code = rmtRunner.Run("seq", "1", "3")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)
	fmt.Println()

	fmt.Println(`Run the "seq 1 3 | grep 2" command locally.`)
	stdout, stderr, code = locRunner.Shell("seq 1 3 | grep 2")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)
	fmt.Println()

	fmt.Println(`Run the "seq 1 3 | grep 2" command remotely.`)
	stdout, stderr, code = rmtRunner.Shell("seq 1 3 | grep 2")
	fmt.Printf("Stdout = %q\n", stdout)
	fmt.Printf("Stderr = %q\n", stderr)
	fmt.Printf("Exit code = %d\n", code)
}
```

Output:
```
Run the "seq 1 3" command locally.
2019/02/28 14:12:10 seq 1 3
Stdout = "1\n2\n3\n"
Stderr = ""
Exit code = 0

Run the "seq 1 3" command remotely.
2019/02/28 14:12:10 ssh andrew@localhost seq 1 3
Stdout = "1\n2\n3\n"
Stderr = ""
Exit code = 0

Run the "seq 1 3 | grep 2" command locally.
2019/02/28 14:12:10 /bin/sh -c "seq 1 3 | grep 2"
Stdout = "2\n"
Stderr = ""
Exit code = 0

Run the "seq 1 3 | grep 2" command remotely.
2019/02/28 14:12:10 ssh andrew@localhost /bin/sh -c "seq 1 3 | grep 2"
Stdout = "2\n"
Stderr = ""
Exit code = 0
```

License
-------

The go-logrun package is available under the [MITLicense](https://mit-license.org/).

Thanks
------

Thanks to [Secure64](https://secure64.com/company/) for
contributing this code.
