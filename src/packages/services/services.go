package services

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
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
		log.Print(err)
		return response.UserAPIResponse{}, err
	}
	defer userResp.Body.Close()
	str, err := io.ReadAll(userResp.Body)
	if err != nil {
		log.Print(err)
		return response.UserAPIResponse{}, err
	}
	err = json.Unmarshal(str, &user)
	if err != nil {
		log.Print(err)
		return response.UserAPIResponse{}, err
	}
	return user, nil
}

func getUserRepoData() ([]response.UserRepoApiResponse, error) {
	var usrRep []response.UserRepoApiResponse
	userProjectResp, err := http.Get(utils.GITHUB_USER_LBRYANT + "/repos")
	if err != nil {
		log.Print(err)
		return []response.UserRepoApiResponse{}, err
	}
	str, err := io.ReadAll(userProjectResp.Body)
	if err != nil {
		log.Print(err)
		return []response.UserRepoApiResponse{}, err
	}
	err = json.Unmarshal(str, &usrRep)
	if err != nil {
		log.Print(err)
		return []response.UserRepoApiResponse{}, err
	}
	return usrRep, nil
}

func RenderSite(w io.Writer) error {
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

	tmpl := template.Must(template.ParseFS(indexHTML, "resources/index.html"))
	return tmpl.Execute(w, data)
}

func CreateDeploySite(location string) {
	user, repos, err := GetData()
	if err != nil {
		log.Fatal("Failed to get data:", err)
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
		log.Fatal("Cannot create file:", err)
	}
	defer file.Close()

	tmpl := template.Must(template.ParseFS(indexHTML, "resources/index.html"))
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal("Cannot execute template:", err)
	}
}
