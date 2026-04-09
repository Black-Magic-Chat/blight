//go:build windows

package updater

import (
	"strings"
	"syscall"
)

func isInstallerAsset(name string) bool {
	return strings.Contains(name, "setup") && strings.HasSuffix(name, ".exe")
}

func installerTempName() string {
	return "blight-setup.exe"
}

func installerSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{HideWindow: false}
}
