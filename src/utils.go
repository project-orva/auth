package main

import (
	"strings"
	"github.com/dgrijalva/jwt-go"
	"errors"

	"fmt"
)

type JtwClaim struct {
	UID string `json:"uid"`
	Permissions string `json:"permissions"`
	jwt.StandardClaims
}

func createJWTToken(id string, secret string, permissions string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JtwClaim{
		id,
		permissions,
		jwt.StandardClaims{
			Issuer:    "standard_auth",
		},
	})

	return token.SignedString([]byte(secret))
}

func isJWTValid(signedToken string, secret string) bool {
	token, _ := jwt.ParseWithClaims(signedToken, &JtwClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return token.Valid
}

var InvalidTokenErr = errors.New("jtw: invalid jtw token")

func parseJWT(signedToken string, secret string) (*JtwClaim, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JtwClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(*JtwClaim)
	
	fmt.Println(token.Valid, err)
	if !token.Valid || !ok {
		return nil, InvalidTokenErr
	}
	
	fmt.Println(claims)
	return claims, nil
}

func parseRemote(address string) string {
	if strings.Contains(address, "[::1]") {
		return "local"
	}

	return strings.Split(address, ":")[0]
}

func Includes(s []string, current string) bool {
	for _, a := range s {
		if a == current {
			return true
		}
	}

	return false
}

func matchPermissions(client string, resource string) bool {
	clientPermissions := strings.Split(client, ",") 
	resourcePermissions := strings.Split(resource, ",")

	for _, cp := range clientPermissions {
		if !Includes(resourcePermissions, cp) {
			return false
		}
	}

	return true
}