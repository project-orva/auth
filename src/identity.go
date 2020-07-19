package main

type Identity struct {
	Token string
	Issued uint64
}

func (creds *DbCreds) createIdentityTable() error{
	db := CreateSession(creds)
	defer db.Close()

	sqlQuery := `CREATE TABLE IF NOT EXISTS identity (
		TOKEN CHAR(36) PRIMARY KEY UNIQUE NOT NULL,
		ISSUED BIGINT NOT NULL
	)`
	_, err := db.Exec(sqlQuery)
	if err != nil {
		panic(err)
	}

	return err
}


func (creds *DbCreds) findIdentity(token string) (*Identity, error) {
	db := CreateSession(creds)
	defer db.Close()

	qry := `select * from identity where token = $1`
	row := db.QueryRow(qry, token)

	identity := &Identity{}
	if err := row.Scan(&identity.Token, &identity.Issued); err != nil {
		return nil, err
	}

	return identity, nil
}


func (creds *DbCreds) deleteIndentity(token string) (error) {
	db := CreateSession(creds)

	qry := `DELETE FROM identity where token = $1`
	_, err := db.Exec(qry, token)

	return err
}