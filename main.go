package main

import (
	// "fmt"
	"log"
	"encoding/json"
	"net/http"
)

func heath(w http.ResponseWriter, r *http.Request) {
	resp := struct{
		Healty bool
	}{true}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(heath))
	mux.Handle("/dispatch", handlePost(http.HandlerFunc(dispatch)))
	mux.Handle("/validate", handlePost(http.HandlerFunc(validate)))
	mux.Handle("/register-resource", handlePost(http.HandlerFunc(registerResource)))
	mux.Handle("/register-client", handlePost(http.HandlerFunc(registerClient)))
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}