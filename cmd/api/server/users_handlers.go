package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/asankov/gira/internal/auth"
	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
	"github.com/hashicorp/go-multierror"
)

var (
	errInvalidToken             = errors.New("invalid token")
	errUserIsRequired           = errors.New("user is required")
	errExpectedToken            = errors.New("expected token in header")
	errEmailRequired            = errors.New("'email' is required field")
	errUsernameRequired         = errors.New("'username' is required field")
	errPasswordRequired         = errors.New("'password' is required field")
	errParsingBody              = errors.New("error while parsing request body")
	errHashedPasswordNotAllowed = errors.New("'hashedPassword' is not allowed field")
)

func (s *Server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			if err == io.EOF {
				s.respondError(w, r, errUserIsRequired, http.StatusBadRequest)
				return
			}
			s.Log.Errorf("Error while parsing body: %v", err)
			s.respondError(w, r, errParsingBody, http.StatusBadRequest)
			return
		}

		if err := validateUser(&user); err != nil {
			s.respondError(w, r, err, http.StatusBadRequest)
			return
		}

		userResponse, err := s.UserModel.Insert(&user)
		if err != nil {
			if err == postgres.ErrEmailAlreadyExists || err == postgres.ErrUsernameAlreadyExists {
				s.respondError(w, r, err, http.StatusBadRequest)
				return
			}
			s.Log.Errorf("Error while inserting user into the DB: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respond(w, r, userResponse, http.StatusOK)
	}
}

func (s *Server) handleUserGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-auth-token")
		if token == "" {
			s.respondError(w, r, errExpectedToken, http.StatusBadRequest)
			return
		}

		if _, err := s.Authenticator.DecodeToken(token); err != nil {
			if errors.Is(err, auth.ErrInvalidSignature) || errors.Is(err, auth.ErrTokenExpired) {
				s.respondError(w, r, errInvalidToken, http.StatusUnauthorized)
				return
			}
			s.Log.Errorf("Error while authenticating user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		user, err := s.UserModel.GetUserByToken(token)
		if err != nil {
			s.respondError(w, r, errInvalidToken, http.StatusUnauthorized)
			return
		}

		s.respond(w, r, &models.UserResponse{User: user}, http.StatusOK)
	}
}

func validateUser(user *models.User) error {
	var err *multierror.Error
	if user.ID != "" {
		err = multierror.Append(err, errIDNotAllowed)
	}

	if user.Username == "" {
		err = multierror.Append(err, errUsernameRequired)
	}

	if user.Email == "" {
		err = multierror.Append(err, errEmailRequired)
	}

	if user.Password == "" {
		err = multierror.Append(err, errPasswordRequired)
	}

	if user.HashedPassword != nil {
		err = multierror.Append(err, errHashedPasswordNotAllowed)
	}

	return err.ErrorOrNil()
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

		usr, err := s.UserModel.Authenticate(user.Email, user.Password)
		if err != nil {
			// TODO: JSON Error
			http.Error(w, "Wrong email/password", http.StatusUnauthorized)
			return
		}

		token, err := s.Authenticator.NewTokenForUser(usr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: persist the token, so we can invalidate it
		if err := s.UserModel.AssociateTokenWithUser(usr.ID, token); err != nil {
			s.Log.Errorf("Error while associating token with user %s: %v", usr.ID, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		s.respond(w, r, &models.UserLoginResponse{Token: token}, http.StatusOK)
	}
}

func (s *Server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := tokenFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.UserModel.InvalidateToken(user.ID, token); err != nil {
			return
		}

		s.respond(w, r, nil, http.StatusOK)
	}
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.Log.Errorf("Error while encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) respondError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()}); err != nil {
		s.Log.Errorf("Error while encoding error response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
