package controllers

import (
	"fmt"
	"net/http"

	"github.com/liuminhaw/wrestic-brw/views"
)

type Repositories struct {
	Templates struct {
		New Template
	}
}

func (rep Repositories) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		JsFiles    []string
		FormInputs []views.RepositoryConfig
	}
	formInputs, err := views.NewRepositoryConfigs()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	data.JsFiles = append(data.JsFiles, "/static/js/new-repository.js")
	data.FormInputs = formInputs

	rep.Templates.New.Execute(w, r, data)
}
