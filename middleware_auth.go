package main

import (
	"fmt"
	"net/http"

	auth "github.com/altar-12/rssagg/internal"
	"github.com/altar-12/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) authenticate(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, badRequest, fmt.Sprintf("Error while fetching API key: %v", err))
			return
		}
		user, err := apiCfg.DB.FetchUserWithApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, badRequest, fmt.Sprintf("Error while fetching the user %v", err))
			return
		}
		handler(w, r, user)
	}

}
