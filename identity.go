package main

struct Identity {
	Token string,
	Issued uint64,
}


func (creds *DbCreds) createIdentityTable() error{
	b := pgdb.CreateSession(req.Creds)
	defer db.Close()

	sqlQuery := `CREATE TABLE identity (
		TOKEN CHAR(36) PRIMARY KEY UNIQUE NOT NULL,
		ISSUED BIGINT NOT NULL,
	)`
	_, err := db.Exec(sqlQuery)

	return err
}


func (creds *DbCreds) findIdentity(token string) (*Client, error) {
	db := CreateSession(creds)
	defer db.Close()

	qry := `select * from identity where token = $1`
	row := db.QueryRow(qry, id)

	c := &Client{}

	if err := row.Scan(&c.Token, &c.Issued); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return d, nil
}


func (creds *DbCreds) deleteIndentity(token string) (error) {
	db := CreateSession(creds)

	qry := `DELETE FROM identity where token = $1`
	_, err := db.Exec(sqlQuery, token)

	return err
}