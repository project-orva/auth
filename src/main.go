package main

import (
	"os"
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
	err := godotenv.Load("../.env")
	if err != nil {
		panic("env isn't being set correctly")
	}

	ctx := &RequestContext{
		JWTSecret: os.Getenv("JWT_SECRET"),
		Creds: &DbCreds{
			Host: os.Getenv("PG_HOST"),
			User: os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			Dbname: os.Getenv("PG_DBNAME"),
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