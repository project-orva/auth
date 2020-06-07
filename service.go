package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

struct DispatchResponse {
	IdentityToken string `json:"identity_token"`
	IAT uint64 `json:"iat"`
	EAT uint64 `json:"eat"`
}

func dispatch(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(req.Body)

    resource := &Resource{}
	err := decoder.Decode(&resource)
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate resource
	res, err := findResource(resource.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(res.Key), []byte(resource.Key))
	// output if correct
	if bcryptErr != nil {

	}

	json.NewEncoder(w).Encode(&DispatchResponse{
		IdentityToken: 0,
		IAT: 0
	})
}

func validate(w http.ResponseWriter, r *http.Request) {

}

func registerResource(w http.ResponseWriter, r *http.Request) {

}

func registerClient(w http.ResponseWriter, r *http.Request) {
	
}