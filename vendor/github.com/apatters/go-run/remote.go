// Copyright 2019 Secure64 Software Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package run

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	defaultSSHPort        = 22
	defaultSSHHostname    = "localhost"
	defaultSSHKeyfileName = "id_rsa"
)

// Credentials contains needed credentials to SSH to a host. It can
// use either a password or SSH private key.
type Credentials struct {
	// Hostname is either the hostname or IP of the remote host.
	Hostname string

	// Port is the port used to connect to the ssh server on the
	// remote host.
	Port int

	// Username is the account name used to authenticate on the
	// remote host.
	Username string

	// Password is password used to authenticate on the remote
	// host. Not needed if using PrivateKeyFilename.
	Password string

	// PrivateKeyFilename is the full path the SSH private key
	// used to authenticate with the remote host.  Not used if
	// Password is specified. You must use ssh-agent or something
	// similar to provide the passphrase if the key is passphrase
	// protected.
	PrivateKeyFilename string
}

func defaultPrivateKeyFilename(username string) (string, error) {
	user, err := user.Lookup(username)
	if err != nil {
		return "", err
	}

	return filepath.Join(user.HomeDir, ".ssh", defaultSSHKeyfileName), nil
}

// RemoteConfig contains configuration data used in the Remote
// constructor.
type RemoteConfig struct {
	// ShellExecutable is the path to the shell on the remote host
	// used to run commands in Shell() methods. See Remote for
	// details.
	ShellExecutable string

	// Stdin specifies the process's standard input. See Remote for
	// details.
	Stdin io.Reader

	// Stdout specifies the process's standard output. See Remote
	// for details.
	Stdout io.Writer

	// Stderr specifies the process's standard error. See Remote
	// for details.
	Stderr io.Writer

	// Credentials used to authenticate on the remote system.
	Credentials Credentials
}

// Remote wraps ssh.Client to make running commands over SSH on a
// remote host relatively as easy as when running them in shell
// script.
type Remote struct {
	// ShellExecutable is the full path to the shell on the remote
	// host to be run when executing shell commands.
	ShellExecutable string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, the descriptor's output is captured and
	// is returned in the Run and Shell functions.
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// Credentials are used to authenticate with the remote host.
	Credentials Credentials

	sshSession *ssh.Session
}

// NewRemote is the constructor for Remote. It takes a RemoteConfig
// object to configure it. The following configuration options are set
// if the default RemoteConfig constructor, RemoteConfig{}, is used:
//
//     ShellExecutable = DefaultShellExecutable
//     Stdin = nil  // Discard stdin.
//     Stdout = nil // Capture stdout.
//     Stderr = nil // Capture stderr,
//     Credentials.Hostname = "localhost"
//     Credentials.Port = 22
//     Credentials.Username = Current user
//     Credentials.Password = ""
//     Credentials.PrivateKeyFilename = Current users default private RSA
//     keyfile ($HOME/.ssh/id_rsa) if present.
func NewRemote(config RemoteConfig) (*Remote, error) {
	r := new(Remote)
	if len(config.ShellExecutable) == 0 {
		r.ShellExecutable = DefaultShellExecutable
	} else {
		r.ShellExecutable = config.ShellExecutable
	}
	r.Stdin = config.Stdin
	r.Stdout = config.Stdout
	r.Stderr = config.Stderr
	r.Credentials = config.Credentials
	if r.Credentials.Hostname == "" {
		r.Credentials.Hostname = defaultSSHHostname
	}
	if r.Credentials.Port == 0 {
		r.Credentials.Port = defaultSSHPort
	}
	if r.Credentials.Username == "" {
		user, err := user.Current()
		if err != nil {
			return nil, err
		}
		r.Credentials.Username = user.Username
	}
	if r.Credentials.Password == "" && r.Credentials.PrivateKeyFilename == "" {
		keyFilename, err := defaultPrivateKeyFilename(r.Credentials.Username)
		if err != nil {
			return nil, err
		}
		r.Credentials.PrivateKeyFilename = keyFilename
	}

	return r, nil
}

func (r *Remote) getSSHAuths() ([]ssh.AuthMethod, error) {
	var auths []ssh.AuthMethod
	if r.Credentials.Password != "" {
		auths = []ssh.AuthMethod{ssh.Password(r.Credentials.Password)}
	} else {
		sshAuthSockEnv := os.Getenv("SSH_AUTH_SOCK")
		if sshAuthSockEnv != "" {
			sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
			if err != nil {
				return nil, err
			}
			agent := agent.NewClient(sock)
			signers, err := agent.Signers()
			if err != nil {
				return nil, err
			}
			auths = []ssh.AuthMethod{ssh.PublicKeys(signers...)}

			return auths, nil
		}
		keyBuf, err := ioutil.ReadFile(r.Credentials.PrivateKeyFilename)
		if err != nil {
			return nil, fmt.Errorf(
				"run: could not read private key file '%s': %s",
				r.Credentials.PrivateKeyFilename,
				err)
		}
		key, err := ssh.ParsePrivateKey(keyBuf)
		if err != nil {
			return nil, fmt.Errorf(
				"run: could not use private key file '%s': %s",
				r.Credentials.PrivateKeyFilename,
				err)
		}
		auths = []ssh.AuthMethod{ssh.PublicKeys(key)}
	}

	return auths, nil
}

