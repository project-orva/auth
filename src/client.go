package main

type Client struct {
	Key string
	IPAddress string
	Permissions string
}

func (creds *DbCreds) createClientTable() error{
	db := CreateSession(creds)
	defer db.Close()

	sqlQuery := `CREATE TABLE IF NOT EXISTS client (
		IPADDRESS TEXT PRIMARY KEY NOT NULL,
		KEY CHAR(36) UNIQUE NOT NULL,
		PERMISSIONS TEXT NOT NULL
	)`
	_, err := db.Exec(sqlQuery)

	if err != nil {
		panic(err)
	}

	return err
}


func (creds *DbCreds) findClient(ip string) (*Client, error) {
	db := CreateSession(creds)
	defer db.Close()

	qry := `select * from client where ipaddress = $1`
	row := db.QueryRow(qry, ip)

	c := &Client{}

	if err := row.Scan(&c.IPAddress, &c.Key, &c.Permissions); err != nil {
		return nil, err
	}

	return c, nil
}


func (creds *DbCreds) insertUpdateClient(c *Client) error {
	db := CreateSession(creds)
	defer db.Close()

	sqlQuery := `
	INSERT INTO CLIENT (
		KEY,
		PERMISSIONS,
		IPADDRESS
	) VALUES (
		$1, $2, $3
	) ON CONFLICT (IPADDRESS) DO UPDATE 
		SET PERMISSIONS = $2,
			KEY = $1
	`
	_, err := db.Exec(sqlQuery, c.Key, c.Permissions, c.IPAddress)

	return err
}