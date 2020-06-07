package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DispatchResponse struct {
	IdentityToken string `json:"identity_token"`
	IAT uint64 `json:"iat"`
}

func (ctx *RequestContext) dispatch(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

    resource := &Resource{}
	err := decoder.Decode(&resource)
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate resource
	res, err := ctx.Creds.findResource(resource.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(res.Key), []byte(resource.Key))

	if bcryptErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	token, signErr := createJWTToken(resource.ID, ctx.JWTSecret)
	if signErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(&DispatchResponse{
		IdentityToken: token,
		IAT: uint64(time.Now().Unix()),
	})
}

type ValidationRequest struct {
	ClientKey string `json:"client_key"`
	IdentityToken string `json:"identity_token"`
}

type ValidationResponse struct {
	Valid bool `json:"valid"`
}

func  (ctx *RequestContext) validate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

    request := &ValidationRequest{}
	err := decoder.Decode(&request)
		
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify that the client is a valid one
	_, findErr := ctx.Creds.findClient(request.ClientKey)
	if findErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// verify the identity token
	json.NewEncoder(w).Encode(&ValidationResponse{
		Valid: isJWTValid(request.IdentityToken, ctx.JWTSecret),
	})
}

func registerResource(w http.ResponseWriter, r *http.Request) {

}

func registerClient(w http.ResponseWriter, r *http.Request) {
	
}