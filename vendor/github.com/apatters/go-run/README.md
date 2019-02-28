go-run
===========

Run is a Go (golang) package that wraps the standard Go os/exec and
golang.org/x/crypto/ssh packages to run commands either locally or
over ssh while capturing stdout, stderr, and exit codes.

[![GoDoc](https://godoc.org/github.com/apatters/go-run?status.svg)](https://godoc.org/github.com/apatters/go-run)


Features
-------

* Run commands either locally or remotely over SSH.
* Run commands in a shell or directly ala glibc's exec().
* Capture stdout, stderr, and exit code.
* Output can be redirected to any Writer.

Documentation
-------------

Documentation can be found at [GoDoc](https://godoc.org/github.com/apatters/go-run)


Installation
------------

Install wordwrap using the "go get" command:

```bash
$ go get github.com/apatters/go-run
```

The Go distribution is run's only dependency.


Examples
--------

### Local

Local is used to run commands on the local host.

``` golang
package main

import (
	"fmt"

	"github.com/apatters/go-run"
)

func main() {
	// Initialize Local object using defaults.
	runner := run.NewLocal(run.LocalConfig{})

	fmt.Println("Run ls command.")
	stdout, stderr, code, _ := runner.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command with an expected error.")
	stdout, stderr, code, _ = runner.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command after changing directory.")
	runner = run.NewLocal(run.LocalConfig{Dir: "/bin"})
	stdout, stderr, code, _ = runner.Run(
		"/bin/ls",
		"-1",
		"true",
		"false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command using shell.")
	runner = run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, _ = runner.Shell("/bin/ls -1 /bin/true /bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command using shell with an expected error.")
	stdout, stderr, code, _ = runner.Shell("/bin/ls -1 /bin/true /bin/false /xyzzy")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run complex shell command.")
	runner = run.NewLocal(run.LocalConfig{})
	stdout, stderr, code, _ = runner.Shell("cd /bin && /bin/ls -1 true false | head -n 1")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()
}
```

Output:

```
Run ls command.
stdout = "/bin/false\n/bin/true\n"
stderr = ""
exit code = 0

Run ls command with an expected error.
stdout = "/bin/false\n/bin/true\n"
stderr = "/bin/ls: cannot access /xyzzy: No such file or directory\n"
exit code = 2

Run ls command after changing directory.
stdout = "false\ntrue\n"
stderr = ""
exit code = 0

Run ls command using shell.
stdout = "/bin/false\n/bin/true\n"
stderr = ""
exit code = 0

Run ls command using shell with an expected error.
stdout = "/bin/false\n/bin/true\n"
stderr = "/bin/ls: cannot access /xyzzy: No such file or directory\n"
exit code = 2

Run complex shell command.
stdout = "false\n"
stderr = ""
exit code = 0
```

### Remote

Remote is used to run commands on remote hosts using SSH. It defaults
to using the current user name and the user's public SSH key for
authentication. Ssh-agent or something similar must be used to provide
the pass-phrase if the key is pass-phrase protected.

```golang
package main

import (
	"fmt"

	"github.com/apatters/go-run"
)

func main() {
	// Initialize Remote object using defaults.
	runner, _ := run.NewRemote(run.RemoteConfig{
		Credentials: run.Credentials{
			Hostname: "localhost"},
	})

	fmt.Println("Run ls command.")
	stdout, stderr, code, _ := runner.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command with an expected error.")
	stdout, stderr, code, _ = runner.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false",
		"/xyzzy")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	runner, _ = run.NewRemote(run.RemoteConfig{})

	fmt.Println("Run ls command using shell.")
	stdout, stderr, code, _ = runner.Shell("/bin/ls -1 /bin/true /bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run ls command using shell with an expected error.")
	stdout, stderr, code, _ = runner.Shell("/bin/ls -1 /bin/true /bin/false /xyzzy")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run complex shell command.")
	runner, _ = run.NewRemote(run.RemoteConfig{})
	stdout, stderr, code, _ = runner.Shell("cd /bin && /bin/ls -1 true false | head -n 1")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()
}
```

Output:

```
Run ls command.
stdout = "/bin/false\n/bin/true\n"
stderr = ""
exit code = 0

Run ls command with an expected error.
stdout = "/bin/false\n/bin/true\n"
stderr = "/bin/ls: cannot access '/xyzzy': No such file or directory\n"
exit code = 2

Run ls command using shell.
stdout = "/bin/false\n/bin/true\n"
stderr = ""
exit code = 0

Run ls command using shell with an expected error.
stdout = "/bin/false\n/bin/true\n"
stderr = "/bin/ls: cannot access '/xyzzy': No such file or directory\n"
exit code = 2

Run complex shell command.
stdout = "false\n"
stderr = ""
exit code = 0
```

### Standard runner

The standard runner can be used to run local commands (only) if you do
not need to use a customized constructor saving the extra line or two
of code needed to call the constructor.

``` golang
package main

import (
	"fmt"

	"github.com/apatters/go-run"
)

func main() {
	// There is no need for a constructor when running local
	// commands when using the standard runner.
	fmt.Println("Run ls command.")
	stdout, stderr, code, _ := run.Run(
		"/bin/ls",
		"-1",
		"/bin/true",
		"/bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	fmt.Println("Run id command.")
	stdout, stderr, code, _ = run.Run("/bin/id", "root")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	stdout, stderr, code, _ = run.Shell("/bin/ls -1 /bin/true /bin/false")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()

	// Print out the command to run, then run it.
	fmt.Println("Run id command.")
	stdout, stderr, code, _ = run.Shell("/bin/id root | cut -f1 -d ' '")
	fmt.Printf("stdout = %q\n", stdout)
	fmt.Printf("stderr = %q\n", stderr)
	fmt.Printf("exit code = %d\n", code)
	fmt.Println()
}
```

Output:

```
Run ls command.
stdout = "/bin/false\n/bin/true\n"
stderr = ""
exit code = 0

Run id command.
stdout = "uid=0(root) gid=0(root) groups=0(root)\n"
stderr = ""
exit code = 0

stdout = "/bin/false\n/bin/true\n"
stderr = ""
exit code = 0

Run id command.
stdout = "uid=0(root)\n"
stderr = ""
exit code = 0
```

License
-------

The go-run package is available under the [MITLicense](https://mit-license.org/).


Thanks
------

Thanks to [Secure64](https://secure64.com/company/) for
contributing this code.




