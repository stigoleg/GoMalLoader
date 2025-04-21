package utils

import "testing"

func TestCheckMutex_NoPanic(t *testing.T) {
	CheckMutex("test")
}

func TestCheckSandbox_NoPanic(t *testing.T) {
	CheckSandbox()
}

func TestCheckSleepSkew_NoPanic(t *testing.T) {
	CheckSleepSkew()
}

func TestCheckMutex_ConcurrentInstances(t *testing.T) {
	calls := 0
	exitCalled := false
	origImpl := checkMutexImpl
	checkMutexImpl = func(name string) {
		calls++
		if calls > 1 {
			exitCalled = true
		}
	}
	defer func() { checkMutexImpl = origImpl }()

	CheckMutex("test-mutex") // first instance
	CheckMutex("test-mutex") // second instance (should trigger exit)

	if !exitCalled {
		t.Error("expected exit to be called on second instance")
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}
