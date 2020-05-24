package server

import (
	"net/http"
	"text/template"
)

// TemplateRenderer implements Renderer and is used to render templates
type TemplateRenderer struct{}

// Render implements Renderer
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
