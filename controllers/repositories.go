package controllers

import "net/http"

type Repositories struct {
	Templates struct {
		New Template
	}
}

func (rep Repositories) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		JsFiles []string
	}
	data.JsFiles = append(data.JsFiles, "/static/js/new-repository.js")
	rep.Templates.New.Execute(w, r, data)
}
