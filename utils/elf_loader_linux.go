//go:build linux

package utils

/*
#include <dlfcn.h>
#include <stdlib.h>

void* go_dlopen(const char* path) {
    return dlopen(path, RTLD_NOW | RTLD_LOCAL);
}
*/
import "C"
import (
	"fmt"
	"os"
	"testing"
	"unsafe"
)

func DlopenFromFd(fd int) error {
	path := fmt.Sprintf("/proc/self/fd/%d", fd)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.go_dlopen(cpath) == nil {
		return fmt.Errorf("dlopen failed for %s", path)
	}
	return nil
}

func TestDlopenFromFd_InvalidFd(t *testing.T) {
	err := DlopenFromFd(-1)
	if err == nil {
		t.Error("expected error for invalid fd")
	}
}

func TestDlopenFromFd_NonELFFile(t *testing.T) {
	f, err := os.Open("/dev/null")
	if err != nil {
		t.Fatalf("failed to open /dev/null: %v", err)
	}
	defer f.Close()
	err = DlopenFromFd(int(f.Fd()))
	if err == nil {
		t.Error("expected error for non-ELF file")
	}
}
