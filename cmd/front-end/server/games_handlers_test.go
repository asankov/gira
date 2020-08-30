package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asankov/gira/cmd/front-end/server"
	"github.com/asankov/gira/internal/fixtures"
	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
)

func TestHandleHome(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	renderer := fixtures.NewRendererMock(ctrl)

	srv := &server.Server{
		Log:      logrus.StandardLogger(),
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
