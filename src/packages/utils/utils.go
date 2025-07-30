package utils

import (
	"path/filepath"
	"runtime"
)

func GetProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
}

func GetDeployableSiteURI() string {
	return filepath.Join(GetProjectRoot(), "docs", "index.html")
}
