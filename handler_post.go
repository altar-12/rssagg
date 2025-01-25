package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/altar-12/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerFetchUserPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Invalid page number passed: %v", err))
		return
	}
	entriesPerPage, err := strconv.Atoi(r.URL.Query().Get("entries_per_page"))
	if err != nil {
		respondWithError(w, badRequest, fmt.Sprintf("Invalid 'entries per page' passed: %v", err))
		return
	}
	skipCount := (page - 1) * entriesPerPage
	rawPosts, err := apiCfg.DB.FetchUserPosts(r.Context(), database.FetchUserPostsParams{
		UserID: user.ID,
		Offset: int32(skipCount),
		Limit:  int32(entriesPerPage),
	})
	if err != nil {
		respondWithError(w, internalServerError, fmt.Sprintf("Error fetching posts from database: %v", err))
		return
	}
	respondWithJSON(w, ok, ConvertDBPostsToPosts(rawPosts))
}
