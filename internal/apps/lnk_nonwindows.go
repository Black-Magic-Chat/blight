//go:build !windows

package apps

func ResolveLnkTarget(lnkPath string) string {
	return ""
}

func FindAppIcon(targetPath string) string {
	return ""
}
