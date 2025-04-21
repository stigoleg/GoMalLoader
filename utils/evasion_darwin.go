//go:build darwin

package utils

var checkMutexImpl = func(name string) {}

func CheckMutex(name string) {
	checkMutexImpl(name)
}

func CheckSandbox()   {}
func CheckSleepSkew() {}
