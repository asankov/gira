package templates

import (
	"net/http"
	"text/template"

	"github.com/asankov/gira/cmd/front-end/server"
)

// Renderer implements Renderer and is used to render templates
type Renderer struct{}

// NewRenderer returns new Renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Render implements Renderer
func (t *Renderer) Render(w http.ResponseWriter, r *http.Request, d interface{}, p server.Page) error {
	tt, err := template.ParseFiles("./ui/html/"+string(p), "./ui/html/base.layout.tmpl")
	if err != nil {
		return err
	}

	if err := tt.Execute(w, d); err != nil {
		return err
	}
	return nil
}
