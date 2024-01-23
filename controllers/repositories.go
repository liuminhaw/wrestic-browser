package controllers

import (
	"fmt"
	"net/http"

	"github.com/liuminhaw/wrestic-brw/models"
	"github.com/liuminhaw/wrestic-brw/restic"
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
	repoType := r.FormValue("type")

	var repository restic.Repository
	switch repoType {
	case "local":
		repository = &restic.LocalRepository{
			Password:    r.FormValue("password"),
			Destination: r.FormValue("destination"),
		}
	case "s3":
		repository = &restic.S3Repository{
			Password:        r.FormValue("password"),
			Destination:     r.FormValue("destination"),
			AccessKeyId:     r.FormValue("access-key"),
			SecretAccessKey: r.FormValue("secret-key"),
		}
	case "sftp":
		repository = &restic.SftpRepository{
			Password:    r.FormValue("password"),
			Destination: r.FormValue("destination"),
			User:        r.FormValue("sftp-user"),
			Host:        r.FormValue("sftp-host"),
			Pem:         r.FormValue("sftp-pem"),
		}
	default:
		// TODO: direct back to current page and show error message
		message := fmt.Sprintf("Repository type: %s not supported", repoType)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, message)
		return
	}

	// Check repository connection
	if err := repository.Connect(); err != nil {
		fmt.Printf("Connection failed: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Failed to connect to repository")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Connection test success")
}
