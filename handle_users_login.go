package main

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"
	"log"

	"github.com/Alody/Go-HTTP-server/internal/auth"
	"github.com/Alody/Go-HTTP-server/internal/database"
	"database/sql"
)

// handlerUsersLogin handles user login requests.
func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {

	// request body
	type parameters struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
	}

	// response body
	type response struct {
		User
		Token string `json:"token,omitempty"` // Optional token field for responses
		RefreshToken string `json:"refresh_token,omitempty"` // Optional refresh token field for responses
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

	token, err := auth.MakeJWT(user.ID, cfg.JWT_key, time.Hour)
	log.Println("Token received:", token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	log.Println("Refresh token received:", refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour) // 60 days
	now := time.Now()

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
	Token:     refreshToken,
	UserID:    user.ID,
	CreatedAt: now,
	UpdatedAt: now,
	ExpiresAt: expiresAt,
	RevokedAt: sql.NullTime{Valid: false}, // no revocation time yet
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
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
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// extract auth header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header", nil)
		return
	}
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	// look up token in the database
	tokenInfo, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// check if valid and not expired
	if tokenInfo.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	if tokenInfo.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}

	// create new JWT token
	token, err := auth.MakeJWT(tokenInfo.UserID, cfg.JWT_key, time.Hour)
	log.Println("Token received:", token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new JWT", err)
		return
	}

	// return new token
	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	// extract auth header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header", nil)
		return
	}
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")
	

	// revoke the token in the database
	err := cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
	Token:     refreshToken,
	RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Refresh token not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}