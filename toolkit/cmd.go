package toolkit

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

/*
  执行shell命令
*/
type CmdExecutor interface {
	Execute() error
	Stderr() string
	Stdout() string
	//RetCode() int
}

type CommandExecutor struct {
	cmd     string
	stdout  bytes.Buffer
	stderr  bytes.Buffer
	retCode int
}

func NewCommandExecutor(cmd string) CmdExecutor {
	return &CommandExecutor{
		cmd: cmd,
	}
}

func (ce *CommandExecutor) Execute() error {
	if strings.Trim(ce.cmd, " ") == "" {
		return errors.New("cmd should not be empty")
	}

	cmd := exec.Command("/bin/sh", "-c", ce.cmd)
	cmd.Stdout = &ce.stdout
	cmd.Stderr = &ce.stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (ce *CommandExecutor) Stderr() string {
	return ce.stderr.String()
}

func (ce *CommandExecutor) Stdout() string {
	return ce.stdout.String()
}

type CommandExecutorWithTimeout struct {
	cmd     string
	stdout  bytes.Buffer
	stderr  bytes.Buffer
	retCode int

	timeout time.Duration
}

// with timeout
func NewCommandExecutorWithTimeout(cmd string, timeout time.Duration) CmdExecutor {
	return &CommandExecutorWithTimeout{
		cmd:     cmd,
		timeout: timeout,
	}
}

func (ct *CommandExecutorWithTimeout) Execute() error {
	if strings.Trim(ct.cmd, " ") == "" {
		return errors.New("cmd should not be empty")
	}
	if ct.timeout == 0 {
		ct.timeout = 60 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), ct.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", ct.cmd)
	cmd.Stdout = &ct.stdout
	cmd.Stderr = &ct.stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (ct *CommandExecutorWithTimeout) Stderr() string {
	return ct.stderr.String()
}

func (ct *CommandExecutorWithTimeout) Stdout() string {
	return ct.stdout.String()
}
