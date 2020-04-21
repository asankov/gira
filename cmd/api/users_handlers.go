package main

import (
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
	"gopkg.in/square/go-jose.v2/json"
)

var (
	errEmailRequired            = errors.New("'email' is required field")
	errUsernameRequired         = errors.New("'username' is required field")
	errPasswordRequired         = errors.New("'password' is required field")
	errHashedPasswordNotAllowed = errors.New("'hashedPassword' is not allowed field")
)

func (s *server) createUserHandler() http.HandlerFunc {
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

		if _, err := s.userModel.Insert(&user); err != nil {
			if err == postgres.ErrEmailAlreadyExists || err == postgres.ErrUsernameAlreadyExists {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			s.log.Printf("error while inserting user into the DB: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
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
