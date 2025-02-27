package controllers

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"ubiquitous-eye/packages/model/response"
	"ubiquitous-eye/packages/utils"
)

func RootController(w http.ResponseWriter, r *http.Request) {
	user, repos, err := getData()
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}

	data := struct {
		User  response.UserAPIResponse
		Repos []response.UserRepoApiResponse
	}{
		User:  user,
		Repos: repos,
	}

	tmpl := template.Must(template.ParseFiles(utils.INDEX_HTML_PATH))
	_ = tmpl.Execute(w, data)
}

func getData() (response.UserAPIResponse, []response.UserRepoApiResponse, error) {
	user, err := getUserData()
	if err != nil {
		return response.UserAPIResponse{}, []response.UserRepoApiResponse{}, err
	}
	user_repo, err := getUserRepoData()
	if err != nil {
		return response.UserAPIResponse{}, []response.UserRepoApiResponse{}, err
	}
	if utils.DEBUG {
		printUserData(user, user_repo)
	}
	return user, user_repo, nil
}

func getUserData() (response.UserAPIResponse, error) {
	var user response.UserAPIResponse
	user_resp, err := http.Get(utils.GITHUB_USER_LBRYANT)
	if err != nil {
		log.Fatal(err)
		return response.UserAPIResponse{}, err
	}
	defer user_resp.Body.Close()
	str, err := io.ReadAll(user_resp.Body)
	if err != nil {
		log.Fatal(err)
		return response.UserAPIResponse{}, err
	}
	err = json.Unmarshal(str, &user)
	if err != nil {
		log.Fatal(err)
		return response.UserAPIResponse{}, err
	}
	return user, nil
}

func getUserRepoData() ([]response.UserRepoApiResponse, error) {
	var usrRep []response.UserRepoApiResponse
	user_project_resp, err := http.Get(utils.GITHUB_USER_LBRYANT + "/repos")
	if err != nil {
		log.Fatal(err)
		return []response.UserRepoApiResponse{}, err
	}
	str, err := io.ReadAll(user_project_resp.Body)
	if err != nil {
		log.Fatal(err)
		return []response.UserRepoApiResponse{}, err
	}
	err = json.Unmarshal(str, &usrRep)
	if err != nil {
		log.Fatal(err)
		return []response.UserRepoApiResponse{}, err
	}
	return usrRep, nil
}

func printUserData(user response.UserAPIResponse, repos []response.UserRepoApiResponse) {
	log.Printf("User: %v\n", user)
	for _, repo := range repos {
		log.Printf("Repo: %+v\n", repo)
	}
}
