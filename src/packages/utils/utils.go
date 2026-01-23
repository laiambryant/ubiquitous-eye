package utils

import (
	"os"
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

func GetBuildSHA() string {
	sha := os.Getenv("GITHUB_SHA")
	if sha == "" {
		return "local"
	}
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}
