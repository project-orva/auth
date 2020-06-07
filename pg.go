package pg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // comment justifing it.. lol
	"github.com/truauth/truauth/pkg/settings"
)

type DbCreds struct {
	Host     string
	User     string
	Password string
	Dbname   string
}

func CreateSession(creds *DbCreds) *sql.DB {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", creds.Host, creds.User, creds.Password, creds.Dbname)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err.Error())
	}

	return db
}

func FetchCreds() *DbCreds {
	creds := &DbCreds{}
	settings.Init(creds, "aws-rds")

	return creds
}