func (r *Remote) open() error {
	auths, err := r.getSSHAuths()
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User:            r.Credentials.Username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint: gosec
	}
	client, err := ssh.Dial("tcp",
		fmt.Sprintf("%s:%d", r.Credentials.Hostname, r.Credentials.Port),
		config)
	if err != nil {
		return fmt.Errorf("run: connection to %s@%s failed: %s",
			r.Credentials.Username,
			r.Credentials.Hostname,
			err)
	}
	r.sshSession, err = client.NewSession()
	if err != nil {
		return err
	}

	return nil
}

func (r *Remote) close() error {
	if r.sshSession != nil {
		err := r.sshSession.Close()
		r.sshSession = nil
		return err
	}

	return nil
}

func (r *Remote) exec(args ...string) (string, string, int, error) {
	err := r.open()
	if err != nil {
		return "", "", 0, err
	}
	defer r.close() // nolint
	if r.sshSession == nil {
		panic("Session == nil")
	}

	// Hook up standard files.
	r.sshSession.Stdin = r.Stdin
	var stdoutPipe io.Reader
	if r.Stdout == nil {
		stdoutPipe, err = r.sshSession.StdoutPipe()
		if err != nil {
			return "", "", 0, err
		}
	} else {
		r.sshSession.Stdout = r.Stdout
	}
	var stderrPipe io.Reader
	if r.Stderr == nil {
		stderrPipe, err = r.sshSession.StderrPipe()
		if err != nil {
			return "", "", 0, err
		}
	} else {
		r.sshSession.Stderr = r.Stderr
	}

	code := 0
	cmdLine := strings.Join(args, " ")
	err = r.sshSession.Run(cmdLine)
	if err != nil {
		switch err.(type) {
		case *ssh.ExitError:
			// Extract exit code from error message.
			re := regexp.MustCompile("^Process exited with status ([0-9]+)$")
			match := re.FindStringSubmatch(err.Error())
			if match != nil {
				code, err = strconv.Atoi(match[1])
				if err != nil {
					return "", "", 0, err
				}
			}
		default:
			return "", "", 0, err
		}
	}

	// Process the I/O.
	var stdoutBuf []byte
	if r.Stdout == nil {
		stdoutBuf, err = ioutil.ReadAll(stdoutPipe)
		if err != nil {
			return "", "", 0, err
		}
	}
	var stderrBuf []byte
	if r.Stderr == nil {
		stderrBuf, err = ioutil.ReadAll(stderrPipe)
		if err != nil {
			return "", "", 0, err
		}
	}

	return string(stdoutBuf), string(stderrBuf), code, err
}

// Run runs a command like glibc's exec() call. It returns the
// standard out, standard error, and exit code of the command when it
// completes.
func (r *Remote) Run(cmd string, args ...string) (string, string, int, error) {
	cmdLine := cmd + " " + strings.Join(args, " ")
	stdout, stderr, code, err := r.exec(cmdLine)

	return stdout, stderr, code, err
}

// FormatRun returns a string representation of the what command would
// be run using Run(). Useful for logging commands.
func (r *Remote) FormatRun(cmd string, args ...string) string {
	s := fmt.Sprintf(`ssh %s@%s %s %s`,
		r.Credentials.Username,
		r.Credentials.Hostname,
		cmd,
		strings.Join(args, " "))

	return strings.TrimSpace(s)
}

// Shell runs a command in a shell. The command is passed to the shell
// as the -c option, so just about any shell code that can be used on
// the command-line will be passed to it. It returns the standard out,
// standard error, and exit code of the command when it completes.
func (r *Remote) Shell(cmd string) (string, string, int, error) {
	cmdLine := fmt.Sprintf(`%s -c "%s"`, r.ShellExecutable, cmd)
	stdout, stderr, code, err := r.exec(cmdLine)

	return stdout, stderr, code, err
}

// FormatShell returns a string representation of the what command
// would be run using Shell().  Useful for logging commands.
func (r *Remote) FormatShell(cmd string) string {
	s := fmt.Sprintf(`ssh %s@%s %s -c "%s"`,
		r.Credentials.Username,
		r.Credentials.Hostname,
		r.ShellExecutable,
		cmd)

	return strings.TrimSpace(s)
}
