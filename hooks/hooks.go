// Package hooks provides functions to invoke a directory of executable hooks,
// used to provide arbitrary handling of significant events.
package hooks

import (
	"fmt"
	deos "github.com/hlandau/degoutils/os"
	"github.com/hlandau/xlog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Log site.
var log, Log = xlog.New("acme.hooks")

// The default hook path is the path at which executable hooks are looked for.
// On POSIX-like systems, this is usually "/usr/lib/acme/hooks" (or
// "/usr/libexec/acme/hooks" if /usr/libexec exists).
var DefaultPath string

func init() {
	// Allow overriding at build time.
	p := DefaultPath
	if p == "" {
		p = "/usr/lib/acme/hooks"
	}

	if _, err := os.Stat("/usr/libexec"); strings.HasPrefix(p, "/usr/lib/") && err == nil {
		p = "/usr/libexec" + p[8:]
	}

	DefaultPath = p
}

// Notifies hook programs that a live symlink has been updated.
//
// If hookDirectory is "", DefaultHookPath is used. stateDirectory and
// hostnames are passed as information to the hooks.
func NotifyLiveUpdated(hookDirectory, stateDirectory string, hostnames []string) error {
	if len(hostnames) == 0 {
		return nil
	}

	hostnameList := strings.Join(hostnames, "\n") + "\n"
	_, err := runParts(hookDirectory, stateDirectory, []byte(hostnameList), "live-updated")
	if err != nil {
		return err
	}

	return nil
}

// Invokes HTTP challenge start hooks.
//
// installed indicates whether at least one hook script indicated success. err
// could still be returned in this case if an error occurs while executing some
// other hook.
func ChallengeHTTPStart(hookDirectory, stateDirectory, hostname, targetFileName, token, ka string) (installed bool, err error) {
	return runParts(hookDirectory, stateDirectory, []byte(ka),
		"challenge-http-start", hostname, targetFileName, token)
}

func ChallengeHTTPStop(hookDirectory, stateDirectory, hostname, targetFileName, token, ka string) error {
	_, err := runParts(hookDirectory, stateDirectory, []byte(ka),
		"challenge-http-stop", hostname, targetFileName, token)
	return err
}

// Implements functionality similar to the "run-parts" command on many distros.
// Implementations vary, so it is reimplemented here.
func runParts(directory, stateDirectory string, stdinData []byte, args ...string) (anySucceeded bool, err error) {
	if directory == "" {
		directory = DefaultPath
	}

	fi, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			// Not an error if the directory doesn't exist; nothing to do.
			return false, nil
		}

		return false, err
	}

	// Probably shouldn't propagate this to all child processes, but it's the
	// easiest way to not replace the entire environment when calling.
	err = os.Setenv("ACME_STATE_DIR", stateDirectory)
	if err != nil {
		return false, err
	}

	// Do not execute a world-writable directory.
	if (fi.Mode() & 02) != 0 {
		return false, fmt.Errorf("refusing to execute hooks, directory is world-writable: %s", directory)
	}

	ms, err := filepath.Glob(filepath.Join(directory, "*"))
	if err != nil {
		return false, err
	}

	for _, m := range ms {
		fi, err := os.Stat(m)
		if err != nil {
			log.Errore(err, "hook: ", m)
			continue
		}

		// Yes, this is vulnerable to race conditions; it's just to stop people
		// from shooting themselves in the foot.
		if (fi.Mode() & 02) != 0 {
			log.Errorf("refusing to execute world-writable hook script: %s", m)
			continue
		}

		var cmd *exec.Cmd
		if shouldSudoFile(m, fi) {
			log.Debugf("calling hook script (with sudo): %s", m)
			args2 := []string{"-n", "--", m}
			args2 = append(args2, args...)
			cmd = exec.Command("sudo", args2...)
		} else {
			log.Debugf("calling hook script: %s", m)
			cmd = exec.Command(m, args...)
		}

		cmd.Dir = "/"

		pipeR, pipeW, err := os.Pipe()
		if err != nil {
			return anySucceeded, err
		}

		defer pipeR.Close()
		go func() {
			defer pipeW.Close()
			pipeW.Write([]byte(stdinData))
		}()

		cmd.Stdin = pipeR
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		logFailedExecution(m, err)
		if err == nil {
			anySucceeded = true
		}
	}

	return anySucceeded, nil
}

func logFailedExecution(hookPath string, err error) {
	if err == nil {
		return
	}

	exitCode, err2 := deos.GetExitCode(err)
	if err2 != nil {
		// Not an error code. ???
		log.Errore(err2, "hook script: ", hookPath)
		return
	}

	switch exitCode {
	case 42:
		// Unsupported event type for this hook. Don't log anything; this is OK.
	default:
		log.Errore(err, "hook script: ", hookPath)
	}
}

// © 2015—2016 Hugo Landau <hlandau@devever.net>    MIT License
