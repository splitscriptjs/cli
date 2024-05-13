package utils

import (
	"path/filepath"
	"strings"

	"splitscript/config"
)

func GenerateDevFileName(conf config.Config, path string) (string, error) {
	rel, err := filepath.Rel("./", path)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(rel)
	if ext == ".ts" {
		rel = strings.TrimSuffix(rel, ext) + ".js"
	}
	return filepath.Join(conf.Dev, rel), nil
}
