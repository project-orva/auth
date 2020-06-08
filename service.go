package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
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
	client, findErr := ctx.Creds.findClient(request.ClientKey)
	if findErr != nil || len(client.IPAddress) == 0 {
		w.WriteHeader(http.StatusUnauthorized)		
	}
	
	// @@@ match the client IP with the incoming IP.

	json.NewEncoder(w).Encode(&ValidationResponse{
		// verify the identity token
		Valid: isJWTValid(request.IdentityToken, ctx.JWTSecret),
	})
}

type RegisterResourceRequest struct {
	ClientKey string `json:"client_key"`
	ResourceID string `json:"resource_id"`
	ResourceKey string `json:"resource_key"`
}

type RegisterResourceResponse struct {
	Valid bool `json:"valid"`
	ResourceKey string `json:"resource_key"`
}

func (ctx *RequestContext) registerResource(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	request := &RegisterResourceRequest{}
	err := decoder.Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	client, err := ctx.Creds.findClient(request.ClientKey)
	if err != nil || len(client.IPAddress) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
	}
	// @@@ match the client IP with the incoming IP.
	
	var resourceKey string
	if len(request.ResourceKey) == 0 {
		key := uuid.New()
		resourceKey = key.String()
	} else {
		resourceKey = request.ResourceKey
	}
	
	resource := &Resource{
		ID: request.ResourceID,
		Key: request.ResourceKey,
	}

	insertErr := ctx.Creds.insertUpdateResource(resource)
	if insertErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(&RegisterResourceResponse{
		Valid: true,
		ResourceKey: resourceKey,
	})
}

func registerClient(w http.ResponseWriter, r *http.Request) {
	
}