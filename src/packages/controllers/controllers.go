package controllers

import (
	"html/template"
	"net/http"
	"ubiquitous-eye/packages/model/response"
	"ubiquitous-eye/packages/services"
	"ubiquitous-eye/packages/utils"
)

func RootController(w http.ResponseWriter, r *http.Request) {
	user, repos, err := services.GetData()
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
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
