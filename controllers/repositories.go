package controllers

import "net/http"

type Repositories struct {
	Templates struct {
		New Template
	}
}

func (rep Repositories) New(w http.ResponseWriter, r *http.Request) {
	rep.Templates.New.Execute(w, r, nil)
}
