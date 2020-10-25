package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asankov/gira/pkg/client"
)

func (s *Server) handleFranchisesAddPost() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		franchiseName := r.PostForm.Get("franchise")
		if franchiseName == "" {
			http.Error(w, "'franchise' is required", http.StatusBadRequest)
			return
		}

		resp, err := s.Client.CreateFranchise(&client.CreateFranchiseRequest{Name: franchiseName}, token)
		if err != nil {
			if errors.Is(err, client.ErrNoAuthorization) {
				w.Header().Add("Location", "/users/login")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/games/new?selectedFranchise=%s", resp.Franchise.ID))
		w.WriteHeader(http.StatusSeeOther)
	}
}
