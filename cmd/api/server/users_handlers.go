package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
)

var (
	errEmailRequired            = errors.New("'email' is required field")
	errUsernameRequired         = errors.New("'username' is required field")
	errPasswordRequired         = errors.New("'password' is required field")
	errHashedPasswordNotAllowed = errors.New("'hashedPassword' is not allowed field")
)

type userResponse struct {
	Token string `json:"token"`
}

func (s *Server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateUser(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userResponse, err := s.UserModel.Insert(&user)
		if err != nil {
			if err == postgres.ErrEmailAlreadyExists || err == postgres.ErrUsernameAlreadyExists {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			s.Log.Printf("error while inserting user into the DB: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.respond(w, r, userResponse, http.StatusOK)
	}
}

func validateUser(user *models.User) error {
	if user.ID != "" {
		return errIDNotAllowed
	}

	if user.Username == "" {
		return errUsernameRequired
	}

	if user.Email == "" {
		return errEmailRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	if user.HashedPassword != nil {
		return errHashedPasswordNotAllowed
	}

	return nil
}

func (s *Server) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Email == "" {
			http.Error(w, "'email' is required on login", http.StatusBadRequest)
			return
		}

		if user.Password == "" {
			http.Error(w, "'password' is required on login", http.StatusBadRequest)
			return
		}

		_, err := s.UserModel.Authenticate(user.Email, user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := s.Auth.NewTokenForUser(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: persist the token, so we can invalidate it

		s.respond(w, r, &userResponse{Token: token}, http.StatusOK)
	}
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.Log.Printf("error while encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
