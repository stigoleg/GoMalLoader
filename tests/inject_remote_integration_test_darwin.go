//go:build darwin
// +build darwin

package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

func TestInjectRemoteIntegration_Positive_Darwin(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "gomalloader_test_inject_darwin_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Start a dummy process (sleep 10)
	dummyCmd := exec.Command("sleep", "10")
	if err := dummyCmd.Start(); err != nil {
		t.Fatalf("failed to start dummy process: %v", err)
	}
	defer dummyCmd.Process.Kill()
	targetPID := dummyCmd.Process.Pid

	// Write safe shellcode (ret)
	payloadPath := filepath.Join(tempDir, "shellcode.bin")
	if err := ioutil.WriteFile(payloadPath, []byte{0xc3}, 0755); err != nil {
		t.Fatalf("failed to write payload: %v", err)
	}

	// Write config.json for inject_remote mode
	configPath := filepath.Join(tempDir, "config.json")
	configObj := map[string]interface{}{
		"mode":           "inject_remote",
		"source":         "file",
		"path":           payloadPath,
		"aes_key":        "0123456789abcdef",
		"target_process": strconv.Itoa(targetPID),
		"obfuscated":     false,
		"self_delete":    false,
	}
	configBytes, err := json.MarshalIndent(configObj, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	if err := ioutil.WriteFile(configPath, configBytes, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Build the loader binary
	loaderPath := filepath.Join(tempDir, "loader_test_bin")
	cmd := exec.Command("go", "build", "-o", loaderPath, "../main.go")
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build loader: %v\n%s", err, string(out))
	}

	// Run the loader as a subprocess
	runCmd := exec.Command(loaderPath)
	runCmd.Dir = tempDir
	output, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("loader run failed: %v\nOutput: %s", err, string(output))
	}
}
