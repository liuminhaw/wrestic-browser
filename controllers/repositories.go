package controllers

import (
	"fmt"
	"net/http"

	"github.com/liuminhaw/wrestic-brw/context"
	"github.com/liuminhaw/wrestic-brw/restic"
	"github.com/liuminhaw/wrestic-brw/views"
)

type Repositories struct {
	Templates struct {
		New   Template
		Index Template
	}

	RepositoryService       *restic.RepositoryService
	RepositoryStatusService *restic.RepositoryStatusService
}

func (rep Repositories) Index(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Settings []struct {
			Name       string
			Status     string
			LastBackup string
			Owner      string
		}
	}
	statuses, err := rep.RepositoryStatusService.List()
	if err != nil {
		fmt.Printf("Repository indices: %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	for _, status := range statuses {
		data.Settings = append(data.Settings, struct {
			Name       string
			Status     string
			LastBackup string
			Owner      string
		}{
			Name:       status.Name,
			Status:     status.Status,
			LastBackup: status.LastBackup,
			Owner:      status.Owner,
		})
	}

	rep.Templates.Index.Execute(w, r, data)
}

func (rep Repositories) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		JsFiles    []string
		FormInputs []views.RepositoryConfig
		RepoTypes  []string
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
	userId := context.User(r.Context()).ID

	switch repoType {
	case "local":
		rep.RepositoryService.Repository = &restic.LocalRepository{
			Name:        r.FormValue("name"),
			Password:    r.FormValue("password"),
			Destination: r.FormValue("destination"),
			Encryption:  &restic.LocalRepositoryEnc{},
		}
	case "s3":
		rep.RepositoryService.Repository = &restic.S3Repository{
			Name:            r.FormValue("name"),
			Password:        r.FormValue("password"),
			Destination:     r.FormValue("destination"),
			AccessKeyId:     r.FormValue("access-key"),
			SecretAccessKey: r.FormValue("secret-key"),
			Region:          r.FormValue("aws-region"),
			Encryption:      &restic.S3RepositoryEnc{},
		}
	case "sftp":
		rep.RepositoryService.Repository = &restic.SftpRepository{
			Name:        r.FormValue("name"),
			Password:    r.FormValue("password"),
			Destination: r.FormValue("destination"),
			User:        r.FormValue("sftp-user"),
			Host:        r.FormValue("sftp-host"),
			Pem:         r.FormValue("sftp-pem"),
			Encryption:  &restic.SftpRepositoryEnc{},
		}
	default:
		// TODO: direct back to current page and show error message
		message := fmt.Sprintf("Repository type: %s not supported", repoType)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, message)
		return
	}
	// Generate encrypted data
	if err := rep.RepositoryService.Repository.GenEnc(rep.RepositoryService.EncKey); err != nil {
		fmt.Printf("Generate repository encrypted data failed: %s\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "Server Error")
	}

	// Check repository connection
	if err := rep.RepositoryService.Connect(); err != nil {
		// TODO: Direct back to repository new page and show error message
		fmt.Printf("Connection failed: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Failed to connect to repository")
		return
	}
	if err := rep.RepositoryService.Create(userId); err != nil {
		fmt.Printf("Create new repository failed: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Failed to create new repository")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Connection test success")
}
