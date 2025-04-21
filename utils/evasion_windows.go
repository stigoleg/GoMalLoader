//go:build windows

package utils

import (
	"log"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var checkMutexImpl = func(name string) {
	handle, _, _ := kernel32.NewProc("CreateMutexW").Call(0, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))))
	if handle == 0 {
		log.Println("mutex exists, exiting")
		syscall.Exit(1)
	}
}

func CheckMutex(name string) {
	checkMutexImpl(name)
}

func CheckSleepSkew() {
	start := time.Now()
	syscall.Sleep(2000)
	if time.Since(start) < 1900*time.Millisecond {
		log.Println("sleep skew detected")
		syscall.Exit(1)
	}
}

func CheckSandbox() {
	buf := make([]uint16, 128)
	n, err := syscall.GetComputerName(&buf[0], new(uint32))
	if err != nil {
		return
	}
	name := syscall.UTF16ToString(buf[:n])
	if strings.Contains(strings.ToLower(name), "sandbox") {
		log.Println("sandbox detected")
		syscall.Exit(1)
	}
}
