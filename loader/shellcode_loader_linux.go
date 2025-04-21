//go:build linux

package loader

func NewShellcodeLoader() ShellcodeLoader {
	return &linuxShellcodeLoader{}
}
