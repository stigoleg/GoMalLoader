//go:build darwin
// +build darwin

package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestReflectiveLoaderIntegration_Positive_Darwin(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "gomalloader_test_reflective_darwin_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Build a minimal valid .dylib shared library
	libSource := `
		#include <stdio.h>
		__attribute__((constructor)) void test_entry() { printf(\"Hello from test_entry\\n\"); }
	`
	libCPath := filepath.Join(tempDir, "libdummy.c")
	if err := ioutil.WriteFile(libCPath, []byte(libSource), 0644); err != nil {
		t.Fatalf("failed to write C source: %v", err)
	}
	libDylibPath := filepath.Join(tempDir, "libdummy.dylib")
	cmd := exec.Command("clang", "-dynamiclib", "-o", libDylibPath, libCPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build shared library: %v\n%s", err, string(out))
	}

	// Write config.json for dll_reflective mode
	configPath := filepath.Join(tempDir, "config.json")
	config := `{
		"mode": "dll_reflective",
		"source": "file",
		"path": "` + libDylibPath + `",
		"aes_key": "0123456789abcdef",
		"obfuscated": false,
		"self_delete": false
	}`
	if err := ioutil.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Build the loader binary
	loaderPath := filepath.Join(tempDir, "loader_test_bin")
	cmd = exec.Command("go", "build", "-o", loaderPath, "../main.go")
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
