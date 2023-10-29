package controllers

import (
	"net/http"
)

// StaticHandler takes in a controllers.Template interface
// and retrun http.HandlerFunc which interpret static html template content
func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}
