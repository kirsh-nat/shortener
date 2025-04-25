package app

import "database/sql"

func DBConnect(ps string) *sql.DB {
	DB, err := sql.Open("pgx", ps)
	if err != nil {
		Sugar.Fatalw(err.Error(), "event", err)
		return nil
	}

	return DB
}
