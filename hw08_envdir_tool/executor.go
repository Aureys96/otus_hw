package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...)

	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout

	c.Env = func() []string {
		for k, v := range env {
			if v.NeedRemove {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v.Value)
			}
		}
		return os.Environ()
	}()

	if err := c.Run(); err != nil {
		returnCode = c.ProcessState.ExitCode()
	}

	return
}
