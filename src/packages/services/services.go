package services

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/laiambryant/ubiquitous-eye/packages/model/response"
	"github.com/laiambryant/ubiquitous-eye/packages/utils"
)

//go:embed resources/index.html resources/site.css resources/site.js
var indexHTML embed.FS

func GetData() (*response.UserAPIResponse, []response.UserRepoApiResponse, error) {
	user, _ := getUserData()
	userRepo, _ := getUserRepoData()
	return &user, userRepo, nil
}

func getUserData() (response.UserAPIResponse, error) {
	var user response.UserAPIResponse
	userResp, err := http.Get(utils.GITHUB_USER_LBRYANT)
	if err != nil {
		slog.Warn(err.Error())
		return response.UserAPIResponse{}, err
	}
	defer userResp.Body.Close()
	str, err := io.ReadAll(userResp.Body)
	if err != nil {
		slog.Warn(err.Error())
		return response.UserAPIResponse{}, err
	}
	err = json.Unmarshal(str, &user)
	if err != nil {
		slog.Warn(err.Error())
		return response.UserAPIResponse{}, err
	}
	return user, nil
}

func getUserRepoData() ([]response.UserRepoApiResponse, error) {
	var usrRep []response.UserRepoApiResponse
	userProjectResp, err := http.Get(utils.GITHUB_USER_LBRYANT + "/repos")
	if err != nil {
		slog.Warn(err.Error())
		return []response.UserRepoApiResponse{}, err
	}
	str, err := io.ReadAll(userProjectResp.Body)
	if err != nil {
		slog.Warn(err.Error())
		return []response.UserRepoApiResponse{}, err
	}
	err = json.Unmarshal(str, &usrRep)
	if err != nil {
		slog.Warn(err.Error())
		return []response.UserRepoApiResponse{}, err
	}
	return usrRep, nil
}

func RenderSite(w io.Writer, location string) error {
	user, repos, err := GetData()
	if err != nil {
		return err
	}
	data := struct {
		User     *response.UserAPIResponse
		Repos    []response.UserRepoApiResponse
		Year     int
		BuildSHA string
	}{
		User:     user,
		Repos:    repos,
		Year:     time.Now().Year(),
		BuildSHA: utils.GetBuildSHA(),
	}

	tmpl := template.Must(template.ParseFS(indexHTML, location))
	return tmpl.Execute(w, data)
}

func CreateDeploySite(location string) error {
	user, repos, err := GetData()
	if err != nil {
		slog.Error("Failed to get data:", "error", err)
		return err
	}
	data := struct {
		User     *response.UserAPIResponse
		Repos    []response.UserRepoApiResponse
		Year     int
		BuildSHA string
	}{
		User:     user,
		Repos:    repos,
		Year:     time.Now().Year(),
		BuildSHA: utils.GetBuildSHA(),
	}
	file, err := os.Create(location)
	if err != nil {
		slog.Error("Cannot create file:", "error", err)
		return err
	}
	defer file.Close()

	err = writeAssets(filepath.Dir(location))
	if err != nil {
		slog.Error("Cannot write assets:", "error", err)
		return err
	}

	tmpl := template.Must(template.ParseFS(indexHTML, "resources/index.html"))
	err = tmpl.Execute(file, data)
	if err != nil {
		slog.Error("Cannot execute template:", "error", err)
		return err
	}
	return nil
}

func writeAssets(outputDir string) error {
	assetsDir := filepath.Join(outputDir, "assets")
	err := os.MkdirAll(assetsDir, 0o755)
	if err != nil {
		return err
	}

	err = writeAssetFile(filepath.Join(assetsDir, "site.css"), "resources/site.css")
	if err != nil {
		return err
	}

	return writeAssetFile(filepath.Join(assetsDir, "site.js"), "resources/site.js")
}

func writeAssetFile(outputPath, sourcePath string) error {
	asset, err := indexHTML.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, asset, 0o644)
}
