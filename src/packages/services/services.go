package services

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/laiambryant/ubiquitous-eye/packages/model/response"
	"github.com/laiambryant/ubiquitous-eye/packages/utils"
)

//go:embed resources/index.html
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
		User  *response.UserAPIResponse
		Repos []response.UserRepoApiResponse
	}{
		User:  user,
		Repos: repos,
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
		User  *response.UserAPIResponse
		Repos []response.UserRepoApiResponse
	}{
		User:  user,
		Repos: repos,
	}
	file, err := os.Create(location)
	if err != nil {
		slog.Error("Cannot create file:", "error", err)
		return err
	}
	defer file.Close()

	tmpl := template.Must(template.ParseFS(indexHTML, "resources/index.html"))
	err = tmpl.Execute(file, data)
	if err != nil {
		slog.Error("Cannot execute template:", "error", err)
		return err
	}
	return nil
}
