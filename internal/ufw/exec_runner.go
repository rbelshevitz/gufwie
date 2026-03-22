package ufw

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
)

type ExecRunner struct{}

func NewExecRunner() *ExecRunner { return &ExecRunner{} }

func (r *ExecRunner) Run(ctx context.Context, name string, args ...string) (stdout string, stderr string, exitCode int, err error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var outB, errB bytes.Buffer
	cmd.Stdout = &outB
	cmd.Stderr = &errB

	err = cmd.Run()
	stdout = outB.String()
	stderr = errB.String()

	if err == nil {
		return stdout, stderr, 0, nil
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		exitCode = exitCodeFromExitError(ee)
		return stdout, stderr, exitCode, nil
	}

	return stdout, stderr, -1, fmt.Errorf("exec %s: %w", name, err)
}

func exitCodeFromExitError(ee *exec.ExitError) int {
	if runtime.GOOS == "windows" {
		// Go sets ProcessState.ExitCode() on Windows.
		return ee.ProcessState.ExitCode()
	}
	if status, ok := ee.Sys().(syscall.WaitStatus); ok {
		return status.ExitStatus()
	}
	return ee.ProcessState.ExitCode()
}

