package utils

import (
	"os"
	"os/exec"
	"testing"
)

func TestSelfDelete_Mock(t *testing.T) {
	// Save original functions
	origExecutable := osExecutable
	origCommand := execCommand
	defer func() {
		osExecutable = origExecutable
		execCommand = origCommand
	}()

	// Mock os.Executable
	osExecutable = func() (string, error) {
		return "/tmp/fakebinary", nil
	}
	// Mock exec.Command
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "mocked")
	}

	// Should not panic or error
	SelfDelete()
}

// --- Mocks for testing ---
var osExecutable = os.Executable
var execCommand = exec.Command

func SelfDelete() {
	path, _ := osExecutable()
	cmd := execCommand("cmd.exe", "/C", "ping 127.0.0.1 -n 3 > NUL & del \""+path+"\"")
	cmd.Start()
}
