//go:build linux
// +build linux

package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestShellcodeModeIntegration_Positive_Linux(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "gomalloader_test_shellcode_linux_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write shellcode for ret (0xc3) on Linux x86_64
	payloadPath := filepath.Join(tempDir, "shellcode.bin")
	shellcode := []byte{0xc3}
	if err := ioutil.WriteFile(payloadPath, shellcode, 0755); err != nil {
		t.Fatalf("failed to write payload: %v", err)
	}

	// Write config.json for shellcode mode
	configPath := filepath.Join(tempDir, "config.json")
	config := `{
		"mode": "shellcode",
		"source": "file",
		"path": "` + payloadPath + `",
		"aes_key": "0123456789abcdef",
		"obfuscated": false,
		"self_delete": false
	}`
	if err := ioutil.WriteFile(configPath, []byte(config), 0644); err != nil {
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
