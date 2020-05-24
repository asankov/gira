package server

import (
	"net/http"
	"text/template"

	"github.com/asankov/gira/pkg/models"
)

type TemplateRenderer struct{}

func (t *TemplateRenderer) Render(w http.ResponseWriter, r *http.Request, d interface{}, p Page) error {
	tt, err := template.ParseFiles("./ui/html/"+string(p), "./ui/html/base.layout.tmpl")
	if err != nil {
		return err
	}

	if err := tt.Execute(w, d); err != nil {
		return err
	}
	return nil
}
