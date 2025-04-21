//go:build darwin

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

func DlopenFromPath(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.go_dlopen(cpath) == nil {
		return fmt.Errorf("dlopen failed for %s", path)
	}
	return nil
}

func TestDlopenFromPath_InvalidPath(t *testing.T) {
	err := DlopenFromPath("/nonexistent/path/to/file.dylib")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestDlopenFromPath_NonDylibFile(t *testing.T) {
	f, err := os.Open("/dev/null")
	if err != nil {
		t.Fatalf("failed to open /dev/null: %v", err)
	}
	defer f.Close()
	err = DlopenFromPath("/dev/null")
	if err == nil {
		t.Error("expected error for non-dylib file")
	}
}
