package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"fmt"
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
	if (err != nil || res == nil) {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(res.Key), []byte(resource.Key))

	if bcryptErr != nil {
		fmt.Println(bcryptErr, res, resource)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, signErr := createJWTToken(
		resource.ID,
		ctx.JWTSecret,
		res.Permissions,
	)
	if signErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jData, marshalErr := json.Marshal(&DispatchResponse{
		IdentityToken: token,
		IAT: uint64(time.Now().Unix()),
	})

	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write(jData)
}

type ValidationRequest struct {
	ClientKey string `json:"client_key"`
	IdentityToken string `json:"identity_token"`
}

type ValidationResponse struct {
	Valid bool `json:"valid"`
}

func (ctx *RequestContext) validate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

    request := &ValidationRequest{}
	err := decoder.Decode(&request)
		
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify that the client is a valid one
	remoteAddr := parseRemote(r.RemoteAddr)
	client, findErr := ctx.Creds.findClient(remoteAddr)

	if findErr != nil || client.Key != request.ClientKey {
		fmt.Println("request key mismatch")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// match the client IP with the incoming IP.
	claims, err := parseJWT(request.IdentityToken, ctx.JWTSecret)
	if err != nil || !matchPermissions(
		client.Permissions,
		claims.Permissions,
	) {
		fmt.Println("permission mismatch")

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(&ValidationResponse{
		Valid: true,
	})
}

type RegisterResourceRequest struct {
	ClientKey string `json:"client_key"`
	ResourceID string `json:"resource_id"`
	ResourceKey string `json:"resource_key"`
	Permissions []string `json:"permissions"`
}

type RegisterResourceResponse struct {
	Valid bool `json:"valid"`
	ResourceKey string `json:"resource_key"`
}

func (ctx *RequestContext) registerResource(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// @@ verify that all fields are present in the request.
	request := &RegisterResourceRequest{}

	err := decoder.Decode(&request)
	
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	remoteAddr := parseRemote(r.RemoteAddr)

	client, err := ctx.Creds.findClient(remoteAddr)
	if err != nil || len(client.IPAddress) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if client.Key != request.ClientKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	var resourceKey string
	if len(request.ResourceKey) == 0 {
		key := uuid.New()
		resourceKey = key.String()
	} else {
		resourceKey = request.ResourceKey
	}

	keyHash, hashErr := bcrypt.GenerateFromPassword([]byte(resourceKey), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println(hashErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resource := &Resource{
		ID: request.ResourceID,
		Key: string(keyHash),
		Permissions: strings.Join(request.Permissions[:], ","),
	}

	insertErr := ctx.Creds.insertUpdateResource(resource)
	if insertErr != nil {
		fmt.Println(insertErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&RegisterResourceResponse{
		Valid: true,
		ResourceKey: resourceKey,
	})
}

type RegisterClientRequest struct {
	Permissions []string `json:"permissions"`
}

type RegisterClientResponse struct {
	Key string `json:"key"`
}

func (ctx *RequestContext) registerClient(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	request := &RegisterClientRequest{}
	err := decoder.Decode(&request)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := uuid.New()
	resourceKey := key.String()

	client := &Client{
		Key: resourceKey,
		Permissions: strings.Join(request.Permissions[:], ","),
		IPAddress: parseRemote(r.RemoteAddr),
	}

	insertErr := ctx.Creds.insertUpdateClient(client)
	if insertErr != nil {
		fmt.Println(insertErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&RegisterClientResponse{
		Key: resourceKey,
	})
}