//go:build linux

package utils

import (
	"syscall"
	"testing"
)

type linuxMemoryOps struct{}

func NewMemoryOps() MemoryOps {
	return &linuxMemoryOps{}
}

func (m *linuxMemoryOps) AllocRWX(size int) uintptr {
	addr, _, err := syscall.Syscall6(syscall.SYS_MMAP, 0, uintptr(size), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_ANON|syscall.MAP_PRIVATE, 0, 0)
	if err != 0 {
		panic("mmap failed")
	}
	return addr
}

func (m *linuxMemoryOps) CreateThread(startAddr uintptr) uintptr {
	// For now, just return the address; shellcode will be called in current thread
	return startAddr
}

func (m *linuxMemoryOps) WaitForThread(thread uintptr) {
	// No-op for now
}

func TestLinuxMemoryOps_AllocRWX(t *testing.T) {
	mem := &linuxMemoryOps{}
	addr := mem.AllocRWX(4096)
	if addr == 0 {
		t.Errorf("AllocRWX returned zero address")
	}
}

func TestLinuxMemoryOps_CreateThread_WaitForThread(t *testing.T) {
	mem := &linuxMemoryOps{}
	addr := uintptr(0x1234)
	if mem.CreateThread(addr) != addr {
		t.Errorf("CreateThread should return input address")
	}
	// Should not panic
	mem.WaitForThread(addr)
}
