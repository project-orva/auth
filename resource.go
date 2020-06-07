package main

struct Resource {
	ID string `json:"resource_id"`
	Key string `json:"resource_key"`
}

func (creds *DbCreds) createResourceTable() error{
	b := pgdb.CreateSession(req.Creds)
	defer db.Close()

	sqlQuery := `CREATE TABLE resource (
		ID CHAR(36) PRIMARY KEY UNIQUE NOT NULL,
		KEY TEXT NOT NULL,
	)`
	_, err := db.Exec(sqlQuery)

	return err
}


func (creds *DbCreds) findResource(id string) (*Resource, error) {
	db := CreateSession(creds)
	defer db.Close()

	qry := `select * from resource where id=$1`
	row := db.QueryRow(qry, id)

	c := &Client{}

	if err := row.Scan(&c.ID, &c.Key, &c.Type); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return d, nil
}


func (creds *DbCreds) insertUpdateResource(id string) (*Client, error) {
	b := pgdb.CreateSession(req.Creds)
	defer db.Close()

	sqlQuery := `insert into resource VALUES ($1, $2, $3) ON DUPLICATE KEY UPDATE KEY=$2, TYPE=$3`
	_, err := db.Exec(sqlQuery, c.Key, c.Permissions, c.IPAddress)

	return err
}