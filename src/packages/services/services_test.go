package services

import (
	"os"
	"testing"

	"github.com/laiambryant/ubiquitous-eye/packages/utils"
)

func TestCreateDeploySite(t *testing.T) {
	wd, _ := os.Getwd()
	file, err := os.Stat(wd + "/" + utils.DEPLOYABLE_SITE_URI)
	if err != nil {
		t.Errorf("File %s/%s should always exist %v", wd, utils.DEPLOYABLE_SITE_URI, err)
		return
	}
	CreateDeploySite()
	if updatedFile, err := os.Stat(utils.DEPLOYABLE_SITE_URI); err != nil || updatedFile.ModTime().Equal(file.ModTime()) {
		t.Errorf("The file should exist")
		return
	}
}

func TestRenderSite(t *testing.T) {
	if err := RenderSite(os.Stdout); err != nil {
		t.Errorf("RenderSite should not return an error: %v", err)
		return
	}
}
