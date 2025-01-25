package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	internalServerError = 500
	ok                  = 200
	badRequest          = 400
	created             = 201
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		fmt.Printf("Responding with code %v and message %v\n", code, msg)
	}
	type response struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, response{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("failed to convert the value to JSON format %v\n", payload)
		w.WriteHeader(internalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
