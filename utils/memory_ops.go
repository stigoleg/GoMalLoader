package utils

type MemoryOps interface {
	AllocRWX(size int) uintptr
	CreateThread(startAddr uintptr) uintptr
	WaitForThread(thread uintptr)
}
