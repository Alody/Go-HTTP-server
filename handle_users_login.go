package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Alody/Go-HTTP-server/internal/auth"
)

// handlerUsersLogin handles user login requests.
func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {

	// request body
	type parameters struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds *int64 `json:"expires_in_seconds"`
	}

	// response body
	type response struct {
		User
		Token string `json:"token,omitempty"` // Optional token field for responses
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(user.PasswordHash, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	if params.Expires_in_seconds == nil {
		defaultValue := int64(3600)               // Default to 1 hour
		params.Expires_in_seconds = &defaultValue // Default to 1 hour if not specified
	} else if *params.Expires_in_seconds > 3600 {
		cappedVal := int64(3600)               // Cap to 1 hour
		params.Expires_in_seconds = &cappedVal // Cap to 1 hour
	}
	token, err := auth.MakeJWT(user.ID, cfg.JWT_key, time.Duration(*params.Expires_in_seconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
