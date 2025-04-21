//go:build windows

package loader

func NewShellcodeLoader() ShellcodeLoader {
	return &windowsShellcodeLoader{}
}
