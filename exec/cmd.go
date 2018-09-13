package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Cmd cmd struct
func Cmd(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// CmdOut out bytes
func CmdOut(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	out = bytes.TrimSpace(out)
	if err != nil {
		return nil, fmt.Errorf(err.Error() + ":" + string(out))
	}
	return out, nil
}

// CmdOutWithTimeout timeout
func CmdOutWithTimeout(ctx context.Context, command string, timeout int) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	done := make(chan error)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	cmd.Start()
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		err := cmd.Process.Kill() // timeout
		if err != nil {
			return stdout.Bytes(), fmt.Errorf(stderr.String() + err.Error() + ctx.Err().Error())
		}
		return stdout.Bytes(), fmt.Errorf(stderr.String() + ctx.Err().Error())
	case <-time.After(time.Duration(timeout) * time.Second):
		err := cmd.Process.Kill() // timeout
		if err != nil {
			return stdout.Bytes(), fmt.Errorf(stderr.String() + err.Error())
		}
		return stdout.Bytes(), fmt.Errorf(stderr.String()+`TIMEOUT %d,Process "%s" has been killed`, timeout, strings.Join(cmd.Args, " "))
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf(err.Error() + ":" + stderr.String())
		}
		return stdout.Bytes(), err
	}
}
