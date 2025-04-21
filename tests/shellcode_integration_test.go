package tests

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestShellcodeModeIntegration(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "gomalloader_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write dummy shellcode payload (just a NOP sled for test)
	payloadPath := filepath.Join(tempDir, "shellcode.bin")
	if err := ioutil.WriteFile(payloadPath, []byte{0x90, 0x90, 0x90, 0x90}, 0644); err != nil {
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
	if err == nil {
		t.Fatalf("expected loader to fail, but it succeeded. Output: %s", string(output))
	}
	if !bytes.Contains(output, []byte("unexpected fault address")) && !bytes.Contains(output, []byte("segmentation violation")) {
		t.Errorf("unexpected error output: %s", string(output))
	}

	// Optionally, check output for expected log lines
	if len(output) == 0 {
		t.Errorf("expected some output from loader, got none")
	}
}
