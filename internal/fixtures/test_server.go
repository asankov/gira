package fixtures

import (
	http "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/gira-games/api/pkg/models"
)

type ServerBuilder struct {
	path         string
	method       string
	token        string
	data         interface{}
	responseCode int
	query        string
	t            *testing.T
}

func NewTestServer(t *testing.T) ServerBuilder {
	return ServerBuilder{
		path:         "/",
		method:       http.MethodGet,
		responseCode: http.StatusOK,
	}
}

func (s ServerBuilder) Path(path string) ServerBuilder {
	s.path = path
	return s
}

func (s ServerBuilder) Method(method string) ServerBuilder {
	s.method = method
	return s
}

func (s ServerBuilder) Token(token string) ServerBuilder {
	s.token = token
	return s
}

func (s ServerBuilder) Data(data interface{}) ServerBuilder {
	s.data = data
	return s
}

func (s ServerBuilder) Return(responseCode int) ServerBuilder {
	s.responseCode = responseCode
	return s
}

func (s ServerBuilder) Query(query string) ServerBuilder {
	s.query = query
	return s
}

func (s ServerBuilder) Build() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != s.path || r.Method != s.method {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if s.token != "" {
			if r.Header.Get(models.XAuthToken) != s.token {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		if s.query != "" {
			if !strings.Contains(r.URL.RawQuery, "excludeAssigned=true") {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// TODO: assert body

		w.WriteHeader(s.responseCode)
		if s.data != nil {
			if _, err := w.Write(MarshalBytes(s.t, s.data)); err != nil {
				s.t.Fatalf("error while writing response - %v", err)
			}
		}
	}))
}
