//go:build !windows

package files

import "syscall"

// HiddenCmd returns process attributes for launching helper commands detached on non-Windows targets.
func HiddenCmd(name string, args ...string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
