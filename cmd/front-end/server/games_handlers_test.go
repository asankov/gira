package server_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asankov/gira/cmd/front-end/server"
	"github.com/asankov/gira/internal/fixtures"

	"github.com/golang/mock/gomock"
)

func TestHandleHome(t *testing.T) {
	ctrl := gomock.NewController(t)

	renderer := fixtures.NewRendererMock(ctrl)

	srv := &server.Server{
		Log:      log.New(os.Stdout, "", 0),
		Renderer: renderer,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	renderer.EXPECT().
		Render(gomock.Eq(w), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	srv.ServeHTTP(w, r)

	got, expected := w.Code, http.StatusOK
	if got != expected {
		t.Errorf("Got (%d) for status code, expected (%d)", got, expected)
	}
}
