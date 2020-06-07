package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
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

	token, jwtErr := jwt.Parse(res.Key, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return ctx.JWTSecret, nil
	})

	if jwtErr != nil {
		w.WriteHeader(http.StatusInternalServerError)		
	}
	
	signedStr, signErr := token.SigningString()

	if signErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(&DispatchResponse{
		IdentityToken: signedStr,
		IAT: uint64(time.Now().Unix()),
	})
}

func  (ctx *RequestContext) validate(w http.ResponseWriter, r *http.Request) {
	
}

func registerResource(w http.ResponseWriter, r *http.Request) {

}

func registerClient(w http.ResponseWriter, r *http.Request) {
	
}