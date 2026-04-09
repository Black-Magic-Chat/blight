//go:build !windows

package apps

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func Launch(app AppEntry) error {
	target := app.Path
	if target == "" {
		target = app.LnkPath
	}

	var cmd *exec.Cmd
	lower := strings.ToLower(target)
	switch {
	case strings.HasSuffix(lower, ".app"):
		cmd = exec.Command("open", "-a", target)
	case strings.HasSuffix(lower, ".desktop"):
		cmd = exec.Command("xdg-open", target)
	default:
		if filepath.IsAbs(target) {
			cmd = exec.Command(target)
		} else {
			cmd = exec.Command("sh", "-lc", target)
		}
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch %s: %w", app.Name, err)
	}
	return nil
}
