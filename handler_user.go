package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/altar-12/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to decode the parameters: %v", err))
		return
	}
	currentTime := time.Now().UTC()
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to add user record: %v", err))
		return
	}
	respondWithJSON(w, created, ConvertDBUserToUser(user))
}

/*
func (apiCfg *apiConfig) handlerFetchUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := apiCfg.DB.FetchUsers(r.Context())
	if err != nil {
		respondWithError(w, internalServerError, fmt.Sprintf("Unable to fetch user records: %v", err))
		return
	}
	users := []User{}
	for _, user := range dbUsers {
		users = append(users, ConvertDBUserToUser(user))
	}
	respondWithJSON(w, ok, users)
}
*/

func (apiCfg *apiConfig) handlerFetchUserWithApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, ok, ConvertDBUserToUser(user))
}

func (apiCfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")
	ID, err := uuid.Parse(rawID)
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Invalid ID specified, %v", err))
		return
	}
	user, err := apiCfg.DB.DeleteUser(r.Context(), ID)
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to delete the user record: %v", err))
		return
	}
	respondWithJSON(w, ok, ConvertDBUserToUser(user))
}
