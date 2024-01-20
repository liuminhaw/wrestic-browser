package controllers

import (
	"fmt"
	"net/http"

	"github.com/liuminhaw/wrestic-brw/models"
	"github.com/liuminhaw/wrestic-brw/views"
)

type Repositories struct {
	Templates struct {
		New Template
	}

	RepositoryService *models.RepositoryService
}

func (rep Repositories) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		JsFiles    []string
		FormInputs []views.RepositoryConfig
		RepoTypes  []models.RepositoryTypes
	}
	formInputs, err := views.NewRepositoryConfigs()
	if err != nil {
		fmt.Printf("New repository: form inputs: %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	repoTypes, err := rep.RepositoryService.Types()
	if err != nil {
		fmt.Printf("New repository: repository types: %s\n", err)
		http.Error(w, "Internal service error", http.StatusInternalServerError)
	}

	data.JsFiles = append(data.JsFiles, "/static/js/new-repository.js")
	data.FormInputs = formInputs
	data.RepoTypes = repoTypes

	rep.Templates.New.Execute(w, r, data)
}

func (rep Repositories) Create(w http.ResponseWriter, r *http.Request) {
	message := "POST /repositories for creating new repository connection"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, message)
}
