package services

import (
	"os"
	"testing"

	"github.com/laiambryant/ubiquitous-eye/packages/utils"
)

func TestCreateDeploySite(t *testing.T) {
	targetPath := utils.GetDeployableSiteURI()

	file, err := os.Stat(targetPath)
	if err != nil {
		t.Errorf("File %s should always exist %v", targetPath, err)
		return
	}

	err = CreateDeploySite(targetPath)
	if err != nil {
		t.Errorf("CreateDeploySite should not return an error: %v", err)
		return
	}

	if updatedFile, err := os.Stat(targetPath); err != nil || updatedFile.ModTime().Equal(file.ModTime()) {
		t.Errorf("The file should exist and be updated")
	}
}

func TestRenderSite(t *testing.T) {
	if err := RenderSite(os.Stdout, utils.DEPLOYABLE_SITE_URI_TEST); err != nil {
		t.Errorf("RenderSite should not return an error: %v", err)
	}
}
