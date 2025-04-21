//go:build darwin

package utils

import (
	"syscall"
)

type darwinMemoryOps struct{}

func NewMemoryOps() MemoryOps {
	return &darwinMemoryOps{}
}

func (m *darwinMemoryOps) AllocRWX(size int) uintptr {
	addr, _, err := syscall.Syscall6(syscall.SYS_MMAP, 0, uintptr(size), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_ANON|syscall.MAP_PRIVATE, 0, 0)
	if err != 0 {
		panic("mmap failed")
	}
	return addr
}

func (m *darwinMemoryOps) CreateThread(startAddr uintptr) uintptr {
	// For now, just return the address; shellcode will be called in current thread
	return startAddr
}

func (m *darwinMemoryOps) WaitForThread(thread uintptr) {
	// No-op for now
}
