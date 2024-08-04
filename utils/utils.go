package utils

import (
	"path/filepath"
	"strings"

	"github.com/splitscriptjs/cli/config"
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
func Includes[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
