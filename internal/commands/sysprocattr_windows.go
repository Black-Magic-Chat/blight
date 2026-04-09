//go:build windows

package commands

import (
	"os/exec"
	"syscall"
)

func applyHiddenProcessAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
