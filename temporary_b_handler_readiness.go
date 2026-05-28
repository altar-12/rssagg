package main

import "net/http"

func handlerReadinessa(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
