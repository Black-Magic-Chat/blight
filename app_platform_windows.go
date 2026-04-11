//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"
)

var procShellExecute = syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW")

func configureSettingsCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
}

// shellOpen opens a file with its default handler via ShellExecuteW — no cmd.exe flash.
func shellOpen(path string) {
	verb, _ := syscall.UTF16PtrFromString("open")
	file, _ := syscall.UTF16PtrFromString(path)
	procShellExecute.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), 0, 0, 1)
}

// explorerSelect opens Windows Explorer with the file selected, without spawning a console.
func explorerSelect(path string) {
	arg, _ := syscall.UTF16PtrFromString("/select," + path)
	explorer, _ := syscall.UTF16PtrFromString("explorer.exe")
	procShellExecute.Call(0, 0, uintptr(unsafe.Pointer(explorer)), uintptr(unsafe.Pointer(arg)), 0, 1)
}

var (
	wtOnce sync.Once
	wtPath string
)

func detectWindowsTerminal() string {
	wtOnce.Do(func() {
		local, _ := os.UserCacheDir()
		candidate := filepath.Join(local, "Microsoft", "WindowsApps", "wt.exe")
		if _, err := os.Stat(candidate); err == nil {
			wtPath = candidate
			return
		}
		if p, err := exec.LookPath("wt"); err == nil {
			wtPath = p
		}
	})
	return wtPath
}

func openInTerminal(dir string) {
	wt := detectWindowsTerminal()
	if wt != "" {
		wtPtr, _ := syscall.UTF16PtrFromString(wt)
		argsPtr, _ := syscall.UTF16PtrFromString(fmt.Sprintf(`-d "%s"`, dir))
		procShellExecute.Call(0, 0, uintptr(unsafe.Pointer(wtPtr)), uintptr(unsafe.Pointer(argsPtr)), 0, 1)
		return
	}
	cmdPtr, _ := syscall.UTF16PtrFromString("cmd.exe")
	argsPtr, _ := syscall.UTF16PtrFromString(fmt.Sprintf(`/k cd /d "%s"`, dir))
	procShellExecute.Call(0, 0, uintptr(unsafe.Pointer(cmdPtr)), uintptr(unsafe.Pointer(argsPtr)), 0, 1)
}

func blightInstallDir() string {
	local, err := os.UserCacheDir() // returns %LocalAppData% on Windows
	if err != nil {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "blight")
	}
	return filepath.Join(local, "blight")
}

func runAsAdmin(path string) error {
	verb, _ := syscall.UTF16PtrFromString("runas")
	exe, _ := syscall.UTF16PtrFromString(path)
	cwd, _ := syscall.UTF16PtrFromString(filepath.Dir(path))

	ret, _, _ := procShellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(exe)),
		0,
		uintptr(unsafe.Pointer(cwd)),
		1, // SW_SHOWNORMAL
	)
	if ret <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", ret)
	}
	return nil
}
