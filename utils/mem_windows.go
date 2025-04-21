//go:build windows

package utils

import (
	"syscall"
	"testing"
)

type windowsMemoryOps struct{}

func NewMemoryOps() MemoryOps {
	return &windowsMemoryOps{}
}

func (m *windowsMemoryOps) AllocRWX(size int) uintptr {
	virtAlloc := GetProcAddr("kernel32.dll", "VirtualAlloc")
	return CallProc(virtAlloc, 0, uintptr(size), 0x3000, 0x40)
}

func (m *windowsMemoryOps) CreateThread(startAddr uintptr) uintptr {
	createThread := GetProcAddr("kernel32.dll", "CreateThread")
	return CallProc(createThread, 0, 0, startAddr, 0, 0, 0)
}

func (m *windowsMemoryOps) WaitForThread(thread uintptr) {
	waitForSingleObject := GetProcAddr("kernel32.dll", "WaitForSingleObject")
	CallProc(waitForSingleObject, thread, 0xFFFFFFFF)
}

func GetProcAddr(dll, proc string) uintptr {
	lib, _ := syscall.LoadLibrary(dll)
	addr, _ := syscall.GetProcAddress(lib, proc)
	return uintptr(addr)
}

func CallProc(addr uintptr, args ...uintptr) uintptr {
	ret, _, _ := syscall.SyscallN(addr, args...)
	return ret
}

func TestWindowsMemoryOps_AllocRWX(t *testing.T) {
	mem := &windowsMemoryOps{}
	addr := mem.AllocRWX(4096)
	if addr == 0 {
		t.Errorf("AllocRWX returned zero address")
	}
}

func TestWindowsMemoryOps_CreateThread_WaitForThread(t *testing.T) {
	mem := &windowsMemoryOps{}
	addr := uintptr(0x1234)
	thread := mem.CreateThread(addr)
	if thread == 0 {
		t.Errorf("CreateThread returned zero thread handle")
	}
	// Should not panic
	mem.WaitForThread(thread)
}
