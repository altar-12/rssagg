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

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to decode the request body: %v", err))
		return
	}
	currentTime := time.Now().UTC()
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to create a feed_follow record: %v", err))
		return
	}
	respondWithJSON(w, created, ConvertDBFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerFetchUserFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, internalServerError, fmt.Sprintf("Unable to fetch user feeds: %v", err))
		return
	}
	respondWithJSON(w, ok, ConvertDBFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDString := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDString)
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to parse the ID: %v", feedFollowID))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to delete the feed follow record: %v", err))
		return
	}
	respondWithJSON(w, ok, struct {
		Success bool `json:"success"`
	}{Success: true})
}
