package main

import (
	"database/sql"
)

type Client struct {
	Key string
	Permissions string
	IPAddress string
}

func (creds *DbCreds) createClientTable() error{
	db := CreateSession(creds)
	defer db.Close()

	sqlQuery := `CREATE TABLE client (
		KEY CHAR(36) PRIMARY KEY UNIQUE NOT NULL,
		PERMISSIONS TEXT NOT NULL,
		IPADDRESS TEXT NOT NULL
	)`
	_, err := db.Exec(sqlQuery)

	return err
}


func (creds *DbCreds) findClient(key string) (*Client, error) {
	db := CreateSession(creds)
	defer db.Close()

	qry := `select * from client where key = $1`
	row := db.QueryRow(qry, key)

	c := &Client{}

	if err := row.Scan(&c.Key, &c.Permissions, &c.IPAddress); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return c, nil
}


func (creds *DbCreds) insertUpdateClient(c *Client) error {
	db := CreateSession(creds)
	defer db.Close()

	sqlQuery := `insert into client VALUES ($1, $2, $3) ON DUPLICATE KEY UPDATE PERMISSIONS=$2, IPADDRESS=$3`
	_, err := db.Exec(sqlQuery, c.Key, c.Permissions, c.IPAddress)

	return err
}