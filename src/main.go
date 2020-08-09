package main

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/joho/godotenv"
)

func heath(w http.ResponseWriter, r *http.Request) {
	resp := struct{
		Healty bool
	}{true}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(resp)
}

type RequestContext struct {
	JWTSecret string
	Creds *DbCreds
}

func main() {
	godotenv.Load("../.env")

	ctx := &RequestContext{
		JWTSecret: validateEnvVar("JWT_SECRET"),
		Creds: &DbCreds{
			Host: validateEnvVar("PG_HOST"),
			User: validateEnvVar("PG_USER"),
			Password: validateEnvVar("PG_PASSWORD"),
			Dbname: validateEnvVar("PG_DBNAME"),
		},
	}

	ctx.Creds.createClientTable()
	ctx.Creds.createResourceTable()	
	ctx.Creds.createIdentityTable()

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(heath))
	mux.Handle("/dispatch", handlePost(http.HandlerFunc(ctx.dispatch)))
	mux.Handle("/validate", handlePost(http.HandlerFunc(ctx.validate)))
	mux.Handle("/register-resource", handlePost(http.HandlerFunc(ctx.registerResource)))
	mux.Handle("/register-client", handlePost(http.HandlerFunc(ctx.registerClient)))
	
	log.Fatal(http.ListenAndServe(":5258", mux))
}