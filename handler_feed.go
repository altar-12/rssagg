package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/altar-12/rssagg/internal/database"
	"github.com/google/uuid"
)

// this is a dummy message, nothing else to mention here.
func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Unable to decode request body: %v", err))
		return
	}
	currentTime := time.Now().UTC()
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, internalServerError, fmt.Sprintf("Unable to create a feed record: %v", err))
		return
	}
	respondWithJSON(w, created, ConvertDBFeedToFeed(feed))
}

/*
func (apiCfg *apiConfig) handlerGetUserFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	rawFeeds, err := apiCfg.DB.FetchUserFeeds(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, internalServerError, fmt.Sprintf("Unable to fetch user feeds: %v", err))
		return
	}
	feeds := []Feed{}
	for _, feed := range rawFeeds {
		feeds = append(feeds, ConvertDBFeedToFeed(feed))
	}
	respondWithJSON(w, ok, feeds)
}
*/

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	rawFeeds, err := apiCfg.DB.FetchFeeds(r.Context())
	if err != nil {
		respondWithError(w, internalServerError, fmt.Sprintf("Unable to fetch feeds: %v", err))
	}
	respondWithJSON(w, ok, ConvertDBFeedsToFeeds(rawFeeds))
}
