package main

import (
	"strings"
	"github.com/dgrijalva/jwt-go"
)

func createJWTToken(id string, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"uid": id,
	}

	return token.SignedString([]byte(secret))
}

func isJWTValid(signedToken string, secret string) bool {
	claims := jwt.MapClaims{} 
	token, _ := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return token.Valid
}

func parseRemote(address string) string {
	if strings.Contains(address, "[::1]") {
		return "local"
	}

	return address
}