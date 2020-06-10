package main

import (
	"database/sql"
)

type Resource struct {
	ID string `json:"resource_id"`
	Key string `json:"resource_key"`
}

func (creds *DbCreds) createResourceTable() error{
	db := CreateSession(creds)
	defer db.Close()

	sqlQuery := `CREATE TABLE IF NOT EXISTS resource (
		ID CHAR(36) PRIMARY KEY UNIQUE NOT NULL,
		KEY TEXT NOT NULL
	)`
	_, err := db.Exec(sqlQuery)
	if err != nil {
		panic(err)
	}

	return err
}


func (creds *DbCreds) findResource(id string) (*Resource, error) {
	db := CreateSession(creds)
	defer db.Close()

	qry := `select * from resource where id=$1`
	row := db.QueryRow(qry, id)

	resource := &Resource{}
	if err := row.Scan(&resource.ID, &resource.Key); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return resource, nil
}


func (creds *DbCreds) insertUpdateResource(resource *Resource) error {
	db := CreateSession(creds)
	defer db.Close()

	sqlQuery := `INSERT INTO RESOURCE (
		ID,
		KEY
	) VALUES (
		$1, $2
	) ON CONFLICT (ID) DO UPDATE
		SET KEY = $2
	`
	_, err := db.Exec(sqlQuery, resource.ID, resource.Key)

	return err
}