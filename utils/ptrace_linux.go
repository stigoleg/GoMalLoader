//go:build linux

package utils

import (
	"os"
	"syscall"
	"unsafe"
)

type userRegsStruct struct {
	R15      uint64
	R14      uint64
	R13      uint64
	R12      uint64
	Rbp      uint64
	Rbx      uint64
	R11      uint64
	R10      uint64
	R9       uint64
	R8       uint64
	Rax      uint64
	Rcx      uint64
	Rdx      uint64
	Rsi      uint64
	Rdi      uint64
	Orig_rax uint64
	Rip      uint64
	Cs       uint64
	Eflags   uint64
	Rsp      uint64
	Ss       uint64
	Fs_base  uint64
	Gs_base  uint64
	Ds       uint64
	Es       uint64
	Fs       uint64
	Gs       uint64
}

func GetRegs(pid int) (*userRegsStruct, error) {
	regs := &userRegsStruct{}
	_, _, errno := syscall.Syscall6(syscall.SYS_PTRACE, syscall.PTRACE_GETREGS, uintptr(pid), 0, uintptr(unsafe.Pointer(regs)), 0, 0)
	if errno != 0 {
		return nil, os.NewSyscallError("ptrace_getregs", errno)
	}
	return regs, nil
}

func SetRegs(pid int, regs *userRegsStruct) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_PTRACE, syscall.PTRACE_SETREGS, uintptr(pid), 0, uintptr(unsafe.Pointer(regs)), 0, 0)
	if errno != 0 {
		return os.NewSyscallError("ptrace_setregs", errno)
	}
	return nil
}

// RemoteMmap64 performs an mmap syscall in the target process and returns the allocated address
func RemoteMmap64(pid int, size int) (uintptr, error) {
	regs, err := GetRegs(pid)
	if err != nil {
		return 0, err
	}
	backup := *regs

	// Setup registers for mmap syscall
	regs.Rax = 9 // __NR_mmap
	regs.Rdi = 0 // addr
	regs.Rsi = uint64(size)
	regs.Rdx = syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	regs.R10 = syscall.MAP_ANON | syscall.MAP_PRIVATE
	regs.R8 = 0   // fd
	regs.R9 = 0   // offset
	regs.Rip -= 2 // back up to re-execute last instruction (safe for syscalls)

	if err := SetRegs(pid, regs); err != nil {
		return 0, err
	}
	// Single step
	if err := syscall.PtraceSingleStep(pid); err != nil {
		return 0, err
	}
	if err := Wait(pid); err != nil {
		return 0, err
	}
	regs2, err := GetRegs(pid)
	if err != nil {
		return 0, err
	}
	addr := uintptr(regs2.Rax)
	// Restore original regs
	SetRegs(pid, &backup)
	return addr, nil
}

func PtraceAttach(pid int) error {
	return syscall.PtraceAttach(pid)
}

func PtraceDetach(pid int) error {
	return syscall.PtraceDetach(pid)
}

func Wait(pid int) error {
	var ws syscall.WaitStatus
	_, err := syscall.Wait4(pid, &ws, 0, nil)
	return err
}

// WriteMemory uses process_vm_writev to write data to another process's memory
func WriteMemory(pid int, remoteAddr uintptr, data []byte) error {
	localIov := syscall.Iovec{
		Base: &data[0],
		Len:  uint64(len(data)),
	}
	remoteIov := syscall.Iovec{
		Base: (*byte)(unsafe.Pointer(remoteAddr)),
		Len:  uint64(len(data)),
	}
	_, _, errno := syscall.Syscall6(310,
		uintptr(pid),
		uintptr(unsafe.Pointer(&localIov)), 1,
		uintptr(unsafe.Pointer(&remoteIov)), 1,
		0)
	if errno != 0 {
		return os.NewSyscallError("process_vm_writev", errno)
	}
	return nil
}
