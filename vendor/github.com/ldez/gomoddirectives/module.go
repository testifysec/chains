package gomoddirectives

import (
<<<<<<< HEAD
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ldez/grignotin/gomod"
	"golang.org/x/mod/modfile"
)

// GetModuleFile gets module file.
// It's better to use [GetGoModFile] instead of this function.
func GetModuleFile() (*modfile.File, error) {
	info, err := gomod.GetModuleInfo()
	if err != nil {
		return nil, err
	}

	if info[0].GoMod == "" {
		return nil, errors.New("working directory is not part of a module")
	}

	return parseGoMod(info[0].GoMod)
}

func parseGoMod(goMod string) (*modfile.File, error) {
	raw, err := os.ReadFile(filepath.Clean(goMod))
=======
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/mod/modfile"
)

type modInfo struct {
	Path      string `json:"Path"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
	Main      bool   `json:"Main"`
}

// GetModuleFile gets module file.
func GetModuleFile() (*modfile.File, error) {
	// https://github.com/golang/go/issues/44753#issuecomment-790089020
	cmd := exec.Command("go", "list", "-m", "-json")

	raw, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command go list: %w: %s", err, string(raw))
	}

	var v modInfo
	err = json.NewDecoder(bytes.NewBuffer(raw)).Decode(&v)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling error: %w: %s", err, string(raw))
	}

	if v.GoMod == "" {
		return nil, errors.New("working directory is not part of a module")
	}

	raw, err = os.ReadFile(v.GoMod)
>>>>>>> 70e0318b1 ([WIP] add archivista storage backend)
	if err != nil {
		return nil, fmt.Errorf("reading go.mod file: %w", err)
	}

	return modfile.Parse("go.mod", raw, nil)
}
