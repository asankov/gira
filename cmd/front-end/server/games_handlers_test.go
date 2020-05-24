package server

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asankov/gira/internal/fixtures"

	"github.com/golang/mock/gomock"
)

func TestHandleHome(t *testing.T) {
	ctrl := gomock.NewController(t)

	renderer := fixtures.NewRendererMock(ctrl)

	srv := &Server{
		Log:      log.New(os.Stdout, "", 0),
		Renderer: renderer,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	srv.ServeHTTP(w, r)

	// fix template not found error

	got, expected := r.Response.StatusCode, http.StatusOK
	if r.Response.StatusCode != expected {
		t.Errorf("Got (%d) for status code, expected (%d)", got, expected)
	}
}
